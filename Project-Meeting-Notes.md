# k8gb.io Project Meeting Notes & Agenda

This document will capture the agenda and meeting notes and links for the recurring k8gb.io office hours meeting. 

Its latest home was [Proton Docs](https://drive.proton.me/urls/MHENT68VER#AM8e6BqdFvmv). Its first home was [Google Docs](https://docs.google.com/document/d/1YdpEVhtyCKvwFtXR7cn1Kn2Xc_tdNskoFnhHmFPbtA4/edit?tab=t.0#heading=h.b62b237a7ky). 

It lives here now, in [GitHub](https://github.com/k8gb-io/k8gb), to allow the greatest number of folks to see it and interact with it.

### Schedule

The meeting runs every other Wednesday from [13:00 CET](https://dateful.com/convert/prague-czechia?t=13) to [13:30 CET](https://dateful.com/convert/prague-czechia?t=1330) ([**calendar**](https://zoom-lfx.platform.linuxfoundation.org/meetings/k8gb?view=month)):

Join the [Zoom Meeting](https://zoom-lfx.platform.linuxfoundation.org/meeting/92572060749?password=645f8346-1952-44fa-bd9b-45208260fc10) 

**Old** Meeting recordings: [YouTube](https://www.youtube.com/channel/UCwvtktvdZu_pg-t-INvuW5g/featured) 

**New** Meeting recordings will be on CNCF site - links follow per date

### Links

- Project website: [k8gb.io](http://k8gb.io)
- GitHub Repo:  [https://github.com/k8gb-io/k8gb](https://github.com/k8gb-io/k8gb)
- Slack: [#k8gb](https://cloud-native.slack.com/archives/C021P656HGB)
- Mailing-list: [cncf-k8gb-maintainers@lists.cncf.io](mailto:cncf-k8gb-maintainers@lists.cncf.io)
- LinkedIn: [https://www.linkedin.com/company/k8gb/](https://www.linkedin.com/company/k8gb/) 
- Twitter / X: [https://x.com/k8gb\_io](https://x.com/k8gb_io) 
- Medium: [https://medium.com/@kubernetesglobalbalancer](https://medium.com/@kubernetesglobalbalancer) 
- YouTube: [https://youtube.com/@k8gb823](https://youtube.com/@k8gb823)

<details><summary><strong>Backlog</strong></summary>

## Backlog:

- Community user reporting docker hub rate limiting during k8gb installation - republish to github
- [https://github.com/k8gb-io/k8gb/issues/1314](https://github.com/k8gb-io/k8gb/issues/1314) - split brain documentation request

</details>

<details open><summary><strong>Feb 4, 2026 #85</strong></summary>

## Feb 4, 2026 #85

Zoom Recording: https://zoom.us/rec/share/VoC_axlmGGixnrhlu0Wa4S9B69GRzsHr5vSZCUnkX4a7ame_b_YOoubxUg0pV1Ld.4Y-OjjxtdMoR3lSc

On YouTube: []()

**Attendees**

- Yury Tsarev
- 

**Agenda**

- [k8gb project board](https://github.com/orgs/k8gb-io/projects/2/views/1)
- **News**
    - Incubation
      - https://github.com/cncf/toc/pull/2020 Governance review
        - https://github.com/k8gb-io/k8gb/issues/2210 Related Project Lead Election
      - https://github.com/cncf/toc/pull/2029 Tech review
- **Discuss**
  - https://github.com/k8gb-io/k8gb/issues/2195 Switch to vendor neutral OCI registry and repository
  - https://github.com/k8gb-io/k8gb/issues/2180 Switch to vendor-neutral API group
    - this one is tricky to make non-intrusive https://github.com/k8gb-io/k8gb/pull/2203 - testing this approach
  - Release this week?
- **PR review**
    - Trivy implementation [https://github.com/k8gb-io/k8gb/pull/2179](https://github.com/k8gb-io/k8gb/pull/2179) 🙏 @itsfarhan - merged, needs some follow up fixes
    - Mkdocs versioning [https://github.com/k8gb-io/k8gb/pull/2178](https://github.com/k8gb-io/k8gb/pull/2178) 🙏 @itsfarhan
    - Incubation DTR [https://github.com/k8gb-io/k8gb/pull/1909](https://github.com/k8gb-io/k8gb/pull/1909) 🙏 @itsfarhan
    - GSLB reconciliation should not fail on hostnames outside delegated zones in referenced resources [https://github.com/k8gb-io/k8gb/issues/2183](https://github.com/k8gb-io/k8gb/issues/2183) 🙏 @Piroddi - merged
    - https://github.com/k8gb-io/k8gb/pull/2204 Add new k8gb_gslb_healthy_local_records prom metrics 🙏 @Piroddi - merged
    - https://github.com/k8gb-io/k8gb/pull/2184 🙏 @Piroddi - merged
- **Issue Review**
    - Good first issues: https://github.com/k8gb-io/k8gb/issues?q=is%3Aissue%20state%3Aopen%20label%3A%22good%20first%20issue%22 
- **Community Update**
    - WIP: "what is k8gb" and "k8gb getting started" videos
- **Other**
- **Action Items**


<details><summary><strong>Feb 4, 2026 #85</strong></summary>

## Jan 21, 2026 #84

Zoom Recording: https://zoom.us/rec/share/VoC_axlmGGixnrhlu0Wa4S9B69GRzsHr5vSZCUnkX4a7ame_b_YOoubxUg0pV1Ld.4Y-OjjxtdMoR3lSc

On YouTube: [https://youtu.be/uIAZB8DFzqo](https://youtu.be/uIAZB8DFzqo)

**Attendees**

- Farhan Ahmed
- Yury Tsarev
- Bradley Andersen

**Agenda**

- [k8gb project board](https://github.com/orgs/k8gb-io/projects/2/views/1)
- **News**
    - KubeCon
      - co-located talk waitlisted: _Building Unified Global Load Balancing for the Edge With k8gb_
      - [Maintainer Summit talk](https://sched.co/2EF6x) accepted
      - kiosk approved  
      - no lightning talk this time
      - submitted to Rejekts:
        - _Building Unified Global Load Balancing for the Edge With k8gb_
        - _Community Manager Speedrun: Sandbox to Incubation_
        - _Brea-k8gb-ing Good: Say My (Domain) Name_
    - Incubation
      - ADOPTER reviews WIP
      - need to make sure everything on the [Incubation Application](https://github.com/cncf/toc/issues/1472) is up2date
        - especially: licensing and governance
          - licensing is handled by FOSSA
          - how do other projects handle governance? do we need a steering committee?
      - Incubating DD Review
        - [https://github.com/k8gb-io/k8gb/issues/1906](https://github.com/k8gb-io/k8gb/issues/1906) 
        - [https://github.com/k8gb-io/k8gb/pull/2095](https://github.com/k8gb-io/k8gb/pull/2095) 
- **Discuss**
  - k8gb-kb?
    - https://github.com/k8gb-io/k8gb/issues/2022
    - https://github.com/k8gb-io/k8gb/issues/1872
- **PR review**
    - Trivy implementation [https://github.com/k8gb-io/k8gb/pull/2179](https://github.com/k8gb-io/k8gb/pull/2179) 🙏 @itsfarhan
    - Mkdocs versioning [https://github.com/k8gb-io/k8gb/pull/2178](https://github.com/k8gb-io/k8gb/pull/2178) 🙏 @itsfarhan
    - Incubation DTR [https://github.com/k8gb-io/k8gb/pull/1909](https://github.com/k8gb-io/k8gb/pull/1909) 🙏 @itsfarhan
    - GSLB reconciliation should not fail on hostnames outside delegated zones in referenced resources [https://github.com/k8gb-io/k8gb/issues/2183](https://github.com/k8gb-io/k8gb/issues/2183) 🙏 @Piroddi
- **Issue Review**
    - 46 (-21) after Yury's heroic review 💪
- **Community Update**
    - WIP: "what is k8gb" and "k8gb getting started" videos
- **Other**
- **Action Items**

</details><details><summary><strong>Jan 07, 2026 #83</strong></summary>

## Jan 07, 2026 #83

Recording: https://zoom.us/rec/share/ukYqcJtxJH2wkDFmaK0xo09aYNal3yHdzsW1-FU1rmuL8JKyd4T5TXUmeIwhklgl.OK4fbf5aodiPF3z4

Attendees:

- Tomáš Boros
- Farhan Ahmed
- Yury Tsarev
- Bradley Andersen

Agenda:

- **News**
    - blog added: https://www.k8gb.io/blog/
    - agenda moved to markdown
    - Community Meetings now on [YouTube](https://www.youtube.com/@k8gb823)
    - KubeCon
      - co-located talk waitlisted: "Building Unified Global Load Balancing for the Edge With k8gb" (Yur, Bradley)
      - lightning talk: to be delivered by Nuno
      - non-accepted talks: to be / submitted to Rejekts and possibly made into blog posts
      - kiosk: waiting
    - Incubation kickoff with TOC held
      - 3-month plan
        - ADOPTER reviews (2 already WIP) and TOC availability
      - need to make sure everything on the [Incubation Application](https://github.com/cncf/toc/issues/1472) is up2date
        - especially: licensing and governance
          - licensing is handled by FOSSA
          - how do other projects handle governance? do we need a steering
          committee?
    - Christmas Release v0.17.0 with Gateway API support
      - 4 new contributors!
        - 🙇 @WesleyKlop, @angelbarrera92, @actionjax, @mattwelke
      - Supported resources: HTTPRoute, GRPCRoute, TCPRoute, UDPRoute, TLSRoute
      - See [resource reference examples](https://www.k8gb.io/resource_ref/)
      - Full [release notes](https://github.com/k8gb-io/k8gb/releases/tag/v0.17.0)
- **Discuss**
- **PR review**
    - Merged / closed since last time:
      - (m) fix: geodatafilepath and geodatafield are lowercase in coredns [https://github.com/k8gb-io/k8gb/pull/2103](https://github.com/k8gb-io/k8gb/pull/2103) 
      - (m) fix(docs): fix dig commands for each CoreDNS instance in local tutorial [https://github.com/k8gb-io/k8gb/pull/1832](https://github.com/k8gb-io/k8gb/pull/1832)
      - (c) Add support for ExternalName service health checks [https://github.com/k8gb-io/k8gb/pull/1888](https://github.com/k8gb-io/k8gb/pull/1888) 
      - (m) feat: OCI Registry support fix [https://github.com/k8gb-io/k8gb/pull/2089](https://github.com/k8gb-io/k8gb/pull/2089) 
      - (m) fix(docs): added more rollback procedures [https://github.com/k8gb-io/k8gb/pull/2105](https://github.com/k8gb-io/k8gb/pull/2105)
    - No change since last time:
      - DynamicZones [https://github.com/k8gb-io/k8gb/pull/2102](https://github.com/k8gb-io/k8gb/pull/2102) / Dinar 
      - fix: restore extraServerBlocks as global setting [https://github.com/k8gb-io/k8gb/pull/2121](https://github.com/k8gb-io/k8gb/pull/2121)
      - feat: Add health checking for ingress controllers [https://github.com/k8gb-io/k8gb/pull/2110](https://github.com/k8gb-io/k8gb/pull/2110)
      - test(DTR): Testing different K8S versions for DTR (incubation)  [https://github.com/k8gb-io/k8gb/pull/2095](https://github.com/k8gb-io/k8gb/pull/2095) 
- **Issue Review**
    - 67 (+4) currently open - need to review
    - André / Peishu: In-cluster DNS Issue [https://github.com/k8gb-io/k8gb/issues/2022](https://github.com/k8gb-io/k8gb/issues/2022) <-- close
    - Yury / Farhan: Switch Helm Repo to OCI [https://github.com/k8gb-io/k8gb/issues/1973](https://github.com/k8gb-io/k8gb/issues/1973)
    - Incubating DD Review [https://github.com/k8gb-io/k8gb/issues/1906](https://github.com/k8gb-io/k8gb/issues/1906) 
- **Community Update**
    - WIP: "what is k8gb" and "k8gb getting started" videos
    - LinkedIn vanity stats
      - Followers +10.7% last week, +60 since last meeting
      - Page visitors +21.4% last week
      - Search appearances +106.3% last week
      - Posting
        - 5 posts since last community meeting:
          - CTR 6.62% - cta, join community meeting (post)
          - CTR 2.35% - community meeting highlights (repost)
          - CTR 3.6% - announcing youtube uploads (post)
          - CTR 5% - gw api release announcement (post)
          - CTR 9.03% - cta, add to adopters (post)
        - CTR = clicks / impressions
        - impression = >= 50% of the post is visible on a member's screen for >= 300ms
    - 11 new [stars](https://github.com/k8gb-io/k8gb/stargazers) since last community meeting 📈
- **Other**
- **Action Items**

</details><details><summary><strong>Dec 10, 2025 #82</strong></summary>

## Dec 10, 2025 #82

Recording: [https://zoom.us/rec/share/ryP\_sAYA5uV8qeyIQJ0dwP0JNnXtLWcK9tt0B\_RpGFa\_uAmojUZiW\_61BeFK9erc.EZqEdwsjNg25hiSL](https://zoom.us/rec/share/ryP_sAYA5uV8qeyIQJ0dwP0JNnXtLWcK9tt0B_RpGFa_uAmojUZiW_61BeFK9erc.EZqEdwsjNg25hiSL)

Attendees:

- Bradley
- Yury
- Dinar
- Farhan

Agenda:

- **News**
    - 
- **Discuss**
    - Christmas Release with GateWay API support
- **PR review**
    - [x] Add support for GatewayAPI's TCPRoute [https://github.com/k8gb-io/k8gb/pull/2116](https://github.com/k8gb-io/k8gb/pull/2116) <-- relevant to xmas release / Andre
    - [x] DynamicZones [https://github.com/k8gb-io/k8gb/pull/2102](https://github.com/k8gb-io/k8gb/pull/2102) / Dinar 
    - fix: restore extraServerBlocks as global setting [https://github.com/k8gb-io/k8gb/pull/2121](https://github.com/k8gb-io/k8gb/pull/2121)
    - fix: geodatafilepath and geodatafield are lowercase in coredns [https://github.com/k8gb-io/k8gb/pull/2103](https://github.com/k8gb-io/k8gb/pull/2103) 
    - [x] feat: Add health checking for ingress controllers [https://github.com/k8gb-io/k8gb/pull/2110](https://github.com/k8gb-io/k8gb/pull/2110)
    - fix(docs): fix dig commands for each CoreDNS instance in local tutorial [https://github.com/k8gb-io/k8gb/pull/1832](https://github.com/k8gb-io/k8gb/pull/1832)
    - Add support for ExternalName service health checks [https://github.com/k8gb-io/k8gb/pull/1888](https://github.com/k8gb-io/k8gb/pull/1888) 
    - [x] feat: OCI Registry support fix [https://github.com/k8gb-io/k8gb/pull/2089](https://github.com/k8gb-io/k8gb/pull/2089) 
    - [x] test(DTR): Testing different K8S versions for DTR (incubation)  [https://github.com/k8gb-io/k8gb/pull/2095](https://github.com/k8gb-io/k8gb/pull/2095) 
    - [x] fix(docs): added more rollback procedures [https://github.com/k8gb-io/k8gb/pull/2105](https://github.com/k8gb-io/k8gb/pull/2105)
- **Issue Review**
    - 63 currently open - need to review
    - André / Peishu: In-cluster DNS Issue [https://github.com/k8gb-io/k8gb/issues/2022](https://github.com/k8gb-io/k8gb/issues/2022) <-- close
    - Yury / Farhan: Switch Helm Repo to OCI [https://github.com/k8gb-io/k8gb/issues/1973](https://github.com/k8gb-io/k8gb/issues/1973)
    - Incubating DD Review [https://github.com/k8gb-io/k8gb/issues/1906](https://github.com/k8gb-io/k8gb/issues/1906) 
- **Community Update**
    - WIP: "what is k8gb" and "k8gb getting started" videos
    - Incubating app: [security self-assessment](https://github.com/cncf/toc/issues/1472#issuecomment-3602317184) merged
        - invite to KubeCon maintainer track panel (TAG Security)
- **Other**
    - KubeCon (co-located): "Building Unified Global Load Balancing for the Edge With k8gb"
    - KubeCon lightning talk: any volunteers? Nuno!
- **Action Items**
    - bradley: to prevent some folks from being blocked: consider moving agenda to markdown on github, consider something other than zoom (google meet is blocked)

</details><details><summary><strong>Nov 26, 2025 #81</strong></summary>

## Nov 26, 2025 #81

Recording: [https://zoom.us/rec/share/Z46B3TJM-PoKDS\_XtoJBz2fE1iE6FYfoOTWH6wstbQCMvEUDPb1455JS-m4QMOMz.s7ujQqquoDUhxrnE](https://zoom.us/rec/share/Z46B3TJM-PoKDS_XtoJBz2fE1iE6FYfoOTWH6wstbQCMvEUDPb1455JS-m4QMOMz.s7ujQqquoDUhxrnE)

Attendees:

- Bradley
- Yury
- Farhan
- André

Agenda:

- **News**
    - KubeCon: [https://www.linkedin.com/posts/k8gb\_k8gb-kubecon-cncf-activity-7397197050511261696-GkBF/](https://www.linkedin.com/posts/k8gb_k8gb-kubecon-cncf-activity-7397197050511261696-GkBF/) 
        - maintainer summit
        - lightning talk
        - 2x kiosk
        - community boost
    - Ingress / GW API: [https://www.linkedin.com/posts/k8gb\_kubernetes-k8gb-gatewayapi-activity-7396423054643335168-v1EM/](https://www.linkedin.com/posts/k8gb_kubernetes-k8gb-gatewayapi-activity-7396423054643335168-v1EM/) 
        - ingress-nginx controller retires 3/26
        - ingress api stays / frozen 
        - christmas release :tada:
- **Issue Review**
    - [https://github.com/k8gb-io/k8gb/issues/2022](https://github.com/k8gb-io/k8gb/issues/2022) anything to add here? should we close?
        - andre will check
    - [https://github.com/k8gb-io/k8gb/issues/1973](https://github.com/k8gb-io/k8gb/issues/1973) is implemented - double check the OCI repo state, then we can close
        - yury will take a look at merging
    - [https://github.com/k8gb-io/k8gb/issues/1906](https://github.com/k8gb-io/k8gb/issues/1906) Due Diligence Technical Review readiness
        - Important for Incubating
        - Farhan is already working on subitems
        - WIP 
- **Community Update**
    - LinkedIn :rocket:
        - > 150 followers (+ 12% last 7 days) / + 50% last 30 days
        - 3,175 post impressions last 30 days  (+ 182% last 7 days)
        - + 10% search appearances last 7 days
        - + 216% page views last 30 days
        - + 280% unique visitors last 30 days
    - Incubating app
- **Actions**
    - @bradley:
        - kubecon medium post
        - incubating
            - ~~take another look at what is required by&#32;[https://github.com/k8gb-io/k8gb/issues/1907](https://github.com/k8gb-io/k8gb/issues/1907)
            - ~~mv incubating app security review ticket to new repo~~

</details><details><summary><strong>Oct 29, 2025 #80</strong></summary>

## Oct 29, 2025 #80

Recording: [https://zoom.us/rec/share/ssNtsRt72jx0IHzvdbGN0UmmsnDl1bPI4\_U2MTPc7UR8R1tkyiGWdC\_HBWlrLsBF.s2VmP3-t9fzyH6J7](https://zoom.us/rec/share/ssNtsRt72jx0IHzvdbGN0UmmsnDl1bPI4_U2MTPc7UR8R1tkyiGWdC_HBWlrLsBF.s2VmP3-t9fzyH6J7) 

Attendees:

- Yury
- Bradley
- Farhan
- André
- Dinar

Agenda:

- **News**
    - v0.16.0 released: 
        - [https://github.com/k8gb-io/k8gb/releases/tag/v0.16.0](https://github.com/k8gb-io/k8gb/releases/tag/v0.16.0) 🔥
        - [https://www.linkedin.com/posts/yurytsarev\_release-v0160-k8gb-iok8gb-activity-7386324075910184961-K3mV/](https://www.linkedin.com/posts/yurytsarev_release-v0160-k8gb-iok8gb-activity-7386324075910184961-K3mV/) 
            - delivers comprehensive multi-cloud DNS automation, with automated zone delegation now supporting AWS Route53, Azure DNS, and newly added GCP Cloud DNS.

This milestone enables seamless global load balancing across any combination of cloud environments without manual DNS configuration.
🚀 Key highlights:
 • GCP Cloud DNS support — full integration with automated zone delegation and end-to-end testing
 • LoadBalancer Service support — extending global load balancing to Layer 4, complementing Ingress
 • Upstream external-dns migration — ensuring feature parity and maintainability across all major providers
 • 90+ additional improvements including enhanced Infoblox integration, CoreDNS hot-reload, and dependency updates (Go 1.25.2, k8s API Machinery v0.34.1)
- **Discuss**
    - 
- **Issue Review**
    - [https://github.com/k8gb-io/k8gb/issues/2022](https://github.com/k8gb-io/k8gb/issues/2022) anything to add here? should we close?
        - andre will check
    - [https://github.com/k8gb-io/k8gb/issues/1973](https://github.com/k8gb-io/k8gb/issues/1973) is implemented - double check the OCI repo state, then we can close
        - yury will take a look at merging
    - [https://github.com/k8gb-io/k8gb/issues/1906](https://github.com/k8gb-io/k8gb/issues/1906) Due Diligence Technical Review readiness
        - Important for Incubating
        - Farhan is already working on subitems
        - WIP 
    - Sudhamsh question in slack [https://cloud-native.slack.com/archives/C021P656HGB/p1760512987294269](https://cloud-native.slack.com/archives/C021P656HGB/p1760512987294269)
        - 2 ingresses suggested
- **Community Update**
    - 
- **Other**
    - absa case study
- **Action Items**
    - 

</details><details><summary><strong>Oct 15, 2025 #79</strong></summary>

## Oct 15, 2025 #79

Recording: [https://zoom.us/rec/share/lllg5DBDRQWE59-3Tlffc-pl-XQr009HmNYEfsWFyy40VcCBeHzo7s60vZ2yKqgo.GKYs14IpCrIQObtB](https://zoom.us/rec/share/lllg5DBDRQWE59-3Tlffc-pl-XQr009HmNYEfsWFyy40VcCBeHzo7s60vZ2yKqgo.GKYs14IpCrIQObtB)

Attendees:

- Yury
- Bradley
- Michal
- Marcus
- Farhan

Agenda:

- **News**
    - [Document and test GCP Cloud DNS Provider Integration for K8GB](https://github.com/k8gb-io/k8gb/pull/2065) - merged, we can safely declare big3 in upcoming release
    - Movement in Incubating application
        - [https://github.com/orgs/cncf/projects/27/views/9](https://github.com/orgs/cncf/projects/27/views/9) 
        - 
        - Due Diligence assigned! [https://github.com/cncf/toc/issues/1472#issuecomment-3402773646](https://github.com/cncf/toc/issues/1472#issuecomment-3402773646)
        - 
    - [https://www.k8gb.io/](https://www.k8gb.io/) got new fancy star button, it already helped getting more stars, thanks Bradley!
    - 
    - CLO Monitor dropped to 95, Bradley added LXF trademark footer [https://github.com/k8gb-io/k8gb/pull/2077](https://github.com/k8gb-io/k8gb/pull/2077) 

    - Now we are back to 100 [https://clomonitor.io/projects/cncf/k8gb](https://clomonitor.io/projects/cncf/k8gb) 
- **Discuss**
    - next release - planned 20. Oct [https://github.com/k8gb-io/k8gb/issues/1992](https://github.com/k8gb-io/k8gb/issues/1992)  should we prepare release this week?
- **Issue Review**
    - [https://github.com/k8gb-io/k8gb/issues/2022](https://github.com/k8gb-io/k8gb/issues/2022) anything to add here? should we close?
    - [https://github.com/k8gb-io/k8gb/issues/1973](https://github.com/k8gb-io/k8gb/issues/1973) is implemented - double check the OCI repo state, then we can close
    - [https://github.com/k8gb-io/k8gb/issues/1906](https://github.com/k8gb-io/k8gb/issues/1906) Due Diligence Technical Review readiness
        - Important for Incubating
        - Farhan is already working on subitems
    - Sudhamsh question in slack [https://cloud-native.slack.com/archives/C021P656HGB/p1760512987294269](https://cloud-native.slack.com/archives/C021P656HGB/p1760512987294269)
        - 2 ingresses suggested
- **Community Update**
    - several talks submitted for KubeCon EU 🤞
- **Action Items**
    - mbcp case study promotion - intro video 
    - release

</details><details><summary><strong>Oct 1, 2025 #78</strong></summary>

## Oct 1, 2025 #78

Recording: [https://zoom.us/rec/share/hOAjb0ylEzCQZUFOhro7Y20rU2uFAG7\_p9SELmst9wMzMxxWwJ9-Q8aYnYg6JKI-.ncvJ4qR8n8DENCat](https://zoom.us/rec/share/hOAjb0ylEzCQZUFOhro7Y20rU2uFAG7_p9SELmst9wMzMxxWwJ9-Q8aYnYg6JKI-.ncvJ4qR8n8DENCat) 

Attendees:

- André
- Yury
- Bradley

Agenda:

- **News**
    - stuff from last time
    - [Document and test GCP Cloud DNS Provider Integration for K8GB](https://github.com/k8gb-io/k8gb/pull/2065)
- **Discuss**
    - 
- **Issue Review**
    - [https://github.com/k8gb-io/k8gb/issues/1872](https://github.com/k8gb-io/k8gb/issues/1872)
- **Community Update**
    - 
- **Action Items**
    - stars
    - issue review
    - mbcp case study promotion - intro video 
    - next release - 15. Oct [https://github.com/k8gb-io/k8gb/issues/1992](https://github.com/k8gb-io/k8gb/issues/1992) 
    - [https://github.com/cncf/toc/issues/1472](https://github.com/cncf/toc/issues/1472) was submitted nearly one year ago - it needs a general "another look" 
        - added / updated:
            - release info, adopters
            - Additional Information section:
                - [millennium-bcp case study](https://www.cncf.io/case-studies/millennium-bcp/)
                - [presentations-featuring-k8gb](https://www.k8gb.io/#presentations-featuring-k8gb)
                - [online-publications-featuring-k8gb](https://www.k8gb.io/#online-publications-featuring-k8gb)
                - [books-featuring-k8gb](https://www.k8gb.io/#and-even-books-featuring-k8gb)

</details><details><summary><strong>Sep 17, 2025 #77</strong></summary>

## Sep 17, 2025 #77

Agenda:

- **News**
    - 🎉 **MBCP case study** live [https://www.linkedin.com/posts/k8gb\_millennium-bcp-activity-7373817349780475904-szXu/](https://www.linkedin.com/posts/k8gb_millennium-bcp-activity-7373817349780475904-szXu/) 🎉
        - zero downtime in regional failover tests
        - 70% faster incident response
        - 99.99% uptime for banking apps
    - **Global Blue/Green Deployments with Crossplane v2 and k8gb** [https://tinyurl.com/ysd2y8jb](https://tinyurl.com/ysd2y8jb) - originally from KubeCon Hong Kong - Works across geo-distributed clusters. Automatic failover.
        - ✅ Namespace-scoped resources
        - ✅ Embedded KCL functions
        - ✅ Auto management policy switching based on Gslb resource health 
    - **2 k8gb presentations at ContainerDays**: [https://tinyurl.com/4b5ts48n](https://tinyurl.com/4b5ts48n)
        - **Kubermatic** "Evaluating Global Load Balancing Options for Kubernetes in Practice”
            - real-world use cases / how industry leaders are tackling multi-cluster networking.
        - **Prodyna** "Ensuring high availability with global load balancing in Kubernetes"
    - KubeCon NA Atlanta k8gb approved activities
        - Contribfest | Lightning Talk | Kiosk
    - 🎉**&#32;Community Infoblox PR** [https://github.com/k8gb-io/k8gb/pull/2058](https://github.com/k8gb-io/k8gb/pull/2058) 🎉
    - anything else from last time?
- **Community Update**
    - next in line for assignment [https://github.com/orgs/cncf/projects/27/views/9](https://github.com/orgs/cncf/projects/27/views/9) (cncf toc board-moving levels)
        - yury: [https://github.com/cncf/tag-security/pull/1446](https://github.com/cncf/tag-security/pull/1446) (k8gb security self-assessment questions)

</details><details><summary><strong>Sep 3, 2025 #76</strong></summary>

## Sep 3, 2025 #76

Recording: [https://zoom.us/rec/share/xprgYvNjRfW-G6Ne69ui0WvK0OylzWtlN\_3\_Hxm1-U8TvQCbwZOoW7YyrIbyDVev.n\_gUozC5BmVeYAil](https://zoom.us/rec/share/xprgYvNjRfW-G6Ne69ui0WvK0OylzWtlN_3_Hxm1-U8TvQCbwZOoW7YyrIbyDVev.n_gUozC5BmVeYAil) 

Attendees:

- Andre
- Bradley
- Florian
- Yury

Agenda:

- **News**
    - Load balancer service integration: [https://github.com/k8gb-io/k8gb/pull/2029](https://github.com/k8gb-io/k8gb/pull/2029)
    - Deprecated GSLB configuration via annotations: [https://github.com/k8gb-io/k8gb/pull/2043](https://github.com/k8gb-io/k8gb/pull/2043)
    - Concluded migration to external dns upstream helm chart:
        - Azure: [https://github.com/k8gb-io/k8gb/pull/2046](https://github.com/k8gb-io/k8gb/pull/2046)
        - Cloudflare: [https://github.com/k8gb-io/k8gb/pull/2045](https://github.com/k8gb-io/k8gb/pull/2045)
        - NS1: [https://github.com/k8gb-io/k8gb/pull/2042](https://github.com/k8gb-io/k8gb/pull/2042)
        - RFC2136: [https://github.com/k8gb-io/k8gb/pull/2047](https://github.com/k8gb-io/k8gb/pull/2047)
    - **Follow-up**
    - Review of [https://github.com/k8gb-io/k8gb/pull/2019](https://github.com/k8gb-io/k8gb/pull/2019) helm oci repo
    - KubeCon NA Atlanta k8gb approved activities
        - Contribfest
        - Lightning Talk
        - Kiosk
    - anything else from last time?
- **Discuss**
    - 
- **Community Update**
    - Linux Foundation case study with MBCP in final review
    - next in line for assignment [https://github.com/orgs/cncf/projects/27/views/9](https://github.com/orgs/cncf/projects/27/views/9) (cncf toc board-moving levels)
        - yury: [https://github.com/cncf/tag-security/pull/1446](https://github.com/cncf/tag-security/pull/1446) (k8gb security self-assessment questions)
- **Action Items**
    - 
</details><details><summary><strong>Aug 20, 2025 #75</strong></summary>

## Aug 20, 2025 #75

Recording: [https://zoom.us/rec/share/CbRTBkqcnsUp6rGof\_fUj5hhggCk5UBBk1s0YhdtFI\_IHFiT3QXvnzdDIT8h1KF9.D0hWEUTZb\_dGJtY1](https://zoom.us/rec/share/CbRTBkqcnsUp6rGof_fUj5hhggCk5UBBk1s0YhdtFI_IHFiT3QXvnzdDIT8h1KF9.D0hWEUTZb_dGJtY1)

Attendees:

- Andre
- Yury

Agenda:

- **News**
    - k8gb project table at KubeCon India was a blast [https://www.linkedin.com/posts/yurytsarev\_great-afternoon-at-the-k8gb-project-table-activity-7359639781858779138-HtQy](https://www.linkedin.com/posts/yurytsarev_great-afternoon-at-the-k8gb-project-table-activity-7359639781858779138-HtQy?utm_source=share&utm_medium=member_desktop&rcm=ACoAAASmnkkBA2oCkuBDmaZLvqDmY57S7LCBuh8)
    - KubeCon NA Atlanta k8gb approved activities
        - Contribfest
        - Lightning Talk
        - Kiosk
    - next in line for assignment [https://github.com/orgs/cncf/projects/27/views/9](https://github.com/orgs/cncf/projects/27/views/9) (cncf toc board-moving levels)
        - yury: [https://github.com/cncf/tag-security/pull/1446](https://github.com/cncf/tag-security/pull/1446) (k8gb security self-assessment questions)
- **Follow-up**
    - Linux Foundation case study with MBCP in final review
    - anything else from last time?
    - [https://github.com/k8gb-io/k8gb/pull/2017](https://github.com/k8gb-io/k8gb/pull/2017) Add helm support for geo data fields in coredns cm
        - Merged! Thanks Kelvin!
- **Discuss**
    - Terratest to Chainsaw migration
        - Migrate more test scenarios to chainsaw
    - Next release
        - Service integration
            - hostname question
        - externaldns helm chart values
            - in v0.15.0 we migrated route53
            - Migrate the rest of the values to official helm chart style
        - Gateway support? most probably postponed to next-next release given limited capacity
- **Community Update**
    - 
- **Action Items**
    - Review [https://github.com/k8gb-io/k8gb/pull/2019](https://github.com/k8gb-io/k8gb/pull/2019) helm oci repo

</details><details><summary><strong>Aug 6, 2025 #74</strong></summary>

## Aug 6, 2025 #74

Recording: [https://zoom.us/rec/share/oP\_773GXp9BYyU9lKAIFXXgBOms5OJychdD0l0bTGLIuTWIAR\_S4vpMRpa3yQZNR.0R5jww1eO274U5bd](https://zoom.us/rec/share/oP_773GXp9BYyU9lKAIFXXgBOms5OJychdD0l0bTGLIuTWIAR_S4vpMRpa3yQZNR.0R5jww1eO274U5bd) 

Attendees:

- Andre
- Farhan

Agenda:

- **News**
    - Webinar on [v0.15 release](https://github.com/k8gb-io/k8gb/releases/tag/v0.15.0): [https://www.youtube.com/watch?v=Jvro15FkudY](https://www.youtube.com/watch?v=Jvro15FkudY)
    - New Adopter! [https://github.com/k8gb-io/k8gb/pull/2014](https://github.com/k8gb-io/k8gb/pull/2014)
    - Yury in Kubecon India
- **Follow-up**
    - First time contributor, add tolerations: [https://github.com/k8gb-io/k8gb/pull/2009](https://github.com/k8gb-io/k8gb/pull/2009)
    - CoreDNS now exposed as load balancer service in local setup: [https://github.com/k8gb-io/k8gb/pull/1828](https://github.com/k8gb-io/k8gb/pull/1828)
    - Configured renovate to bump more versions: [https://github.com/k8gb-io/k8gb/pull/2007](https://github.com/k8gb-io/k8gb/pull/2007)
- **Discuss**
    - new issue - [https://github.com/k8gb-io/k8gb/issues/2015](https://github.com/k8gb-io/k8gb/issues/2015) - not possible to configure geoip databases
    - Andre's Priorities:
        - migrate more tests to chainsaw to reduce testing time to under 10 minutes, then add API Gateway integration
        - CoreDNS serving NS and glue records
- **Community Update**
    - 
- **Action Items**
    - 

</details><details><summary><strong>Jul 23, 2025 #73</strong></summary>

## Jul 23, 2025 #73

Recording: [https://zoom.us/rec/share/ucm5Vho6-gghLObhpb5lo0-Cl\_nN8IXkbNQF0-UFfhO-Qq00-fNQmtI8b-FzUgsY.sE9BQxzE2gMRfvJe](https://zoom.us/rec/share/ucm5Vho6-gghLObhpb5lo0-Cl_nN8IXkbNQF0-UFfhO-Qq00-fNQmtI8b-FzUgsY.sE9BQxzE2gMRfvJe) 

Attendees:

- Bradley
- Yury
- Andre
- Michal
- Dinar

Agenda:

- **News**
    - [v0.15 released](https://github.com/k8gb-io/k8gb/releases/tag/v0.15.0) 🎉 … starting quarterly release cycle 🎉
        - changes from this release
        - social posts done … blog coming (maybe also webinar)
    - New website - community contribution 🎉
    - Bradley unavailable in August and November 12
        - need someone to do these calls
            - ✅ 6 August - Andre taking
            - ✅ 20 August - Yury taking
            - ✅ 12 November - Yury
    - moved these notes from google docs to proton docs
- **Follow-up**
    - ✅ FOSSA license scan issue
    - Incubation application: [https://github.com/orgs/cncf/projects/27/views/9](https://github.com/orgs/cncf/projects/27/views/9) 
        - waiting for assignment
        - security self-assessment in progress [https://github.com/cncf/tag-security/pull/1446#issuecomment-3091833870](https://github.com/cncf/tag-security/pull/1446#issuecomment-3091833870) (from [https://github.com/cncf/tag-security/issues/1441#issuecomment-2621683616](https://github.com/cncf/tag-security/issues/1441#issuecomment-2621683616)) 
    - ✅ Linux Foundation case study with MBCP in final review

- **Discuss**
    - 
- **Community Update**
    - 
- **Action Items**
    - 

</details><details><summary><strong>Jul 9, 2025 #72</strong></summary>

## Jul 9, 2025 #72

Recording: [https://zoom.us/rec/share/Ecg3J6p5vpzIDvFZ0nflOOZJt8SX-drS\_TAXMzPiZBB4ODJd0s3qRLNnsBtxkjH1.DDap105UZoYfHUrN](https://zoom.us/rec/share/Ecg3J6p5vpzIDvFZ0nflOOZJt8SX-drS_TAXMzPiZBB4ODJd0s3qRLNnsBtxkjH1.DDap105UZoYfHUrN) 

Attendees:

- Bradley
- Farhan
- Yury

Agenda:

- **News**
    - Yury has applied to KubeCon NA: lightning talk, kiosk, and Contribfest
    - Bradley unavailable in August and November 12
        - need someone to do these calls
- **Follow-up**
    - FOSSA license scan issue: [https://cncfservicedesk.atlassian.net/servicedesk/customer/portal/1/CNCFSD-2819](https://cncfservicedesk.atlassian.net/servicedesk/customer/portal/1/CNCFSD-2819) 
        - officially false positive, but having trouble marking it as such
    - Incubation application: [https://github.com/orgs/cncf/projects/27/views/9](https://github.com/orgs/cncf/projects/27/views/9) 
        - waiting for assignment
    - Linux Foundation case study with MBCP
        - pinged again today

- **Discuss**
    - Anything is blocking v0.15 release?

- **Community Update**
    - 

- **Action Items**
    - 

</details><details><summary><strong>Jun 25, 2025 #71</strong></summary>

## Jun 25, 2025 #71

Recording: [https://zoom.us/rec/share/\_DUR3Z6qC03cJSyNYR3tT-jfyoOVWuCvzjN7l15sxUIB69fxbguEuG\_2nT6\_0\_15.jt4UABLOEtRUNivK](https://zoom.us/rec/share/_DUR3Z6qC03cJSyNYR3tT-jfyoOVWuCvzjN7l15sxUIB69fxbguEuG_2nT6_0_15.jt4UABLOEtRUNivK) 

Attendees:

- Yury
- Michal
- Andre

Agenda:

- **News**
    - KubeCon China Crossplane + k8gb talk 
        - Crossplane Configuration for global app with k8gb ( [Yury Tsarev](mailto:xnullz@gmail.com))  reference example
        - [https://www.youtube.com/watch?v=L9mRWljLnzw](https://www.youtube.com/watch?v=L9mRWljLnzw) recording
        - [https://github.com/k8gb-io/k8gb/tree/master/docs/examples/crossplane/globalapp](https://github.com/k8gb-io/k8gb/tree/master/docs/examples/crossplane/globalapp) 
        - Both are linked on the frontpage 
[https://www.k8gb.io/#presentations-featuring-k8gb](https://www.k8gb.io/#presentations-featuring-k8gb) [https://www.k8gb.io/docs/examples/crossplane/globalapp/](https://www.k8gb.io/docs/examples/crossplane/globalapp/) 
    - KCD Czech & Slovak
        - [https://docs.google.com/presentation/d/1ExMDY3tFAFi8RT7d-LDyxcb0ueyIRsLNbLLPu-7btCE/edit?usp=sharing](https://docs.google.com/presentation/d/1ExMDY3tFAFi8RT7d-LDyxcb0ueyIRsLNbLLPu-7btCE/edit?usp=sharing) deck
- **Follow-up**
    - Action Items from last(-ish) time
        - [license scan issue](https://app.fossa.com/projects/custom%2B162%2Fgithub.com%2Fk8gb-io%2Fk8gb/refs/branch/master/5b6f7df90d72afdbfb194a6f977c269d869a4732/preview) → [ticket](https://cncfservicedesk.atlassian.net/servicedesk/customer/portal/1/CNCFSD-2819)
        - Due Diligence review issues from [the DTR](https://docs.google.com/document/d/1nJB-A2RbWEcfxpINa1n1bpuYsQTuvavtDl1pNtNTZs4/edit?tab=t.0#heading=h.iu3tzv2c4euc) are [github issues](https://github.com/k8gb-io/k8gb/issues/1906)**&#32;…** Governance WIP
    - Other
        - [https://github.com/cncf/toc/issues/1472](https://github.com/cncf/toc/issues/1472) still “[ready for assignment](https://github.com/orgs/cncf/projects/27/views/1?filterQuery=status%3A%22Ready+for+assignment%22)” - 2 in progress now, 1 more in front of us
        - did not win [end-user submission for KubeCon Japan](https://www.linkedin.com/posts/cloud-native-computing-foundation_kubecon-activity-7319353673535348736-slHY/?utm_source=share&utm_medium=member_android&rcm=ACoAABaGd30BZimWi9WsuWmSzyrjtQD9mwDhgUg), but Linux Foundation folks want to do a case study with MBCP 🚀0.15-rc1 release

- **Discuss**
    - (missing) Documentation for release and new features

- **Community Update**
    - Working on Medium post / short intro video “What is k8gb?”
    - Working on basic audience growth

- **Action Items**
    - Session to pass through issues? Is everything still valid / relevant?
    - Make 0.15 Release ASAP
    - Regular (quarterly?) release schedule proposal

</details><details><summary><strong>May 28, 2025 #70</strong></summary>

## May 28, 2025 #70

Recording: [https://zoom.us/rec/share/UIY\_WtVR70bdi59KBrF6Rrcwz0MkniO3XypN95Im0xugB6xFOWd\_JBWbDPRyfXfq.5RTBi5JE0pCHOTMD](https://zoom.us/rec/share/UIY_WtVR70bdi59KBrF6Rrcwz0MkniO3XypN95Im0xugB6xFOWd_JBWbDPRyfXfq.5RTBi5JE0pCHOTMD) 

Attendees:

- Bradley A
- Yury T
- Michal K
- Dinar V
- Andre A

Agenda:

- **Follow-up**
    - Action Items from last(-ish) time
        - ~~[x] do we still need a call to discuss gslb resource in the context of adding gw api support?~~
        - [license scan issue](https://app.fossa.com/projects/custom%2B162%2Fgithub.com%2Fk8gb-io%2Fk8gb/refs/branch/master/5b6f7df90d72afdbfb194a6f977c269d869a4732/preview) → [ticket](https://cncfservicedesk.atlassian.net/servicedesk/customer/portal/1/CNCFSD-2819)
        - Due Diligence review issues from [the DTR](https://docs.google.com/document/d/1nJB-A2RbWEcfxpINa1n1bpuYsQTuvavtDl1pNtNTZs4/edit?tab=t.0#heading=h.iu3tzv2c4euc) are [github issues](https://github.com/k8gb-io/k8gb/issues/1906)**&#32;…** Governance WIP
    - Other
        - [https://github.com/cncf/toc/issues/1472](https://github.com/cncf/toc/issues/1472) still “[ready for assignment](https://github.com/orgs/cncf/projects/27/views/1?filterQuery=status%3A%22Ready+for+assignment%22)” - 3 in progress now, 4 more in front of us
        - did not win [end-user submission for KubeCon Japan](https://www.linkedin.com/posts/cloud-native-computing-foundation_kubecon-activity-7319353673535348736-slHY/?utm_source=share&utm_medium=member_android&rcm=ACoAABaGd30BZimWi9WsuWmSzyrjtQD9mwDhgUg), but Linux Foundation folks want to do a case study with MBCP 🚀
    - 0.15-rc1 release

- **Discuss**
    - (missing) Documentation for release and new features
    - k8gb Project Table at KubeCon India. Anybody joins?
    - Crossplane+k8gb talk at KubeCon China [https://kccncchn2025.sched.com/event/1x5jS/resilient-multiregion-global-control-planes-with-crossplane-and-k8gb-yury-tsarev-steven-borrelli-upbound?iframe=yes&w=100%&sidebar=yes&bg=no](https://kccncchn2025.sched.com/event/1x5jS/resilient-multiregion-global-control-planes-with-crossplane-and-k8gb-yury-tsarev-steven-borrelli-upbound?iframe=yes&w=100%&sidebar=yes&bg=no) - I would love to base the demo on 0.15 release(Yury)
    - First feedback on breaking changes [https://github.com/k8gb-io/k8gb/discussions/364#discussioncomment-13273055](https://github.com/k8gb-io/k8gb/discussions/364#discussioncomment-13273055) 




- **Community Update**
    - Working on Medium post / short intro video “What is k8gb?”
    - Working on basic audience growth

- **Action Items**
    - Bradley make slack post about release
    - Session to pass through issues? Is everything still valid / relevant?

</details><details><summary><strong>May 14, 2025 #69</strong></summary>

## May 14, 2025 #69

Recording: [https://zoom.us/rec/share/-a5TpXZtn3WTeWyWiA8SU\_hRPX4dZCvDEimxuQjS5jnUQmzCSsGsy5TUUOAzvHY.MHthdn\_nUsL5Lxsv](https://zoom.us/rec/share/-a5TpXZtn3WTeWyWiA8SU_hRPX4dZCvDEimxuQjS5jnUQmzCSsGsy5TUUOAzvHY.MHthdn_nUsL5Lxsv) 

Attendees:

- Bradley A
- Andre A
- Dinar V
- Victor L
- Michal K

Agenda:

- **Follow-up**
    - Action Items from last time
        - ⏳ set up call to discuss gslb resource in the context of adding gw api support
        - ⏳ look in to [license scan issue](https://app.fossa.com/projects/custom%2B162%2Fgithub.com%2Fk8gb-io%2Fk8gb/refs/branch/master/5b6f7df90d72afdbfb194a6f977c269d869a4732/preview)
            - [created a ticket](https://cncfservicedesk.atlassian.net/servicedesk/customer/portal/1/CNCFSD-2819)
        - ✅ convert DD review issues from [the DTR](https://docs.google.com/document/d/1nJB-A2RbWEcfxpINa1n1bpuYsQTuvavtDl1pNtNTZs4/edit?tab=t.0#heading=h.iu3tzv2c4euc) to github issues: [https://github.com/k8gb-io/k8gb/issues/1906](https://github.com/k8gb-io/k8gb/issues/1906) 
        - ✅ better / additional blog platform?
            - seems most relevant are medium and TNS
            - [https://thenewstack.io/](https://thenewstack.io/) ← submitted to be a contributor
- **Discuss**
    - [Kubenet](https://docs.google.com/document/u/0/d/1tn0h_1sDlzmEBQ7TkwouMUwy4IMPgJDixQ3q9Pk-AFY/mobilebasic) +/- k8gb
    - bug discovered + squashed yesterday 🎉
    - 0.15-rc1 release ASAP (maybe EOW) 🚀




- **Community Update**
    - no change since [last time](https://docs.google.com/document/d/1YdpEVhtyCKvwFtXR7cn1Kn2Xc_tdNskoFnhHmFPbtA4/edit?tab=t.0#bookmark=id.4831hz2h8l2w), except agreed to create customer case study with mbcp 🎉

- **Action Items**
    - 



</details><details><summary><strong>Apr 30, 2025 #68</strong></summary>

## Apr 30, 2025 #68

Recording: [https://zoom.us/rec/share/5dHZ\_ibeQ-yicNEgp1tCBu\_gwEg0PVw1P1yc2Ww7s0npCtYwRgEQdTJiZMVo3MPG.h-pRVeFk39TdX1MX](https://zoom.us/rec/share/5dHZ_ibeQ-yicNEgp1tCBu_gwEg0PVw1P1yc2Ww7s0npCtYwRgEQdTJiZMVo3MPG.h-pRVeFk39TdX1MX) 

Attendees:

- Bradley A
- Dinar V
- Victor L
- Yury T
- Michal K
- Bilal J

Agenda:

- **News**
    - Domain xfer complete 😅
    - GitHub ⭐ > 1000 🚀




- **Discuss**
    - CRD change for 1.0 ← vendor neutrality
    - Support for Gateway API - [https://github.com/k8gb-io/k8gb/issues/954](https://github.com/k8gb-io/k8gb/issues/954) 
        - question from community member
        - external DNS first
    - Gslb resource - is complicated, will be more so with GW API support
        - separate call to discuss




- **Issues**
    - relevant from last time?
        - license scan still failing: [https://app.fossa.com/projects/custom%2B162%2Fgithub.com%2Fk8gb-io%2Fk8gb/refs/branch/master/5b6f7df90d72afdbfb194a6f977c269d869a4732/preview](https://app.fossa.com/projects/custom%2B162%2Fgithub.com%2Fk8gb-io%2Fk8gb/refs/branch/master/5b6f7df90d72afdbfb194a6f977c269d869a4732/preview) 
            - seems to be false … need to look into it




-  **PRs** 
    - relevant from last time?
        - [https://github.com/k8gb-io/k8gb/pull/1889](https://github.com/k8gb-io/k8gb/pull/1889) 
        - [https://github.com/k8gb-io/k8gb/pull/1888](https://github.com/k8gb-io/k8gb/pull/1888) 
    - anything new?




- **Community Update**
    - [https://github.com/cncf/toc/issues/1472](https://github.com/cncf/toc/issues/1472) has been moved up to “[ready for assignment](https://github.com/orgs/cncf/projects/27/views/1?filterQuery=status%3A%22Ready+for+assignment%22)” - progress! (3 in progress now, 4 more in front of us?)
    - Due Diligence reviews have begun:
        - [Technical](https://docs.google.com/document/d/1nJB-A2RbWEcfxpINa1n1bpuYsQTuvavtDl1pNtNTZs4/edit?tab=t.0#heading=h.iu3tzv2c4euc) 24 April
        - Governance next TAG meeting on 8 May
        - 👆once we are assigned, reviews will go quickly because we have done these things early.
    - Social posts - last 2 weeks:
        - LI
            - [1000 stars](https://www.linkedin.com/posts/k8gb_github-k8gb-iok8gb-a-cloud-native-kubernetes-activity-7318268843762593794-Es4g/) 
            - > 2,700 impressions (^ 475%) - 77 clicks - 23 reactions
            - +21 followers (^ 24%)
    - ✅ Working on [end-user submission for KubeCon Japan](https://www.linkedin.com/posts/cloud-native-computing-foundation_kubecon-activity-7319353673535348736-slHY/?utm_source=share&utm_medium=member_android&rcm=ACoAABaGd30BZimWi9WsuWmSzyrjtQD9mwDhgUg) - submitted today - winner announced 7 May
        - this will also become a blog post
    - Working on Medium post “What is k8gb?”
    - Working on basic audience growth

- **Action Items**
    - set up call to discuss gslb [Bradley Andersen](mailto:bradley.d.andersen@gmail.com)
    - look in to license scan issue [https://app.fossa.com/projects/custom%2B162%2Fgithub.com%2Fk8gb-io%2Fk8gb/refs/branch/master/5b6f7df90d72afdbfb194a6f977c269d869a4732/preview](https://app.fossa.com/projects/custom%2B162%2Fgithub.com%2Fk8gb-io%2Fk8gb/refs/branch/master/5b6f7df90d72afdbfb194a6f977c269d869a4732/preview)  [Bradley Andersen](mailto:bradley.d.andersen@gmail.com)
    - need to convert DD review issues from [https://docs.google.com/document/d/1nJB-A2RbWEcfxpINa1n1bpuYsQTuvavtDl1pNtNTZs4/edit?tab=t.0#heading=h.iu3tzv2c4euc](https://docs.google.com/document/d/1nJB-A2RbWEcfxpINa1n1bpuYsQTuvavtDl1pNtNTZs4/edit?tab=t.0#heading=h.iu3tzv2c4euc) to github issues here [https://github.com/k8gb-io/k8gb/issues](https://github.com/k8gb-io/k8gb/issues)  [Bradley Andersen](mailto:bradley.d.andersen@gmail.com)
    - ✅ better / additional blog platform? [Bradley Andersen](mailto:bradley.d.andersen@gmail.com)
        - seems most relevant are medium and TNS
        - [https://thenewstack.io/](https://thenewstack.io/) ← submitted to be a contributor

</details><details><summary><strong>Apr 16, 2025 #67</strong></summary>

## Apr 16, 2025 #67

Recording: [https://zoom.us/rec/share/161iEONShSThMfrzq0N09sZ5dbQkc2q0iepjzsGOYpa-4c0cAL90BQfUoaTK7-E.Rr\_Zgp880InufJo0](https://zoom.us/rec/share/161iEONShSThMfrzq0N09sZ5dbQkc2q0iepjzsGOYpa-4c0cAL90BQfUoaTK7-E.Rr_Zgp880InufJo0) 

Attendees:

- Bradley
- Yury
- Michal
- Andre
- Victor Lu
    - gw api ← on roadmap
    - threat modeling 
        - might be useful: [https://github.com/cncf/tag-security/pull/1446/files](https://github.com/cncf/tag-security/pull/1446/files) 

Agenda:

- **News**
    - Domain xfer
        - for now, probably register / use alt like k8gb.dev [Yury Tsarev](mailto:xnullz@gmail.com)
    - CRD change for 1.0 ← vendor neutrality … but … domain 
    - KubeCon recap [https://www.linkedin.com/posts/k8gb\_kubecon-europe-2025-activity-7314962455611273216-FvC\_](https://www.linkedin.com/posts/k8gb_kubecon-europe-2025-activity-7314962455611273216-FvC_) and [https://medium.com/@kubernetesglobalbalancer/kubecon-europe-2025-e381016fb3ce](https://medium.com/@kubernetesglobalbalancer/kubecon-europe-2025-e381016fb3ce) 
        - rejekts talk [https://www.linkedin.com/posts/k8gb\_cloud-native-rejekts-europe-2025-the-nash-activity-7313111304569868288-zVD5](https://www.linkedin.com/posts/k8gb_cloud-native-rejekts-europe-2025-the-nash-activity-7313111304569868288-zVD5) 
        - project lightning talk [https://www.youtube.com/watch?v=YMyrcqZ2sbU](https://www.youtube.com/watch?v=YMyrcqZ2sbU) 
        - booth [https://www.linkedin.com/posts/k8gb\_kubecon-cloudnativecon-europe-2025-project-activity-7313484847719579648-02RV](https://www.linkedin.com/posts/k8gb_kubecon-cloudnativecon-europe-2025-project-activity-7313484847719579648-02RV) 
        - 👆 socials (X, LinkedIn, Medium)
    - Stars
        - 999 … swag for 1000th?
- **Issues**
    - relevant from last time?: [https://docs.google.com/document/d/1YdpEVhtyCKvwFtXR7cn1Kn2Xc\_tdNskoFnhHmFPbtA4/edit?tab=t.0#bookmark=id.u7du1bzev411](https://docs.google.com/document/d/1YdpEVhtyCKvwFtXR7cn1Kn2Xc_tdNskoFnhHmFPbtA4/edit?tab=t.0#bookmark=id.u7du1bzev411) 
    - anything new?
-  **PRs** 
    - relevant from last time?: [https://docs.google.com/document/d/1YdpEVhtyCKvwFtXR7cn1Kn2Xc\_tdNskoFnhHmFPbtA4/edit?tab=t.0#bookmark=id.m926mm3gq77](https://docs.google.com/document/d/1YdpEVhtyCKvwFtXR7cn1Kn2Xc_tdNskoFnhHmFPbtA4/edit?tab=t.0#bookmark=id.m926mm3gq77)
    - External DNS for aws route53 integration: [https://github.com/k8gb-io/k8gb/pull/1856](https://github.com/k8gb-io/k8gb/pull/1856)
    - Deprecating dnsZone and edgeDnsZone: [https://github.com/k8gb-io/k8gb/pull/1876](https://github.com/k8gb-io/k8gb/pull/1876) (fixes [https://github.com/k8gb-io/k8gb/issues/1858](https://github.com/k8gb-io/k8gb/issues/1858))
    - anything new?
- **Community Update**
    - [https://github.com/cncf/toc/issues/1472](https://github.com/cncf/toc/issues/1472) has been moved up to “ready for assignment” - progress!
- **Action Items**
    - domain xfer - Yury 

</details><details><summary><strong>Apr 2, 2025 #66</strong></summary>

## Apr 2, 2025 #66

Recording: [https://zoom.us/rec/share/Edc-il8vKgBKWla117P3cD52-IIKEJK9i5bYwPIRIbCpiTbG8e8q-k2bx0pV1Aqd.r5LNKLB0XA-tO4IM](https://zoom.us/rec/share/Edc-il8vKgBKWla117P3cD52-IIKEJK9i5bYwPIRIbCpiTbG8e8q-k2bx0pV1Aqd.r5LNKLB0XA-tO4IM) 

Attendees:

- Andre Aguas
- Michal
- Dinar

Agenda:

- **News**
    - Bradley, Nuno and Yury at Kubecon!
        - Project Lightning Talk: [https://sched.co/1tcvf](https://sched.co/1tcvf) - Tuesday 1. April 12.20 - 12.25
            - slides started: [https://docs.google.com/presentation/d/143ZAoWAoqlD4C0X1KHmiVr4JRdoZiQbYjMhLHM47LuA/edit?usp=sharing](https://docs.google.com/presentation/d/143ZAoWAoqlD4C0X1KHmiVr4JRdoZiQbYjMhLHM47LuA/edit?usp=sharing) 
        - Project Pavilion: Thursday Afternoon 3. April Kiosk 4B 14:00 - 17:00

- **Issues**
    - Coredns ip address lookup:**&#32;**[https://github.com/k8gb-io/k8gb/pull/1864#issuecomment-2771509977](https://github.com/k8gb-io/k8gb/pull/1864#issuecomment-2771509977)
    - dnsZones breaking change: [https://github.com/k8gb-io/k8gb/issues/1858](https://github.com/k8gb-io/k8gb/issues/1858)
    - GSLB name reverted after change: [https://github.com/k8gb-io/k8gb/issues/1875](https://github.com/k8gb-io/k8gb/issues/1875)
    - K8gb with cert manager: [https://github.com/k8gb-io/k8gb/issues/1872](https://github.com/k8gb-io/k8gb/issues/1872) -> [https://github.com/AbsaOSS/cert-manager-webhook-externaldns](https://github.com/AbsaOSS/cert-manager-webhook-externaldns)
    - Reported issue: k8gb does not watch endpoints, the health of the application is only verified at reconciliation interval; Idea: can we use dynamic watchers? How to find the GSLB resource from an endpoint?

-  **PRs** 
    - External-dns refactoring and update (waiting for review)
        - [https://github.com/k8gb-io/k8gb/pull/1856](https://github.com/k8gb-io/k8gb/pull/1856) 
        - [https://github.com/k8gb-io/k8gb/pull/1762](https://github.com/k8gb-io/k8gb/pull/1762)
    - 
 

- **Discussion**
    - V1.0?




- **Follow-up**
    - 




- **Community Update**
    - Holding pattern on [https://github.com/cncf/toc/issues/1472](https://github.com/cncf/toc/issues/1472). 




- **Action Items**
    - 

</details><details><summary><strong>Mar 19, 2025 #65</strong></summary>

## Mar 19, 2025 #65

Recording: [https://zoom.us/rec/share/iwjXm4jjB1CVeADKzlxm4yiJWL0NE5z3j0ar-ors8cMaJjw4gm-DrmJGLUg6YW04.UZXrvtcU-8HXJrv3](https://zoom.us/rec/share/iwjXm4jjB1CVeADKzlxm4yiJWL0NE5z3j0ar-ors8cMaJjw4gm-DrmJGLUg6YW04.UZXrvtcU-8HXJrv3) 

Attendees:

- Bradley Andersen	
- Dinar Valeev	
- Michal Kuřítka	
- Yury Tsarev

Agenda:

- **News**
    - Kubecon in 2 weeks!
        - Project Lightning Talk: [https://sched.co/1tcvf](https://sched.co/1tcvf) - Tuesday 1. April 12.20 - 12.25
            - slides started: [https://docs.google.com/presentation/d/143ZAoWAoqlD4C0X1KHmiVr4JRdoZiQbYjMhLHM47LuA/edit?usp=sharing](https://docs.google.com/presentation/d/143ZAoWAoqlD4C0X1KHmiVr4JRdoZiQbYjMhLHM47LuA/edit?usp=sharing) 
        - Project Pavilion: Thursday Afternoon 3. April Kiosk 4B 14:00 - 17:00
    - Community coverage by Gerardo: [https://www.cncf.io/blog/2025/02/19/exploring-multi-cluster-fault-tolerance-with-k8gb/](https://www.cncf.io/blog/2025/02/19/exploring-multi-cluster-fault-tolerance-with-k8gb/) 🎉
    - docs revamp: [https://www.youtube.com/watch?v=8eROmJg71gw](https://www.youtube.com/watch?v=8eROmJg71gw) 

- **Issues**
    - 




-  **PRs** 
    - External-dns refactoring and update
        - [https://github.com/k8gb-io/k8gb/pull/1856](https://github.com/k8gb-io/k8gb/pull/1856) 
        - [https://github.com/k8gb-io/k8gb/pull/1762](https://github.com/k8gb-io/k8gb/pull/1762) 
    - 
 

- **Discussion**
    - Skip next community meeting (2. April) in favor of Kubecon?
        - bradley and yury in london - others may still run it if they want to 👍
    - V1.0?




- **Follow-up**
    - 




- **Community Update**
    - Holding pattern on [https://github.com/cncf/toc/issues/1472](https://github.com/cncf/toc/issues/1472). 




- **Action Items**
    - 

</details><details><summary><strong>Mar 5, 2025 #64</strong></summary>

## Mar 5, 2025 #64

Recording: [https://zoom.us/rec/share/DW3rMirpXiNBxCFqbornbxwfsWc1Xx33T9E\_luZNIn-GZl-gTXUgepdoDESNdas.peXLdZOGD\_mUIItO](https://zoom.us/rec/share/DW3rMirpXiNBxCFqbornbxwfsWc1Xx33T9E_luZNIn-GZl-gTXUgepdoDESNdas.peXLdZOGD_mUIItO) 

Attendees:

- Yury
- Dinar
- Michal

Agenda:

- **News**
    - Community coverage from Scott, author of ‘Kubernetes - An Enterprise Guide’ book [https://www.linkedin.com/posts/scottsurovich\_kubernetes-an-enterprise-guide-third-edition-activity-7302863671758655489-TTJK/?utm\_source=share&utm\_medium=member\_desktop&rcm=ACoAAASmnkkBA2oCkuBDmaZLvqDmY57S7LCBuh8](https://www.linkedin.com/posts/scottsurovich_kubernetes-an-enterprise-guide-third-edition-activity-7302863671758655489-TTJK/?utm_source=share&utm_medium=member_desktop&rcm=ACoAAASmnkkBA2oCkuBDmaZLvqDmY57S7LCBuh8) 

- **Issues**
    - [Create test setup for upstream DNS providers :: 1772](https://github.com/k8gb-io/k8gb/issues/1772) Create test setup for upstream DNS providers - WIP - will be Crossplane-based
    - [Revamp website :: 1778](https://github.com/k8gb-io/k8gb/issues/1778) Revamp website - WIP - Gerardo has a PoC:
        - [https://www.youtube.com/watch?v=8eROmJg71gw](https://www.youtube.com/watch?v=8eROmJg71gw)
    - [Allow to override DNS record targets :: 1800](https://github.com/k8gb-io/k8gb/issues/1800#issuecomment-2588261322) partially solved, unclear what to do with CNAME part
    - [HA of k8gb service - questions :: 1035](https://github.com/k8gb-io/k8gb/issues/1035#issuecomment-2582936204)
    - New issues from the community
        - [https://github.com/k8gb-io/k8gb/issues/1833](https://github.com/k8gb-io/k8gb/issues/1833) Ignore 'mesh' gateway when counting referenced gateways in VirtualService
        - [https://github.com/k8gb-io/k8gb/issues/1837](https://github.com/k8gb-io/k8gb/issues/1837) k8gb incorrectly updates Nameserver A record TTL when multiple GSLB objects exist with different TTLs
        - [https://github.com/k8gb-io/k8gb/issues/1840](https://github.com/k8gb-io/k8gb/issues/1840) I do not hava public dns server,I just want to test locally with 3cluster ,how to instalk




-  **PRs** 
    - [Switch external-dns back to official upstream :: 1762](https://github.com/k8gb-io/k8gb/pull/1762) Manual regression tests ongoing but the pipeline is red (Yury)
        - release after regression testing this issue 👍
        - Even trivial environment exposure like AWS\_DEFAULT\_REGION is problematic.
        - We agreed with Andre to refactor external-dns templating or switch to the upstream chart, see discussion at [https://cloud-native.slack.com/archives/C021P656HGB/p1741031785519899?thread\_ts=1741002945.112709&cid=C021P656HGB](https://cloud-native.slack.com/archives/C021P656HGB/p1741031785519899?thread_ts=1741002945.112709&cid=C021P656HGB) 
            - Webhook based providers as an option (outside of big3) 
    - [Add support for multiple zones :: 1774](https://github.com/k8gb-io/k8gb/pull/1774) [kuritka@gmail.com](mailto:kuritka@gmail.com) is working on refactoring
        - PR is merged [https://github.com/k8gb-io/k8gb/pull/1845](https://github.com/k8gb-io/k8gb/pull/1845) 
        - Related [https://github.com/k8gb-io/k8gb/pull/1848](https://github.com/k8gb-io/k8gb/pull/1848) is merged
        - [kuritka@gmail.com](mailto:kuritka@gmail.com) to present the changes in community meeting
 

- **Discussion**
    - Newest multizone support and related challenges, presented by Michal
    - [https://github.com/k8gb-io/k8gb/discussions/1839](https://github.com/k8gb-io/k8gb/discussions/1839) k8gb CoreDNS exposed over hostnetwork ingress - programmatic lookup

- **Follow-up**
    - [Fossa is failing :: 1797](https://github.com/k8gb-io/k8gb/issues/1797) FOSSA fails again, we need to sort out the access to see the actual error [Bradley Andersen](mailto:bradley.d.andersen@gmail.com) look in to it
        - Yury got the access, not yet fully sure how to get properly solve the flagged dependency
        - [Create test setup for Azure DNS integration :: 1773](https://github.com/k8gb-io/k8gb/pull/1773#issuecomment-2606083136) Azure test setup and alternative implementation with Crossplane
            - I used Crossplane for route53 regression test too. Tests need polishing and documentation before I send them for review(Yury)
            - Related [https://github.com/k8gb-io/k8gb/pull/1849](https://github.com/k8gb-io/k8gb/pull/1849) for AWS from Andre




- **Community Update**
    - Making solid progress toward incubating 🎉




- **Action Items**
    - Go through the [https://github.com/k8gb-io/k8gb/issues](https://github.com/k8gb-io/k8gb/issues) and [https://github.com/k8gb-io/k8gb/discussions](https://github.com/k8gb-io/k8gb/discussions) queue, provide initial response 
    - Refactor external-dns deployment template (Andre + Yury)

</details><details><summary><strong>Feb 19, 2025 #63</strong></summary>

## Feb 19, 2025 #63

Recording: 

Attendees:

- Bradley
- Michal

Agenda:

- **News**
    - KubeCon London (attending: Bradley, Yury)
        - project pavilion
        - lightning talk
        - rejekts talk accepted 

- **Issues**
    - [Create test setup for upstream DNS providers :: 1772](https://github.com/k8gb-io/k8gb/issues/1772) Create test setup for upstream DNS providers - WIP - will be Crossplane-based
    - [Revamp website :: 1778](https://github.com/k8gb-io/k8gb/issues/1778) Revamp website - WIP - Gerardo has a PoC:
        - [https://www.youtube.com/watch?v=8eROmJg71gw](https://www.youtube.com/watch?v=8eROmJg71gw)
    - [Allow to override DNS record targets :: 1800](https://github.com/k8gb-io/k8gb/issues/1800#issuecomment-2588261322) partially solved, unclear what to do with CNAME part
    - [HA of k8gb service - questions :: 1035](https://github.com/k8gb-io/k8gb/issues/1035#issuecomment-2582936204)




-  **PRs** 
    - [Switch external-dns back to official upstream :: 1762](https://github.com/k8gb-io/k8gb/pull/1762) Manual regression tests ongoing but the pipeline is red (Yury)
        - release after regression testing this issue 👍
    - [Add support for multiple zones :: 1774](https://github.com/k8gb-io/k8gb/pull/1774) [kuritka@gmail.com](mailto:kuritka@gmail.com) is working on refactoring
        - look for PR soon 🎉
 

- **Discussion**
    - 

- **Follow-up**
    - [Fossa is failing :: 1797](https://github.com/k8gb-io/k8gb/issues/1797) FOSSA fails again, we need to sort out the access to see the actual error [Bradley Andersen](mailto:bradley.d.andersen@gmail.com) look in to it
        - in progress: 
            - ✅ team setup issues - new invites sent ← once this is worked out, we more or less self-manage
            - ✅ team members vs org relationships issues - (ex: Yury moving account to k8gb)
            - associating team with correct repo (absa/k8gb vs k8gb/k8gb)
        - [Create test setup for Azure DNS integration :: 1773](https://github.com/k8gb-io/k8gb/pull/1773#issuecomment-2606083136) Azure test setup and alternative implementation with Crossplane
            - I used Crossplane for route53 regression test too. Tests need polishing and documentation before I send them for review(Yury)




- **Community Update**
    - Making solid progress toward incubating 🎉




- **Action Items**
    -  

</details><details><summary><strong>Feb 5, 2025 #62</strong></summary>

## Feb 5, 2025 #62

Recording: [https://zoom.us/rec/share/NsUVRAcNDwJH3KFDejSCm0bPOKf6\_Y58gOwe8OgevE\_hcfO9BJSU1wRVmhkOozKa.0H3jIpoZG4mvZIae](https://zoom.us/rec/share/NsUVRAcNDwJH3KFDejSCm0bPOKf6_Y58gOwe8OgevE_hcfO9BJSU1wRVmhkOozKa.0H3jIpoZG4mvZIae) 

Attendees:

- Yury
- Dinar

Agenda:

- **News**
    - [https://github.com/k8gb-io/k8gb/pull/1824](https://github.com/k8gb-io/k8gb/pull/1824) Darede as a new official Adopter!
    - CNCF Network TAG Incubating proposal preso was a great success!
        - slides: [https://cloud-native.slack.com/files/U069X07MHEW/F08BMHSD8PQ/k8gb-io-incubation-preso.pdf](https://cloud-native.slack.com/files/U069X07MHEW/F08BMHSD8PQ/k8gb-io-incubation-preso.pdf) 
        - recording: [https://www.youtube.com/watch?v=neWnJad-IxI](https://www.youtube.com/watch?v=neWnJad-IxI) 
    - KubeCon London (attending: bradley, maybe yury)
        - project pavilion
        - lightning talk
        - rejekts talks still open(?) - notifications 10.2 - schedule 17.2 

- **Issues**
    - [https://github.com/k8gb-io/k8gb/issues/1772](https://github.com/k8gb-io/k8gb/issues/1772) Create test setup for upstream DNS providers - WIP - will be Crossplane-based
    - [https://github.com/k8gb-io/k8gb/issues/1778](https://github.com/k8gb-io/k8gb/issues/1778) Revamp website - WIP - Gerardo works on PoC
    - [https://github.com/k8gb-io/k8gb/issues/1800#issuecomment-2588261322](https://github.com/k8gb-io/k8gb/issues/1800#issuecomment-2588261322) partially solved, unclear what to do with CNAME part
    - [https://github.com/k8gb-io/k8gb/issues/1035#issuecomment-2582936204](https://github.com/k8gb-io/k8gb/issues/1035#issuecomment-2582936204) HA questions




-  **PRs** 
    - [https://github.com/k8gb-io/k8gb/pull/1820](https://github.com/k8gb-io/k8gb/pull/1820) Security Self-Assesment Non-goald for review
    - [https://github.com/k8gb-io/k8gb/pull/1762](https://github.com/k8gb-io/k8gb/pull/1762) Manual regression tests ongoing but the pipeline is red (Yury)
    - [https://github.com/k8gb-io/k8gb/pull/1774](https://github.com/k8gb-io/k8gb/pull/1774) [kuritka@gmail.com](mailto:kuritka@gmail.com) is working on refactoring
 

- **Discussion**
    - 

- **Follow-up**
    - [https://github.com/k8gb-io/k8gb/issues/1797](https://github.com/k8gb-io/k8gb/issues/1797) FOSSA fails again, we need to sort out the access to see the actual error [Bradley Andersen](mailto:bradley.d.andersen@gmail.com) look in to it
        - in progress: 
            - team setup issues - new invites sent ← once this is worked out, we more or less self-manage
            - team members vs org relationships issues - (ex: Yury moving account to k8gb)
            - associating team with correct repo (absa/k8gb vs k8gb/k8gb)
        - Docs PoC**&#32;**[https://github.com/k8gb-io/k8gb/issues/1778](https://github.com/k8gb-io/k8gb/issues/1778) 
        - [https://github.com/k8gb-io/k8gb/pull/1762#issuecomment-2606075921](https://github.com/k8gb-io/k8gb/pull/1762#issuecomment-2606075921) switch to upstream external-dns
            - Azure regression test is green
            - Route53 test is green
        - [https://github.com/k8gb-io/k8gb/pull/1773#issuecomment-2606083136](https://github.com/k8gb-io/k8gb/pull/1773#issuecomment-2606083136) Azure test setup and alternative implementation with Crossplane
            - I used Crossplane for route53 regression test too. Tests need polishing and documentation before I send them for review(Yury)




- **Community Update**
    - Making solid progress toward incubating 🎉
        - [https://github.com/cncf/toc/issues/1472#issuecomment-2605486438](https://github.com/cncf/toc/issues/1472#issuecomment-2605486438) 




- **Action Items**
    - Release after regression testing of [https://github.com/k8gb-io/k8gb/pull/1762](https://github.com/k8gb-io/k8gb/pull/1762) 

</details><details><summary><strong>Jan 22, 2025 #61</strong></summary>

## Jan 22, 2025 #61

Recording: [https://zoom.us/rec/share/JMmJgi9kqo2GNkJFxeDaqsBdlV2sFLzq7y5-hlCa2aJo30tZ5JdjCdrW-8Gzph5z.bHvyJx4k5jkdnvuW](https://zoom.us/rec/share/JMmJgi9kqo2GNkJFxeDaqsBdlV2sFLzq7y5-hlCa2aJo30tZ5JdjCdrW-8Gzph5z.bHvyJx4k5jkdnvuW) 

Attendees:

- Bradley Andersen (@elohmrow)
- Michal
- Dinar
- Andre
- Yury

Agenda:

- **News**
    - KubeCon London (attending: bradley, maybe yury)
        - project pavilion
        - lightning talk
        - rejekts talks still open
    - [https://www.youtube.com/live/tKUNI6E1\_7c](https://www.youtube.com/live/tKUNI6E1_7c) ChatLoopBackOff - Episode 42 (K8gb)

- **Issues**
    - [https://github.com/k8gb-io/k8gb/issues/1797](https://github.com/k8gb-io/k8gb/issues/1797) FOSSA fails again, we need to sort out the access to see the actual error [Bradley Andersen](mailto:bradley.d.andersen@gmail.com) look in to it
    - Docs PoC**&#32;**[https://github.com/k8gb-io/k8gb/issues/1778](https://github.com/k8gb-io/k8gb/issues/1778) 

-  **PRs** 
    - [https://github.com/k8gb-io/k8gb/pull/1762#issuecomment-2606075921](https://github.com/k8gb-io/k8gb/pull/1762#issuecomment-2606075921) switch to upstream external-dns
        - Azure regression test is green
 

- **Discussion**
    - [https://github.com/k8gb-io/k8gb/pull/1773#issuecomment-2606083136](https://github.com/k8gb-io/k8gb/pull/1773#issuecomment-2606083136) Azure test setup and alternative implementation with Crossplane

- **Follow-up**
    - 
- **Community Update**
    - Vanity Metrics Highlights
        - GitHub: 922 stars - history: [https://emanuelef.github.io/daily-stars-explorer/#/k8gb-io/k8gb](https://emanuelef.github.io/daily-stars-explorer/#/k8gb-io/k8gb) 
        - Slack: 106 members
        - ADOPTERS: 5
- **Action Items**
    - **CNCF presentation (tomorrow Jan 23rd)**
    - **Release after regression testing of&#32;[https://github.com/k8gb-io/k8gb/pull/1762](https://github.com/k8gb-io/k8gb/pull/1762)&#160;**

</details><details><summary><strong>Jan 8, 2025 #60</strong></summary>

## Jan 8, 2025 #60

Recording: [https://zoom.us/rec/share/KdAOg09uyVV8yAhJ24x2Plo3KCZUwrYPhpbBz7DcHlrFmPUu7C1W3DHkO-xgqogZ.Kw49ZCQt8ZjGyfiP](https://zoom.us/rec/share/KdAOg09uyVV8yAhJ24x2Plo3KCZUwrYPhpbBz7DcHlrFmPUu7C1W3DHkO-xgqogZ.Kw49ZCQt8ZjGyfiP) 

Attendees:

- Yury (@ytsarev)
- Andre (@abaguas)
- Nuno Guedes (@infbase)

Agenda:

- **News**
    - KubeCon London: lightning talk is accepted! Bradley will take the stage!
    - contribfest is unfortunately rejected
    - We should learn more next week Jan 13 on the rest of the applications to the main track. For colocated events cfp notifications are scheduled for  21 January
    - Medium blog post on k8gb experience at KubeCon Salt Lake City [https://medium.com/@kubernetesglobalbalancer/k8gb-rocks-at-kubecon-na-2024-4acb01721560](https://medium.com/@kubernetesglobalbalancer/k8gb-rocks-at-kubecon-na-2024-4acb01721560) 
    - K8gb was mentioned in another KubeCon SLC presentation [https://www.youtube.com/watch?v=cpkKinqdwqA](https://www.youtube.com/watch?v=cpkKinqdwqA) 
    - 

- [https://www.youtube.com/live/tKUNI6E1\_7c](https://www.youtube.com/live/tKUNI6E1_7c) ChatLoopBackOff - Episode 42 (K8gb) scheduled for Jan 16

- **Issues**
    - [https://github.com/k8gb-io/k8gb/issues/1797](https://github.com/k8gb-io/k8gb/issues/1797) Fossa pipe is failing
    - [https://github.com/k8gb-io/k8gb/issues/1765](https://github.com/k8gb-io/k8gb/issues/1765) Updating incorrect Ingress annotation value does not update Gslb
    - [https://github.com/k8gb-io/k8gb/issues/1778](https://github.com/k8gb-io/k8gb/issues/1778) revamp website

-  **PRs** 
    - \*\* K8GB behind reverse proxy [https://github.com/k8gb-io/k8gb/pull/1710](https://github.com/k8gb-io/k8gb/pull/1710)
        - [https://github.com/k8gb-io/k8gb/pull/1710#pullrequestreview-2490242371](https://github.com/k8gb-io/k8gb/pull/1710#pullrequestreview-2490242371) final review bits
    - Test setup for Azure integration [https://github.com/k8gb-io/k8gb/pull/1773](https://github.com/k8gb-io/k8gb/pull/1773) (AWS and GCP also prepared)
        - [Yury Tsarev](mailto:xnullz@gmail.com) will test it e2e this week
    - [https://github.com/k8gb-io/k8gb/pull/1796](https://github.com/k8gb-io/k8gb/pull/1796) badges on frontpage fix
    - [https://github.com/k8gb-io/k8gb/pull/1762](https://github.com/k8gb-io/k8gb/pull/1762) should we revive this one after the external-dns upstream accepted and [released](https://github.com/kubernetes-sigs/external-dns/releases/tag/v0.15.1) changes by Andre? Use [https://github.com/k8gb-io/k8gb/pull/1773](https://github.com/k8gb-io/k8gb/pull/1773) for testing
    - [https://github.com/k8gb-io/k8gb/pull/1774](https://github.com/k8gb-io/k8gb/pull/1774) from [donovan.muller@gmail.com](mailto:donovan.muller@gmail.com)
 

- **Discussion**
    - Should we plan for the next release?
        - Switch to upstream external-dns + reverse proxy support merge, then we can release

- **Follow-up**
    - [Bradley Andersen](mailto:bradley.d.andersen@gmail.com) make sure [https://github.com/cncf/foundation/blob/main/project-maintainers.csv](https://github.com/cncf/foundation/blob/main/project-maintainers.csv) is up2date 
- **Community Update**
    - <high prio> 23 Jan preso to Network TAG as part of Incubation application 🚀
        - [Yury Tsarev](mailto:xnullz@gmail.com) prepare the slides 
    - Vanity Metrics Highlights
        - GitHub: 911 stars
        - Slack: 104 members
        - LinkedIn: 23 followers
        - Twitter / X: 8 followers
        - ADOPTERS: 5
- **Action Items**
    - **Release**
    - **CNCF presentation**
    - **Docs PoC**

</details><details><summary><strong>Dec 25, 2024</strong></summary>

## Dec 25, 2024 skipping due to Holidays

</details><details><summary><strong>Dec 11, 2024 #59</strong></summary>

## Dec 11, 2024 #59

Recording: [https://zoom.us/rec/share/KdAOg09uyVV8yAhJ24x2Plo3KCZUwrYPhpbBz7DcHlrFmPUu7C1W3DHkO-xgqogZ.Kw49ZCQt8ZjGyfiP](https://zoom.us/rec/share/KdAOg09uyVV8yAhJ24x2Plo3KCZUwrYPhpbBz7DcHlrFmPUu7C1W3DHkO-xgqogZ.Kw49ZCQt8ZjGyfiP) 

Attendees:

- Bradley Andersen (@elohmrow) 
- Yury (@ytsarev)
- Andre (@abaguas)

Agenda:

- **News**
    - No meeting next time? (25th December) … next one 8. Jan.

- **Issues**
    - Potential Infoblox bug: [https://cloud-native.slack.com/archives/C021P656HGB/p1732226260669799](https://cloud-native.slack.com/archives/C021P656HGB/p1732226260669799)
        - [https://www.k8gb.io/docs/deploy\_infoblox.html](https://www.k8gb.io/docs/deploy_infoblox.html) 
    - ✅We need to refresh the coredns plugin: [https://github.com/k8gb-io/coredns-crd-plugin](https://github.com/k8gb-io/coredns-crd-plugin)
        - To fix vulnerabilities: [https://artifacthub.io/packages/helm/k8gb/k8gb](https://artifacthub.io/packages/helm/k8gb/k8gb)

-  **PRs** 
    - \*\* Add CNAME support to dig [https://github.com/k8gb-io/k8gb/pull/1783/files](https://github.com/k8gb-io/k8gb/pull/1783/files)
    - \*\* Use upstream coreDNS chart [https://github.com/k8gb-io/k8gb/pull/1776](https://github.com/k8gb-io/k8gb/pull/1776)
    - ✅ Small fix to reduce health checks [https://github.com/k8gb-io/k8gb/pull/1777](https://github.com/k8gb-io/k8gb/pull/1777)
    - \*\* K8GB behind reverse proxy [https://github.com/k8gb-io/k8gb/pull/1710](https://github.com/k8gb-io/k8gb/pull/1710)
    - Test setup for Azure integration [https://github.com/k8gb-io/k8gb/pull/1773](https://github.com/k8gb-io/k8gb/pull/1773) (AWS and GCP also prepared)
    - Chainsaw as new e2e testing framework: [https://github.com/k8gb-io/k8gb/pull/1758](https://github.com/k8gb-io/k8gb/pull/1758)
    - Add flag to support ClusterIP exposed CoreDNS: [https://github.com/k8gb-io/k8gb/pull/1788](https://github.com/k8gb-io/k8gb/pull/1788) 

- **Discussion**
    - 

- **Follow-up**
    - Prepare contribfest and booth for Kubecon Europe
    - Apply for KubeCon co-located events?
    - ? [https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/features-add-ons/maintainer-summit/](https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/features-add-ons/maintainer-summit/) 
    - FOSDEM?
        - Bradley cannot attend - conflict
- **Community Update**
    - <high prio> 9 Jan preso to Network TAG as part of Incubation application 🚀 
    - <low prio> change bluesky handle to @k8gb.io
        - need help from someone with access to k8gb dns (to add a TXT record) or the web server (to create a text file)
    - <info> k8gb Community Manager talk submitted to KubeCon EU 🤞
    - <info> k8gb Project lightning talk submitted to KubeCon EU
    - <info> k8gb Contribfest talk submitted to KubeCon EU
    - <info> k8gb Kiosk submitted to KubeCon EU
    - Vanity Metrics Highlights
        - GitHub: 899 stars
        - Slack: 100 members
        - LinkedIn: 23 followers
        - Twitter / X: 6 followers
        - ADOPTERS: 5
- **Action Items**
    - [Bradley Andersen](mailto:bradley.d.andersen@gmail.com) make sure [https://github.com/cncf/foundation/blob/main/project-maintainers.csv](https://github.com/cncf/foundation/blob/main/project-maintainers.csv) is up2date
    - [Bradley Andersen](mailto:bradley.d.andersen@gmail.com) follow-up, incubating app + watch recorded presos for guidance
    - [Andre Aguas](mailto:andre.b.aguas@gmail.com) - apply for FOSDEM with Gateway API talk
    - [https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/co-located-events/cfp-colocated-events/](https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/co-located-events/cfp-colocated-events/) 
        - Kubernetes on Edge Day?
            - [https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/co-located-events/kubernetes-on-edge-day/](https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/co-located-events/kubernetes-on-edge-day/) 

</details><details><summary><strong>Nov 27, 2024 #58</strong></summary>

## Nov 27, 2024 #58

Recording: [https://zoom.us/rec/share/ztKGoAJPDxxOO0WjT6vnPiwp9VhU1vNIkQGYXz0sdt6yrZJj1MkcuHAHw18H2srY.YMx0sWqnnYa6T9Ic](https://zoom.us/rec/share/ztKGoAJPDxxOO0WjT6vnPiwp9VhU1vNIkQGYXz0sdt6yrZJj1MkcuHAHw18H2srY.YMx0sWqnnYa6T9Ic) 

Attendees:

- Bradley Andersen (@elohmrow)
- Yury Tsarev (@ytsarev)
- Andre Aguas (@abaguas)

Agenda:

- **News**
    - Kubecon NA lightning ([https://www.youtube.com/watch?v=vCzl15AIoU0&t=6s](https://www.youtube.com/watch?v=vCzl15AIoU0&t=6s)) + booth + contribfest; There was a great turnout at the booth 🎉

- **Issues**
    - Potential Infoblox bug: [https://cloud-native.slack.com/archives/C021P656HGB/p1732226260669799](https://cloud-native.slack.com/archives/C021P656HGB/p1732226260669799)
        - [https://www.k8gb.io/docs/deploy\_infoblox.html](https://www.k8gb.io/docs/deploy_infoblox.html) 
    - We need to refresh the coredns plugin: [https://github.com/k8gb-io/coredns-crd-plugin](https://github.com/k8gb-io/coredns-crd-plugin)
        - To fix vulnerabilities: [https://artifacthub.io/packages/helm/k8gb/k8gb](https://artifacthub.io/packages/helm/k8gb/k8gb)

-  **PRs** 
    - \*\* Add CNAME support to dig [https://github.com/k8gb-io/k8gb/pull/1783/files](https://github.com/k8gb-io/k8gb/pull/1783/files)
    - \*\* Use upstream coreDNS chart [https://github.com/k8gb-io/k8gb/pull/1776](https://github.com/k8gb-io/k8gb/pull/1776)
    - Small fix to reduce health checks [https://github.com/k8gb-io/k8gb/pull/1777](https://github.com/k8gb-io/k8gb/pull/1777)
    - \*\* K8GB behind reverse proxy [https://github.com/k8gb-io/k8gb/pull/1710](https://github.com/k8gb-io/k8gb/pull/1710)
    - Test setup for Azure integration [https://github.com/k8gb-io/k8gb/pull/1773](https://github.com/k8gb-io/k8gb/pull/1773) (AWS and GCP also prepared)
    - Chainsaw as new e2e testing framework: [https://github.com/k8gb-io/k8gb/pull/1758](https://github.com/k8gb-io/k8gb/pull/1758)

- **Discussion**
    - 

- **Follow-up**
    - Prepare contribfest, lightning talk and booth for Kubecon Europe
    - Apply for KubeCon co-located events?
    - ? [https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/features-add-ons/maintainer-summit/](https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/features-add-ons/maintainer-summit/) 
    - FOSDEM?
- **Community Update**
    - <high prio> 12/12 preso to Network TAG as part of Incubation application 🚀 
    - <low prio> change bluesky handle to @k8gb.io
        - need help from someone with access to k8gb dns (to add a TXT record) or the web server (to create a text file)
    - <info> k8gb Community Manager talk submitted to KubeCon EU 🤞
- **Action Items**
    - [Bradley Andersen](mailto:bradley.d.andersen@gmail.com) make sure [https://github.com/cncf/foundation/blob/main/project-maintainers.csv](https://github.com/cncf/foundation/blob/main/project-maintainers.csv) is up2date
    - 
    - [Bradley Andersen](mailto:bradley.d.andersen@gmail.com) follow-up, incubating app + watch recorded presos for guidance
    - @ everyone - think about KubeCon EU proj lightning talk … 1m at end is updates - Bradley can deliver it
    - [Andre Aguas](mailto:andre.b.aguas@gmail.com) - apply for FOSDEM with Gateway API talk
    - [https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/co-located-events/cfp-colocated-events/](https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/co-located-events/cfp-colocated-events/) 
        - Kubernetes on Edge Day?
            - [https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/co-located-events/kubernetes-on-edge-day/](https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/co-located-events/kubernetes-on-edge-day/) 

</details><details><summary><strong>Nov 13, 2024</strong></summary>

## Nov 13, 2024 skipping due to KubeCon

</details><details><summary><strong>Oct 30, 2024 #57</strong></summary>

## Oct 30, 2024 #57

Recording: [https://zoom.us/rec/share/H58AgQY5rWICcQLk8eekS3qHyKIhhjMwWlnzniGfA0lgQRP6-8V4VwQhGq\_hdRU2.T7sebQcmYrbpJcwh](https://zoom.us/rec/share/H58AgQY5rWICcQLk8eekS3qHyKIhhjMwWlnzniGfA0lgQRP6-8V4VwQhGq_hdRU2.T7sebQcmYrbpJcwh) 

Attendees:

- Bradley Andersen (@elohmrow)
- Yury Tsarev (@ytsarev)
- Andre Aguas (@abaguas)
- Michal Kuřítka

Agenda:

- **News**
    - k8gb Incubation Application submitted: [https://github.com/cncf/toc/issues/1472](https://github.com/cncf/toc/issues/1472) 🎉
    - New adopter PagBank: [https://github.com/k8gb-io/k8gb/pull/1755](https://github.com/k8gb-io/k8gb/pull/1755)

- **Issues**
    - Strategy not updated if initially set with a wrong value: [https://github.com/k8gb-io/k8gb/issues/1765](https://github.com/k8gb-io/k8gb/issues/1765)
        - might be useful for experimenting: [https://playcel.undistro.io/](https://playcel.undistro.io/) 
    - Health checks of ingress controller: [https://github.com/k8gb-io/k8gb/issues/1754](https://github.com/k8gb-io/k8gb/issues/1754)

-  **PRs** 
    - anything [from last time](https://docs.google.com/document/d/1YdpEVhtyCKvwFtXR7cn1Kn2Xc_tdNskoFnhHmFPbtA4/edit?tab=t.0#bookmark=id.7bli0m25ht36)?
        - [https://github.com/k8gb-io/k8gb/pull/1743](https://github.com/k8gb-io/k8gb/pull/1743) → Chainsaw experiments
    - Switch external DNS back to upstream: [https://github.com/k8gb-io/k8gb/pull/1762](https://github.com/k8gb-io/k8gb/pull/1762)
    - Chainsaw PoC: [https://github.com/k8gb-io/k8gb/pull/1758](https://github.com/k8gb-io/k8gb/pull/1758)
        - External DNS flapping for test parallelism: [https://github.com/k8gb-io/k8gb/pull/1767](https://github.com/k8gb-io/k8gb/pull/1767)

- **Discussion**
    - The API group of the CRD is k8gb.absa.oss, does it meet the incubation application requirements? [Bradley Andersen](mailto:bradley.d.andersen@gmail.com) check if it \_must\_ change (vendor neutrality)

- **Follow-up**
    - Content for Medium posts
- **Community Update**
    - Vanity Metrics Highlights (+change 1 month)
        - GitHub: +19 stars
        - Slack: +3 members
        - LinkedIn: +11 followers
        - Twitter / X: +28 followers
        - ADOPTERS: +2
- **Action Items**
    - [Bradley Andersen](mailto:bradley.d.andersen@gmail.com) social posts around Kubecon
    - Yury / Andre prep for Kubecon
    - Andre investigates validation of annotation values

</details><details><summary><strong>Oct 16, 2024 #56</strong></summary>

## Oct 16, 2024 #56

Recording: [https://zoom.us/rec/share/IM8GFXJroDwQysy1ghIu7yNUxtP3rL3qtVxhY8249cArmfBmmMk7v5iH9FSmAQY9.5fH\_eXwTS2-3ME-p](https://zoom.us/rec/share/IM8GFXJroDwQysy1ghIu7yNUxtP3rL3qtVxhY8249cArmfBmmMk7v5iH9FSmAQY9.5fH_eXwTS2-3ME-p) 

Attendees:

- Andre Aguas (@abaguas)
- Nuno Guedes (@infbase)

Agenda:

- **News**
    - Open Systems is now an official adopter: [https://github.com/k8gb-io/k8gb/pull/1753](https://github.com/k8gb-io/k8gb/pull/1753)

- **Issues**
    - [https://github.com/k8gb-io/k8gb/issues/1741](https://github.com/k8gb-io/k8gb/issues/1741) CoreDNS AWS NLB health check not getting healthy

-  **PRs** 
    - [https://github.com/k8gb-io/k8gb/pull/1743](https://github.com/k8gb-io/k8gb/pull/1743) use upstream coredns chart instead of fork
        - All features seem to be there and resolution is working in local cluster but integration test setup is not yet working
    - [https://github.com/k8gb-io/k8gb/pull/1710](https://github.com/k8gb-io/k8gb/pull/1710) Support k8gb behind a reverse proxy
        - PR ready to be merged after a round of feedback
    - [https://github.com/k8gb-io/k8gb/issues/1662](https://github.com/k8gb-io/k8gb/issues/1662) Encubation proposal close to completion. Kudos Bradley

- **Discussion**
    - 
- **Follow-up**
    - Medium posts
    - Intro video
- **Community Update**
    - 868 github stars**!&#160;**
- **Action Items**
    - 

</details><details><summary><strong>Oct 2, 2024 #55</strong></summary>

## Oct 2, 2024 #55

Recording: [https://zoom.us/rec/share/uzOCWcKBrVeutDfGJD44Yy8dFXBSyiOwAzPf6-mdy3hhDr07RhX3gJC\_evdvwUNv.e7jyztGpJW6eolBR](https://zoom.us/rec/share/uzOCWcKBrVeutDfGJD44Yy8dFXBSyiOwAzPf6-mdy3hhDr07RhX3gJC_evdvwUNv.e7jyztGpJW6eolBR) 

Attendees:

- Yury Tsarev (@ytsarev)
- Nuno Guedes (@infbase)
- Andre Aguas (@abaguas)

Agenda:

- **News**
    - Nuno and Diego presentation in [KCD Porto](https://community.cncf.io/events/details/cncf-kcd-porto-presents-kcd-porto-2024/) delivered 🎉
    - Jiri and Yury [OSS Summit Vienna](https://osseu2024.sched.com/event/1ej7D?iframe=no) delivered 🎉
        - Recording is available: [https://www.youtube.com/watch?v=5eLX4kMgo8Q](https://www.youtube.com/watch?v=5eLX4kMgo8Q) 
    - For the renovate process, we are pinning releases of github actions from now on, example [https://github.com/k8gb-io/k8gb/commit/eb1f4dd8009514c7d5dd4189cf73e03eab561af8](https://github.com/k8gb-io/k8gb/commit/eb1f4dd8009514c7d5dd4189cf73e03eab561af8)

- **Issues**
    - [https://github.com/k8gb-io/k8gb/issues/1741](https://github.com/k8gb-io/k8gb/issues/1741) CoreDNS AWS NLB health check not getting healthy
    - [Support reading external metrics to assess Service readiness · Issue #1745 · k8gb-io/k8gb (github.com)](https://github.com/k8gb-io/k8gb/issues/1745)

-  **PRs** 
    - [https://github.com/k8gb-io/k8gb/pull/1743](https://github.com/k8gb-io/k8gb/pull/1743) use upstream coredns chart instead of fork
    - [https://github.com/k8gb-io/k8gb/pull/1710](https://github.com/k8gb-io/k8gb/pull/1710) Support k8gb behind a reverse proxy

- **Discussion**
    - Direct support for K8s services of type LoadBalancer (using GSLBReferenceResolver type?)
        - [Support Service of type LoadBalancer to enable global load balancing on L4 · Issue #147 · k8gb-io/k8gb (github.com)](https://github.com/k8gb-io/k8gb/issues/147)
        - [K8GB for none K8S Workload externally · Issue #1140 · k8gb-io/k8gb (github.com)](https://github.com/k8gb-io/k8gb/issues/1140)
        - [K8GB for the service of Type ExternalName/Loadbalancer · Issue #1212 · k8gb-io/k8gb (github.com)](https://github.com/k8gb-io/k8gb/issues/1212)
        - [add support for svc backend by tanujd11 · Pull Request #1363 · k8gb-io/k8gb (github.com)](https://github.com/k8gb-io/k8gb/pull/1363)
    - Revised 1.0 release plans
        - ​​[https://github.com/k8gb-io/k8gb/issues/52#issuecomment-2387967660](https://github.com/k8gb-io/k8gb/issues/52#issuecomment-2387967660) 
            - Should we include Service type LoadBalancer in 1.0?
            - Should we include[&#32;Gateway API support&#32;](https://github.com/k8gb-io/k8gb/issues/954)in 1.0?
- **Follow-up**
    - 
- **Community Update**
    - 864 github stars**!**
- **Action Items**
    - 

</details><details><summary><strong>Sep 18, 2024 #54</strong></summary>

## Sep 18, 2024 #54

Recording: [https://zoom.us/rec/share/Q0A-na6wmDh\_yPUfdQNu7CwsAPWSX9u-x4RLBcHPlidDNR30qRfQ2U87Rd46x2lM.wRsJ8OPTvOqnO7xi](https://zoom.us/rec/share/Q0A-na6wmDh_yPUfdQNu7CwsAPWSX9u-x4RLBcHPlidDNR30qRfQ2U87Rd46x2lM.wRsJ8OPTvOqnO7xi)

Attendees:

- Bradley Andersen (@elohmrow)
- Andre Aguas (@abaguas)

Agenda:

- **News**

- Yury @ [OSS Summit Vienna](https://osseu2024.sched.com/event/1ej7D?iframe=no) today
- Next week: K8gb session at [CNCF KCD Porto 2024](https://community.cncf.io/events/details/cncf-kcd-porto-presents-kcd-porto-2024/)
- [0.14.0 released](https://github.com/k8gb-io/k8gb/releases/tag/v0.14.0)
    - [https://github.com/k8gb-io/k8gb/issues/552](https://github.com/k8gb-io/k8gb/issues/552) 🎉 - Gateway API next - edging toward feature-completeness
- anything else?

- **Issues**

- 

- **Discussion**

- [1.0 from last time](https://docs.google.com/document/d/1YdpEVhtyCKvwFtXR7cn1Kn2Xc_tdNskoFnhHmFPbtA4/edit#bookmark=id.nzbn36qmqlsa)
    - some blockers, no dates yet 👍

- **Follow-up**

- Intro video
- [Ingress Controller health status issue](https://cloud-native.slack.com/archives/C021P656HGB/p1722973213080159) - testing in progress
- Blog
    - perhaps medium?

- **Community Update**

- Started making social posts social posts - plan to do more, more regular
- Vanity Metrics (+change 1 month)
    - GitHub Repo:  [https://github.com/k8gb-io/k8gb](https://github.com/k8gb-io/k8gb) - 17 watchers (+1), 92 forks (+0), 856 stars (+8), +25 commits, +24 PRs, 32 releases (+1)
    - Slack: [#k8gb](https://cloud-native.slack.com/archives/C021P656HGB) - 94 members (+3), +26 messages
    - LinkedIn: [https://www.linkedin.com/company/k8gb/](https://www.linkedin.com/company/k8gb/) - 11 followers (+1), 5 posts (+3), 281 post views (+108), 16 clicks (+10), 11 (+6) reactions
    - Twitter / X: [https://x.com/k8gb\_io](https://x.com/k8gb_io) - 3 posts (+1), 7 followers (+4), 92 views, 7 reactions, 3 re-posts
- WIP CNCF Incubating application: [https://github.com/k8gb-io/k8gb/issues/1662](https://github.com/k8gb-io/k8gb/issues/1662) - will work on finishing it up by first week of October.

- **Action Items**
**&#160;**
- [Bradley Andersen](mailto:bradley.d.andersen@gmail.com) to post to socials [0.14.0 release](https://github.com/k8gb-io/k8gb/releases/tag/v0.14.0)
- [Bradley Andersen](mailto:bradley.d.andersen@gmail.com) to look into blog stack
- [Bradley Andersen](mailto:bradley.d.andersen@gmail.com) to work on [Incubating app](https://github.com/k8gb-io/k8gb/issues/1662)

</details><details><summary><strong>Sep 4, 2024 #53</strong></summary>

## Sep 4, 2024 #53

Recording: [https://zoom.us/rec/share/jR82XtAiMvqawCnCbsIz-zlUiGNZVMCs4V43W\_o10JMnylBu5fJy82JStkQ0Bw6z.AudTu70f0zdrEBf3](https://zoom.us/rec/share/jR82XtAiMvqawCnCbsIz-zlUiGNZVMCs4V43W_o10JMnylBu5fJy82JStkQ0Bw6z.AudTu70f0zdrEBf3)

Attendees:

- Andre Aguas (@abaguas)
- Yury Tsarev (@ytsarev)

Agenda:

- **News**

- 

- **Issues**

- [https://github.com/k8gb-io/k8gb/issues/1275](https://github.com/k8gb-io/k8gb/issues/1275)  - reverse proxy support / annotation-based IP list control
    - simple implementation wise

- **Discussion**

- 1.0 prereqs
    - get rid of external-dns fork
        - It looks like NS record support is there in upstream
            - Will need retesting of AWS and Azure
        - Add GCP support on top of latest external-dns
    - Update [https://github.com/k8gb-io/coredns-helm](https://github.com/k8gb-io/coredns-helm) to latest upstream if possible
    - [https://github.com/k8gb-io/coredns-crd-plugin](https://github.com/k8gb-io/coredns-crd-plugin) check the state and deps
- [https://github.com/k8gb-io/k8gb/pull/1672#discussion\_r1737285468](https://github.com/k8gb-io/k8gb/pull/1672#discussion_r1737285468) shape of API

- **Follow-up**

- Action items from last time
    - Nuno will start speaking with folks for testimonials related to the intro video
        - 
    - Nuno will create an issue around this (Ingress Controller health status) discussion: [https://cloud-native.slack.com/archives/C021P656HGB/p1722973213080159](https://cloud-native.slack.com/archives/C021P656HGB/p1722973213080159) 
        - testing in progress
    - Yury will speak with Donovan about moving the repo from this (Helm chart) discussion: [https://cloud-native.slack.com/archives/C021P656HGB/p1722186772957459](https://cloud-native.slack.com/archives/C021P656HGB/p1722186772957459) 
        - Done
    - keep thinking about blog 
        - Still in progress

- **Community Update**

- 

- **Action Items**
**&#160;**
- Vienna talk prep
- 0.14.0 release before the conference 

</details><details><summary><strong>Aug 21, 2024 #52</strong></summary>

## Aug 21, 2024 #52

Recording: [https://zoom.us/rec/share/5RtNmAfG2f7QyAduMWufPE8ktUcw\_7jWflEiUOE4jiPw1Ru4Wf3zVrufQ9MrE3H\_.ZM7dmDGwZR4ea9z4](https://zoom.us/rec/share/5RtNmAfG2f7QyAduMWufPE8ktUcw_7jWflEiUOE4jiPw1Ru4Wf3zVrufQ9MrE3H_.ZM7dmDGwZR4ea9z4) 

Attendees:

- Bradley Andersen (@elohmrow)
- Nuno Guedes (@infbase) 
- Andre Aguas (@abaguas)

Agenda:

- **News**

- 

- **Issues**

- [https://github.com/k8gb-io/k8gb/issues/1275](https://github.com/k8gb-io/k8gb/issues/1275)  - reverse proxy support / annotation-based IP list control
    - simple implementation wise

- **Discussion**

- 

- **Follow-up**

- Action items from last time
    - Nuno will start speaking with folks for testimonials related to the intro video
        - 
    - Nuno will create an issue around this (Ingress Controller health status) discussion: [https://cloud-native.slack.com/archives/C021P656HGB/p1722973213080159](https://cloud-native.slack.com/archives/C021P656HGB/p1722973213080159) 
        - testing in progress
    - Yury will speak with Donovan about moving the repo from this (Helm chart) discussion: [https://cloud-native.slack.com/archives/C021P656HGB/p1722186772957459](https://cloud-native.slack.com/archives/C021P656HGB/p1722186772957459) 

- **Community Update**

- We now have @k8gb.io email domain
- We now have a Twitter / X account 
- We now use the CNCF Calendar and Zoom - new link is in the Schedule section, above 👆
    - this required setting up a CNCF Service Desk account
- Social posts to be made:
    - K8gb session at [CNCF KCD Porto 2024](https://community.cncf.io/events/details/cncf-kcd-porto-presents-kcd-porto-2024/) - in one month
    - K8gb session at [OSS Summit Vienna](https://osseu2024.sched.com/event/1ej7D?iframe=no) - in three  weeks
- Vanity Metrics
    - GitHub Repo:  [https://github.com/k8gb-io/k8gb](https://github.com/k8gb-io/k8gb) - 16 watchers, 92 forks, 847 stars, [4 commits, 12 PRs since last meeting]
    - Slack: [#k8gb](https://cloud-native.slack.com/archives/C021P656HGB) - 91 members - [9 messages since last meeting]
    - LinkedIn: [https://www.linkedin.com/company/k8gb/](https://www.linkedin.com/company/k8gb/) - 10 followers, [2 posts, 73 post views, 6 clicks, 5 reactions since the last meeting]
    - Twitter / X: [https://x.com/k8gb\_io](https://x.com/k8gb_io) - [2 posts, 3 followers since last meeting]
- WIP CNCF Incubating application: [https://github.com/k8gb-io/k8gb/issues/1662](https://github.com/k8gb-io/k8gb/issues/1662) - will work on finishing it up this week.

- **Action Items**
**&#160;**
- keep thinking about blog 

</details><details><summary><strong>Aug 7, 2024 #51</strong></summary>

## Aug 7, 2024 #51

Recording: [https://youtu.be/yy53PgAlx7o](https://youtu.be/yy53PgAlx7o) 

Attendees:

- Bradley Andersen (@elohmrow)
- Nuno Guedes (@infbase)
- Yury Tsarev (@ytsarev)
- Michal Kuritka (@kuritka)

Agenda:

- News
    - K8gb session at [CNCF KCD Porto 2024](https://community.cncf.io/events/details/cncf-kcd-porto-presents-kcd-porto-2024/) 🎉
    - K8gb session at [OSS Summit Vienna](https://osseu2024.sched.com/event/1ej7D?iframe=no) - Jiri agreed to co-present, we will do it together! 🎉
    - [https://github.com/k8gb-io/k8gb/pull/1682](https://github.com/k8gb-io/k8gb/pull/1682) and [https://github.com/k8gb-io/k8gb/pull/1684](https://github.com/k8gb-io/k8gb/pull/1684) - glorious terratest fixes by Andre
- Issues:
    - [https://github.com/k8gb-io/k8gb/issues/1275](https://github.com/k8gb-io/k8gb/issues/1275)  - reverse proxy support / annotation-based IP list control
- Discussion
    - Slack-based discussion [https://cloud-native.slack.com/archives/C021P656HGB/p1722973213080159](https://cloud-native.slack.com/archives/C021P656HGB/p1722973213080159) - should we track ingress controller health status?
    - Andre: coreDNS helm chart issue: [https://cloud-native.slack.com/archives/C021P656HGB/p1722186772957459](https://cloud-native.slack.com/archives/C021P656HGB/p1722186772957459). It is quite low effort and it would be important to keep all our repos together and up-to-date.
- Follow-up
    - Michal: k8gb.io domain ownership, MX records needed
    - Nuno: 2-minute intro video
    - Yury: issue clean up [https://github.com/orgs/k8gb-io/projects/2](https://github.com/orgs/k8gb-io/projects/2)
        - [https://github.com/orgs/k8gb-io/projects/2/views/2?sortedBy%5Bdirection%5D=asc&sortedBy%5BcolumnId%5D=Status](https://github.com/orgs/k8gb-io/projects/2/views/2?sortedBy%5Bdirection%5D=asc&sortedBy%5BcolumnId%5D=Status)
    - Bradley: WIP CNCF Incubating application: [https://github.com/k8gb-io/k8gb/issues/1641](https://github.com/k8gb-io/k8gb/issues/1641)
        - Zoom
- Action Items
    - Nuno will start speaking with folks for testimonials related to the intro video
    - Nuno will create an issue around this (Ingress Controller health status) discussion: [https://cloud-native.slack.com/archives/C021P656HGB/p1722973213080159](https://cloud-native.slack.com/archives/C021P656HGB/p1722973213080159) 
    - Yury will speak with Donovan about moving the repo from this (Helm chart) discussion: [https://cloud-native.slack.com/archives/C021P656HGB/p1722186772957459](https://cloud-native.slack.com/archives/C021P656HGB/p1722186772957459) 
    - ~~Michal will further investigate the transfer of the k8gb.io domain~~
    - ~~Bradley / Yury to send out new Zoom link to help meet requirements under&#32;[https://github.com/k8gb-io/k8gb/issues/1661](https://github.com/k8gb-io/k8gb/issues/1661)&#160;~~
    - ~~Bradley to make social posts  around K8gb session at&#32;[CNCF KCD Porto 2024](https://community.cncf.io/events/details/cncf-kcd-porto-presents-kcd-porto-2024/)&#32;and K8gb session at&#32;[OSS Summit Vienna](https://osseu2024.sched.com/event/1ej7D?iframe=no)

</details><details><summary><strong>Jul 24, 2024 #50</strong></summary>

## Jul 24, 2024 #50

Recording: [https://youtu.be/walfel6rijE](https://youtu.be/walfel6rijE) 

Attendees:

- Nuno Guedes (@infbase)
- Andre Aguas (@abaguas)

Agenda:

- News
    - Community Management update
        - WIP CNCF Incubating application: [https://github.com/k8gb-io/k8gb/issues/1641](https://github.com/k8gb-io/k8gb/issues/1641)
        - social
            - Nuno: [2-minute intro video](https://github.com/k8gb-io/k8gb/issues/1640)
                - Nuno will do the editing, we need to record testimonials and the narrator track
            - Michal: k8gb.io domain ownership, MX records needed

- Issues:
    - 

- PRs:
    - 

- Discussion
    - k8gb.io domain ownership, MX records needed
    - 2-minute intro video
    - issue clean up [https://github.com/orgs/k8gb-io/projects/2](https://github.com/orgs/k8gb-io/projects/2)
        - [https://github.com/orgs/k8gb-io/projects/2/views/2?sortedBy%5Bdirection%5D=asc&sortedBy%5BcolumnId%5D=Status](https://github.com/orgs/k8gb-io/projects/2/views/2?sortedBy%5Bdirection%5D=asc&sortedBy%5BcolumnId%5D=Status) Yury performed clean up, lot of issues are closed, some of the realigned. 1.0 looks tangible

- Action Items
    - 


</details><details><summary><strong>Jul 10, 2024 #49</strong></summary>

## Jul 10, 2024 #49

Recording: [https://www.youtube.com/watch?v=h\_fclJwhPmE](https://www.youtube.com/watch?v=h_fclJwhPmE) 

Attendees:

- Bradley Andersen (@elohmrow)
- Andre Aguas (@abaguas)
- Michal Kuritka (@kuritka)
- Yury Tsarev (@ytsarev)

Agenda:

- News
    - Community Management update
        - Bradley not here next time due to work trip 👍
        - BIG early goals:
            - to spread awareness
            - make sure the docs are tight
            - roadmap cleanup [https://github.com/orgs/k8gb-io/projects/2](https://github.com/orgs/k8gb-io/projects/2)
            - start the sandbox --> incubation process
        - quick doc fixes ([1633](https://github.com/k8gb-io/k8gb/issues/1633) [1638](https://github.com/k8gb-io/k8gb/pull/1638)) to be followed by [docu maintenance](https://github.com/k8gb-io/k8gb/issues/1643)
        - social
            - everything set up except [twitter](https://github.com/k8gb-io/k8gb/issues/1642)
            - need a [2-minute intro video](https://github.com/k8gb-io/k8gb/issues/1640)
                - Nuno will do the editing, we need to record testimonials and the narrator track
            - will start announcing more upcoming talks, etc. on socials
        - [issue](https://github.com/k8gb-io/k8gb/issues/1641) created to start thinking about CNCF incubation 
        - latest release [announced on LinkedIn](https://www.linkedin.com/posts/k8gb_release-v0130-k8gb-iok8gb-activity-7214901953359704065-UgMu?utm_source=share&utm_medium=member_desktop) by Yury, and on k8gb LinkedIn page
            - [0.13.0 release on linkedin](https://www.linkedin.com/feed/update/urn:li:activity:7214901953359704065)
            - [today’s community meeting on linkedin](https://www.linkedin.com/feed/update/urn:li:activity:7216382497704570881)
        - Bradley created LinkedIn org [https://www.linkedin.com/company/k8gb/](https://www.linkedin.com/company/k8gb/) 
            - ‘Hire’ yourself there as ‘Core Maintainer’ or ‘Contributor’ for visibility! 

- Issues:
    - [https://github.com/orgs/k8gb-io/projects/2](https://github.com/orgs/k8gb-io/projects/2) needs a lot of ❤️

- PRs:
    - [https://github.com/k8gb-io/k8gb/pull/1639](https://github.com/k8gb-io/k8gb/pull/1639) small post-release fix

- Discussion
    - k8gb.io domain ownership, MX records needed

- Action Items
    - Yury to attempt first mass roadmap cleanup, targeting clean release 1.0 plan
        - Cleanup performed, lot of issues closed and realigned [https://github.com/orgs/k8gb-io/projects/2/views/2?sortedBy%5Bdirection%5D=asc&sortedBy%5BcolumnId%5D=Status](https://github.com/orgs/k8gb-io/projects/2/views/2?sortedBy%5Bdirection%5D=asc&sortedBy%5BcolumnId%5D=Status) 
    - Michal to review [https://github.com/k8gb-io/k8gb/pull/1639](https://github.com/k8gb-io/k8gb/pull/1639) 
    - Michal will check k8gb.io domain status in Absa
    - Yury to create chainsaw PoC Issue, Michal to review contents from terratest perspective
        - Done: ​​[https://github.com/k8gb-io/k8gb/issues/1660](https://github.com/k8gb-io/k8gb/issues/1660) 
    - Yury and Bradley to work on CNCF incubation
        - See WIP application: [https://github.com/k8gb-io/k8gb/issues/1641](https://github.com/k8gb-io/k8gb/issues/1641) 


</details><details><summary><strong>Jun 26, 2024 #48</strong></summary>

## Jun 26, 2024 #48

Recording: [https://www.youtube.com/watch?v=xHXlqAhdjcM](https://www.youtube.com/watch?v=xHXlqAhdjcM) 

Attendees:

- Yury Tsarev (@ytsarev)
- Bradley Andersen(@elohmrow)
- Theo Chatzimichos (@tampakrap)
- Andre Aguas (@abaguas)
- Nuno Guedes (@infbase)
- Dinar Valeev (@k0da)

Agenda:

- News
    - k8gb Community Manager joins the team 🎉
    - **Last call**: K8gb talk accepted to Open Source Summit Vienna.  Does anybody want to be a co-speaker? Ping Jiri

- Issues:
    - [https://github.com/k8gb-io/k8gb/issues/1566](https://github.com/k8gb-io/k8gb/issues/1566) scorecard pipeline is failing in master branch
        - [https://github.com/k8gb-io/k8gb/issues/1566#issuecomment-2155197678](https://github.com/k8gb-io/k8gb/issues/1566#issuecomment-2155197678) some hints found

- PRs:
    - Great contribution from Andre
        - [https://github.com/k8gb-io/k8gb/pull/1557](https://github.com/k8gb-io/k8gb/pull/1557) Decouple gslb from the kubernetes Ingress resource
            - Under review/testing, terratest suite to be extended
                - Terratest done 👍
                - Regression testing done 👍
                - Resource reference namespace isolation concern [https://github.com/k8gb-io/k8gb/pull/1557#pullrequestreview-2136723205](https://github.com/k8gb-io/k8gb/pull/1557#pullrequestreview-2136723205) 🟥 
    - [https://github.com/k8gb-io/k8gb/pull/1612](https://github.com/k8gb-io/k8gb/pull/1612) goreleaser to buildx (merged)
    - [https://github.com/k8gb-io/k8gb/pull/1610](https://github.com/k8gb-io/k8gb/pull/1610) renovate: Do not group upgrades of dependencies on major version 0 (merged)
    - Dependency bump PRs from Michal ( [kuritka@gmail.com](mailto:kuritka@gmail.com) ) - everything is merged
        - [https://github.com/k8gb-io/k8gb/pull/1597](https://github.com/k8gb-io/k8gb/pull/1597) Unit Tests: limit warnings, fix racing
        - [https://github.com/k8gb-io/k8gb/pull/1598](https://github.com/k8gb-io/k8gb/pull/1598) Bump ControllerGen, CRD
        - [https://github.com/k8gb-io/k8gb/pull/1599](https://github.com/k8gb-io/k8gb/pull/1599) Bump golangci to v1.59.1
        - [https://github.com/k8gb-io/k8gb/pull/1600](https://github.com/k8gb-io/k8gb/pull/1600) Bump mocks 
- Discussion
    - When should we make a release?
        - Proposal:
            - Immediately: release will include deps bump + Azure support
            - Next one: after [https://github.com/k8gb-io/k8gb/pull/1557](https://github.com/k8gb-io/k8gb/pull/1557) is merged, make a good announcement about it - Yury will proceed with release asap
            - Next next one: GCP support
                - Nuno can provide a test environment


</details><details><summary><strong>Jun 12, 2024 #47</strong></summary>

## Jun 12, 2024 #47

Recording: [https://youtu.be/QP6q6qFYFoo](https://youtu.be/QP6q6qFYFoo) 

Attendees:

- Yury Tsarev (@ytsarev)
- Michal Kuritka (@kuritka)

Agenda:

- News
    - K8gb talk accepted to Open Source Summit Vienna.  Anybody want to be co-speaker?
    - We didn’t release k8gb for a while. Let’s do it!

- Issues:
    - [https://github.com/k8gb-io/k8gb/issues/1566](https://github.com/k8gb-io/k8gb/issues/1566) scorecard pipeline is failing in master branch
        - [https://github.com/k8gb-io/k8gb/issues/1566#issuecomment-2155197678](https://github.com/k8gb-io/k8gb/issues/1566#issuecomment-2155197678) some hints found

- PRs:
    - [https://github.com/k8gb-io/k8gb/pull/1587](https://github.com/k8gb-io/k8gb/pull/1587) - helm OLM publish fix - needs review
    - [https://github.com/k8gb-io/k8gb/pull/1525](https://github.com/k8gb-io/k8gb/pull/1525) Documented Azure DNS deployment  - merged 
        - [https://github.com/k8gb-io/k8gb/pull/1593](https://github.com/k8gb-io/k8gb/pull/1593) Azure secret reference refactoring and enhancement - merged
    - Great contribution from Andre
        - [https://github.com/k8gb-io/k8gb/pull/1557](https://github.com/k8gb-io/k8gb/pull/1557) Decouple gslb from the kubernetes Ingress resource
            - Under review/testing, terratest suite to be extended
    - New PRs from Michal ( [kuritka@gmail.com](mailto:kuritka@gmail.com) ) - need review
        - [https://github.com/k8gb-io/k8gb/pull/1597](https://github.com/k8gb-io/k8gb/pull/1597) Unit Tests: limit warnings, fix racing
        - [https://github.com/k8gb-io/k8gb/pull/1598](https://github.com/k8gb-io/k8gb/pull/1598) Bump ControllerGen, CRD
        - [https://github.com/k8gb-io/k8gb/pull/1599](https://github.com/k8gb-io/k8gb/pull/1599) Bump golangci to v1.59.1
        - [https://github.com/k8gb-io/k8gb/pull/1600](https://github.com/k8gb-io/k8gb/pull/1600) Bump mocks 


</details><details><summary><strong>May 29, 2024 #46</strong></summary>

## May 29, 2024 #46

Recording: [https://youtu.be/Pvyw2jWA3P4](https://youtu.be/Pvyw2jWA3P4) 

Attendees:

- Yury Tsarev (@ytsarev)
- Nuno Guedes (@infbase) 
- Michal Kuritka (@kuritka)

Agenda:

- News
    - Updates on [https://www.k8gb.io/](https://www.k8gb.io/) 
    - > 800 GitHub stars achieved
    - Cool feedback

    - New infoblox webhook - [https://github.com/AbsaOSS/external-dns-infoblox-webhook](https://github.com/AbsaOSS/external-dns-infoblox-webhook)
        - External dns community is moving to webhooks based approach
        - Good chance to deprecate our own external dns fork in future
        - Easier to maintain, lesser blast radius
        - Possibly use [https://github.com/googleapis/release-please](https://github.com/googleapis/release-please) in future

- Issues:
    - [https://github.com/k8gb-io/k8gb/issues/1566](https://github.com/k8gb-io/k8gb/issues/1566) scorecard pipeline is failing in master branch
    - [https://github.com/k8gb-io/k8gb/issues/1348](https://github.com/k8gb-io/k8gb/issues/1348) Document the case of when a load balancing configuration is not deployed in every cluster
    - googleapis/release-please-action@v4

- PRs:
    - [https://github.com/k8gb-io/k8gb/pull/1525](https://github.com/k8gb-io/k8gb/pull/1525) Documented Azure DNS deployment 
        - Review is done, waiting for [nuno.guedes@millenniumbcp.pt](mailto:nuno.guedes@millenniumbcp.pt) feedback
        - [https://github.com/k8gb-io/k8gb/pull/1525#issuecomment-2137198955](https://github.com/k8gb-io/k8gb/pull/1525#issuecomment-2137198955) agreed with way to proceed
    - Great contribution from Andre
        - [https://github.com/k8gb-io/k8gb/pull/1557](https://github.com/k8gb-io/k8gb/pull/1557) Decouple gslb from the kubernetes Ingress resource
            - Under review/testing, terratest suite to be extended
    - Doc: [https://github.com/k8gb-io/k8gb/pull/1578](https://github.com/k8gb-io/k8gb/pull/1578) RFC2136 to front page


</details><details><summary><strong>May 15, 2024 #45</strong></summary>

## May 15, 2024 #45

Recording: [2024-05-15 k8gb community meeting](https://www.youtube.com/watch?v=U4dcUiH8x9I)

Attendees:

- Yury Tsarev (@ytsarev)
- Andre Aguas (@abaguas)
- Michal Kuritka (@kuritka)

Agenda:

- News
    - Great coverage at [https://oilbeater.com/en/2024/04/18/k8gb-best-cloudnative-gslb/](https://oilbeater.com/en/2024/04/18/k8gb-best-cloudnative-gslb/) and nice associated github start bump [https://star-history.com/#k8gb-io/k8gb&Date](https://star-history.com/#k8gb-io/k8gb&Date) 



- Issues:
    - Question around bind in discussions [https://github.com/k8gb-io/k8gb/discussions/364#discussioncomment-9303278](https://github.com/k8gb-io/k8gb/discussions/364#discussioncomment-9303278)
        - We probably should make it more obvious in documentation
- PRs:
    - [https://github.com/k8gb-io/k8gb/pull/1525](https://github.com/k8gb-io/k8gb/pull/1525) Documented Azure DNS deployment 
        - Review is done, waiting for Nuno feedback
    - [https://github.com/k8gb-io/k8gb/pull/1510](https://github.com/k8gb-io/k8gb/pull/1510) Contribfest hero  - review from Michal and Dinar required
    - Great contributions from Andre
        - [https://github.com/k8gb-io/k8gb/pull/1549](https://github.com/k8gb-io/k8gb/pull/1549) Queue reconciliation of all GSLBs that reference the same endpoint - merged
        - [https://github.com/k8gb-io/k8gb/pull/1548](https://github.com/k8gb-io/k8gb/pull/1548) Helm supports extra env, volumes and volume mounts for externaldns - merged
        - [https://github.com/k8gb-io/k8gb/pull/1557](https://github.com/k8gb-io/k8gb/pull/1557) Decouple gslb from the kubernetes Ingress resource
            - this is huge! review and thorough testing is planned


</details><details><summary><strong>Apr 17, 2024 #44</strong></summary>

## Apr 17, 2024 #44

Recording: [https://www.youtube.com/watch?v=JNJ2k7mcHXc](https://www.youtube.com/watch?v=JNJ2k7mcHXc) 

Attendees:

- Yury Tsarev (@ytsarev)
- Andre Aguas (@abaguas)
- Theo Chatzimichos (@tampakrap)

Agenda:

- News
    - K8gb is covered in French from last 2023 KubeCon [https://blog.ippon.fr/2023/05/02/retour-sur-les-conferences-de-la-kubecon-2023-partie-2-3/](https://blog.ippon.fr/2023/05/02/retour-sur-les-conferences-de-la-kubecon-2023-partie-2-3/) 
    - [cncf-k8gb-maintainers] FYI: Social Engineering Attempts on Open Source Projects - important email from Chris
- Issues:
    - New interest in Gateway (and specifically Istio gateway) support [https://github.com/k8gb-io/k8gb/issues/552#issuecomment-2026348447](https://github.com/k8gb-io/k8gb/issues/552#issuecomment-2026348447)
        - [https://github.com/k8gb-io/k8gb/issues/552#issuecomment-2041232957](https://github.com/k8gb-io/k8gb/issues/552#issuecomment-2041232957) some ideas from Yury on ingress decoupling. Comments are welcome! 
        - Andre will try to contribute Ingress decoupling! Thanks a lot!
- PRs:
    - [https://github.com/k8gb-io/k8gb/pull/1525](https://github.com/k8gb-io/k8gb/pull/1525) Documented Azure DNS deployment - Yury will prioritize review
    - [https://github.com/k8gb-io/k8gb/pull/1510](https://github.com/k8gb-io/k8gb/pull/1510) Contribfest hero  - review from Michal and Dinar required
    - [https://github.com/k8gb-io/k8gb/pull/1363](https://github.com/k8gb-io/k8gb/pull/1363) add support for svc backend - no feedback from PR author
        - Ingress decoupling strategy can also help here 


</details><details><summary><strong>Apr 3, 2024 #43</strong></summary>

## Apr 3, 2024 #43

Recording: [https://www.youtube.com/watch?v=AVfJnnSFkmw](https://www.youtube.com/watch?v=AVfJnnSFkmw) 

Attendees:

- Yury Tsarev (@ytsarev)
- Nuno Guedes (@infbase)
- Theo Chatzimichos (@tampakrap)

Agenda:

- Post KubeCon 
    - Yury did k8gb lightning talk
        - Slides: [https://docs.google.com/presentation/d/1nhCcG-JknDU\_YCipZTk4xrcakpc0ZUxT/edit?usp=sharing&ouid=103063253599686332438&rtpof=true&sd=true](https://docs.google.com/presentation/d/1nhCcG-JknDU_YCipZTk4xrcakpc0ZUxT/edit?usp=sharing&ouid=103063253599686332438&rtpof=true&sd=true) 
        - Recording: [https://www.youtube.com/watch?v=MsQ0E7SYNPo](https://www.youtube.com/watch?v=MsQ0E7SYNPo) 
    - Nuno and Yury did k8gb Contribfest
        - Slides: [https://docs.google.com/presentation/d/1eMPPEj5E1nCEiVTP6SZJ45tF-6bDz3Kz/edit?usp=sharing&ouid=103063253599686332438&rtpof=true&sd=true](https://docs.google.com/presentation/d/1eMPPEj5E1nCEiVTP6SZJ45tF-6bDz3Kz/edit?usp=sharing&ouid=103063253599686332438&rtpof=true&sd=true) 
        - No recording has been found yet - no confidence it was done for Contribfest sessions
- PRs:
    - [https://github.com/k8gb-io/k8gb/pull/1525](https://github.com/k8gb-io/k8gb/pull/1525) Documented Azure DNS deployment
    - [https://github.com/k8gb-io/k8gb/pull/1510](https://github.com/k8gb-io/k8gb/pull/1510) Contribfest hero :) 
    - [https://github.com/k8gb-io/k8gb/pull/1363](https://github.com/k8gb-io/k8gb/pull/1363) add support for svc backend - no feedback from PR author


</details><details><summary><strong>Mar 6, 2024 #42</strong></summary>

## Mar 6, 2024 #42

Attendees:

- Yury Tsarev (@ytsarev)
- Nuno Guedes (@infbase)

Agenda:

- KubeCon 
    - Nuno and Yury sync up for preparation to lightning talk and Contribfest this/next week
    - Azure public DNS support release goal before KubeCon
- [https://github.com/k8gb-io/k8gb/issues/175](https://github.com/k8gb-io/k8gb/issues/175) - disaster / service resilience tests are done internally by Nuno, we plan to document and incorporate them into the codebase later on


</details><details><summary><strong>Feb 21, 2024 #41</strong></summary>

## Feb 21, 2024 #41

Attendees:

- Yury Tsarev (@ytsarev)
- Vitor Esteves (@v-esteves)

Agenda:

- Lightning talk approved and submitted for KubeCon
- Contribfest preparation
- [nuno.guedes@millenniumbcp.pt](mailto:nuno.guedes@millenniumbcp.pt) update on Azure Private DNS
- [jiri.kremser@gmail.com](mailto:jiri.kremser@gmail.com) share the FOSDEM experience
- Status of open PRs:
    - [https://github.com/k8gb-io/k8gb/pull/1363](https://github.com/k8gb-io/k8gb/pull/1363) add support for svc backend
    - [https://github.com/k8gb-io/k8gb/pull/1422](https://github.com/k8gb-io/k8gb/pull/1422) Add support of secret based AuthN/Z for Route53 


</details><details><summary><strong>Feb 7, 2024 #40</strong></summary>

## Feb 7, 2024 #40

Attendees:

- Nuno Guedes (@infbase)
- Vitor Esteves (@v-esteves)

Agenda:

- Nuno Guedes is preparing PR for Azure Private DNS 


</details><details><summary><strong>Jan 24, 2024 #39</strong></summary>

## Jan 24, 2024 #39

Recording:

Attendees:

- Yury Tsarev (@ytsarev)
- Jiri Kremser (@jkremser)

Agenda:

- [https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/program/project-opportunities/](https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/program/project-opportunities/) 
    - Lightning-talk on March Tuesday 19 March
    - **K8gb Contribfest WAS ACCEPTED!&#32;[https://kccnceu2024.sched.com/event/1Yheq/contribfest-k8gb-contribfest?iframe=no&w=100%&sidebar=yes&bg=no](https://kccnceu2024.sched.com/event/1Yheq/contribfest-k8gb-contribfest?iframe=no&w=100%&sidebar=yes&bg=no)&#32;&#160;&#32;-** everybody welcome to join! CO-SPEAKERS wanted
- [**https://github.com/k8gb-io/k8gb/issues/974](https://github.com/k8gb-io/k8gb/issues/974)&#32;**last flaky test to fix
- [https://github.com/k8gb-io/k8gb/pull/1363](https://github.com/k8gb-io/k8gb/pull/1363) add support for svc backend - any help needed?
- [jiri.kremser@gmail.com](mailto:jiri.kremser@gmail.com) is giving a talk at FOSDEM! [https://fosdem.org/2024/schedule/event/fosdem-2024-1793-k8gb-meets-cluster-api/](https://fosdem.org/2024/schedule/event/fosdem-2024-1793-k8gb-meets-cluster-api/) ! 


</details><details><summary><strong>Jan 10, 2024 #38</strong></summary>

## Jan 10, 2024 #38

Recording: [https://www.youtube.com/watch?v=icVx9vB8OAU](https://www.youtube.com/watch?v=icVx9vB8OAU) 

Attendees:

- Yury Tsarev (@ytsarev)
- Tanuj Dwivedi (@tanujd11)
- Vitor Esteves (@v-esteves)

Agenda:

- [https://github.com/k8gb-io/k8gb/releases/tag/v0.12.2](https://github.com/k8gb-io/k8gb/releases/tag/v0.12.2) release and release pipeline fixes
- [https://github.com/cncf/toc/pull/1153](https://github.com/cncf/toc/pull/1153) The annual review is merged
- [https://github.com/k8gb-io/k8gb/pull/1374/](https://github.com/k8gb-io/k8gb/pull/1374/) The last flaky terratest is fixed
- Request to rebase good old [https://github.com/kubernetes-sigs/external-dns/pull/2835#issuecomment-1864710554](https://github.com/kubernetes-sigs/external-dns/pull/2835#issuecomment-1864710554) 
- KubeCon call [https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/program/project-opportunities/](https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/program/project-opportunities/?utm_medium=email&_hsmi=287304328&_hsenc=p2ANqtz-_cUg6aa7gb1rtYc23jJhDG54R3iyrDsTyVukwa2atokRNHI0rZ9qnoIF9wicqdJRTbAqYOS-7KZRzLxJRRYHpoETd-UA&utm_content=287304328&utm_source=hs_email)
    - Are you attending the conference?
    - If yes, are you willing to participate in k8gb booth and/or contribfest?
- [https://github.com/k8gb-io/k8gb/pull/1363](https://github.com/k8gb-io/k8gb/pull/1363) add support for svc backend #1363 
- [https://github.com/k8gb-io/k8gb/issues/1314](https://github.com/k8gb-io/k8gb/issues/1314) document split brain scenario, Tanuj will contribute


</details><details><summary><strong>Dec 13, 2023 #37</strong></summary>

## Dec 13, 2023 #37

Attendees:

- Nuno Guedes (@infbase)
- Yury Tsarev (@ytsarev)
- Tanuj Dwivedi (@tanujd11)
- 

Agenda:

- Flaky terratests fix [https://github.com/k8gb-io/k8gb/pull/1340](https://github.com/k8gb-io/k8gb/pull/1340) . It got much more stable, still not 100% pass rate. 
    - 
    - [https://github.com/k8gb-io/k8gb/issues/1345](https://github.com/k8gb-io/k8gb/issues/1345) last problematic one identified
- Planned v0.12.0 release with Cloudflare support. Do we want to include anything else in the release?
- [https://github.com/k8gb-io/k8gb/pull/1064](https://github.com/k8gb-io/k8gb/pull/1064) can be rebased and rerun with the more stable pipeline
- PrivateDNS scenario discussion
- We can apply for Incubating without waiting for Annual review!


</details><details><summary><strong>Nov 29, 2023 #36</strong></summary>

## Nov 29, 2023 #36

Attendees:

- Nuno Guedes (@infbase)

Agenda:

- Nuno Guedes is wrapping up updating Azure Private DNS example to support managed identities. PR to follow
- Vitor Esteves still going thru PR 1064 for rfc2136
- KubeCon:
    - [Take It to the Edge: Creating a Globally Distributed Ingress with Istio & K8gb - Jimmi Dyson, D2iQ](https://www.youtube.com/watch?v=4qJDkw5YGqM)
    - Submitted session “Multi-Cloud Global Content Distribution at Cloud Native Speeds” to KubeCon in Paris, for main track as well as k8s on edge pre-day





</details><details><summary><strong>Nov 15, 2023 #35</strong></summary>

## Nov 15, 2023 #35

Attendees:

- Yury Tsarev (@ytsarev)
- Michal Kuritka (@kuritka)
- Vitor Esteves (@v-esteves)

Agenda:

- CNCF Security Slam award [https://www.linkedin.com/posts/knight1776\_k8gb-activity-7130192715794825216-t10\_](https://www.linkedin.com/posts/knight1776_k8gb-activity-7130192715794825216-t10_) 
- Cloudflare support is ready for review [https://github.com/k8gb-io/k8gb/pull/1278](https://github.com/k8gb-io/k8gb/pull/1278) 
- Pipelines are still failing in [https://github.com/k8gb-io/k8gb/pull/1064](https://github.com/k8gb-io/k8gb/pull/1064) help is needed to [vitor.esteves@millenniumbcp.pt](mailto:vitor.esteves@millenniumbcp.pt) to sort it out 
    - Multicloud - rfc2136 as the glue DNS between the clouds
    - Crossplane Configuration for global app with k8gb ( [Yury Tsarev](mailto:xnullz@gmail.com))
    - [https://github.com/k8gb-io/k8gb/blob/master/Makefile#L446-L455](https://github.com/k8gb-io/k8gb/blob/master/Makefile#L446-L455) as a starting point for the investigation
- The final steps for full Azure support [nuno.guedes@millenniumbcp.pt](mailto:nuno.guedes@millenniumbcp.pt) - what is blocking us currently from claiming the Azure support?
    - Nuno will work on public DNS PR for Azure this week
- [https://github.com/k8gb-io/coredns-crd-plugin/pull/56](https://github.com/k8gb-io/coredns-crd-plugin/pull/56) community contributions review ([kuritka@gmail.com](mailto:kuritka@gmail.com))


</details><details><summary><strong>Nov 1, 2023 #34</strong></summary>

## Nov 1, 2023 #34

Attendees:

- Yury Tsarev (@ytsarev)
- Jiri Kremser (@jkremser)
- Michal Kuritka (@kuritka)

Agenda:

- 
- We won the security slam again. Super huge big kudos to [jiri.kremser@gmail.com](mailto:jiri.kremser@gmail.com)!!!
    - Look into other repos like coredns-crd-plugin and external-dns fork in future!
- New contributor [https://github.com/k8gb-io/coredns-crd-plugin/pull/55](https://github.com/k8gb-io/coredns-crd-plugin/pull/55) from D2IQ!
- [https://github.com/k8gb-io/k8gb/issues/1314](https://github.com/k8gb-io/k8gb/issues/1314) - split brain documentation request from Redhat
- 
- should we ping the reviewer about [https://github.com/cncf/toc/pull/1153](https://github.com/cncf/toc/pull/1153) ?
    - [Yury Tsarev](mailto:xnullz@gmail.com) will look at the schedule and ping cncf toc mailing lists
    - The pace of other annual review PRs deviates, apparently 2 approvers are required
    - [jiri.kremser@gmail.com](mailto:jiri.kremser@gmail.com) will try to ping cncf mailing list to [cncf-toc@lists.cncf.io](mailto:cncf-toc@lists.cncf.io) 
- eddie asked about contributing our self-assessment.md into [https://github.com/cncf/tag-security/tree/main/assessments/projects](https://github.com/cncf/tag-security/tree/main/assessments/projects) 
    - We track it at [https://github.com/k8gb-io/k8gb/issues/1322](https://github.com/k8gb-io/k8gb/issues/1322) 
- 2 quality of life improvements for renovate & testing:
    - [https://github.com/k8gb-io/k8gb/pull/1320](https://github.com/k8gb-io/k8gb/pull/1320)
    - [https://github.com/k8gb-io/k8gb/pull/1321](https://github.com/k8gb-io/k8gb/pull/1321)
    - [https://github.com/k8gb-io/k8gb/issues/1323](https://github.com/k8gb-io/k8gb/issues/1323) - follow up for renovate exposure


</details><details><summary><strong>Oct 18, 2023 #33</strong></summary>

## Oct 18, 2023 #33

Attendees:

- Yury Tsarev (@ytsarev)
- Jiri Kremser (@jkremser)
- Vitor Esteves (@v-esteves)
- Michal Kuritka (@kuritka)
- Nuno Guedes (@infbase)

Agenda:

- [https://github.com/k8gb-io/k8gb/pull/1293](https://github.com/k8gb-io/k8gb/pull/1293) - New adopter! 
    - Makes us eligible for Incubation(!), reflected here [https://github.com/cncf/toc/pull/1153#issuecomment-1754648129](https://github.com/cncf/toc/pull/1153#issuecomment-1754648129) 
- Cloudflare support [https://github.com/k8gb-io/k8gb/pull/1278](https://github.com/k8gb-io/k8gb/pull/1278)
- [https://opensource.microsoft.com/azure-credits/](https://opensource.microsoft.com/azure-credits/) azure credits
- Security Slam
    - **K8gb is used by the US Space Force.!!!**
- OLM automation for community-operators-prod repo 
    - [https://github.com/k8gb-io/k8gb/pull/1308](https://github.com/k8gb-io/k8gb/pull/1308)
    - [https://github.com/redhat-openshift-ecosystem/community-operators-prod/pull/3472](https://github.com/redhat-openshift-ecosystem/community-operators-prod/pull/3472)
- k8gb lite migration plan
    - [kuritka@gmail.com](mailto:kuritka@gmail.com) created a set of GH issues [https://github.com/k8gb-io/k8gb/issues?q=is%3Aissue+is%3Aopen+label%3Ak8gb-lite](https://github.com/k8gb-io/k8gb/issues?q=is%3Aissue+is%3Aopen+label%3Ak8gb-lite) 👍
- Failed helm publish pipeline - [https://github.com/k8gb-io/k8gb/actions/runs/6146868686/job/16677091875](https://github.com/k8gb-io/k8gb/actions/runs/6146868686/job/16677091875)
- Hacktoberfest participation - [jiri.kremser@gmail.com](mailto:jiri.kremser@gmail.com) volunteered to bootstrap! 👍
- [https://github.com/k8gb-io/k8gb/pull/1064](https://github.com/k8gb-io/k8gb/pull/1064) pipes are failing; investigation needed [vitor.esteves@millenniumbcp.pt](mailto:vitor.esteves@millenniumbcp.pt)


</details><details><summary><strong>Sep 6, 2023 #32</strong></summary>

## Sep 6, 2023 #32

Attendees:

- Yury Tsarev (@ytsarev)
- Michal Kuritka (@kuritka)

Agenda:

- [https://github.com/k8gb-io/toc/pull/3](https://github.com/k8gb-io/toc/pull/3) maintainers, please review the CNCF Annual report before we send it to CNCF TOC [Dinar “k0da” Valeev](mailto:dinarv@gmail.com) [kuritka@gmail.com](mailto:kuritka@gmail.com) [jiri.kremser@gmail.com](mailto:jiri.kremser@gmail.com) [donovan.muller@gmail.com](mailto:donovan.muller@gmail.com) - done -  The final [https://github.com/cncf/toc/pull/1153](https://github.com/cncf/toc/pull/1153) was created and email to CNCF TOC was sent with the required associated announcement. (@ytsarev)  - we got positive feedback -**&#32;all set** [https://github.com/cncf/toc/pull/1153#pullrequestreview-1598567037](https://github.com/cncf/toc/pull/1153#pullrequestreview-1598567037) 
- [https://github.com/orgs/k8gb-io/projects/2](https://github.com/orgs/k8gb-io/projects/2) roadmap cleanup
- k8gb lite migration plan
    - [kuritka@gmail.com](mailto:kuritka@gmail.com) will create GH issue with the migration plan  
- Cloudflare support [https://github.com/k8gb-io/k8gb/issues/884](https://github.com/k8gb-io/k8gb/issues/884) - interest and willingness to test from the community member - [Yury Tsarev](mailto:xnullz@gmail.com) will prepare the initial PR for the testing
- Rancheer fleet support - https://github.com/k8gb-io/k8gb/pull/1255


</details><details><summary><strong>Aug 9, 2023 #31</strong></summary>

## Aug 9, 2023 #31

Attendees:

- Michal Kuritka (@kuritka)
- Yury Tsarev (@ytsarev)
- Jiri Kremser (@jkremser)
- Dinar Valeev (@k0da)
- Vitor Esteves (@v-esteves)

Agenda:

- [https://github.com/k8gb-io/k8gb/pull/1064](https://github.com/k8gb-io/k8gb/pull/1064) and [https://github.com/k8gb-io/k8gb/pull/1065](https://github.com/k8gb-io/k8gb/pull/1065) Looks like 1064 should be deprecated and we should proceed with 1065 only? [vitor.esteves@millenniumbcp.pt](mailto:vitor.esteves@millenniumbcp.pt)
- [https://www.youtube.com/watch?v=vrDCUIVyc4g&list=PLo4lFffE9Ct9rvNjSOOL64VTs\_qVnrSnI&index=43](https://www.youtube.com/watch?v=vrDCUIVyc4g&list=PLo4lFffE9Ct9rvNjSOOL64VTs_qVnrSnI&index=43) k8gb presentation at KCD Bengaluru. Thanks a ton, Tamil!
- [https://github.com/k8gb-io/k8gb/issues/1206](https://github.com/k8gb-io/k8gb/issues/1206) CNCF TOC Annual review
    - We should also consider as part of the review [https://github.com/cncf/toc/blob/main/process/graduation\_criteria.md#incubating-stage](https://github.com/cncf/toc/blob/main/process/graduation_criteria.md#incubating-stage) 
    - [Yury Tsarev](mailto:xnullz@gmail.com) will create an annual review PR(done)
- k8gb lite migration plan
    - [kuritka@gmail.com](mailto:kuritka@gmail.com) will create GH issue with the migration plan 
- Recent release [https://github.com/k8gb-io/k8gb/releases/tag/v0.11.2](https://github.com/k8gb-io/k8gb/releases/tag/v0.11.2) 
    - [https://github.com/k8gb-io/coredns-crd-plugin/pull/54](https://github.com/k8gb-io/coredns-crd-plugin/pull/54) 


</details><details><summary><strong>Jul 26, 2023 #30</strong></summary>

## Jul 26, 2023 #30

Attendees:

- Vitor Esteves (@v-esteves)
- Nuno Guedes (@infbase)
- K Tamil Vanan
- Michal Kuritka (@kuritka)

Agenda:

- [https://github.com/k8gb-io/k8gb/pull/106](https://github.com/k8gb-io/k8gb/pull/1065)4
- [https://github.com/k8gb-io/k8gb/pull/1065](https://github.com/k8gb-io/k8gb/pull/1065)
- updates on Azure support
- Other questions

Action:

	[jiri.kremser@gmail.com](mailto:jiri.kremser@gmail.com) could you look at the PRs above?


</details><details><summary><strong>Jun 28, 2023 #29</strong></summary>

## Jun 28, 2023 #29

Attendees:

- Michal Kuritka (@kuritka)
- Jiri Kremser (@jkremser)
- Yury Tsarev (@ytsarev)
- Nuno Guedes (@infbase)

Agenda:

- k8gb-lite
- [https://github.com/ossf/scorecard-action](https://github.com/ossf/scorecard-action) scorecard action setup
    - [https://github.com/k8gb-io/k8gb/pull/1193](https://github.com/k8gb-io/k8gb/pull/1193) first iteration
- Any updates on Azure support?
- Pipelines are still frequently producing false negatives


</details><details><summary><strong>May 31, 2023 #28</strong></summary>

## May 31, 2023 #28

Attendees:

- Vitor Esteves (@v-esteves)
- Nuno Guedes (@infbase)
- Michal Kuritka (@kuritka)
- Jiri Kremser (@jkremser)
- Yury Tsarev (@ytsarev)

Agenda:

- 600 stars on GitHub celebration! (610 already!) 
- Jiri’s great contributions :D
    - [https://github.com/k8gb-io/k8gb/pull/1186](https://github.com/k8gb-io/k8gb/pull/1186) 
    - [https://github.com/k8gb-io/k8gb/pull/1185](https://github.com/k8gb-io/k8gb/pull/1185) 
- Azure Public DNS support [https://github.com/k8gb-io/k8gb/pull/912](https://github.com/k8gb-io/k8gb/pull/912) 
    - Nuno will own the PR delivery
- Further steps with infra support from MS
    - Connection with Open Source Initiative manager? Nuno will take care
    - Yury will evaluate CNCF Service Desk for Sandbox to align 
- Oh no, CLOMonitor degraded from 100 to 95 :D Apparently, we need to implement new [https://clomonitor.io/docs/topics/checks/#openssf-scorecard-badge](https://clomonitor.io/docs/topics/checks/#openssf-scorecard-badge) 
    - Yury will take care
    - Create Issue on ssf-scorecard linting pipeline in k8gb
    - Check badge levels (potentially upgrade to Silver or Gold)
        - silver: [https://bestpractices.coreinfrastructure.org/en/projects/4866?criteria\_level=1](https://bestpractices.coreinfrastructure.org/en/projects/4866?criteria_level=1)
        - gold: [https://bestpractices.coreinfrastructure.org/en/projects/4866?criteria\_level=2](https://bestpractices.coreinfrastructure.org/en/projects/4866?criteria_level=2)
- [https://github.com/k8gb-io/k8gb/issues/175](https://github.com/k8gb-io/k8gb/issues/175) test that splitbrain(horizontal network partitioning) is not required in Azure/Route53/in general. 
    - Instead of the complex logic, be verbose in case of horizontal network partitioning and split-brain
        - Clarify logging message
        - Potentially propagate to Gslb status
    - Test failure scenarios to actually deeply understand what will happen


</details><details><summary><strong>May 17, 2023 #27</strong></summary>

## May 17, 2023 #27

Attendees:

- Yury Tsarev (@ytsarev)

Nobody else was able to join, so transferring agenda to the next one


</details><details><summary><strong>May 3, 2023 #26</strong></summary>

## May 3, 2023 #26

Attendees:

- Vitor Esteves (@v-esteves)
- Nuno Guedes (@infbase)
- Michal Kuritka (@kuritka)

Agenda:

- Kubecon
    - Over 600 attendees showed interest in the session
    - Positive feedback from session attendees
    - Video: [https://www.youtube.com/watch?v=U46hlF0Z3xs](https://www.youtube.com/watch?v=U46hlF0Z3xs) 
    - Demo code: [https://dev.azure.com/infbase/k8gb-kubeconeu2023](https://dev.azure.com/infbase/k8gb-kubeconeu2023) 
- PRs open
    - Vitor’s contributions
        - [https://github.com/k8gb-io/k8gb/pull/1065](https://github.com/k8gb-io/k8gb/pull/1065) 
        - [https://github.com/k8gb-io/k8gb/pull/1064](https://github.com/k8gb-io/k8gb/pull/1064) 
- [Roadmap](https://github.com/orgs/k8gb-io/projects/2/views/1)
    - Azure Public DNS support
- WeightedRoundRobin merged


</details><details><summary><strong>Apr 5, 2023 #25</strong></summary>

## Apr 5, 2023 #25

Attendees:

- Yury Tsarev (@ytsarev)
- Vitor Esteves (@v-esteves)
- Nuno Guedes (@infbase)

Agenda:

- PRs open
    - Vitor’s contributions
        - [https://github.com/k8gb-io/k8gb/pull/1065](https://github.com/k8gb-io/k8gb/pull/1065) 
        - [https://github.com/k8gb-io/k8gb/pull/1064](https://github.com/k8gb-io/k8gb/pull/1064) 
    - Trivial from Yury
        - [https://github.com/k8gb-io/k8gb/pull/1111](https://github.com/k8gb-io/k8gb/pull/1111) 
- PRs merged
    - Michal’s 
        - [https://github.com/k8gb-io/k8gb/pull/1101](https://github.com/k8gb-io/k8gb/pull/1101)
    - We welcome new contributor [https://github.com/eliasbokreta](https://github.com/eliasbokreta)  to the project!
        - [https://github.com/k8gb-io/k8gb/pull/1112](https://github.com/k8gb-io/k8gb/pull/1112) 
- [Roadmap](https://github.com/orgs/k8gb-io/projects/2/views/1)
    - Anything to target specifically before the kubecon?
        - External-dns NS fork and announce Azure support?
        - Project visibility/cosmetics: Pipelines, badges
- Kubecon
    - Yury and Nuno need to crack the preparation for [https://kccnceu2023.sched.com/event/1HyW3](https://kccnceu2023.sched.com/event/1HyW3) 
    - Azure / GCP multicloud test? VPN between clusters
    - Create Crossplane based Azure/Azure GlobalEKS abstraction ++ Thing to try first ++
    - Demo env from MBCP? Cool. Nuno will set it up.
    - We will meet with Nuno next week. Vitor will aslo help(thanks!)
    - Global 2 clusters ‘global by default’ pattern to share


</details><details><summary><strong>Mar 23, 2023 #24</strong></summary>

## Mar 23, 2023 #24

Attendees

- Michal Kuritka (@kuritka)
- Nuno Guedes (@infbase)
- Vitor Esteves (@v-esteves)

Agenda

- Pipelines are flaky and unstable - should we remove heavy terratest pipelines behind the PR comment?
- [https://github.com/crossplane/crossplane/blob/master/ADOPTERS.md](https://github.com/crossplane/crossplane/blob/master/ADOPTERS.md) can we get Absa and MBCP there? This is important to enter Incubating level process.
- [https://kccnceu2023.sched.com/event/1HyW3](https://kccnceu2023.sched.com/event/1HyW3) who joins KubeCon?
    - Nuno
    - Jiri
    - Yury
    - ( Michal, Dinar, Vito ) virtually

Notes

- we've gone through each of the points from last time. 
- - Vitor is working on DNS related PR and incorporating comments
- - We have successfully added adopters. The plan is to go to v1.0.0 ASAP
- - k8gb-lite integration will come later (after v1.0.0 release)
- - FOSSA issues - skipped at this meeting
- - OCI repo - pre-GHCR
- -simplifying Github pipelines: started by removing dependabot alerts and cleaning it up a bit. The PR is open but needs clarification in the PR conversation 

Action Items

- Ask Vitor to present his PRs at next meeting - Windows DNS related - guys are focusing on this PR’s 

- Make pipelines more lightweight?
    - Dep updates are special - scope to go.mod?
    - Gate with PR comment?
    - Create an issue for pipe optimization
- Absa to Adopters.md double check with Don
    - Yury to create PR, Michal(Absa) to approve it [https://github.com/k8gb-io/k8gb/pull/1080](https://github.com/k8gb-io/k8gb/pull/1080)
    - MBCP goes next
- [https://github.com/kubernetes-sigs/external-dns/pull/2835#issuecomment-1419453279](https://github.com/kubernetes-sigs/external-dns/pull/2835#issuecomment-1419453279) progress there, still waiting for the merge
- OCI repo for k8gb - it still absaoss ? we should move somewhere neutral 
    - Dinar; move to ghcr ? github registry
    - Jiri: cosign might need some changes
- FOSSA:  “License scan found 2 issues”
    - [https://app.fossa.com/projects/custom%2B162%2Fgithub.com%2Fk8gb-io%2Fk8gb/refs/branch/master/dc9d74afe914f9d78480a42b4e12bb2fc881ff81/preview](https://app.fossa.com/projects/custom%2B162%2Fgithub.com%2Fk8gb-io%2Fk8gb/refs/branch/master/dc9d74afe914f9d78480a42b4e12bb2fc881ff81/preview) we do not have access anymore(?)
    - https://app.fossa.com/projects/custom%2B162%2Fgithub.com%2Fk8gb-io%2Fk8gb/refs/branch/master/dc9d74afe914f9d78480a42b4e12bb2fc881ff81/browse/dependencies/?flagged=true:
    - 
- K8gb light into main 
    - Code is ready
    - Pipelines and additional project stuff is not
    - See changes: [https://github.com/k8gb-io/k8gb-lite/issues/1](https://github.com/k8gb-io/k8gb-lite/issues/1)


</details><details><summary><strong>Mar 8, 2023 #23</strong></summary>

## Mar 8, 2023 #23

Attendees

- Jiri Kremser (@jkremser)
- Michal Kuritka (@kuritka)
- Yury Tsarev(@ytsarev)

Agenda

- Pipelines are flaky and unstable - should we remove heavy terratest pipelines behind the PR comment?
- [https://github.com/crossplane/crossplane/blob/master/ADOPTERS.md](https://github.com/crossplane/crossplane/blob/master/ADOPTERS.md) can we get Absa and MBCP there? This is important to enter Incubating level process.
- [https://kccnceu2023.sched.com/event/1HyW3](https://kccnceu2023.sched.com/event/1HyW3) who joins KubeCon?

Notes

- 

Action Items

- Ask Vitor to present his PRs at next meeting - Windows DNS related
- Make pipelines more lightweight?
    - Dep updates are special - scope to go.mod?
    - Gate with PR comment?
    - Create an issue for pipe optimization
- Absa to Adopters.md double check with Don
    - Yury to create PR, Michal(Absa) to approve it [https://github.com/k8gb-io/k8gb/pull/1080](https://github.com/k8gb-io/k8gb/pull/1080)
    - MBCP goes next
- Kubecon
    - Jiri - yes
- [https://github.com/kubernetes-sigs/external-dns/pull/2835#issuecomment-1419453279](https://github.com/kubernetes-sigs/external-dns/pull/2835#issuecomment-1419453279) progress there, still waiting for the merge
- OCI repo for k8gb - it still absaoss ? we should move somewhere neutral 
    - Dinar; move to ghcr ? github registry
    - Jiri: cosign might need some changes
- FOSSA:  “License scan found 2 issues”
    - [https://app.fossa.com/projects/custom%2B162%2Fgithub.com%2Fk8gb-io%2Fk8gb/refs/branch/master/dc9d74afe914f9d78480a42b4e12bb2fc881ff81/preview](https://app.fossa.com/projects/custom%2B162%2Fgithub.com%2Fk8gb-io%2Fk8gb/refs/branch/master/dc9d74afe914f9d78480a42b4e12bb2fc881ff81/preview) we do not have access anymore(?)
    - https://app.fossa.com/projects/custom%2B162%2Fgithub.com%2Fk8gb-io%2Fk8gb/refs/branch/master/dc9d74afe914f9d78480a42b4e12bb2fc881ff81/browse/dependencies/?flagged=true:
    - 
- K8gb light into main 
    - Code is ready
    - Pipelines and additional project stuff is not
    - See changes: [https://github.com/k8gb-io/k8gb-lite/issues/1](https://github.com/k8gb-io/k8gb-lite/issues/1)





</details><details><summary><strong>Jan 17, 2023 #22</strong></summary>

## Jan 17, 2023 #22

Attendees

- Timofey Ilinykh (@somaritane)
- Michal Kuritka(@kuritka)
- Yury Tsarev (@ytsarev)
- Vitor Esteves (@v-esteves)
- Nuno Guedes (@infbase)

Agenda

- Action Items
- Backporting
- K8gb light (getting rid of GSLB, use standard primitives with annotations)
- Project roadmap: [https://github.com/orgs/k8gb-io/projects/2/views/2](https://github.com/orgs/k8gb-io/projects/2/views/2) 

Notes

- K8gb light (getting rid of GSLB, use standard primitives with annotations)
    - Nuno: might definitely help end-uses with onboarding, adopting new CRD is hard

Action Items

- Revisit >2 cluster strategy support ([https://github.com/k8gb-io/k8gb/pull/815](https://github.com/k8gb-io/k8gb/pull/815)) (@jkremser) not done for failover yet (comment on issue)
Michal: Have some solution in case of 3 (more clusters?)
AP on Michal: follow-up to discuss the architecture
[Michal K](mailto:kuritka@gmail.com): see 
    - [https://github.com/k8gb-io/k8gb-light/blob/main/terratest/test/resources/ingress\_fo3\_ordered1.yaml

](https://github.com/k8gb-io/k8gb-light/blob/main/terratest/test/resources/ingress_fo3_ordered1.yaml)
    - [https://github.com/k8gb-io/k8gb-light/blob/main/terratest/test/resources/ingress\_fo3\_ordered2.yaml](https://github.com/k8gb-io/k8gb-light/blob/main/terratest/test/resources/ingress_fo3_ordered2.yaml)
    - [https://github.com/k8gb-io/k8gb-light/blob/main/terratest/test/fo\_lifecycle\_3\_clusters\_ordered\_test.go#L31-L32

](https://github.com/k8gb-io/k8gb-light/blob/main/terratest/test/fo_lifecycle_3_clusters_ordered_test.go#L31-L32)(this is covered by the test and terrateston three clusters )
All failover ordering is defined in k8gb.io/primary-geotag. If k8gb.io/primary-geotag does not contain all available clusters , the remaining clusters are sorted alphabetically.

- For example, we have the clusters "eu, us, uk, za", 
- for "k8gb.io/primary-geotag: uk,us" the failover order is "uk, us, eu, za"
- for "k8gb.io/primary-geotag: za" the failover order is "za, eu, uk, us"
- for "k8gb.io/primary-geotag: eu,us,za" the failover is "eu, us, za, uk"




- [https://github.com/k8gb-io/k8gb/issues/642](https://github.com/k8gb-io/k8gb/issues/642) (Azure Support)
    - Fork external-dns and merge [https://github.com/kubernetes-sigs/external-dns/issues/2826](https://github.com/kubernetes-sigs/external-dns/issues/2826) 
Forks done, didn’t update env for that yet
        - [https://github.com/kubernetes-sigs/external-dns/pull/2835](https://github.com/kubernetes-sigs/external-dns/pull/2835) 
        - Need to incorporate this PR to the fork
        - Need help to test on public azure
    - Ping Nuno Guedes (@infbase) on the case documentation
    - Vitor: not using public but private dns. Private azuredns doesn’t support NS records creation. (using windowsdns as resolver )
    - Happy to share the use-case and workaround
- Cert-manager integration issue: [https://blog.aba](https://blog.abaganon.com/going-global-with-kubernetes-490cf51e2bf8)[https://drive.proton.me/urls/MHENT68VER#AM8e6BqdFvmv](https://drive.proton.me/urls/MHENT68VER#AM8e6BqdFvmv)[ganon.com/goin](https://blog.abaganon.com/going-global-with-kubernetes-490cf51e2bf8)[https://drive.proton.me/urls/MHENT68VER#AM8e6BqdFvmv](https://drive.proton.me/urls/MHENT68VER#AM8e6BqdFvmv)[g-global-with-kubernetes-490cf51e2bf8](https://blog.abaganon.com/going-global-with-kubernetes-490cf51e2bf8) 
    - Create associated ticket for TLS and invite Eric for participation. Provide Vault backend as an example for mitigating that issue. 
    - Create tutorial on related topic. ([Slack discussion](https://cloud-native.slack.com/archives/C021P656HGB/p1650293408364539)) (@k0da, @jkremser)
    - Tested locally, the name’s not resolvable, need to merge the related PR and deploy to aws for proper testing 
    - Refresh with @k0da what was the case 
- WRR, CRD plugin PR: (@kuritka to make a presentation on changes) Presentation is done.&#160;
    - Make it covered by terratest. 
    - Ingress support 
    - K8gb coredns is not deployed yet
    - Fallback for corner cases
    - Needs to be integrated with k8gb (waiting for PR): [https://github.com/k8gb-io/coredns-crd-plugin/pull/45](https://github.com/k8gb-io/coredns-crd-plugin/pull/45) 
- Cloudflare support: +1 ask. Need to implement support for NS and TXT in external-dns provider (???) ([https://github.com/k8gb-io/k8gb/issues/884](https://github.com/k8gb-io/k8gb/issues/884)) 
- usage of k8gb without edgeDNS: provide our view on the topic (@k0da)


</details><details><summary><strong>Nov 29, 2022 #21</strong></summary>

## Nov 29, 2022 #21

Attendees

- Yury Tsarev (@ytsarev)
- Timofey Ilinykh (@somaritane)
- Michal Kuritka(@kuritka)
- Dinar Valeev (@k0da)

Agenda

- K8gb service HA questions: [https://github.com/k8gb-io/k8gb/issues/1035](https://github.com/k8gb-io/k8gb/issues/1035)
- Action Items
- Backporting
- Project roadmap: [https://github.com/orgs/k8gb-io/projects/2/views/2](https://github.com/orgs/k8gb-io/projects/2/views/2) 
- SecuritySlam quotes
- ADOPTERS.md

Notes

- K8gb service HA questions:
    - 1.1 supported
    - 1.2 not supportede
    - 2.1. Yes, but need to update helm to wire edgedns settings for that type
    - 2.2
    - 2.3 Multiple EdgeDNS servers: done by @jkremser, but not tested. Also support for multiple EdgeDNS types is needed ([https://github.com/k8gb-io/k8gb/issues/919](https://github.com/k8gb-io/k8gb/issues/919)) 
    - AP: Need to reflect possible HA scenarios in the documentation
- SecuritySlam quotes:
    - Draft in slack today
- ADOPTERS.md
    - First PR from ABSA as the first adopter (was shared to CNCF 1 year ago)
    - 2nd adopter
    - 3 non-SW vendors as adopters for Incubation project phase
- Jiri: [https://github.com/k8gb-io/k8gb/pull/1021](https://github.com/k8gb-io/k8gb/pull/1021) can we move with this one, CLO Monitor linter, is there a GH app for CLOlinter
- Michal: Thinks to discuss: removal of the GLSB and focusing on Ingress only

Action Items

- Bring the bevvy screen sharing issue to KubeCon/CNCF organizers as feedback (@somaritane) (CNCF service desk ticket + email )
- Revisit >2 cluster strategy support ([https://github.com/k8gb-io/k8gb/pull/815](https://github.com/k8gb-io/k8gb/pull/815)) (@jkremser) not done for failover yet (comment on issue)
- [https://github.com/k8gb-io/k8gb/issues/642](https://github.com/k8gb-io/k8gb/issues/642) (Azure Support)
    - Fork external-dns and merge [https://github.com/kubernetes-sigs/external-dns/issues/2826](https://github.com/kubernetes-sigs/external-dns/issues/2826) 
Forks done, didn’t update env for that yet
    - Ping Nuno Guedes (@infbase) on the case documentation
- Cert-manager integration issue: [https://blog.abaganon.com/going-global-with-kubernetes-490cf51e2bf8](https://blog.abaganon.com/going-global-with-kubernetes-490cf51e2bf8) 
    - Create associated ticket for TLS and invite Eric for participation. Provide Vault backend as example for mitigating that issue. 
    - Create tutorial on related topic. ([Slack discussion](https://cloud-native.slack.com/archives/C021P656HGB/p1650293408364539)) (@k0da, @jkremser)
    - Tested locally, the name’s not resolvable, need to merge the related PR and deploy to aws for proper testing 
- Ping infoblox on v2 release (@somaritane) and check on docker/k3s ways of installation. (there’s flask Python API mock app)
No response from the Infoblox client maintainers. On Hold
- WRR, CRD plugin PR: (@kuritka to make a presentation on changes) Presentation is done.&#160;
    - Make it covered by terratest. 
    - Ingress support 
    - K8gb coredns is not deployed yet
    - Fallback for corner cases
- ~~CFP for KubeCon EU 2023 on WRR (@kuritka)~~
- Cloudflare support: +1 ask. Need to implement support for NS and TXT in external-dns provider (???) ([https://github.com/k8gb-io/k8gb/issues/884](https://github.com/k8gb-io/k8gb/issues/884)) 
- usage of k8gb without edgeDNS: provide our view on the topic (@k0da)


</details><details><summary><strong>15 Nov 2022 #20</strong></summary>

## 15 Nov 2022 #20

The meeting is cancelled due to project members availability


</details><details><summary><strong>1 Nov 2022 #19</strong></summary>

## 1 Nov 2022 #19

Attendees

- Yury Tsarev (@ytsarev)
- Timofey Ilinykh (@somaritane)
- Michal Kuritka(@kuritka)
- Dinar Valeev (@k0da)

Agenda

- Action Items
- KubeCon NA office hours retro: bevy, equipment, steps to undertake
- Project roadmap: [https://github.com/orgs/k8gb-io/projects/2/views/2](https://github.com/orgs/k8gb-io/projects/2/views/2) 
- SecuritySlam quotes
- ADOPTERS.md

Notes

- KubeCon NA office hours retro
    - Use maintainers channel as priority channel (from @kuritka,, does it make sense to involve support by default into thread ? )
    - Use non-corp equipment during the call
    - Have a failover plan (videos, list of people to take over the presentation)
    - Platform is being changed every single year. Switch format to the panel discussion? 
    - K8gb requires an intro as it’s not that known as k8s. Make it short. One slide, failover demo, then discussion, Q&A
    - Feedback:
        - Ticket to CNCF service desk  (AP)
        - Email on the situation, ask to get real test session and try to understand how to mitigate, propose to have a control for screen-sharing. Propose the chance to test upfront not 15 mins before the presentation. (AP)
- SecuritySlam quotes:
    - Draft in slack today
- ADOPTERS.md
    - First PR from ABSA as the first adopter (was shared to CNCF 1 year ago)
    - 2nd adopter
    - 3 non-SW vendors as adopters for Incubation project phase

Action Items

- ~~Check for low-hanging fruit actions to raise the score for KubeCon Security Slam event ([https://clomonitor.io/projects/cncf/k8gb#k8gb\_security](https://clomonitor.io/projects/cncf/k8gb#k8gb_security)) (@ytsarev, @somaritane) (nice talk:&#32;[https://www.youtube.com/watch?v=iZpFtalj4xE](https://www.youtube.com/watch?v=iZpFtalj4xE)) (@jkremser - added myself)~~
- ~~Ask KubeCon organizers on bevvy meeting pre-test (@somaritane,@ytsarev)~~
- Bring the bevvy screen sharing issue to KubeCon/CNCF organizers as feedback (@somaritane) (CNCF service desk ticket + email )
- Revisit >2 cluster strategy support ([https://github.com/k8gb-io/k8gb/pull/815](https://github.com/k8gb-io/k8gb/pull/815)) (@jkremser) not done for failover yet (comment on issue)
- [https://github.com/k8gb-io/k8gb/issues/642](https://github.com/k8gb-io/k8gb/issues/642) (Azure Support)
    - Fork external-dns and merge [https://github.com/kubernetes-sigs/external-dns/issues/2826](https://github.com/kubernetes-sigs/external-dns/issues/2826) 
Forks done, didn’t update env for that yet
    - Ping Nuno Guedes (@infbase) on the case documentation
- Cert-manager integration issue: [https://blog.abaganon.com/going-global-with-kubernetes-490cf51e2bf8](https://blog.abaganon.com/going-global-with-kubernetes-490cf51e2bf8) 
    - Create associated ticket for TLS and invite Eric for participation. Provide Vault backend as example for mitigating that issue. 
    - Create tutorial on related topic. ([Slack discussion](https://cloud-native.slack.com/archives/C021P656HGB/p1650293408364539)) (@k0da, @jkremser)
    - Tested locally, the name’s not resolvable, need to merge the related PR and deploy to aws for proper testing 
- Ping infoblox on v2 release (@somaritane) and check on docker/k3s ways of installation. (there’s flask Python API mock app)
No response from the Infoblox client maintainers. On Hold
- WRR, CRD plugin PR: (@kuritka to make a presentation on changes) Presentation is done.&#160;
    - Make it covered by terratest. 
    - Ingress support 
    - K8gb coredns is not deployed yet
    - Fallback for corner cases
- CFP for KubeCon EU 2023 on WRR (@kuritka)
- Cloudflare support: +1 ask. Need to implement support for NS and TXT in external-dns provider (???) ([https://github.com/k8gb-io/k8gb/issues/884](https://github.com/k8gb-io/k8gb/issues/884)) 
- ~~Update project main page with roadmap and ensure visibility to the community (@somaritane)~~
- usage of k8gb without edgeDNS: provide our view on the topic (@k0da)

***


</details><details><summary><strong>18 Oct 2022 #18</strong></summary>

## 18 Oct 2022 #18

Attendees

- Yury Tsarev (@ytsarev)
- Timofey Ilinykh (@somaritane)
- Jiri Kremser (@jkremser)
- Michal Kuritka(@kuritka)

Agenda

- Action Items
- KubeCon NA prep
- Project roadmap: [https://github.com/orgs/k8gb-io/projects/2/views/2](https://github.com/orgs/k8gb-io/projects/2/views/2) 
- SecuritySlam
- cosign - demo (sboms, signed images, provenance)
- Project office hours timeline

Notes

- 

Action Items

- Check for low-hanging fruit actions to raise the score for KubeCon Security Slam event ([https://clomonitor.io/projects/cncf/k8gb#k8gb\_security](https://clomonitor.io/projects/cncf/k8gb#k8gb_security)) (@ytsarev, @somaritane) (nice talk: [https://www.youtube.com/watch?v=iZpFtalj4xE](https://www.youtube.com/watch?v=iZpFtalj4xE)) (@jkremser - added myself)
- Ask KubeCon organizers on bevvy meeting pre-test (@somaritane)
- Revisit >2 cluster strategy support ([https://github.com/k8gb-io/k8gb/pull/815](https://github.com/k8gb-io/k8gb/pull/815)) (@jkremser) not done for failover yet (comment on issue)
- [https://github.com/k8gb-io/k8gb/issues/642](https://github.com/k8gb-io/k8gb/issues/642) (Azure Support)
    - Fork external-dns and merge [https://github.com/kubernetes-sigs/external-dns/issues/2826](https://github.com/kubernetes-sigs/external-dns/issues/2826) 
Forks done, didn’t update env for that yet
    - Ping Nuno Guedes (@infbase) on the case documentation
- Cert-manager integration issue: [https://blog.abaganon.com/going-global-with-kubernetes-490cf51e2bf8](https://blog.abaganon.com/going-global-with-kubernetes-490cf51e2bf8) 
    - Create associated ticket for TLS and invite Eric for participation. Provide Vault backend as example for mitigating that issue. 
    - Create tutorial on related topic. ([Slack discussion](https://cloud-native.slack.com/archives/C021P656HGB/p1650293408364539)) (@k0da, @jkremser)
- Ping infoblox on v2 release (@somaritane) and check on docker/k3s ways of installation. (there’s flask Python API mock app)
- WRR, CRD plugin PR: (@kuritka to make a presentation on changes)
- Cloudflare support: +1 ask. Need to implement support for NS and TXT in external-dns provider (???) ([https://github.com/k8gb-io/k8gb/issues/884](https://github.com/k8gb-io/k8gb/issues/884)) 
- Update project main page with roadmap and ensure visibility to the community (@somaritane)
- usage of k8gb without edgeDNS: provide our view on the topic (@k0da)

***


</details><details><summary><strong>5 Oct 2022 #17</strong></summary>

## 5 Oct 2022 #17

Attendees

- Yury Tsarev (@ytsarev)
- Timofey Ilinykh (@somaritane)
- Jiri Kremser (@jkremser)
- Michal Kuritka(@kuritka)

Agenda

- Action Items
- KubeCon NA prep
- Project roadmap: [https://github.com/orgs/k8gb-io/projects/2/views/2](https://github.com/orgs/k8gb-io/projects/2/views/2) 
- Presentation from @kuritka on WRR
- We are invited to [https://events.linuxfoundation.org/kubecon-cloudnativecon-north-america/attend/experiences/#security-slam](https://events.linuxfoundation.org/kubecon-cloudnativecon-north-america/attend/experiences/#security-slam) 
    - TL;DR we need to hit 100 in [https://clomonitor.io/projects/cncf/k8gb#k8gb\_security](https://clomonitor.io/projects/cncf/k8gb#k8gb_security). We already have 75, yey
- Documentation improvements
- Grafana dashboards for k8gb

Notes

- 

Action Items

- Check for low-hanging fruit actions to raise the score for KubeCon Security Slam event ([https://clomonitor.io/projects/cncf/k8gb#k8gb\_security](https://clomonitor.io/projects/cncf/k8gb#k8gb_security)) (@ytsarev, @somaritane) (nice talk: [https://www.youtube.com/watch?v=iZpFtalj4xE](https://www.youtube.com/watch?v=iZpFtalj4xE)) (@jkremser - added myself)
- Ask KubeCon organizers on bevvy meeting pre-test (@somaritane)
- Revisit >2 cluster strategy support ([https://github.com/k8gb-io/k8gb/pull/815](https://github.com/k8gb-io/k8gb/pull/815)) (@jkremser) not done for failover yet (comment on issue)
- [https://github.com/k8gb-io/k8gb/issues/642](https://github.com/k8gb-io/k8gb/issues/642) (Azure Support)
    - Fork external-dns and merge [https://github.com/kubernetes-sigs/external-dns/issues/2826](https://github.com/kubernetes-sigs/external-dns/issues/2826) 
Forks done, didn’t update env for that yet
    - Ping Nuno Guedes (@infbase) on the case documentation
- Cert-manager integration issue: [https://blog.abaganon.com/going-global-with-kubernetes-490cf51e2bf8](https://blog.abaganon.com/going-global-with-kubernetes-490cf51e2bf8) 
    - Create associated ticket for TLS and invite Eric for participation. Provide Vault backend as example for mitigating that issue. 
    - Create tutorial on related topic. ([Slack discussion](https://cloud-native.slack.com/archives/C021P656HGB/p1650293408364539)) (@k0da, @jkremser)
- Ping infoblox on v2 release (@somaritane) and check on docker/k3s ways of installation. (there’s flask Python API mock app)
- WRR, CRD plugin PR: (@kuritka to make a presentation on changes)
- Cloudflare support: +1 ask. Need to implement support for NS and TXT in external-dns provider (???) ([https://github.com/k8gb-io/k8gb/issues/884](https://github.com/k8gb-io/k8gb/issues/884)) 
- Update project main page with roadmap and ensure visibility to the community (@somaritane)
- usage of k8gb without edgeDNS: provide our view on the topic (@k0da)




***


</details><details><summary><strong>20 Sept 2022 #16</strong></summary>

## 20 Sept 2022 #16

Attendees

- Yury Tsarev (@ytsarev)
- Timofey Ilinykh (@somaritane)
- Jiri Kremser (@jkremser)
- Michal Kuritka(@kuritka)
- Dinar Valeev (@k0da)

Agenda

- Action Items
- KubeCon NA prep: event structure
- [https://github.com/k8gb-io/k8gb/issues/944#issuecomment-1222096091](https://github.com/k8gb-io/k8gb/issues/944#issuecomment-1222096091) let’s craft the config for Cloudflare and help with the testing
- Project roadmap: [https://github.com/orgs/k8gb-io/projects/2/views/2](https://github.com/orgs/k8gb-io/projects/2/views/2) 
- Yet another Gateway api request [https://cloud-native.slack.com/archives/C021P656HGB/p1663432791213089](https://cloud-native.slack.com/archives/C021P656HGB/p1663432791213089) 
- Presentation from @kuritka on WRR

Notes

- KubeCon NA prep: 
    - Wed 7pm (CEST) (Date & Time: 13:00 - 13:45 (Eastern Time - US & Canada) on Wednesday, October 26, 2022)
    - Duration: 45 mins
    - Bevvy again: need to ask organizers on pre-test
    - Project intro
    - Snapshot on what has changed previously
    - add WRR presentation by @kuritka

Action Items

- K8s crd\_plugin (@kuritka)
- Ask KubeCon organizers on bevvy meeting pre-test (@somaritane)
- Revisit >2 cluster strategy support ([https://github.com/k8gb-io/k8gb/pull/815](https://github.com/k8gb-io/k8gb/pull/815)) (@jkremser) 
- [https://github.com/k8gb-io/k8gb/issues/642](https://github.com/k8gb-io/k8gb/issues/642) (Azure Support)
    - Fork external-dns and merge [https://github.com/kubernetes-sigs/external-dns/issues/2826](https://github.com/kubernetes-sigs/external-dns/issues/2826) 
    - Ping Nuno Guedes (@infbase) on the case documentation
- Cert-manager integration issue: [https://blog.abaganon.com/going-global-with-kubernetes-490cf51e2bf8](https://blog.abaganon.com/going-global-with-kubernetes-490cf51e2bf8) 
    - Create associated ticket for TLS and invite Eric for participation. Provide Vault backend as example for mitigating that issue. 
    - Create tutorial on related topic. ([Slack discussion](https://cloud-native.slack.com/archives/C021P656HGB/p1650293408364539)) (@k0da, @jkremser)
- Ping infoblox on v2 release (@somaritane) and check on docker/k3s ways of installation. (there’s flask Python API mock app)
- WRR, CRD plugin PR: (@kuritka to make a presentation on changes)
- Cloudflare support: +1 ask. Need to implement support for NS and TXT in external-dns provider (???) ([https://github.com/k8gb-io/k8gb/issues/884](https://github.com/k8gb-io/k8gb/issues/884)) 
- Update project main page with roadmap and ensure visibility to the community (@somaritane)
- usage of k8gb without edgeDNS: provide our view on the topic (@k0da)


***


</details><details><summary><strong>6 Sept 2022 #15</strong></summary>

## 6 Sept 2022 #15

Attendees

- Yury Tsarev (@ytsarev)
- Timofey Ilinykh (@somaritane)
- Jiri Kremser (@jkremser)
- Michal Kuritka(@kuritka)

Agenda

- K8gb office hours on KubeCon NA - we need to agree on and book the slot before 9.9. (this Friday)
- [https://github.com/k8gb-io/k8gb/issues/944#issuecomment-1222096091](https://github.com/k8gb-io/k8gb/issues/944#issuecomment-1222096091) let’s craft the config for Cloudflare and help with the testing
- Project roadmap as a new Github Board(CNCF latest review ask)
- [https://github.com/k8gb-io/k8gb/issues/945](https://github.com/k8gb-io/k8gb/issues/945) - usage of k8gb without edgeDNS? Thoughts?
- [https://blog.abaganon.com/going-global-with-kubernetes-490cf51e2bf8](https://blog.abaganon.com/going-global-with-kubernetes-490cf51e2bf8) 
    - *K8gb and TLS - This is probably the main reason why you won’t use K8gb. Provisioning TLS certificates with it becomes a nightmare of a manual process.&#160;&#32;*- how is it solved in Absa? Worth to share the guide?
- Work on k8gb integration with Linode
- external-dns way forward: almost every integration request ends-up with missing NS and TXT record support in external-dns. Need to think about way forward to unblock k8gb.io coverage expansion
- Infoblox v2 upgrade:
[https://github.com/infobloxopen/infoblox-go-client/issues/177](https://github.com/infobloxopen/infoblox-go-client/issues/177) 
[https://github.com/kubernetes-sigs/external-dns/issues/2945#issuecomment-1235732386](https://github.com/kubernetes-sigs/external-dns/issues/2945#issuecomment-1235732386) 
- CI improvements - [https://github.com/k8gb-io/k8gb/pull/949](https://github.com/k8gb-io/k8gb/pull/949)
- Weight Round Robin, CRD plugin PR: [https://github.com/k8gb-io/coredns-crd-plugin/pull/40](https://github.com/k8gb-io/coredns-crd-plugin/pull/40)
- Release stable  [https://github.com/k8gb-io/go-weight-shuffling](https://github.com/k8gb-io/go-weight-shuffling)
- Refactoring SetupManager [https://github.com/k8gb-io/k8gb/issues/932](https://github.com/k8gb-io/k8gb/issues/932)
- Refactor [https://github.com/k8gb-io/k8gb/pull/948](https://github.com/k8gb-io/k8gb/pull/948)
- Mocks test

Notes

- K8gb.io project roadmap: [https://github.com/orgs/k8gb-io/projects/2/views/2](https://github.com/orgs/k8gb-io/projects/2/views/2) 
- [https://blog.abaganon.com/going-global-with-kubernetes-490cf51e2bf8](https://blog.abaganon.com/going-global-with-kubernetes-490cf51e2bf8) work on cert-manager webhook almost there, still waiting for support from external-dns. 
- Work on k8gb integration with Linode: needs support for NS and TXT in external-dns
- external-dns way forward: 
- Become provider owner/maintainer: not an easy way to become provider maintainers as project is under k8s sigs, not an option
- Decision: fork external-dns for k8gb related work on providers with upstream contribution

Action Items

- Cloudflare support: +1 ask. Need to implement support for NS and TXT in external-dns provider (???) ([https://github.com/k8gb-io/k8gb/issues/884](https://github.com/k8gb-io/k8gb/issues/884)) 
- Update project main page with roadmap and ensure visibility to the community (@somaritane)
- usage of k8gb without edgeDNS: provide our view on the topic (@k0da)
- Cert-manager integration issue: [https://blog.abaganon.com/going-global-with-kubernetes-490cf51e2bf8](https://blog.abaganon.com/going-global-with-kubernetes-490cf51e2bf8) Create associated ticket for TLS and invite Eric for participation. Provide Vault backend as example for mitigating that issue. Create tutorial on related topic. ([Slack discussion](https://cloud-native.slack.com/archives/C021P656HGB/p1650293408364539)) (@k0da, @jkremser)
- Ping infoblox on v2 release (@somaritane) and check on docker/k3s ways of installation. (there’s flask Python API mock app)
- WRR, CRD plugin PR: (@kuritka to make a presentation on changes)


***


</details><details><summary><strong>July 26, 2022 #14</strong></summary>

## July 26, 2022 #14

Attendees

- Timofey Ilinykh (@somaritane)
- Dinar Valeev (@k0da)
- Michal Kuritka(@kuritka)
- Nuno Guedes (@infbase)
- Diego Marques (@diego7marques)

Agenda

- Roadmap draft
- Provider onboarding (Azure, linode, …)
- WRR setting 

Notes

- 

Action Items

- 


***


</details><details><summary><strong>July 12, 2022 #13</strong></summary>

## July 12, 2022 #13

Canceled (Project members' unavailability)

***


</details><details><summary><strong>June 28, 2022 #12</strong></summary>

## June 28, 2022 #12

Attendees

- Timofey Ilinykh (@somaritane)
- Jiri Kremser (@jkremser)
- Yury Tsarev(@ytsarev)
- Nuno Guedes (@infbase)
- Diego Marques (@diego7marques)

Agenda

- TOC project review outcomes ([https://youtu.be/VpkYtxzd13E?t=1204](https://youtu.be/VpkYtxzd13E?t=1204)) 
    - [https://github.com/cncf/toc/pull/837#issuecomment-1164576102](https://github.com/cncf/toc/pull/837#issuecomment-1164576102) 
- Tickets triage help needed
- Support for Azure
- [https://github.com/k8gb-io/k8gb/pull/917](https://github.com/k8gb-io/k8gb/pull/917)  Liqo integration

Notes

- 

Action Items

- 


***


</details><details><summary><strong>June 14, 2022 #11</strong></summary>

## June 14, 2022 #11

Attendees

- Timofey Ilinykh (@somaritane)
- Dinar Valeev (@k0da)
- Jiri Kremser (@jkremser)
- Yury Tsarev(@ytsarev)
- Michal Kuritka(@kuritka)

Agenda

- Sandbox review 21th
- Support for Azure
- Should we apply for the k8gb talks at smaller conferences like Kubernetes Days Berlin?
- [https://www.surveymonkey.com/r/CNCF-Maintainers-22-H1](https://www.surveymonkey.com/r/CNCF-Maintainers-22-H1) 
- cluster id ([https://youtu.be/cYFxjZEXucM&t=730s](https://youtu.be/cYFxjZEXucM&t=730s))

Notes

- Sandbox review:
    - No participation is required from our side, and PR for TOC review is already submitted 
- Support for Azure::
    - Looks promising, @k0da is helping with the case
	

Action Items

- Brush out milestone assignment (0.9 -> 0.10): @somaritane
- Contact CNCF support for zoom call time extension: @somaritane
- Add k8gb.io to CNCF meeting calendar
- Fill [https://www.surveymonkey.com/r/CNCF-Maintainers-22-H1](https://www.surveymonkey.com/r/CNCF-Maintainers-22-H1) all :)

***



</details><details><summary><strong>May 31, 2022 #10</strong></summary>

## May 31, 2022 #10

Attendees

- Timofey Ilinykh (@somaritane)
- Dinar Valeev (@k0da)
- Jiri Kremser (@jkremser)
- Michal Kuritka(@kuritka)

Agenda

- News + Ideas from KubeCon 2022 EU
- K8gb + OpenTelemetry demo
- [Michal K](mailto:kuritka@gmail.com) round\_robin issue closed on redundancy
- [Michal K](mailto:kuritka@gmail.com) consistent hashing

Notes

- 

Action Items

- ~~Contact CNCF for k8gb.io office hours on Kubecon 2022: @somaritane~~
- Brush out milestone assignment (0.9 -> 0.10): @somaritane
- ~~Review cncf annual toc review state for k8gb project&#160;~~
- Contact CNCF support for zoom call time extension: @somaritane


***


</details><details><summary><strong>May 17, 2022 #9</strong></summary>

## May 17, 2022 #9

Canceled this one (project members attending KubeCon 2022 EU) 


</details><details><summary><strong>May 3, 2022 #8</strong></summary>

## May 3, 2022 #8

Attendees

- Timofey Ilinykh (@somaritane)
- Dinar Valeev (@k0da)
- Jiri Kremser (@jkremser)
- Yury Tsarev(@ytsarev)
- Michal Kuritka(@kuritka)

Agenda

- [https://github.com/k8gb-io/coredns-crd-plugin/issues/28](https://github.com/k8gb-io/coredns-crd-plugin/issues/28)
- KubeCon
- CNCF annual review (first shot here [https://github.com/jkremser/toc/blob/main/reviews/2022-k8gb-annual.md](https://github.com/jkremser/toc/blob/main/reviews/2022-k8gb-annual.md))
- round-robin for CoreDNS (demo)

Notes

- Ideas for office hours:
    - High-level demo
    - List of milestones,next points
- https://github.com/k8gb-io/k8gb/issues/884
    - Needs additional info from the originator
- 1.0 vs 0.10 release: stick to 0.10

Action Items

- Contact CNCF for k8gb.io office hours on Kubecon 2022: @somaritane
- Brush out milestone assignment (0.9 -> 0.10): @somaritane
- Review cncf annual toc review state for k8gb project (pr [https://github.com/cncf/toc/pull/837](https://github.com/cncf/toc/pull/837))
- Contact CNCF support for zoom call time extension: @somaritane


</details><details><summary><strong>April 19, 2022 #7</strong></summary>

## April 19, 2022 #7

Canceled this one (due to project members availability) 


</details><details><summary><strong>April 5, 2022 #6</strong></summary>

## April 5, 2022 #6

Attendees

- Timofey Ilinykh (@somaritane)
- Dinar Valeev (@k0da)
- Jiri Kremser (@jkremser)
- Yury Tsarev(@ytsarev)
- Michal Kuritka(@kuritka)

Agenda

- Putting/Moving images to another registry (ghcr.io as candidate)
- Become OCI-compatible

Notes

- Jiri: added OCI-compatible labeling for now. Going to dive deeper to produce fully-compatible OCI images.
- Can try incorporate kaniko into pipeline
- Jury: we should check KubeCon schedule for k8gb.io project office hours
- K8gb.io v0.9.0 is released: v1 ingress is now supported
- Michal: working on proper round-robin implementation

Action Items

- Jiri: to update the video section in k8gb.io with new videos
- Check KubeCon 2022 schedule for possible k8gb.io office hours


</details><details><summary><strong>March 22, 2022 #5</strong></summary>

## March 22, 2022 #5

Attendees

- Timofey Ilinykh (@somaritane)
- Dinar Valeev (@k0da)
- Jiri Kremser (@jkremser)
- Yury Tsarev(@ytsarev)
- Michal Kuritka(@kuritka)

Agenda

- Ingress v1 migration ([#847](https://github.com/k8gb-io/k8gb/issues/847))
- Getting rid of Docker Desktop in a dev environment, review and merge [#845](https://github.com/k8gb-io/k8gb/pull/845)
- K8gb in private hosted zones
- Community interest in Slack from OpenShift guys

Notes

- Ingress v1 migration: solve it with major version release and think about proper migration [https://book.kubebuilder.io/multiversion-tutorial/conversion.html](https://book.kubebuilder.io/multiversion-tutorial/conversion.html) 
[https://book.kubebuilder.io/multiversion-tutorial/webhooks.html](https://book.kubebuilder.io/multiversion-tutorial/webhooks.html)
- #845 - the recommended approach is to create the svc as make target vs helm chart, Jiri is about to come with separate PR and we’re about to test
- Dinar: Private zone workaround works (AWS case), with CoreDNS delegation plugin POC
- Dinar: looking for a way to properly implement the DNS glue records. Right now we’re simply creating A records

Action Items

- 


</details><details><summary><strong>March 8, 2022 #4</strong></summary>

## March 8, 2022 #4

Attendees

- Timofey Ilinykh (@somaritane)
- Dinar Valeev (@k0da)
- Jiri Kremser (@jkremser)

Agenda

- Ingress v1 migration ([#847](https://github.com/k8gb-io/k8gb/issues/847))
- Getting rid of Docker Desktop in a dev environment

Notes

- 

Action Items

- Review and Merge [#845](https://github.com/k8gb-io/k8gb/pull/845) 


</details><details><summary><strong>February 22, 2022 #3</strong></summary>

## February 22, 2022 #3

Attendees

- Timofey Ilinykh (@somaritane)
- Dinar Valeev (@k0da)
- Jiri Kremser (@jkremser)

Agenda

- Ingress v1 migration ([#847](https://github.com/k8gb-io/k8gb/issues/847))
- Getting rid of Docker Desktop in a dev environment

Notes

- Getting rid of Docker Desktop:
    - We can enable tcp port for dig, but make it disabled by default in chart values
    - We’re not impacting users with that
- Ingress v1 migration:
    - Suggestion: Sync migration to ingress v1 to next big release (v1) and backport fixes
    - How are we going to upgrade the existing CRs?  (kubectl convert)
    - Food for thought: check the k8gb ingress class
    - kubectl convert can be our friend (https://kubernetes.io/blog/2021/07/14/upcoming-changes-in-kubernetes-1-22/#kubectl-convert)

Action Items

- Review and Merge [#845](https://github.com/k8gb-io/k8gb/pull/845) 















</details><details><summary><strong>February 8, 2022 #2</strong></summary>

## February 8, 2022 #2

Attendees

- Timofey Ilinykh (@somaritane)
- Yury Tsarev (@ytsarev)
- Dinar Valeev (@k0da)
- Jiri Kremser (@jkremser)

Agenda

- The issue with the NS record in AWS Route53 for the private hosted zone
- Groom the GH issues backlog
- Getting rid of Docker Desktop in a dev environment




Notes

-  We need to sort out the usage of GH Project boards vs repo issues lists
- Jiri: Colima has issues with UDP forwarding, the w/a is to use TCP/UDP in tests, plus we have to expose CoreDNS.
- Dinar: k8s doesn’t like when UDP/TCP is exposed on the same port, doesn’t work well with LB
- Question from a customer about k8gb usage in non-loadbalancing case: recommend using ExternalDNS
- Need to investigate means of further promoting k8gb
- Kudos to Jiri for [presentation](https://fosdem.org/2022/schedule/event/container_k8gb_balancer/) on FOSDEM 2022
- Let’s try to use Zoom next time in order to record the video

Action Items

- Create an issue for AWS Route53 private hosted zone (+Azure & GCP)
- Check if two services that expose the pods on the same port, but using different protocol have issue with load balancer (recent change because of colima support - upd not forwarded)







</details><details><summary><strong>January 22, 2022 #1</strong></summary>

## January 22, 2022 #1

Attendees

- Timofey Ilinykh (@somaritane)
- Yury Tsarev (@ytsarev)
- Dinar Valeev (@k0da)
- Jiri Kremser (@jkremser)
- Michal Kuritka (@kuritka)

Agenda

- ingress v1 discussion
- ns records duplicated in the sub-zone or not
- some ideas for our referencial setup
- engagement on the public slack channel

Notes

-  new release 0.8.8 is out \\o/

Action Items

- jkremser: update the PR with colima workaround so that the svc is deployed iif it’s run together with terratests


