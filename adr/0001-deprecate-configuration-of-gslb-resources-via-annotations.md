# ADR-0001: Deprecate configuration of GSLB resources via annotations

## Status

Accepted (removal deferred)

## Date

2025-08-29

## Context

When k8gb only supported ingress integration using Ingress resources, there were two ways the strategy for an application could be configured: creating a GSLB resource or annotating an Ingress resource

With the introduction of new ingress integrations, it triggered the question whether both configuration methods should be supported, or only one of them

## Decision

We will only support configuration using a GSLB resource as the long-term interface. Ingress `k8gb.io/strategy` (and related) annotations are **deprecated**. Prefer an explicit `Gslb` with `spec.resourceRef` (see [Resource References](../docs/resource_ref.md)).

The original decision text targeted removal in k8gb **v0.17**. That removal did **not** ship: as of **v0.20+** the annotation path is still implemented (`controllers/handlers.go`) and emits a deprecation warning when used. Removal remains intended, but the calendar target is deferred to a future major/minor after maintainers announce a migration window — not a silent drop.

## Consequences

### Positive

- Smaller code base, easier to maintain (once removal lands)
- Removed source of bugs since it was not clear which configuration method is the source of truth:
  - If GSLB is the source of truth and someone tweaks an Ingress annotation, should that change be pushed back into the GSLB or should the annotation be overwritten?
  - Conversely, if a GSLB was created from annotations, what takes precedence later — editing the GSLB or editing the annotations?

### Negative

- Configuring k8gb using annotations is very user friendly
- Until removal ships, docs and demos that still show annotations must label them as **deprecated / legacy**

## Amendment (2026-08)

Documented the deferred removal after the v0.20 documentation audit (#2455). Code, ADR, and tutorials must agree: annotations still work, are deprecated, and must not be presented as the preferred path for new installs.
