---
date:
  created: 2026-11-25
authors:
  - elohmrow
categories:
  - Community
  - Releases
  - Integrations
tags:
  - CNCF
  - Incubation
  - Community Growth
  - Open Source
  - Contributors
  - Kubernetes
  - GSLB
  - Reference Architectures
  - KubeCon
  - Project Governance
  - Global Resilience
  - Gateway API
  - AI Inference
  - Sovereign Cloud
  - Software Supply Chain Security
pin: true
draft: true
---

# k8gb Is Growing Up: Beyond Incubation

In July 2026, k8gb became a CNCF Incubating project. The [official announcement](https://www.cncf.io/announcements/2026/08/05/k8gb-becomes-a-cncf-incubating-project/) followed in August, marking an important milestone in a journey that began more than six years ago.

More people are trying k8gb, and more of them are sending changes back.

## Six years in the making

In September, during our 100th Community Meeting, we released [k8gb v1.0.0](https://github.com/k8gb-io/k8gb/releases/tag/v1.0.0).

We did it live. We merged the release PR, watched the automation run, reviewed the generated release notes, and pressed publish together.

The release automation generated provenance attestations and SBOMs, and all of it worked.

Even after six years, we were still welcoming first-time contributors: [@dejoski](https://github.com/dejoski), [@Goutham-Annem](https://github.com/Goutham-Annem), and [@nightcityblade](https://github.com/nightcityblade) all had their first PRs in v1.0.0.

## More people behind the PRs

Since incubation in July, our GitHub stars have grown by nearly 13%, and forks by more than 17%.

The bigger change is in who contributes. Comparing the 130 days before incubation with the 130 days since, twice as many people had PRs merged: 16, up from 8.

- The number of people opening PRs grew from 11 to 21.
- Nine people had their first k8gb PR merged, compared with four in the preceding period.
- Human-authored merged PRs increased from 74 to 83.

These figures exclude bot accounts and are based on the [main k8gb repository](https://github.com/k8gb-io/k8gb).

Some recent contributions:

- [@vicsufer](https://github.com/vicsufer) introduced [Gateway API ListenerSet support](https://github.com/k8gb-io/k8gb/pull/2540). k8gb now accepts routes whose parent is a ListenerSet, part of the Gateway API v1.5 standard. Before, those routes failed because k8gb only looked for a Gateway parent.
- [@ankushkhakale](https://github.com/ankushkhakale) improved [DNS fallback behavior](https://github.com/k8gb-io/k8gb/pull/2625) and [DNS integration testing](https://github.com/k8gb-io/k8gb/pull/2605).
- [@aryasoni98](https://github.com/aryasoni98) has had 13 PRs merged since incubation.
- [Karsten S.](https://www.linkedin.com/in/kubernetes-experte/) shared his insights into [split-brain scenarios](https://lnkd.in/p/ejk6W9RP), including practical monitoring, testing, and architectural considerations.

## Growing the next generation of contributors

We've introduced a formal [Contributor Ladder](https://www.k8gb.io/latest/CONTRIBUTOR_LADDER/), creating a clearer path from Contributor to Reviewer to Maintainer, and recognizing the many ways people can contribute to the project.

In September, [@itsfarhan became our first formally recognized Contributor](https://github.com/k8gb-io/k8gb/issues/2534).

Farhan had attended our community meetings for more than a year, contributed dozens of PRs, introduced Trivy into our build and release pipelines, and even created this very blog!

This fall, his work helped k8gb earn all six badges in the [Fall 2026 CNCF and OpenSSF Security Slam](https://securityslam.com/outcomes/fall-2026/k8gb).

Thanks, Farhan!

We also started getting more low-quality, AI-generated PRs, and reviewing them cost maintainer time.

We introduced an [AI Contribution Policy](https://www.k8gb.io/latest/AI_POLICY/) that welcomes responsible AI-assisted contributions. Contributors must understand their changes well enough to explain them without AI, and AI tools can't appear as `Signed-off-by` or `Co-authored-by` on a commit.

## Building together across the ecosystem

When we entered incubation, our roadmap named expanding the ecosystem integration list as a priority. ListenerSet support was one step. These reference architectures and demos are the next:

- **[NeoNephos](https://github.com/k8gb-io/k8gb/issues/2535):** Global traffic management as part of sovereign cloud infrastructure.
- **[KEDA](https://github.com/k8gb-io/k8gb/issues/2536):** Event-driven scaling and cost-effective multi-region high availability. Jiri Kremser, a KEDA maintainer, is building it, following up on the KCD Porto talk [Highly Available, Barely Running](https://kcd-porto-2026.sessionize.com/session/1300190).
- **[NVIDIA Dynamo](https://github.com/k8gb-io/k8gb/pull/2579):** Global routing and regional failover for AI inference workloads.
- **[Zero-Traffic Region Upgrades with GitOps](https://github.com/k8gb-io/k8gb/issues/2618):** GitOps-driven reference architecture that uses k8gb to safely drain traffic from a Kubernetes cluster region, apply disruptive platform upgrades (such as Kubernetes, ingress, or CNI updates) while the region receives zero traffic, verify its health, and restore traffic before repeating the process for the next region. 

Community members are building all three, and each one shows how k8gb works with other open-source projects to solve a real infrastructure problem.

## Taking the conversation beyond GitHub

Incubation brought CNCF blog coverage, keynote recognition, and more kiosk time at KubeCon.

At KubeCon North America in Salt Lake City, we spoke with more than 60 people at our kiosk in eight hours.

Three topics kept coming up: split-brain prevention, single-cluster failover, and AI inference resilience. Enough people asked about single-cluster failover that we're preparing a blog post and demo on it for the coming weeks.

We also had two opportunities to share our work on stage:

- [**k8gb Is Growing Up**](https://www.youtube.com/watch?v=aUFtEHqHhsY), covering project developments, community growth, AI inference, and sovereign multi-cloud infrastructure.
- [**When the Region Goes Dark**](https://www.youtube.com/watch?v=0LpKYR3gAz4), at Kubernetes on Edge Day, exploring decentralized global load balancing and failover across distributed edge environments.

And at KCD Porto, the conversation continued:

- Nuno Guedes and Filipe Revez presented [**Highly Available, Barely Running**](https://kcd-porto-2026.sessionize.com/session/1300190), demonstrating cost-effective high availability using Flux, KEDA, and k8gb.
- Nuno Guedes and Yury Tsarev presented [**Open Source Is Not Free**](https://kcd-porto-2026.sessionize.com/session/1298811), exploring why enterprises must actively help sustain the open-source projects they depend on.

## We're just getting started

Next up is k8gb v1.1, which will ship the latest security fixes, starting with a security release of the [CoreDNS CRD plugin](https://github.com/k8gb-io/coredns-crd-plugin). We also plan to run the same security scanner [Artifact Hub](https://artifacthub.io/packages/helm/k8gb/k8gb) uses before each release, so we catch issues earlier.

If you run k8gb in production, tell us about it. Millennium bcp, Portugal's largest privately owned bank, shared how k8gb helped them reach 99.99% uptime across critical banking applications. Your story, in Slack or at a community meeting, helps us decide what to build next.

Whether you're contributing code, reviewing PRs, improving documentation, building reference architectures, sharing your production experiences, or simply asking difficult questions, there is a place for you in the k8gb community.

Looking for a first issue? Browse the [issue list](https://github.com/k8gb-io/k8gb/issues), or ask in Slack or at a community meeting. You'll find the meeting schedule and Zoom links on [our website](https://www.k8gb.io/latest/#presentations-and-talks).

Come build with us!

- [GitHub](https://github.com/k8gb-io/k8gb)
- [CNCF Slack #k8gb](https://cloud-native.slack.com/archives/C021P656HGB)
- [Community meetings](https://www.k8gb.io/latest/#presentations-and-talks)
- [Contributor Ladder](https://www.k8gb.io/latest/CONTRIBUTOR_LADDER/)
- [k8gb.io](https://www.k8gb.io/)
