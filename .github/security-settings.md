# Repository security settings

The desired settings for the OSPS Baseline Level 2 controls AC-04.01 and
BR-07.01 are tracked in [actions-permissions.json](actions-permissions.json) and
[secret-protection.json](secret-protection.json). They make the default
`GITHUB_TOKEN` read-only, disable its ability to create and approve pull requests,
and enable GitHub secret scanning and push protection.

These files are inputs to GitHub's admin API. Merging them does not apply the
settings automatically. The rollout below follows the same review-then-apply
approach as [required-checks.md](required-checks.md).

## Preserve automated pull requests

The Helm docs and offline changelog workflows open pull requests with
`CHANGELOG_GH_TOKEN`, independently of `GITHUB_TOKEN`. Their job tokens need only
`contents: read`; write access is supplied only to the pinned
`peter-evans/create-pull-request` step. Checkout does not persist its credentials.

Before merging the workflow changes, verify or rotate the existing
`CHANGELOG_GH_TOKEN` repository secret. Use a fine-grained personal access token
restricted to `k8gb-io/k8gb` with **Contents: read and write** and **Pull requests:
read and write**, approved for the organization where required. It needs no
administration or workflow-editing permission. The token's owner must be able to
push the automation branches and open pull requests in this repository. Keep it
valid and rotate it before expiry.

The presence of a secret name does not establish that its value is valid or has
the necessary permissions. GitHub does not expose its value for inspection.
Do not substitute `GITHUB_TOKEN`: it cannot open these PRs once the repository
setting is disabled. The other existing release credentials are managed
separately.

## Rollout

1. Verify the automation credential above, then merge the workflow and settings
   files after maintainer review and successful CI. Job-level write permissions
   in release, publishing, and other workflows remain explicit and override the
   new read-only default.
2. From a checkout of the merged `master`, use an authenticated maintainer account
   with repository administration access to inspect and save the current values:

   ```sh
   gh api repos/k8gb-io/k8gb/actions/permissions/workflow \
     > /tmp/k8gb-actions-permissions-before.json
   gh api repos/k8gb-io/k8gb \
     --jq '{security_and_analysis}' \
     > /tmp/k8gb-secret-protection-before.json
   ```

   Compare these with the reviewed files and preserve intervening administrator
   changes. An organization policy may prevent repository-level changes.
3. Apply the two reviewed configurations:

   ```sh
   gh api repos/k8gb-io/k8gb/actions/permissions/workflow \
     --method PUT --input .github/actions-permissions.json
   gh api repos/k8gb-io/k8gb \
     --method PATCH --input .github/secret-protection.json \
     --jq '{secret_scanning: .security_and_analysis.secret_scanning,
            push_protection: .security_and_analysis.secret_scanning_push_protection}'
   ```

   These are separate API operations. If either fails, inspect the current state
   before continuing; the other operation may already have succeeded. The
   security payload updates only secret scanning and push protection.
4. Read back both settings:

   ```sh
   gh api repos/k8gb-io/k8gb/actions/permissions/workflow
   gh api repos/k8gb-io/k8gb \
     --jq '{secret_scanning: .security_and_analysis.secret_scanning.status,
            push_protection: .security_and_analysis.secret_scanning_push_protection.status}'
   ```

   Expect `default_workflow_permissions: read`,
   `can_approve_pull_request_reviews: false`, and both security statuses `enabled`.
   Confirm the next Helm docs and offline changelog runs can open or update their
   automation PRs, and review the alerts from the initial secret scan privately.
5. Run a fresh OSPS Baseline evaluation with a scanner containing
   [the QA-03.01 fix](https://github.com/ossf/pvtr-github-repo-scanner/pull/529).
   A local admin-authenticated scan can verify settings that the Actions job token
   cannot observe. Publish through `publish-to-grc-store.yml` once that workflow's
   installed scanner includes the fix, respecting the ten-minute publishing
   interval documented in [required-checks.md](required-checks.md).

These settings address the two additional failures observed in the live
admin-authenticated scan. They do not resolve the four human-review controls;
their evidence remains in
[SECURITY.md](../SECURITY.md#osps-baseline-review-evidence).

## References

- [GitHub Actions workflow permissions](https://docs.github.com/en/rest/actions/permissions#set-default-workflow-permissions-for-a-repository)
- [Repository security settings API](https://docs.github.com/en/rest/repos/repos#update-a-repository)
- [create-pull-request token permissions](https://github.com/peter-evans/create-pull-request#token)
