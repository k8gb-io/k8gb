# ADR-0001: Deprecate configuration of GSLB resources via annotations

## Status

Accepted

## Date

2025-08-29

## Context

When k8gb only supported ingress integration using Ingress resources, there were two ways the strategy for an application could be configured: creating a GSLB resouce or annotating an Ingress resource

With the introduction of new ingress integrations, it triggered the question whether both configuration methods should be supported, or only one of them

## Decision

We will only support configuration using a GSLB resource. Ingress annotations will be deprecated and removed in k8gb v0.17

## Consequences

### Positive

- Smaller code base, easier to maintain
- Removed source of bugs since it was not clear which configuration method is the source of truth:
  - If GSLB is the source of truth and someone tweaks an Ingress annotation, should that change be pushed back into the GSLB or should the annotation be overwritten?
  - Conversely, if a GSLB was created from annotations, what takes precedence later — editing the GSLB or editing the annotations?

### Negative

- Configuring k8gb using annotations is very user friendly
