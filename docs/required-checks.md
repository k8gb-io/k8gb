# Required pull-request checks

The default branch's [active ruleset](https://github.com/k8gb-io/k8gb/rules/26274974)
is the authoritative record of merge requirements. The desired configuration
including CI gates is tracked in [master-ruleset.json](../.github/master-ruleset.json).

Each pull-request workflow reports a stable `CI / …` check. Its final job uses
`always()`, depends on all other jobs in that workflow, and fails if any dependency
failed or was cancelled. It accepts intentional skips. This covers every matrix
entry without putting Kubernetes versions into the required check names.

Path filtering happens inside the workflow through the pinned
`fkirc/skip-duplicate-actions` action. It can reuse successful checks from the
commit history when only ignored files changed. The workflow itself still runs
on every pull request, so documentation changes receive a completed gate rather
than an indefinitely pending check. CodeQL also checks Go dependency changes.

| Required check | What it covers |
| --- | --- |
| `DCO` | Commit sign-offs, reported by the DCO GitHub App |
| `CI / Go` | Change detection, lint, license headers, and Go tests |
| `CI / FOSSA` | Change detection and license scanning |
| `CI / Terratest` | Change detection and the selected Kubernetes test matrix |
| `CI / Chainsaw` | Change detection and the selected Kubernetes test matrix |
| `CI / Upgrade` | Change detection and upgrade testing |
| `CI / CodeQL` | Change detection and CodeQL analysis |
| `CI / Helm schema` | Change detection and Helm values/schema validation |
| `CI / Istio` | Istio tests when the `istio` label is present |
| `CI / Multi-cluster` | Multi-cluster tests when the `heavy-tests` label is present |

The existing `renovate` label exclusions for Terratest and upgrade testing remain
in place. Label-selected Istio and multi-cluster tests rerun on subsequent commits,
not only when a label is first applied. The gate jobs have no token permissions
and do not check out or run repository code.

## Rollout

1. Merge the workflow changes after maintainer review and successful CI. The
   active ruleset already requires DCO, one approval, dismissal of stale approvals,
   an up-to-date branch, linear history, and protection against force pushes and
   deletion. There are no bypass actors.
2. Compare the live ruleset with the tracked configuration, preserving any
   intervening administrator changes. Then enable the CI gates:

   ```sh
   gh api repos/k8gb-io/k8gb/rulesets/26274974 \
     --method PUT --input .github/master-ruleset.json
   gh api repos/k8gb-io/k8gb/rules/branches/master
   ```

3. Update existing PR branches from `master` so they report the new gates. Keep
   the classic branch protection rule during verification; both sets of rules
   apply. The `gh-pages` rule is managed separately.
4. Publish a new OSPS Baseline result:

   ```sh
   gh workflow run publish-to-grc-store.yml --repo k8gb-io/k8gb --ref master
   ```

   Check the workflow artifact and [published result](https://grc.store/targets/k8gb-io/k8gb).
   The scanner inspects checks from the last merged PR and the gate dependencies
   in the default-branch workflows. Check names in settings alone do not establish
   coverage before those workflows are merged.

The publish workflow also runs weekly. Its target accepts one published result
per ten minutes; avoid dispatching repeated runs within that window.
