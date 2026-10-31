# Required pull-request checks

The default branch's [active ruleset](https://github.com/k8gb-io/k8gb/rules/26274974)
is the authoritative record of merge requirements. The desired configuration
including CI requirements is tracked in [master-ruleset.json](master-ruleset.json).

Normal CI runs on every pull request, including documentation changes. There are
no path filters or change-detection jobs. Branch push runs are limited to `master`
to avoid duplicating PR runs. Chainsaw also runs on tag pushes, and CodeQL runs
weekly.

Most existing job checks are required directly. Only Terratest, Chainsaw, and
multi-cluster testing have final gates, so Kubernetes versions do not appear in
the required check names. Each gate uses `always()`, depends on its matrix job,
and fails if that job failed or was cancelled. It accepts intentional skips.
CodeQL's fixed Go analysis job and its separate security-results check are both
required.

| Required check | What it covers |
| --- | --- |
| `DCO` | Commit sign-offs, reported by the DCO GitHub App |
| `Inspect packages` | Lint, license headers, and Go tests |
| `build` | FOSSA license scanning |
| `CI / Terratest` | The selected Kubernetes test matrix |
| `CI / Chainsaw` | The selected Kubernetes test matrix |
| `upgrade-testing` | Upgrade testing |
| `Analyze (go)` | CodeQL analysis job completion |
| `CodeQL` | Security results reported by the GitHub code scanning app |
| `Check values.yaml and its schema in PR` | Helm values/schema validation |
| `terratest-istio` | Istio tests when the `istio` label is present |
| `CI / Multi-cluster` | Multi-cluster tests when the `heavy-tests` label is present |

The existing `renovate` label exclusions for Terratest and upgrade testing remain
in place. Label-selected Istio and multi-cluster tests rerun on subsequent commits,
not only when a label is first applied. The three gate jobs have no token
permissions and do not check out or run repository code.

## Rollout

1. Merge the workflow changes after maintainer review and successful CI. The
   active ruleset already requires DCO, one approval, dismissal of stale approvals,
   an up-to-date branch, linear history, and protection against force pushes and
   deletion. There are no bypass actors.
2. Compare the live ruleset with the tracked configuration, preserving any
   intervening administrator changes. Then enable the additional required checks:

   ```sh
   gh api repos/k8gb-io/k8gb/rulesets/26274974 \
     --method PUT --input .github/master-ruleset.json
   gh api repos/k8gb-io/k8gb/rules/branches/master
   ```

3. Update existing PR branches from `master` so they report all required checks. Keep
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
