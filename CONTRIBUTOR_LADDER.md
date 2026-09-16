# k8gb Contributor Ladder

## Purpose and authority

This document describes a lightweight path for growing contribution and review responsibility in the [`k8gb-io/k8gb`](https://github.com/k8gb-io/k8gb) repository. It recognizes code, documentation, testing, issue investigation, review, release, and community work as valuable contributions.

This ladder complements [`GOVERNANCE.md`](GOVERNANCE.md) and [`CONTRIBUTING.md`](CONTRIBUTING.md). Where this guide conflicts with either document, those documents control. It does not change the Maintainer election process, the pull-request approval requirement, `CODEOWNERS`, or any GitHub permission, branch rule, release credential, or administrative access.

> **Role names and platform access are separate.** Recognition as a Contributor or Reviewer does not grant GitHub organization membership, repository permission, merge authority, release access, or administrative access. Any platform access is granted and removed separately by the people authorized to manage it.

## Ladder at a glance

| Role | Purpose | Scope of authority |
|---|---|---|
| **Contributor** | Participates in improving k8gb through code or non-code work. | May contribute through the normal project process. This role carries no repository permission. |
| **Reviewer** | A Contributor recognized by Maintainers for sound, constructive review in an identified area. | Provides review feedback and recommendations. This role does not satisfy the project’s required Maintainer approval and carries no GitHub permission. |
| **Maintainer** | A governance role responsible for the project’s technical health and community stewardship. | Defined by `CODEOWNERS` and appointed only under `GOVERNANCE.md`. Platform access is separately granted where needed. |

A person may contribute in more than one capacity. Role recognition is based on sustained contribution, judgment, and stewardship rather than on a fixed count of commits or pull requests.

## Shared expectations

Everyone participating under this guide must follow the [Code of Conduct](CODE_OF_CONDUCT.md), follow the contribution process in [`CONTRIBUTING.md`](CONTRIBUTING.md), and engage constructively with the community. Contributors should perform the testing, validation, documentation, and review follow-up appropriate to the change they submit. Security vulnerabilities must be reported through the confidential process in [`SECURITY.md`](SECURITY.md), not through ordinary issue triage.

## Contributor

For the purposes of this guide, a **Contributor** is a person who has made a meaningful contribution to k8gb. Anyone may contribute under [`CONTRIBUTING.md`](CONTRIBUTING.md); this recognition is not a prerequisite for participation and does not grant repository access.

Meaningful contributions may include code changes, documentation or website improvements, tests or test infrastructure, reproducible bug reports, issue investigation or triage, release or project-infrastructure work, review feedback, user support, community facilitation, or other sustained work that advances k8gb.

Contributors may participate in public project discussions, propose fixes and improvements, and be invited or assigned to work where repository settings permit. They are expected to explain their submissions, respond constructively to feedback, and respect the project’s quality, security, licensing, and DCO requirements.

## Reviewer

A **Reviewer** is a Contributor whom Maintainers recognize as a trusted source of constructive technical or documentation review in an identified area. The recognition may be scoped to an area such as core code, Helm charts, documentation, testing, integrations, or release engineering. It is a community-recognition role under this guide, not a Governance role or a GitHub access grant.

### Recognition criteria

Maintainers should consider sustained, high-quality evidence rather than a fixed numerical threshold. Relevant evidence includes meaningful contributions in the area, technically useful reviews of others’ work, familiarity with the applicable architecture and contribution standards, reliable participation over time, constructive collaboration, and willingness to help newer contributors.

A Maintainer may nominate a Contributor for Reviewer recognition by opening a GitHub issue titled `Reviewer nomination: @username`, describing the proposed scope and supporting contributions or reviews. The nominee must confirm their willingness to take on the role. Recognition requires approval from at least two existing Maintainers. The final decision and agreed scope must be recorded in the issue before it is closed. A Contributor may also ask a Maintainer to consider their readiness for this role.

### Responsibilities

Reviewers are expected to review work within their demonstrated area of knowledge, provide timely and respectful feedback, help maintain quality and testing standards, and escalate architectural, security, licensing, or governance questions to Maintainers when appropriate. They should help Contributors understand the contribution process without assuming Maintainer responsibilities.

### Boundaries and privileges

Reviewers may provide review feedback, approving GitHub reviews where the platform permits, and LGTM recommendations. Their review is valuable evidence for a change, but it is not a substitute for the project’s existing approval rule:

> **Reviewer recognition, a GitHub approval, or a Reviewer LGTM does not satisfy the requirement in `CONTRIBUTING.md` for an LGTM from at least one Maintainer listed in `CODEOWNERS`.**

Reviewer recognition grants no automatic GitHub permission. If a separate access decision is ever needed, it must be made, configured, and documented separately from this ladder.

## Maintainer

A **Maintainer** is the project governance role defined by GitHub handles in [`CODEOWNERS`](CODEOWNERS). Maintainers are responsible for the project’s technical direction, health, governance, and sustainability as described in [`GOVERNANCE.md`](GOVERNANCE.md).

### Appointment

This ladder does not create an alternative path to Maintainer status. Under [`GOVERNANCE.md`](GOVERNANCE.md), an existing Maintainer must open a GitHub issue nominating a candidate. The issue must remain open for at least 14 days for public feedback, and existing Maintainers elect a new Maintainer by a supermajority vote. This guide does not alter the Project Lead’s final-vote role for project decisions and escalated conflicts described in `GOVERNANCE.md`.

A candidate will normally have demonstrated sustained and significant contribution, sound technical judgment, reliability, mentorship, and broad understanding of the project. Previous Reviewer recognition can be relevant evidence, but it is not a separate prerequisite or appointment mechanism.

### Responsibilities and authority

The authoritative list of Maintainer responsibilities remains in [`GOVERNANCE.md`](GOVERNANCE.md). It includes issue triage, pull-request review, releases, architecture and roadmap stewardship, community support, CNCF participation and compliance, and project-infrastructure stewardship.

Under [`CONTRIBUTING.md`](CONTRIBUTING.md), a pull request requires an LGTM from at least one Maintainer listed in `CODEOWNERS`. Maintainers participate in project decisions under [`GOVERNANCE.md`](GOVERNANCE.md), including the Project Lead’s stated final-vote role.

Merge, release, and administrative permissions require separate authorization. Under current governance, adding someone to the repository-wide `CODEOWNERS` entry designates them as a Maintainer; Reviewer recognition alone does not justify an entry. Any future scoped ownership or access model requires a separate governance and access decision.

## Inactivity, emeritus status, and return

Maintainer emeritus status and removal are governed by [`GOVERNANCE.md`](GOVERNANCE.md). This ladder does not alter that process.

Maintainers may periodically review Reviewer recognition when a Reviewer is no longer active, no longer wishes to perform the role, or no longer meets the role’s expectations. They should seek to contact the person and record the decision appropriately. Withdrawal of Reviewer recognition does not itself change GitHub access; any access change follows the separate process that granted the access. Returning Contributors are welcome and may be recognized again when current Maintainers determine that the evidence supports it.

## Changes to this ladder

Changes to this guide must be proposed through the normal pull-request process and reviewed and, where required by [`GOVERNANCE.md`](GOVERNANCE.md), voted on by project Maintainers. The guide should remain proportionate to k8gb’s actual community, governance, and repository controls.

## Credits

Sections of this document were adapted from the [CNCF Contributor Ladder template](https://github.com/cncf/project-template/blob/main/CONTRIBUTOR_LADDER.md), the [Kubernetes community membership guide](https://github.com/kubernetes/community/blob/master/community-membership.md), and the [OpenTelemetry membership, roles, and responsibilities guide](https://github.com/open-telemetry/community/blob/main/guides/contributor/membership.md).
