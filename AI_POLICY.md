# k8gb AI Contribution Policy

## Scope and relationship to existing rules

This policy applies to contributions to the [`k8gb-io/k8gb`](https://github.com/k8gb-io/k8gb) repository. It covers AI-assisted code, documentation, tests, issues, pull requests, commit messages, and review comments submitted by Contributors, Reviewers, or Maintainers.

This policy supplements, and does not replace, [`CONTRIBUTING.md`](CONTRIBUTING.md), [`GOVERNANCE.md`](GOVERNANCE.md), [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md), [`SECURITY.md`](SECURITY.md), the project license, and the Developer Certificate of Origin (DCO) requirement. Those documents control if there is a conflict. This policy does not create a different review or merge path, grant GitHub permissions, or change project-maintained automation that is separately configured.

## Core principle

k8gb permits the use of generative artificial-intelligence (AI) tools, including large language models, coding assistants, and agents, to assist with contributions. AI assistance does not change who is accountable: the named human contributor remains responsible for the correctness, security, licensing, quality, and reviewability of what they submit.

A contribution is not rejected solely because it was produced substantially with AI assistance. It is acceptable only when a named human contributor has reviewed it, understands it, can explain and defend it, and accepts responsibility for it under the same contribution requirements that apply to any other submission.

## Owning and validating a contribution

Contributors must understand the purpose and expected impact of each submitted change well enough to explain it without relying on an AI tool. Before requesting review, they must perform appropriate self-review and validation. This includes reviewing generated code, documentation, tests, configuration, citations, and claims; revising them as needed; and running the applicable checks described in [`CONTRIBUTING.md`](CONTRIBUTING.md).

AI-generated tests are useful only when the contributor has reviewed them for relevance and correctness and has validated the behavior they are intended to cover. AI-assisted documentation must be technically accurate, appropriately scoped, and consistent with the repository. Do not add an API, dependency, behavior claim, citation, or configuration solely because an AI tool proposed it; verify it against authoritative project or upstream sources.

Complete an applicable pull-request or issue template where one is provided. Use focused, reviewable changes. For a non-trivial proposed change, start an issue or discussion before implementation unless a Maintainer asks otherwise. Avoid opening multiple related, unreviewed pull requests at once.

## Human communication and review judgment

Issues, pull-request descriptions, replies to review comments, and review comments must reflect the submitting person’s own understanding and judgment. AI may help a contributor improve wording, translate, explore an implementation, or prepare a draft, but the contributor must review and revise the result before posting it.

Do not configure an AI tool to post or reply autonomously on a contributor’s behalf in project issues, pull requests, or reviews. Do not submit an AI-generated review as your own without exercising and expressing your own judgment. This expectation applies equally when a Maintainer uses AI while reviewing or responding to a contribution.

This rule does not prohibit clearly identified project-maintained automation that is separately configured and does not substitute for a contributor’s response to review. In particular, it does not change existing repository workflows that generate release material under their own controls.

## AI-use disclosure

k8gb does not require a contributor to disclose AI use solely because an AI tool was used. A contributor must nevertheless provide information needed to assess a submission, including information necessary to identify relevant third-party rights, comply with applicable law or tool terms, or answer a concrete Maintainer question about correctness, security, dependencies, licensing, or provenance.

The absence of a general disclosure requirement does not excuse inaccurate, incomplete, or unowned issue, pull-request, or review content.

## DCO, licensing, and third-party material

Every contribution remains subject to the project’s Apache-2.0 license and the DCO requirements in [`CONTRIBUTING.md`](CONTRIBUTING.md). For a contributor-authored AI-assisted commit, the human contributor must sign off using their own name and email as required by the DCO process. Do not list an AI system as a `Signed-off-by`, `Co-authored-by`, or similar human-contributor trailer. This provision does not alter project-maintained automation that is separately configured to create signed-off commits or pull requests.

Before submitting AI-assisted material, contributors must ensure that their use of the tool and its output is compatible with the project’s license, intellectual-property requirements, and applicable terms. If output includes or closely reproduces pre-existing third-party copyrighted material, contributors must first confirm that they have permission to use, modify, and contribute it. They must provide the applicable source notice, attribution, and license information in the pull request or accompanying documentation. Do not submit material when its rights cannot be confirmed.

Where available, similarity or licensing features supplied by an AI tool may help identify possible third-party material, but they do not replace the contributor’s review and responsibility.

## Security and confidential reporting

This policy does not change the project’s vulnerability-reporting process. If a contribution, investigation, or AI-assisted analysis identifies a security issue, follow [`SECURITY.md`](SECURITY.md) and use its confidential reporting channel rather than opening a public issue or pull request with vulnerability details.

## Review and enforcement

AI-assisted contributions follow the normal k8gb contribution and review process. In particular, a pull request still requires the Maintainer LGTM specified in [`CONTRIBUTING.md`](CONTRIBUTING.md). AI assistance does not replace required testing, human review, security review where appropriate, or Maintainer judgment.

Maintainers may request clarification or changes, or close contributions that do not meet documented contribution, licensing, security, testing, or review-engagement expectations. They should apply this policy consistently with [`GOVERNANCE.md`](GOVERNANCE.md) and the [Code of Conduct](CODE_OF_CONDUCT.md). A decision must not be based solely on the appearance that a submission may have used AI.

## Related external guidance

The [Linux Foundation guidance on generative AI and open-source development](https://www.linuxfoundation.org/legal/generative-ai) permits AI-generated content while addressing output terms and third-party copyrighted material. This k8gb policy is intended to complement applicable CNCF/Linux Foundation guidance and the project’s existing licensing and contribution requirements. Where an external policy or legal obligation applies, contributors must comply with it.

## Changes to this policy

Changes to this policy must be proposed through the normal pull-request process and reviewed and, where required by [`GOVERNANCE.md`](GOVERNANCE.md), voted on by project Maintainers.

## Credits

Sections of this document were adapted from the [Linux Foundation guidance on generative AI and open-source development](https://www.linuxfoundation.org/legal/generative-ai), the [Crossplane AI Contribution Policy](https://github.com/crossplane/crossplane/blob/main/AI_POLICY.md), the [Microcks AI Policy](https://github.com/microcks/microcks-testcontainers-node/blob/main/AI-POLICY.md), the [Kubewarden AI Contribution Policy](https://github.com/kubewarden/community/blob/main/AI_POLICY.md), and the [Kubernetes pull-request guide](https://www.kubernetes.dev/docs/guide/pull-requests/).
