# ADR-0002: Migrate GSLB API group to vendor-neutral `k8gb.io`

## Status

Accepted

## Date

2026-03-09

## Context

The current GSLB API group `k8gb.absa.oss/v1beta1` contains a vendor-specific domain.
As the project evolves in a neutral governance model, API naming should also be vendor-neutral.

Keeping a vendor-specific group as the long-term public API creates friction for adoption, documentation, and ecosystem integrations.
At the same time, existing users already run production objects with the legacy group and require a safe migration path.

## Migration Options Considered

### Option A: Hard cutover to `k8gb.io` (no compatibility layer)

- Introduce only `k8gb.io/v1beta1` and stop accepting `k8gb.absa.oss/v1beta1`.
- Pros: simplest implementation and lowest ongoing maintenance.
- Cons: immediate breaking change for existing clusters and automation.

### Option B: Dual-write and dual-reconcile both API groups

- Keep both groups fully writable and attempt to keep them in sync in both directions.
- Pros: low short-term friction for users who stay on legacy manifests.
- Cons: two sources of truth, conflict resolution complexity, and long-term ambiguity.

### Option C: Automatic one-way migration bridge to canonical API

- Define `k8gb.io/v1beta1` as the canonical API.
- Keep legacy CRD support during transition.
- Automatically migrate legacy objects to canonical objects.
- Pros: backward compatibility with a single source of truth.
- Cons: larger upgrade blast radius in clusters with many legacy objects.

### Option D: Controlled one-way migration bridge to canonical API (Selected)

- Define `k8gb.io/v1beta1` as the canonical API.
- Keep legacy CRD support during transition.
- Keep legacy reconcile behavior active before migration to avoid behavioral regressions.
- Once `k8gb.io/migration-requested=true` is set, pause legacy runtime reconcile for that object and let migration controller own the transition to avoid dual-writer conflicts.
- Trigger migration per object with explicit label `k8gb.io/migration-requested=true`.
- After migration, treat legacy object as compatibility/read-only source and direct further edits to canonical object.
- Pros: backward compatibility with operator-controlled migration blast radius.
- Cons: temporary coexistence of legacy and canonical logic during transition.

## Decision

We introduce `k8gb.io/v1beta1` as the canonical GSLB API group and migrate all first-party manifests, chart defaults, and controller reconciliation to this group.

Compatibility with `k8gb.absa.oss/v1beta1` is preserved through a migration controller:

- Legacy reconcile stays active before migration to preserve existing behavior.
- When migration is requested, legacy runtime reconcile is paused for that object to avoid concurrent writes with canonical reconcile.
- A deprecation warning event is emitted for legacy objects to steer users to `k8gb.io`.
- Migration runs only when the legacy object has label `k8gb.io/migration-requested=true`.
- Migration logic remains one-way and unchanged in semantics:
  - create canonical `k8gb.io/v1beta1` object (same name/namespace),
  - perform embedded ingress to `resourceRef` conversion,
  - run ownerReference/finalizer safety cleanup,
  - mark legacy object with `k8gb.io/migrated-to-k8gb-io=true`.
- After migration, canonical object is the only write target for new feature work.

Helm continues to support installation of legacy CRDs during the transition window.

### Reconcile Decision Table

| Legacy labels on `k8gb.absa.oss` object | Legacy reconcile mode | Migration action | Notes |
|---|---|---|---|
| no migration labels | full legacy reconcile + deprecation warning event | none | no behavioral downgrade before migration |
| `k8gb.io/migration-requested=true`, not migrated | migration transition mode (legacy runtime paused) | run one-way migration | canonical object becomes write target; avoids dual writers |
| `k8gb.io/migrated-to-k8gb-io=true` | compatibility/read-only + deprecation warning event | none | legacy kept for transition visibility |
| both labels set | compatibility/read-only | none (optional request-label cleanup) | avoids repeated migration attempts |

## Consequences

### Positive

- Public API group is vendor-neutral and aligned with project identity.
- New integrations and examples consistently target `k8gb.io/v1beta1`.
- Existing clusters can migrate without immediate breakage.
- Operators can migrate incrementally instead of cluster-wide big-bang migration.

### Negative

- Additional controller logic is required during the transition period.
- Legacy and canonical logic coexist during the transition period.
- Operators must understand migration trigger labels and write-target handoff.
- Migration events and labels add short-term operational noise during rollout.

### Neutral

- Both API groups may coexist temporarily.
- Documentation must clearly separate canonical and legacy usage until deprecation is complete.
- A future ADR or release note must define when legacy CRD installation defaults change.

## Alternatives

- Keep `k8gb.absa.oss` as the long-term API group.
  - Rejected: does not satisfy vendor-neutral API expectations.
- Rename API group and drop legacy support immediately.
  - Rejected: unacceptable breaking change for current users (Option A).
- Support both groups indefinitely as equal first-class writable APIs.
  - Rejected: duplicates maintenance burden and keeps ambiguity (Option B).
- Keep automatic migration for all legacy objects on reconcile.
  - Rejected: larger operational blast radius during upgrades compared to controlled migration.

## References

- `api/v1beta1io` package and CRD definitions
- `controllers/gslb_migration_controller.go`
- `main.go` legacy migration controller registration
- `chart/k8gb/crd/k8gb.io_gslbs.yaml`
- `chart/k8gb/crd/k8gb.absa.oss_gslbs.yaml`
- `docs/migration_acceptance.md` manual migration acceptance checklist (pre-merge and post-merge)
