# Repository security settings

The desired repository settings are tracked in
[actions-permissions.json](actions-permissions.json) and
[secret-protection.json](secret-protection.json). They make the default
`GITHUB_TOKEN` read-only and enable GitHub secret scanning and push protection.

These files are inputs to GitHub's admin API. Merging them does not apply the
settings automatically. The rollout below follows the same review-then-apply
approach as [required-checks.md](required-checks.md).

## Preserve automated pull requests

The Helm docs and offline changelog workflows use the automatically generated
`GITHUB_TOKEN`, with explicit job-level `contents: write` and
`pull-requests: write` permissions. Those grants override the read-only repository
default. The job tokens are scoped to this repository and expire automatically
when their jobs finish; no manually rotated personal token is needed.

The configuration preserves the existing
`can_approve_pull_request_reviews: true` setting. GitHub combines permission to
create and approve pull requests in this repository setting, so disabling it
would also prevent these workflows from opening their PRs. Changing the default
token permissions to read-only does not require disabling this setting.

## Scanner expectations

Secret scanning and push protection address OSPS-BR-07.01. OSPS-AC-04.01 remains
outstanding: the scanner requires both a read-only default and the combined
PR-creation/approval setting to be disabled. Keeping the job-scoped tokens does
not satisfy that additional policy check, even with a read-only repository
default.

Resolving that policy requirement is separate work. An alternative would be a
dedicated, narrowly scoped GitHub App that issues short-lived installation
tokens, with its signing key protected appropriately. Replacing the job tokens
with a stored personal access token solely to satisfy the scanner introduces a
longer-lived credential and manual rotation duties.

## Rollout

1. Merge the settings files after maintainer review and successful CI. Job-level
   write permissions in release, publishing, and other workflows remain explicit
   and override the new read-only default.
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
   `can_approve_pull_request_reviews: true`, and both security statuses `enabled`.
   Confirm the next Helm docs and offline changelog runs can open or update their
   automation PRs, and review the alerts from the initial secret scan privately.
5. Run a fresh OSPS Baseline evaluation with a scanner containing
   [the QA-03.01 fix](https://github.com/ossf/pvtr-github-repo-scanner/pull/529).
   A local admin-authenticated scan can verify settings that the Actions job token
   cannot observe. Publish through `publish-to-grc-store.yml` once that workflow's
   installed scanner includes the fix, respecting the ten-minute publishing
   interval documented in [required-checks.md](required-checks.md).

The AC-04.01 policy check and the four human-review controls remain outstanding;
the review evidence remains in
[SECURITY.md](../SECURITY.md#osps-baseline-review-evidence).

## References

- [GitHub Actions workflow permissions](https://docs.github.com/en/rest/actions/permissions#set-default-workflow-permissions-for-a-repository)
- [Repository security settings API](https://docs.github.com/en/rest/repos/repos#update-a-repository)
- [GITHUB_TOKEN scope and lifetime](https://docs.github.com/en/actions/concepts/security/github_token)
