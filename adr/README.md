# Architectural Decision Records (ADRs)

This directory contains Architectural Decision Records (ADRs) for the k8gb project. ADRs are documents that capture important architectural decisions made during the project's development, along with their context, consequences, and rationale.

## ADR Index

- [ADR-0000: Template](0000-template.md) - Template for new ADRs
- [ADR-0001: Deprecate Configuration of GSLB resources via annotations](0001-deprecate-configuration-of-gslb-resources-via-annotations.md) - GSLB configuration

## What are ADRs?

Architectural Decision Records are short text documents that capture a single architecture decision. They help teams:

- Understand why certain decisions were made
- Avoid repeating discussions about already-solved problems
- Provide context for future architectural changes
- Onboard new team members to the project's architecture

## ADR Format

Each ADR follows this structure:

- **ADR-0001**: Sequential number for the decision
- **Title**: Short, descriptive title
- **Status**: Current status (Proposed, Accepted, Deprecated, Superseded)
- **Date**: When the decision was made
- **Context**: The situation that led to the decision
- **Decision**: What was decided
- **Consequences**: What happens as a result, both positive and negative
- **Alternatives**: Other options that were considered
- **References**: Links to relevant documentation, discussions, or code

## Creating a New ADR

1. Copy the template from `0000-template.md`
2. Rename it to the next sequential number (e.g., `0001-example-decision.md`)
3. Fill in all sections
4. Submit a pull request for review

## ADR Status Values

- **Proposed**: Decision is under discussion and review
- **Accepted**: Decision has been made and implemented
- **Deprecated**: Decision is no longer relevant
- **Superseded**: Replaced by a newer ADR

## When to Create an ADR

Create an ADR when making decisions about:

- Architecture patterns and designs
- Technology choices
- API design decisions
- Data model changes
- Integration approaches
- Performance optimizations
- Security implementations

## Examples

- ADR-0001: Use Go modules for dependency management
- ADR-0002: Implement controller-runtime pattern for Kubernetes operators
- ADR-0003: Choose CoreDNS as the DNS provider interface

## References

- [ADR GitHub Repository](https://github.com/joelparkerhenderson/architecture_decision_record)
- [ADR Tools](https://adr.github.io/)
- [Documenting Architecture Decisions](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions)
