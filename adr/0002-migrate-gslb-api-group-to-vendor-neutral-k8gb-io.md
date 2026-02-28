# ADR-0002: Migrate GSLB API group to vendor-neutral `k8gb.io`

## Status

Accepted

## Date

2026-02-28

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

### Option C: One-way migration bridge to canonical API (Selected)

- Define `k8gb.io/v1beta1` as the canonical API.
- Keep legacy CRD support during transition.
- Migrate legacy objects to canonical objects and direct further edits to canonical objects.
- Pros: backward compatibility with a single source of truth.
- Cons: temporary migration controller and transition communication burden.

## Decision

We introduce `k8gb.io/v1beta1` as the canonical GSLB API group and migrate all first-party manifests, chart defaults, and controller reconciliation to this group.

Compatibility with `k8gb.absa.oss/v1beta1` is preserved through a migration controller:

- Legacy `Gslb` objects are converted to `k8gb.io/v1beta1` objects with the same name and namespace.
- Legacy objects are marked with `k8gb.io/migrated-to-k8gb-io=true`.
- A warning event is emitted to indicate that users must edit the `k8gb.io` object going forward.

Helm continues to support installation of legacy CRDs during the transition window.

## Consequences

### Positive

- Public API group is vendor-neutral and aligned with project identity.
- New integrations and examples consistently target `k8gb.io/v1beta1`.
- Existing clusters can migrate without immediate breakage.

### Negative

- Additional controller logic is required during the transition period.
- Operators must understand that legacy objects become migration sources, not long-term write targets.
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

## References

- `api/v1beta1io` package and CRD definitions
- `controllers/gslb_migration_controller.go`
- `main.go` legacy migration controller registration
- `chart/k8gb/crd/k8gb.io_gslbs.yaml`
- `chart/k8gb/crd/k8gb.absa.oss_gslbs.yaml`
