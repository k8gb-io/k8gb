# Changelog

## [v0.8.6](https://github.com/k8gb-io/k8gb/tree/v0.8.6) (2021-12-05)

[Full Changelog](https://github.com/k8gb-io/k8gb/compare/v0.8.5...v0.8.6)

**Merged pull requests:**

- RELEASE: v0.8.6 [\#783](https://github.com/k8gb-io/k8gb/pull/783) ([jkremser](https://github.com/jkremser))
- Update Offline Changelog [\#782](https://github.com/k8gb-io/k8gb/pull/782) ([github-actions[bot]](https://github.com/apps/github-actions))
- Update Helm Docs [\#781](https://github.com/k8gb-io/k8gb/pull/781) ([github-actions[bot]](https://github.com/apps/github-actions))
- switch hostAlias to real edgeDNS [\#726](https://github.com/k8gb-io/k8gb/pull/726) ([k0da](https://github.com/k0da))



## [v0.8.5](https://github.com/k8gb-io/k8gb/tree/v0.8.5) (2021-12-01)

[Full Changelog](https://github.com/k8gb-io/k8gb/compare/v0.8.4...v0.8.5)

**Fixed bugs:**

- dev: make command errors [\#770](https://github.com/k8gb-io/k8gb/issues/770)
- deploy-test-apps and deploy-full-local-setup make commands are broken [\#751](https://github.com/k8gb-io/k8gb/issues/751)

**Closed issues:**

- k8gb controller's pid fails to start correctly when deployed by OLM [\#757](https://github.com/k8gb-io/k8gb/issues/757)
- Use pod's dnsConfig for our demo [\#712](https://github.com/k8gb-io/k8gb/issues/712)
- How does this compare to others? [\#689](https://github.com/k8gb-io/k8gb/issues/689)

**Merged pull requests:**

- RELEASE: v0.8.5 [\#780](https://github.com/k8gb-io/k8gb/pull/780) ([jkremser](https://github.com/jkremser))
- Revert "RELEASE: v0.8.5 \(\#776\)" [\#779](https://github.com/k8gb-io/k8gb/pull/779) ([jkremser](https://github.com/jkremser))
- .gitignoring file called 'changes' that's produced and consumed by goreleaser [\#778](https://github.com/k8gb-io/k8gb/pull/778) ([jkremser](https://github.com/jkremser))
- RELEASE: v0.8.5 [\#776](https://github.com/k8gb-io/k8gb/pull/776) ([jkremser](https://github.com/jkremser))
- Reorder terratest Workflow [\#774](https://github.com/k8gb-io/k8gb/pull/774) ([k0da](https://github.com/k0da))
- dev: Update for k3d v5.1.0 | k3d-action@v2 [\#773](https://github.com/k8gb-io/k8gb/pull/773) ([somaritane](https://github.com/somaritane))
- Don't evaluate COREDNS IP too early [\#772](https://github.com/k8gb-io/k8gb/pull/772) ([k0da](https://github.com/k0da))
- Running all the tests on two clusters and only full-rr on 3 clusters [\#769](https://github.com/k8gb-io/k8gb/pull/769) ([jkremser](https://github.com/jkremser))
- Makefile: surround vars with quotes in conditional expressions [\#768](https://github.com/k8gb-io/k8gb/pull/768) ([jkremser](https://github.com/jkremser))
- Invoke the OLM pipeline from release pipeline [\#767](https://github.com/k8gb-io/k8gb/pull/767) ([jkremser](https://github.com/jkremser))
- Do not deploy test-gslb for terratest runs [\#766](https://github.com/k8gb-io/k8gb/pull/766) ([jkremser](https://github.com/jkremser))
- Drop resolv.conf hack [\#764](https://github.com/k8gb-io/k8gb/pull/764) ([k0da](https://github.com/k0da))
- Fix two typos in Makefile [\#762](https://github.com/k8gb-io/k8gb/pull/762) ([jkremser](https://github.com/jkremser))
- LICENSE & README are required by artifacthub.io, so un-.helmignoring [\#761](https://github.com/k8gb-io/k8gb/pull/761) ([jkremser](https://github.com/jkremser))
- Fix helm docs for hostAliases entries [\#760](https://github.com/k8gb-io/k8gb/pull/760) ([jkremser](https://github.com/jkremser))
- Update Helm Docs [\#759](https://github.com/k8gb-io/k8gb/pull/759) ([github-actions[bot]](https://github.com/apps/github-actions))
- Fix 'cannot verify user is non-root in OLM' [\#758](https://github.com/k8gb-io/k8gb/pull/758) ([jkremser](https://github.com/jkremser))
- Fixing vulnerabilities in terratests [\#756](https://github.com/k8gb-io/k8gb/pull/756) ([kuritka](https://github.com/kuritka))
- Add missing displayName + remove existing dir [\#755](https://github.com/k8gb-io/k8gb/pull/755) ([jkremser](https://github.com/jkremser))
- In CONTRIBUTING.md point to a release pr that went smoothly [\#754](https://github.com/k8gb-io/k8gb/pull/754) ([jkremser](https://github.com/jkremser))
- Switch to CR\_TOKEN with the workflow scope [\#753](https://github.com/k8gb-io/k8gb/pull/753) ([ytsarev](https://github.com/ytsarev))
- fix: Broken first user experience commands [\#752](https://github.com/k8gb-io/k8gb/pull/752) ([somaritane](https://github.com/somaritane))
- Allow the EXT\_GSLB\_CLUSTERS\_GEO\_TAGS to contain CLUSTER\_GEO\_TAG [\#750](https://github.com/k8gb-io/k8gb/pull/750) ([jkremser](https://github.com/jkremser))
- cleaning .gitignore [\#740](https://github.com/k8gb-io/k8gb/pull/740) ([kuritka](https://github.com/kuritka))
- Update Offline Changelog [\#736](https://github.com/k8gb-io/k8gb/pull/736) ([github-actions[bot]](https://github.com/apps/github-actions))
- Multicluster setup \(n \> 2\) [\#722](https://github.com/k8gb-io/k8gb/pull/722) ([jkremser](https://github.com/jkremser))



## [v0.8.4](https://github.com/k8gb-io/k8gb/tree/v0.8.4) (2021-11-16)

[Full Changelog](https://github.com/k8gb-io/k8gb/compare/v0.8.3...v0.8.4)

**Implemented enhancements:**

- Consolidate goreleaser builds and local docker builds [\#588](https://github.com/k8gb-io/k8gb/issues/588)
- Create separate README.md file for k8gb helm chart [\#359](https://github.com/k8gb-io/k8gb/issues/359)
- consider DNS provider config consolidation [\#219](https://github.com/k8gb-io/k8gb/issues/219)

**Closed issues:**

- Get rid of `EXTERNAL_DNS_` prefix at Infoblox ENV vars [\#683](https://github.com/k8gb-io/k8gb/issues/683)
- Gh action for pushing the up-to-date image of k8gb-demo-curl [\#651](https://github.com/k8gb-io/k8gb/issues/651)
- Fix the vulnerabilities reported by Artifacthub - part 2 [\#637](https://github.com/k8gb-io/k8gb/issues/637)
- Fix the vulnerabilities reported by Artifacthub - part 1 [\#636](https://github.com/k8gb-io/k8gb/issues/636)

**Merged pull requests:**

- RELEASE: v0.8.4 [\#749](https://github.com/k8gb-io/k8gb/pull/749) ([ytsarev](https://github.com/ytsarev))
- Revert "RELEASE: v0.8.4" [\#747](https://github.com/k8gb-io/k8gb/pull/747) ([ytsarev](https://github.com/ytsarev))
- Extend cut\_release GHA workflow with release revert mechanism [\#745](https://github.com/k8gb-io/k8gb/pull/745) ([ytsarev](https://github.com/ytsarev))
- curldemo: Build the container image only for pushes to master [\#744](https://github.com/k8gb-io/k8gb/pull/744) ([jkremser](https://github.com/jkremser))
- Remove only the first two lines from old changelog when appending [\#742](https://github.com/k8gb-io/k8gb/pull/742) ([jkremser](https://github.com/jkremser))
- Fix onlyLastTag in change log generator [\#741](https://github.com/k8gb-io/k8gb/pull/741) ([k0da](https://github.com/k0da))
- Deploy test-gslb namespace only in local cluster [\#738](https://github.com/k8gb-io/k8gb/pull/738) ([k0da](https://github.com/k0da))
- curldemo: Use full path to Dockerfile [\#737](https://github.com/k8gb-io/k8gb/pull/737) ([jkremser](https://github.com/jkremser))
- Update Helm Docs [\#735](https://github.com/k8gb-io/k8gb/pull/735) ([github-actions[bot]](https://github.com/apps/github-actions))
- RELEASE: v0.8.4 [\#734](https://github.com/k8gb-io/k8gb/pull/734) ([ytsarev](https://github.com/ytsarev))
- Revert "RELEASE: v0.8.4 " [\#733](https://github.com/k8gb-io/k8gb/pull/733) ([ytsarev](https://github.com/ytsarev))
- Update docs about automatic tag during release process + different token [\#732](https://github.com/k8gb-io/k8gb/pull/732) ([jkremser](https://github.com/jkremser))
- Time out and reconcile  aggressively. [\#730](https://github.com/k8gb-io/k8gb/pull/730) ([k0da](https://github.com/k0da))
- Create new tag when version in Chart.yaml is bumped [\#729](https://github.com/k8gb-io/k8gb/pull/729) ([jkremser](https://github.com/jkremser))
- Don't deploy local path storage provisioner [\#728](https://github.com/k8gb-io/k8gb/pull/728) ([k0da](https://github.com/k0da))
- Remove K8GB\_VERSION env var [\#727](https://github.com/k8gb-io/k8gb/pull/727) ([kuritka](https://github.com/kuritka))
- Switch version to the latest tag [\#725](https://github.com/k8gb-io/k8gb/pull/725) ([k0da](https://github.com/k0da))
- RELEASE: v0.8.4  [\#724](https://github.com/k8gb-io/k8gb/pull/724) ([ytsarev](https://github.com/ytsarev))
- Don't keep adding localtargets- prefix for each ext geo tag [\#719](https://github.com/k8gb-io/k8gb/pull/719) ([jkremser](https://github.com/jkremser))
- Refactor external DNS to be provider neutral [\#718](https://github.com/k8gb-io/k8gb/pull/718) ([k0da](https://github.com/k0da))
- Update AWS/Route53 reference setup terraform code [\#717](https://github.com/k8gb-io/k8gb/pull/717) ([ytsarev](https://github.com/ytsarev))
- Fix upgrade-candidate [\#716](https://github.com/k8gb-io/k8gb/pull/716) ([kuritka](https://github.com/kuritka))
- open metrics port for k8gb container [\#714](https://github.com/k8gb-io/k8gb/pull/714) ([k0da](https://github.com/k0da))
- olm: fix paths in the resulting PR [\#711](https://github.com/k8gb-io/k8gb/pull/711) ([jkremser](https://github.com/jkremser))
- Upgrade the controller-gen dependency [\#710](https://github.com/k8gb-io/k8gb/pull/710) ([kuritka](https://github.com/kuritka))
- Push image to docker hub \#651 [\#709](https://github.com/k8gb-io/k8gb/pull/709) ([AugustasV](https://github.com/AugustasV))
- Add KubeCon NA 2021 recording [\#708](https://github.com/k8gb-io/k8gb/pull/708) ([ytsarev](https://github.com/ytsarev))
- yet another olm CI fixes [\#707](https://github.com/k8gb-io/k8gb/pull/707) ([jkremser](https://github.com/jkremser))
- helm: translate the edgeDNSServer -\> edgeDNSServers using a tempate [\#706](https://github.com/k8gb-io/k8gb/pull/706) ([jkremser](https://github.com/jkremser))
- olm: Add installModes and change crd name to fully-qualified one [\#705](https://github.com/k8gb-io/k8gb/pull/705) ([jkremser](https://github.com/jkremser))
- \[FIX\] v0.8.3 panicking error [\#704](https://github.com/k8gb-io/k8gb/pull/704) ([kuritka](https://github.com/kuritka))
- Align exact version in release.yaml [\#703](https://github.com/k8gb-io/k8gb/pull/703) ([kuritka](https://github.com/kuritka))
- Drop socrecard OLM annotation [\#702](https://github.com/k8gb-io/k8gb/pull/702) ([k0da](https://github.com/k0da))
- olm: Add missing annotations + typo in repository [\#701](https://github.com/k8gb-io/k8gb/pull/701) ([jkremser](https://github.com/jkremser))
- Bump dependencies \#1 [\#700](https://github.com/k8gb-io/k8gb/pull/700) ([kuritka](https://github.com/kuritka))
- olm: Fix paths when using version=master + more descriptive PR body [\#699](https://github.com/k8gb-io/k8gb/pull/699) ([jkremser](https://github.com/jkremser))
- Set GO1.16.9 as default GO version for goreleaser [\#698](https://github.com/k8gb-io/k8gb/pull/698) ([kuritka](https://github.com/kuritka))
- Using non-default token for opening the PR [\#696](https://github.com/k8gb-io/k8gb/pull/696) ([jkremser](https://github.com/jkremser))
- typo: GH\_TOKEN -\> GITHUB\_TOKEN [\#695](https://github.com/k8gb-io/k8gb/pull/695) ([jkremser](https://github.com/jkremser))
- Add a way to use non-released changes in repo for producing an olm bundle [\#694](https://github.com/k8gb-io/k8gb/pull/694) ([jkremser](https://github.com/jkremser))
- Update dev documentation for custom SSL cert support \(\#687\) [\#691](https://github.com/k8gb-io/k8gb/pull/691) ([somaritane](https://github.com/somaritane))
- Mass AbsaOSS -\> k8gb-io in code and docs [\#690](https://github.com/k8gb-io/k8gb/pull/690) ([ytsarev](https://github.com/ytsarev))
- Adding large logo to README.md [\#686](https://github.com/k8gb-io/k8gb/pull/686) ([kuritka](https://github.com/kuritka))
- Update Helm Docs [\#684](https://github.com/k8gb-io/k8gb/pull/684) ([github-actions[bot]](https://github.com/apps/github-actions))
- Fix url and add maintainers to Chart.yaml [\#682](https://github.com/k8gb-io/k8gb/pull/682) ([jkremser](https://github.com/jkremser))
- Use k8gb.io png icon asset as helm chart icon [\#681](https://github.com/k8gb-io/k8gb/pull/681) ([somaritane](https://github.com/somaritane))
- Update Helm Docs [\#679](https://github.com/k8gb-io/k8gb/pull/679) ([github-actions[bot]](https://github.com/apps/github-actions))
- olm-bundle: generating the ClusterServiceVerion file + gh action [\#678](https://github.com/k8gb-io/k8gb/pull/678) ([jkremser](https://github.com/jkremser))
- Update ignore command for netlify [\#677](https://github.com/k8gb-io/k8gb/pull/677) ([jkremser](https://github.com/jkremser))
- Bind env variables with ENV-BINDER, clearing tests [\#676](https://github.com/k8gb-io/k8gb/pull/676) ([kuritka](https://github.com/kuritka))
- Revert "Use PNG image as chart icon \(\#671\)" [\#675](https://github.com/k8gb-io/k8gb/pull/675) ([k0da](https://github.com/k0da))
- Generate chart README.md with helm-docs [\#674](https://github.com/k8gb-io/k8gb/pull/674) ([k0da](https://github.com/k0da))
- Move local Prometheus port exposure to less standard ports [\#673](https://github.com/k8gb-io/k8gb/pull/673) ([ytsarev](https://github.com/ytsarev))
- Drop svc migration related values [\#672](https://github.com/k8gb-io/k8gb/pull/672) ([k0da](https://github.com/k0da))
- Use PNG image as chart icon [\#671](https://github.com/k8gb-io/k8gb/pull/671) ([k0da](https://github.com/k0da))
- Switch docker build to be built by goreleaser [\#670](https://github.com/k8gb-io/k8gb/pull/670) ([k0da](https://github.com/k0da))
- Update Offline Changelog [\#669](https://github.com/k8gb-io/k8gb/pull/669) ([github-actions[bot]](https://github.com/apps/github-actions))
- Set onlyLastTag to true to prevent rate limiting issues w/ gh api [\#668](https://github.com/k8gb-io/k8gb/pull/668) ([jkremser](https://github.com/jkremser))
- Add icon URL to k8gb helm chart [\#667](https://github.com/k8gb-io/k8gb/pull/667) ([somaritane](https://github.com/somaritane))
- web-preview: list the files to checkout w/o using the {foo,bar}.md syntax [\#665](https://github.com/k8gb-io/k8gb/pull/665) ([jkremser](https://github.com/jkremser))
- Add automatic deployment preview of PRs when changing the site [\#660](https://github.com/k8gb-io/k8gb/pull/660) ([jkremser](https://github.com/jkremser))



## [v0.8.3](https://github.com/k8gb-io/k8gb/tree/v0.8.3) (2021-10-19)

[Full Changelog](https://github.com/k8gb-io/k8gb/compare/v0.8.2...v0.8.3)

**Implemented enhancements:**

- Helm chart produces empty lines in yamls [\#631](https://github.com/k8gb-io/k8gb/issues/631)
- GoKART integration [\#600](https://github.com/k8gb-io/k8gb/issues/600)
- Provide K8gb without IRSA  Annotation capability [\#593](https://github.com/k8gb-io/k8gb/issues/593)
- Fix terrascan detected issues and make the associated GHA workflow blocking [\#585](https://github.com/k8gb-io/k8gb/issues/585)
- Use separate GH token for changelog\_generator job [\#581](https://github.com/k8gb-io/k8gb/issues/581)
- Create Best Practices Badge and pass [\#479](https://github.com/k8gb-io/k8gb/issues/479)
- Align k8gb logging statements with zerolog best practices [\#469](https://github.com/k8gb-io/k8gb/issues/469)
- core dns load balancer does not get the correct targetPort [\#423](https://github.com/k8gb-io/k8gb/issues/423)
- Expose failover status in prometheus endpoint [\#221](https://github.com/k8gb-io/k8gb/issues/221)
- Detect and log NS servers A records propagation clash between two or more distinct k8gb pairs [\#165](https://github.com/k8gb-io/k8gb/issues/165)
- Opportunity to enhance edgeDNSServer from single instance to multiple [\#154](https://github.com/k8gb-io/k8gb/issues/154)
- Expose advanced metrics [\#124](https://github.com/k8gb-io/k8gb/issues/124)
- Make securityContext of the deployment fully configurable [\#633](https://github.com/k8gb-io/k8gb/pull/633) ([jkremser](https://github.com/jkremser))

**Fixed bugs:**

- podinfo readiness timeout [\#595](https://github.com/k8gb-io/k8gb/issues/595)

**Closed issues:**

- Add Trivy gh action to our workflow [\#638](https://github.com/k8gb-io/k8gb/issues/638)
- k8gb-coredns Pod CrashLoopBack on OpenShift :: listen tcp :53: bind: permission denied [\#623](https://github.com/k8gb-io/k8gb/issues/623)
- Bump GOLANGCI-LINT [\#609](https://github.com/k8gb-io/k8gb/issues/609)
- GH actions should be run against the pull request coming from the outside of our bubble [\#607](https://github.com/k8gb-io/k8gb/issues/607)
- Update diagrams in a Concepts documentation section [\#598](https://github.com/k8gb-io/k8gb/issues/598)
- Remove zone cleanup code [\#547](https://github.com/k8gb-io/k8gb/issues/547)
- Add SECURITY.md and security disclosure process [\#478](https://github.com/k8gb-io/k8gb/issues/478)
- Developer contribution guide [\#100](https://github.com/k8gb-io/k8gb/issues/100)

**Merged pull requests:**

- Fix ldflags argument for goreleaser [\#663](https://github.com/k8gb-io/k8gb/pull/663) ([jkremser](https://github.com/jkremser))
- Bump chart version to 0.8.3 - prepare for release [\#662](https://github.com/k8gb-io/k8gb/pull/662) ([jkremser](https://github.com/jkremser))
- Remove non-used png [\#661](https://github.com/k8gb-io/k8gb/pull/661) ([ytsarev](https://github.com/ytsarev))
- Add \#k8gb slack channel link [\#659](https://github.com/k8gb-io/k8gb/pull/659) ([ytsarev](https://github.com/ytsarev))
- Don't install extdns rbac by default [\#656](https://github.com/k8gb-io/k8gb/pull/656) ([k0da](https://github.com/k0da))
- Expose grafana as NodePort service and open it on k3d [\#655](https://github.com/k8gb-io/k8gb/pull/655) ([jkremser](https://github.com/jkremser))
- Bump coredns with plugin to v0.0.7 [\#654](https://github.com/k8gb-io/k8gb/pull/654) ([somaritane](https://github.com/somaritane))
- Add grafana including example dashboard for 'podinfoes' [\#653](https://github.com/k8gb-io/k8gb/pull/653) ([jkremser](https://github.com/jkremser))
- Spice up README headers a bit [\#652](https://github.com/k8gb-io/k8gb/pull/652) ([ytsarev](https://github.com/ytsarev))
- Fix architecture diagram [\#650](https://github.com/k8gb-io/k8gb/pull/650) ([ytsarev](https://github.com/ytsarev))
- Add parameter denoting how long to wait for k8gbcurl.sh demo script [\#649](https://github.com/k8gb-io/k8gb/pull/649) ([jkremser](https://github.com/jkremser))
- Minor: makefile help indentation [\#648](https://github.com/k8gb-io/k8gb/pull/648) ([jkremser](https://github.com/jkremser))
- Fix terratest CVE-2020-10675 [\#647](https://github.com/k8gb-io/k8gb/pull/647) ([kuritka](https://github.com/kuritka))
- Fix terratest CVE-2021-41103, CVE-2020-27813, CVE-2020-26160 [\#646](https://github.com/k8gb-io/k8gb/pull/646) ([kuritka](https://github.com/kuritka))
- Extend README.md with RedHat link explaining Global Load Balancing [\#644](https://github.com/k8gb-io/k8gb/pull/644) ([kuritka](https://github.com/kuritka))
- github.com/containerd/containerd v1.4.11 [\#640](https://github.com/k8gb-io/k8gb/pull/640) ([kuritka](https://github.com/kuritka))
- dependabot github.com/containerd/containerd voulnerability [\#639](https://github.com/k8gb-io/k8gb/pull/639) ([kuritka](https://github.com/kuritka))
- DNS package test coverage \(3/3\) [\#635](https://github.com/k8gb-io/k8gb/pull/635) ([kuritka](https://github.com/kuritka))
- Remove FAKE\_INFOBLOX \(2/3\) [\#634](https://github.com/k8gb-io/k8gb/pull/634) ([kuritka](https://github.com/kuritka))
- Fix \#631: helm - remove new lines from resulting yaml when using conditionals [\#632](https://github.com/k8gb-io/k8gb/pull/632) ([jkremser](https://github.com/jkremser))
- Mock DNS package, extend testing \(1/3\) [\#630](https://github.com/k8gb-io/k8gb/pull/630) ([kuritka](https://github.com/kuritka))
- Address issues found by terrascan and make it blocking [\#628](https://github.com/k8gb-io/k8gb/pull/628) ([k0da](https://github.com/k0da))
- Drop cleanup code [\#627](https://github.com/k8gb-io/k8gb/pull/627) ([k0da](https://github.com/k0da))
- Run coredns on unpriveleged port [\#626](https://github.com/k8gb-io/k8gb/pull/626) ([k0da](https://github.com/k0da))
- Extend metrics.md by metrics description [\#625](https://github.com/k8gb-io/k8gb/pull/625) ([kuritka](https://github.com/kuritka))
- k8gb\_gslb\_reconciliation\_loops\_total per GSLB [\#624](https://github.com/k8gb-io/k8gb/pull/624) ([kuritka](https://github.com/kuritka))
- k8gb\_runtime\_info [\#622](https://github.com/k8gb-io/k8gb/pull/622) ([kuritka](https://github.com/kuritka))
- K8gbEndpointStatus [\#620](https://github.com/k8gb-io/k8gb/pull/620) ([kuritka](https://github.com/kuritka))
- Link strategy doc on the index page [\#619](https://github.com/k8gb-io/k8gb/pull/619) ([ytsarev](https://github.com/ytsarev))
- Enable path filtering for terrascan [\#618](https://github.com/k8gb-io/k8gb/pull/618) ([ytsarev](https://github.com/ytsarev))
- Optimize kubelinter pipeline config [\#617](https://github.com/k8gb-io/k8gb/pull/617) ([ytsarev](https://github.com/ytsarev))
- Move IP's to constant [\#616](https://github.com/k8gb-io/k8gb/pull/616) ([kuritka](https://github.com/kuritka))
- Fix helm linting error for coredns.serviceType [\#615](https://github.com/k8gb-io/k8gb/pull/615) ([somaritane](https://github.com/somaritane))
- Ability to disable IRSA role association in route53 scenario [\#614](https://github.com/k8gb-io/k8gb/pull/614) ([ytsarev](https://github.com/ytsarev))
- Upload terrascan SARIF file [\#613](https://github.com/k8gb-io/k8gb/pull/613) ([ytsarev](https://github.com/ytsarev))
- Remove accidental newline in recordings table [\#612](https://github.com/k8gb-io/k8gb/pull/612) ([ytsarev](https://github.com/ytsarev))
- Add NS1 INS1GHTS recording [\#611](https://github.com/k8gb-io/k8gb/pull/611) ([ytsarev](https://github.com/ytsarev))
- Bump golangci-lint [\#610](https://github.com/k8gb-io/k8gb/pull/610) ([kuritka](https://github.com/kuritka))
- Run all the static analysis tools, tests, etc. against the pull requests [\#608](https://github.com/k8gb-io/k8gb/pull/608) ([jkremser](https://github.com/jkremser))
- s/edgeDNSServer/edgeDNSServers/g [\#605](https://github.com/k8gb-io/k8gb/pull/605) ([jkremser](https://github.com/jkremser))
- Gokart action [\#604](https://github.com/k8gb-io/k8gb/pull/604) ([kuritka](https://github.com/kuritka))
- Local GoKart [\#603](https://github.com/k8gb-io/k8gb/pull/603) ([kuritka](https://github.com/kuritka))
- Update k8gb design diagrams with clear k8gb controller location [\#602](https://github.com/k8gb-io/k8gb/pull/602) ([ytsarev](https://github.com/ytsarev))
- Update ancient arch statement [\#601](https://github.com/k8gb-io/k8gb/pull/601) ([ytsarev](https://github.com/ytsarev))
- Fix api version for RBAC to not produce a warning [\#599](https://github.com/k8gb-io/k8gb/pull/599) ([jkremser](https://github.com/jkremser))
- contextual logging [\#597](https://github.com/k8gb-io/k8gb/pull/597) ([kuritka](https://github.com/kuritka))
- Fix typo in readme [\#596](https://github.com/k8gb-io/k8gb/pull/596) ([jkremser](https://github.com/jkremser))
- Add @jkremser to CODEOWNERS [\#592](https://github.com/k8gb-io/k8gb/pull/592) ([kuritka](https://github.com/kuritka))
- \[docs\] Fixing couple of typos [\#591](https://github.com/k8gb-io/k8gb/pull/591) ([jkremser](https://github.com/jkremser))
- align go1.16 [\#590](https://github.com/k8gb-io/k8gb/pull/590) ([kuritka](https://github.com/kuritka))
- Release build optimisations [\#589](https://github.com/k8gb-io/k8gb/pull/589) ([kuritka](https://github.com/kuritka))
- Added standard vulnerability response time [\#584](https://github.com/k8gb-io/k8gb/pull/584) ([somaritane](https://github.com/somaritane))
- Extend release process in CONTRIBUTING.md [\#583](https://github.com/k8gb-io/k8gb/pull/583) ([kuritka](https://github.com/kuritka))
- Use CR\_TOKEN secret for changelog generator job [\#582](https://github.com/k8gb-io/k8gb/pull/582) ([k0da](https://github.com/k0da))
- Update Offline Changelog [\#580](https://github.com/k8gb-io/k8gb/pull/580) ([github-actions[bot]](https://github.com/apps/github-actions))
- Initial security policy [\#576](https://github.com/k8gb-io/k8gb/pull/576) ([somaritane](https://github.com/somaritane))
- Add Terrascan GHA workflow [\#574](https://github.com/k8gb-io/k8gb/pull/574) ([ytsarev](https://github.com/ytsarev))

## [v0.8.2](https://github.com/k8gb-io/k8gb/tree/v0.8.2) (2021-08-25)

[Full Changelog](https://github.com/k8gb-io/k8gb/compare/v0.8.1...v0.8.2)

**Implemented enhancements:**

- Upgrade to latest operator-sdk [\#527](https://github.com/k8gb-io/k8gb/issues/527)
- Add favicon to k8gb.io website [\#498](https://github.com/k8gb-io/k8gb/issues/498)
- Add test coverage requirements to contribution flow in CONTRIBUTING.md [\#497](https://github.com/k8gb-io/k8gb/issues/497)
- CodeQL GH workflow should be scoped only for source code changes [\#482](https://github.com/k8gb-io/k8gb/issues/482)
- \[openshift-support\] k8gb does not have the permissions to set the ingress name [\#422](https://github.com/k8gb-io/k8gb/issues/422)
- \[openshift-support\] runAsUser=1000 preventing from running pods [\#421](https://github.com/k8gb-io/k8gb/issues/421)
- Add topology/location based load balancing strategy [\#244](https://github.com/k8gb-io/k8gb/issues/244)
- Try to mitigate DNS protocol limitations with ingress controller custom error [\#134](https://github.com/k8gb-io/k8gb/issues/134)

**Fixed bugs:**

- 'test-round-robin' often fails [\#528](https://github.com/k8gb-io/k8gb/issues/528)

**Closed issues:**

- Updates components diagram [\#569](https://github.com/k8gb-io/k8gb/issues/569)
- Document breifly new geoip strategy and other strategies in general [\#548](https://github.com/k8gb-io/k8gb/issues/548)
- Can't deploy K8GB in a Cluster that has an Existing ExternalDNS Deployment [\#542](https://github.com/k8gb-io/k8gb/issues/542)
- Fix security vulnerability in golang.org/x/crypto [\#539](https://github.com/k8gb-io/k8gb/issues/539)
- Extend terratest suite with http end-to-end test. [\#533](https://github.com/k8gb-io/k8gb/issues/533)
- Enable DCO for all of k8gb \(CNCF requirement\) [\#523](https://github.com/k8gb-io/k8gb/issues/523)
- Update CONTRIBUTING.md with logging recommendations [\#468](https://github.com/k8gb-io/k8gb/issues/468)

**Merged pull requests:**

- Sign-off changelog PR [\#579](https://github.com/k8gb-io/k8gb/pull/579) ([k0da](https://github.com/k0da))
- fix helm publish [\#578](https://github.com/k8gb-io/k8gb/pull/578) ([kuritka](https://github.com/kuritka))
- release v0.8.2 [\#575](https://github.com/k8gb-io/k8gb/pull/575) ([kuritka](https://github.com/kuritka))
- Update contribution flow with code style and logging recommendations [\#573](https://github.com/k8gb-io/k8gb/pull/573) ([somaritane](https://github.com/somaritane))
- Metrics \(4/4\) [\#572](https://github.com/k8gb-io/k8gb/pull/572) ([kuritka](https://github.com/kuritka))
- Mention testing in the contribution flow [\#571](https://github.com/k8gb-io/k8gb/pull/571) ([somaritane](https://github.com/somaritane))
- Fix k8gb-components.svg image [\#570](https://github.com/k8gb-io/k8gb/pull/570) ([k0da](https://github.com/k0da))
- Ignore gh-pages | Jekyll generated output [\#568](https://github.com/k8gb-io/k8gb/pull/568) ([somaritane](https://github.com/somaritane))
- Health status enum [\#564](https://github.com/k8gb-io/k8gb/pull/564) ([kuritka](https://github.com/kuritka))
- Fix external-dns managed records option usage [\#563](https://github.com/k8gb-io/k8gb/pull/563) ([k0da](https://github.com/k0da))
- Fix external-dns securityContext identation [\#562](https://github.com/k8gb-io/k8gb/pull/562) ([k0da](https://github.com/k0da))
- Propogate DNS Zone Negative TTL down to CoreDNS [\#561](https://github.com/k8gb-io/k8gb/pull/561) ([k0da](https://github.com/k0da))
- Make CodeQL workflow to react to Go files change only [\#560](https://github.com/k8gb-io/k8gb/pull/560) ([ytsarev](https://github.com/ytsarev))
- Parametrize security settings for k8gb and externaldns pods [\#559](https://github.com/k8gb-io/k8gb/pull/559) ([ytsarev](https://github.com/ytsarev))
- Openshift support: flagged rbac for the Routes [\#558](https://github.com/k8gb-io/k8gb/pull/558) ([ytsarev](https://github.com/ytsarev))
- create-pull-request creates commit [\#557](https://github.com/k8gb-io/k8gb/pull/557) ([k0da](https://github.com/k0da))
- Add initial strategies document [\#556](https://github.com/k8gb-io/k8gb/pull/556) ([k0da](https://github.com/k0da))
- \[Fix\] busybox HitTestApp [\#555](https://github.com/k8gb-io/k8gb/pull/555) ([kuritka](https://github.com/kuritka))
- Upgrade to latest operator-sdk v1.10.1 [\#554](https://github.com/k8gb-io/k8gb/pull/554) ([ytsarev](https://github.com/ytsarev))
- Update to external-dns v0.9.0 [\#553](https://github.com/k8gb-io/k8gb/pull/553) ([k0da](https://github.com/k0da))
- Use `k8gb` prefix for external dns rbac [\#551](https://github.com/k8gb-io/k8gb/pull/551) ([ytsarev](https://github.com/ytsarev))
- Metrics \(3/4\) [\#550](https://github.com/k8gb-io/k8gb/pull/550) ([kuritka](https://github.com/kuritka))
- http failover, wait for app [\#546](https://github.com/k8gb-io/k8gb/pull/546) ([kuritka](https://github.com/kuritka))
- Handling error code for defer functions [\#545](https://github.com/k8gb-io/k8gb/pull/545) ([kuritka](https://github.com/kuritka))
- logging delegated records [\#544](https://github.com/k8gb-io/k8gb/pull/544) ([kuritka](https://github.com/kuritka))
- Cleanup old NS name format [\#543](https://github.com/k8gb-io/k8gb/pull/543) ([k0da](https://github.com/k0da))
- upgrade to k3d-action v1.5.0 [\#541](https://github.com/k8gb-io/k8gb/pull/541) ([kuritka](https://github.com/kuritka))
- Fix x/crypto vulnerability [\#540](https://github.com/k8gb-io/k8gb/pull/540) ([ytsarev](https://github.com/ytsarev))
- Failover HTTP test [\#538](https://github.com/k8gb-io/k8gb/pull/538) ([kuritka](https://github.com/kuritka))
- Temporarly downgrade GitHub runner to Ubuntu 18.04 [\#537](https://github.com/k8gb-io/k8gb/pull/537) ([kuritka](https://github.com/kuritka))
- Trailing whitespace busting [\#536](https://github.com/k8gb-io/k8gb/pull/536) ([ytsarev](https://github.com/ytsarev))
- simplifying Failover and RoundRobin tests [\#535](https://github.com/k8gb-io/k8gb/pull/535) ([kuritka](https://github.com/kuritka))
- Fix local playground setup [\#532](https://github.com/k8gb-io/k8gb/pull/532) ([k0da](https://github.com/k0da))
- Fix version of podinfo test sample chart [\#531](https://github.com/k8gb-io/k8gb/pull/531) ([ytsarev](https://github.com/ytsarev))
- Label dnsendpoints with strategy label [\#530](https://github.com/k8gb-io/k8gb/pull/530) ([k0da](https://github.com/k0da))
- Install prometheus on local clusters \(2/4\) [\#529](https://github.com/k8gb-io/k8gb/pull/529) ([kuritka](https://github.com/kuritka))
- Metrics configuration \#\(1/4\) [\#525](https://github.com/k8gb-io/k8gb/pull/525) ([kuritka](https://github.com/kuritka))
- Revert "Scan k8gb image by Artifacthub.io \(\#519\)" [\#524](https://github.com/k8gb-io/k8gb/pull/524) ([k0da](https://github.com/k0da))
- Document release process [\#522](https://github.com/k8gb-io/k8gb/pull/522) ([ytsarev](https://github.com/ytsarev))
- terratest abstraction [\#514](https://github.com/k8gb-io/k8gb/pull/514) ([kuritka](https://github.com/kuritka))
- Refactor CoreDNS service [\#453](https://github.com/k8gb-io/k8gb/pull/453) ([k0da](https://github.com/k0da))

## [v0.8.1](https://github.com/k8gb-io/k8gb/tree/v0.8.1) (2021-06-14)

[Full Changelog](https://github.com/k8gb-io/k8gb/compare/v0.8.0...v0.8.1)

**Implemented enhancements:**

- edgeDNSServer is not used for remote k8gb NS record resolution [\#513](https://github.com/k8gb-io/k8gb/issues/513)
- Containerized local authoring environment for k8gb.io website [\#483](https://github.com/k8gb-io/k8gb/issues/483)
- Shorten NS names for zone delegation [\#456](https://github.com/k8gb-io/k8gb/issues/456)

**Fixed bugs:**

- Fix for k8gb.io mixed content warnings [\#488](https://github.com/k8gb-io/k8gb/pull/488) ([somaritane](https://github.com/somaritane))

**Merged pull requests:**

- Prepare for v0.8.1 release [\#521](https://github.com/k8gb-io/k8gb/pull/521) ([ytsarev](https://github.com/ytsarev))
- Point Github Actions based badges to master branch [\#520](https://github.com/k8gb-io/k8gb/pull/520) ([ytsarev](https://github.com/ytsarev))
- Scan k8gb image by Artifacthub.io [\#519](https://github.com/k8gb-io/k8gb/pull/519) ([k0da](https://github.com/k0da))
- Fix DNS query logging message [\#518](https://github.com/k8gb-io/k8gb/pull/518) ([ytsarev](https://github.com/ytsarev))
- refactoring suggestions [\#517](https://github.com/k8gb-io/k8gb/pull/517) ([kuritka](https://github.com/kuritka))
- Use edgeDNSServer for NS name resolution [\#516](https://github.com/k8gb-io/k8gb/pull/516) ([ytsarev](https://github.com/ytsarev))
- fakeDNS testupdate [\#515](https://github.com/k8gb-io/k8gb/pull/515) ([kuritka](https://github.com/kuritka))
- fix local playground, list of local-targets [\#512](https://github.com/k8gb-io/k8gb/pull/512) ([kuritka](https://github.com/kuritka))
- fix number of addresses in local.md [\#511](https://github.com/k8gb-io/k8gb/pull/511) ([kuritka](https://github.com/kuritka))
- Bump coredns plugin version [\#510](https://github.com/k8gb-io/k8gb/pull/510) ([k0da](https://github.com/k0da))
- Fix Github repo links in NS1 docs [\#509](https://github.com/k8gb-io/k8gb/pull/509) ([ytsarev](https://github.com/ytsarev))
- Enhance and simplify NS1 reference deployment example [\#508](https://github.com/k8gb-io/k8gb/pull/508) ([ytsarev](https://github.com/ytsarev))
- refactor \(3/3\): Introducing local FakeDNS [\#507](https://github.com/k8gb-io/k8gb/pull/507) ([kuritka](https://github.com/kuritka))
- Update k8gb curl demo to be usable for real deployments [\#506](https://github.com/k8gb-io/k8gb/pull/506) ([ytsarev](https://github.com/ytsarev))
- refactor \(2/3\): Remove responsibility for target DNS from GSLB assistant [\#505](https://github.com/k8gb-io/k8gb/pull/505) ([kuritka](https://github.com/kuritka))
- refactor \(1/3\): simplify controller tests [\#503](https://github.com/k8gb-io/k8gb/pull/503) ([kuritka](https://github.com/kuritka))
- Add Crossplane Day recording [\#502](https://github.com/k8gb-io/k8gb/pull/502) ([ytsarev](https://github.com/ytsarev))
- Add k8gb CII Best Practices status badge [\#501](https://github.com/k8gb-io/k8gb/pull/501) ([somaritane](https://github.com/somaritane))
- Fix failing `clean-test-apps` make target [\#500](https://github.com/k8gb-io/k8gb/pull/500) ([somaritane](https://github.com/somaritane))
- Remove duplicated hit-testapp-host make function [\#496](https://github.com/k8gb-io/k8gb/pull/496) ([somaritane](https://github.com/somaritane))
- Add override for dev env variables with "dotenv" file [\#495](https://github.com/k8gb-io/k8gb/pull/495) ([somaritane](https://github.com/somaritane))
- Update docs for containerized website authoring [\#494](https://github.com/k8gb-io/k8gb/pull/494) ([somaritane](https://github.com/somaritane))
- Refactor test target [\#492](https://github.com/k8gb-io/k8gb/pull/492) ([kuritka](https://github.com/kuritka))
- refactor: interface rename [\#491](https://github.com/k8gb-io/k8gb/pull/491) ([kuritka](https://github.com/kuritka))
- Shrink NS names [\#490](https://github.com/k8gb-io/k8gb/pull/490) ([kuritka](https://github.com/kuritka))
- Fix for k8gb.io mixed content warnings [\#489](https://github.com/k8gb-io/k8gb/pull/489) ([somaritane](https://github.com/somaritane))
- Bump sigs.k8s.io/external-dns from 0.7.6 to 0.8.0 [\#466](https://github.com/k8gb-io/k8gb/pull/466) ([dependabot[bot]](https://github.com/apps/dependabot))
- Flag enabling SplitBrain [\#465](https://github.com/k8gb-io/k8gb/pull/465) ([kuritka](https://github.com/kuritka))

## [v0.8.0](https://github.com/k8gb-io/k8gb/tree/v0.8.0) (2021-05-13)

[Full Changelog](https://github.com/k8gb-io/k8gb/compare/v0.7.7...v0.8.0)

**Implemented enhancements:**

- Thin down available Infoblox connections [\#463](https://github.com/k8gb-io/k8gb/issues/463)
- Split quickstart focused and developer focused make targets [\#446](https://github.com/k8gb-io/k8gb/issues/446)
- Create governance document [\#436](https://github.com/k8gb-io/k8gb/issues/436)
- automate ingress validation annotation [\#401](https://github.com/k8gb-io/k8gb/issues/401)
- Cover RoundRobin IP list merge with Terratest [\#389](https://github.com/k8gb-io/k8gb/issues/389)
- Switch local setup to newer nginx ingress controller helm chart [\#388](https://github.com/k8gb-io/k8gb/issues/388)
- Upgrade to operator-sdk v1.5.0 [\#376](https://github.com/k8gb-io/k8gb/issues/376)
- Issue when deleting ingress rule or annotations removal doesn't remove the gslb records [\#361](https://github.com/k8gb-io/k8gb/issues/361)
- Reuse/enhance terratest test suite for real cluster validation [\#350](https://github.com/k8gb-io/k8gb/issues/350)
- Automate upgrade testing [\#349](https://github.com/k8gb-io/k8gb/issues/349)
- Rename ohmyterratest module to k8gbterratest [\#348](https://github.com/k8gb-io/k8gb/issues/348)
- Logger Enhancements [\#331](https://github.com/k8gb-io/k8gb/issues/331)
- revisit k8gb service account permissions [\#330](https://github.com/k8gb-io/k8gb/issues/330)
- Add support for `k8gb.io/dns-ttl-seconds` and `k8gb.io/splitbrain-threshold-seconds` strategy annotations [\#316](https://github.com/k8gb-io/k8gb/issues/316)

**Fixed bugs:**

- k8gb allows to load multiple providers [\#448](https://github.com/k8gb-io/k8gb/issues/448)
- Existing DNSEndpoint resources are not re-labeled with dnstype after v0.7.5 upgrade [\#324](https://github.com/k8gb-io/k8gb/issues/324)

**Merged pull requests:**

- Fix base for changelog PR [\#486](https://github.com/k8gb-io/k8gb/pull/486) ([k0da](https://github.com/k0da))
- Fix chart repo url after org move [\#484](https://github.com/k8gb-io/k8gb/pull/484) ([k0da](https://github.com/k0da))
- Unify external-dns deployment [\#481](https://github.com/k8gb-io/k8gb/pull/481) ([k0da](https://github.com/k0da))
- Fix NS1 deployment [\#480](https://github.com/k8gb-io/k8gb/pull/480) ([k0da](https://github.com/k0da))
- Updated CONTRIBUTING documentation [\#477](https://github.com/k8gb-io/k8gb/pull/477) ([somaritane](https://github.com/somaritane))
- Bump github.com/miekg/dns from 1.1.41 to 1.1.42 [\#474](https://github.com/k8gb-io/k8gb/pull/474) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump k8s.io/client-go from 0.20.5 to 0.20.6 [\#472](https://github.com/k8gb-io/k8gb/pull/472) ([dependabot[bot]](https://github.com/apps/dependabot))
- Prepare for 0.8 release [\#471](https://github.com/k8gb-io/k8gb/pull/471) ([k0da](https://github.com/k0da))
- Pass endpoint params for ns1 external-dns provider [\#470](https://github.com/k8gb-io/k8gb/pull/470) ([k0da](https://github.com/k0da))
- Sort delegateTo within infoblox ZoneDelegated [\#467](https://github.com/k8gb-io/k8gb/pull/467) ([k0da](https://github.com/k0da))
- Don't reconcile immediately when creating zone delegation fails [\#462](https://github.com/k8gb-io/k8gb/pull/462) ([kuritka](https://github.com/kuritka))
- Initial version of k8gb Governance document [\#458](https://github.com/k8gb-io/k8gb/pull/458) ([somaritane](https://github.com/somaritane))
- Fix RBAC for k8gb ClusterRole [\#455](https://github.com/k8gb-io/k8gb/pull/455) ([ytsarev](https://github.com/ytsarev))
- golint, check capitalized error strings [\#454](https://github.com/k8gb-io/k8gb/pull/454) ([kuritka](https://github.com/kuritka))
- Enhance terratest suite with ability to be executed against real clusters  [\#452](https://github.com/k8gb-io/k8gb/pull/452) ([ytsarev](https://github.com/ytsarev))
- Terratest timeout 15 min, parallel 12 [\#451](https://github.com/k8gb-io/k8gb/pull/451) ([kuritka](https://github.com/kuritka))
- \[Fix\] Validate when multiple providers are defined [\#450](https://github.com/k8gb-io/k8gb/pull/450) ([kuritka](https://github.com/kuritka))
- bump k3d-action to v1.4.0 [\#449](https://github.com/k8gb-io/k8gb/pull/449) ([kuritka](https://github.com/kuritka))
- Stabilize local setup [\#447](https://github.com/k8gb-io/k8gb/pull/447) ([ytsarev](https://github.com/ytsarev))
- Support for optional Ingress strategy annotations [\#445](https://github.com/k8gb-io/k8gb/pull/445) ([ytsarev](https://github.com/ytsarev))
- Shrink k8gb role to what is really required [\#444](https://github.com/k8gb-io/k8gb/pull/444) ([k0da](https://github.com/k0da))
- FIX: Annotate and Label existing DNSEndpoints [\#443](https://github.com/k8gb-io/k8gb/pull/443) ([k0da](https://github.com/k0da))
- FOSSA scan enabled [\#442](https://github.com/k8gb-io/k8gb/pull/442) ([idvoretskyi](https://github.com/idvoretskyi))
- Update license headers with CNCF recommendations [\#441](https://github.com/k8gb-io/k8gb/pull/441) ([ytsarev](https://github.com/ytsarev))
- TestK8gbBasicRoundRobinExample [\#440](https://github.com/k8gb-io/k8gb/pull/440) ([kuritka](https://github.com/kuritka))
- Service CoreDNS Corefile by k8gb chart [\#439](https://github.com/k8gb-io/k8gb/pull/439) ([k0da](https://github.com/k0da))
- Describe testing setup with k3d config [\#438](https://github.com/k8gb-io/k8gb/pull/438) ([k0da](https://github.com/k0da))
- k8gb playground documentation, update A records for one agent [\#437](https://github.com/k8gb-io/k8gb/pull/437) ([kuritka](https://github.com/kuritka))
- Fix deploy-candidate message [\#435](https://github.com/k8gb-io/k8gb/pull/435) ([kuritka](https://github.com/kuritka))
- FIX: race condition detected [\#432](https://github.com/k8gb-io/k8gb/pull/432) ([kuritka](https://github.com/kuritka))
- Fix possible host name clash in tests [\#430](https://github.com/k8gb-io/k8gb/pull/430) ([k0da](https://github.com/k0da))
- Upgrade testing [\#429](https://github.com/k8gb-io/k8gb/pull/429) ([kuritka](https://github.com/kuritka))
- Fix new line escape [\#428](https://github.com/k8gb-io/k8gb/pull/428) ([k0da](https://github.com/k0da))
- Update CoreDNS chart [\#427](https://github.com/k8gb-io/k8gb/pull/427) ([k0da](https://github.com/k0da))
- Enable Ingress to Gslb Owner Reference [\#426](https://github.com/k8gb-io/k8gb/pull/426) ([ytsarev](https://github.com/ytsarev))
- Extend Gslb CRD with additionalPrinterColumns [\#425](https://github.com/k8gb-io/k8gb/pull/425) ([ytsarev](https://github.com/ytsarev))
- Bump operator SDK to v1.5.0 [\#419](https://github.com/k8gb-io/k8gb/pull/419) ([kuritka](https://github.com/kuritka))
- Add DoK community talk recording [\#418](https://github.com/k8gb-io/k8gb/pull/418) ([ytsarev](https://github.com/ytsarev))
- Migration to networking.k8s.io/v1beta1 [\#417](https://github.com/k8gb-io/k8gb/pull/417) ([kuritka](https://github.com/kuritka))
- bump golic v0.5.0 [\#416](https://github.com/k8gb-io/k8gb/pull/416) ([kuritka](https://github.com/kuritka))
- Rename traces of legacy branding [\#415](https://github.com/k8gb-io/k8gb/pull/415) ([ytsarev](https://github.com/ytsarev))
- upgrade terratest go.mod [\#414](https://github.com/k8gb-io/k8gb/pull/414) ([kuritka](https://github.com/kuritka))
- Improve logging for missing environment variables [\#413](https://github.com/k8gb-io/k8gb/pull/413) ([somaritane](https://github.com/somaritane))
- Enable coredns logging [\#412](https://github.com/k8gb-io/k8gb/pull/412) ([ytsarev](https://github.com/ytsarev))
- Bump github.com/rs/zerolog from 1.20.0 to 1.21.0 [\#411](https://github.com/k8gb-io/k8gb/pull/411) ([dependabot[bot]](https://github.com/apps/dependabot))
- cleaning go.mod from github.com/go-logr/zapr [\#410](https://github.com/k8gb-io/k8gb/pull/410) ([kuritka](https://github.com/kuritka))
- Add AWS Containers from the Couch recording [\#408](https://github.com/k8gb-io/k8gb/pull/408) ([ytsarev](https://github.com/ytsarev))
- Added golangci-lint as pre-requisite to local setup doc [\#407](https://github.com/k8gb-io/k8gb/pull/407) ([somaritane](https://github.com/somaritane))
- log debug, optimization [\#406](https://github.com/k8gb-io/k8gb/pull/406) ([kuritka](https://github.com/kuritka))
- Improve initial logging experience [\#405](https://github.com/k8gb-io/k8gb/pull/405) ([somaritane](https://github.com/somaritane))
- Offline Changelog for v0.7.7 [\#404](https://github.com/k8gb-io/k8gb/pull/404) ([somaritane](https://github.com/somaritane))
- Split changelog PR off helm publish workflow [\#403](https://github.com/k8gb-io/k8gb/pull/403) ([k0da](https://github.com/k0da))
- HTTP ingress rule value is Mandatory [\#402](https://github.com/k8gb-io/k8gb/pull/402) ([kuritka](https://github.com/kuritka))
- Update nginx ingress chart [\#391](https://github.com/k8gb-io/k8gb/pull/391) ([k0da](https://github.com/k0da))
- use gopkg.strings.Format\(\) instead of local utils.ToString\(\) [\#387](https://github.com/k8gb-io/k8gb/pull/387) ([kuritka](https://github.com/kuritka))

## [v0.7.7](https://github.com/k8gb-io/k8gb/tree/v0.7.7) (2021-03-22)

[Full Changelog](https://github.com/k8gb-io/k8gb/compare/v0.7.6...v0.7.7)

**Implemented enhancements:**

- Remove `udp-services` ConfigMap creation from k8gb helm chart [\#352](https://github.com/k8gb-io/k8gb/issues/352)
- migrate off deprecated APIs [\#347](https://github.com/k8gb-io/k8gb/issues/347)
- Deprecate `expose53onWorkers` configuration option [\#323](https://github.com/k8gb-io/k8gb/issues/323)
- Add Arm support [\#243](https://github.com/k8gb-io/k8gb/issues/243)

**Fixed bugs:**

- k8gb CRD is removed during helm chart upgrade [\#345](https://github.com/k8gb-io/k8gb/issues/345)
- Installing on a cluster with an existing udp-services ConfigMap fails [\#164](https://github.com/k8gb-io/k8gb/issues/164)

**Closed issues:**

- Document metrics exposure via Prometheus Operator [\#119](https://github.com/k8gb-io/k8gb/issues/119)

**Merged pull requests:**

- Update relative link in doc [\#400](https://github.com/k8gb-io/k8gb/pull/400) ([ytsarev](https://github.com/ytsarev))
- Switch to relative link in cross reference doc [\#399](https://github.com/k8gb-io/k8gb/pull/399) ([ytsarev](https://github.com/ytsarev))
- Fix github\_changelog\_generator defaults [\#398](https://github.com/k8gb-io/k8gb/pull/398) ([k0da](https://github.com/k0da))
- Include pull-requests into changelog [\#397](https://github.com/k8gb-io/k8gb/pull/397) ([k0da](https://github.com/k0da))
- Fix grammar in NOTES.txt [\#395](https://github.com/k8gb-io/k8gb/pull/395) ([ytsarev](https://github.com/ytsarev))
- Remove kustomize and associated make targets [\#393](https://github.com/k8gb-io/k8gb/pull/393) ([somaritane](https://github.com/somaritane))
- Rollback external-dns to get NS record creation back [\#392](https://github.com/k8gb-io/k8gb/pull/392) ([ytsarev](https://github.com/ytsarev))
- Fix helm chart NOTES.txt [\#390](https://github.com/k8gb-io/k8gb/pull/390) ([ytsarev](https://github.com/ytsarev))
- bump golic version [\#385](https://github.com/k8gb-io/k8gb/pull/385) ([kuritka](https://github.com/kuritka))
- Bump github.com/miekg/dns from 1.1.40 to 1.1.41 [\#383](https://github.com/k8gb-io/k8gb/pull/383) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump k8s.io/client-go from 0.20.4 to 0.20.5 [\#382](https://github.com/k8gb-io/k8gb/pull/382) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/prometheus/client\_golang from 1.9.0 to 1.10.0 [\#381](https://github.com/k8gb-io/k8gb/pull/381) ([dependabot[bot]](https://github.com/apps/dependabot))
- Fix changelog in Release [\#378](https://github.com/k8gb-io/k8gb/pull/378) ([k0da](https://github.com/k0da))
- put license back to test target [\#377](https://github.com/k8gb-io/k8gb/pull/377) ([kuritka](https://github.com/kuritka))
- Generation tools improvements [\#375](https://github.com/k8gb-io/k8gb/pull/375) ([somaritane](https://github.com/somaritane))
- Release v0.7.7 preparation [\#374](https://github.com/k8gb-io/k8gb/pull/374) ([somaritane](https://github.com/somaritane))
- Route53 documentation fixes [\#373](https://github.com/k8gb-io/k8gb/pull/373) ([ytsarev](https://github.com/ytsarev))
- Golic CI [\#372](https://github.com/k8gb-io/k8gb/pull/372) ([kuritka](https://github.com/kuritka))
- Drop linter from terratest action [\#371](https://github.com/k8gb-io/k8gb/pull/371) ([k0da](https://github.com/k0da))
- Licenses to be compatible with vscode editor [\#370](https://github.com/k8gb-io/k8gb/pull/370) ([kuritka](https://github.com/kuritka))
- Use simple log format as default for make run [\#369](https://github.com/k8gb-io/k8gb/pull/369) ([somaritane](https://github.com/somaritane))
- License management with GOLIC [\#368](https://github.com/k8gb-io/k8gb/pull/368) ([kuritka](https://github.com/kuritka))
- Integration zerolog [\#367](https://github.com/k8gb-io/k8gb/pull/367) ([kuritka](https://github.com/kuritka))
- Import image [\#363](https://github.com/k8gb-io/k8gb/pull/363) ([k0da](https://github.com/k0da))
- Update DNSEndpoint CRD [\#360](https://github.com/k8gb-io/k8gb/pull/360) ([k0da](https://github.com/k0da))
- change License icon [\#358](https://github.com/k8gb-io/k8gb/pull/358) ([kuritka](https://github.com/kuritka))
- Remove deploy-gslb-operator-14 make target [\#357](https://github.com/k8gb-io/k8gb/pull/357) ([somaritane](https://github.com/somaritane))
- Logger factory [\#356](https://github.com/k8gb-io/k8gb/pull/356) ([kuritka](https://github.com/kuritka))
- Reduce load on test setup [\#355](https://github.com/k8gb-io/k8gb/pull/355) ([k0da](https://github.com/k0da))
- Remove `udp-services` ConfigMap from k8gb helm chart templates [\#354](https://github.com/k8gb-io/k8gb/pull/354) ([somaritane](https://github.com/somaritane))
- Update apiextensions to v1 [\#353](https://github.com/k8gb-io/k8gb/pull/353) ([k0da](https://github.com/k0da))
- Move crds back to templates folder [\#346](https://github.com/k8gb-io/k8gb/pull/346) ([k0da](https://github.com/k0da))
- Fix the license text [\#344](https://github.com/k8gb-io/k8gb/pull/344) ([ytsarev](https://github.com/ytsarev))
- Add Apache 2 license header to every Go file [\#343](https://github.com/k8gb-io/k8gb/pull/343) ([ytsarev](https://github.com/ytsarev))
- Update Contribution guide after changing the license [\#342](https://github.com/k8gb-io/k8gb/pull/342) ([ytsarev](https://github.com/ytsarev))
- Add links to k8gb presentation recordings [\#341](https://github.com/k8gb-io/k8gb/pull/341) ([ytsarev](https://github.com/ytsarev))
- Add Code of Conduct [\#340](https://github.com/k8gb-io/k8gb/pull/340) ([ytsarev](https://github.com/ytsarev))
- Switch to Apache 2 license [\#339](https://github.com/k8gb-io/k8gb/pull/339) ([ytsarev](https://github.com/ytsarev))
- Logger input Environment variables  [\#338](https://github.com/k8gb-io/k8gb/pull/338) ([kuritka](https://github.com/kuritka))
- bump k3d-action to v1.3.1 [\#337](https://github.com/k8gb-io/k8gb/pull/337) ([kuritka](https://github.com/kuritka))
- Offline v0.7.6 release notes [\#335](https://github.com/k8gb-io/k8gb/pull/335) ([somaritane](https://github.com/somaritane))
- Automate releases [\#334](https://github.com/k8gb-io/k8gb/pull/334) ([k0da](https://github.com/k0da))

## [v0.7.6](https://github.com/k8gb-io/k8gb/tree/v0.7.6) (2021-03-01)

[Full Changelog](https://github.com/k8gb-io/k8gb/compare/v0.7.5...v0.7.6)

**Fixed bugs:**

- edgeDNS external-dns pods are failing to start since v0.7.5 [\#328](https://github.com/k8gb-io/k8gb/issues/328)
- "Split brain TXT record expired the time threshold: \(0s\)"  when `gslb` CR gets recreated [\#317](https://github.com/k8gb-io/k8gb/issues/317)

**Closed issues:**

- \[helm chart\] fsGroup not a valid securityContext field [\#293](https://github.com/k8gb-io/k8gb/issues/293)

**Merged pull requests:**

- Release v0.7.6 preparation [\#333](https://github.com/k8gb-io/k8gb/pull/333) ([somaritane](https://github.com/somaritane))
- Make k8gb demo curl script ready for local invocation [\#332](https://github.com/k8gb-io/k8gb/pull/332) ([ytsarev](https://github.com/ytsarev))
- Bring back external-dns service account [\#329](https://github.com/k8gb-io/k8gb/pull/329) ([k0da](https://github.com/k0da))
- Solve fsGroup issue [\#327](https://github.com/k8gb-io/k8gb/pull/327) ([ytsarev](https://github.com/ytsarev))
- Update absaoss/k8s\_crd CoreDNS plugin to v0.0.2 [\#326](https://github.com/k8gb-io/k8gb/pull/326) ([k0da](https://github.com/k0da))
- Doc crds badge [\#325](https://github.com/k8gb-io/k8gb/pull/325) ([ytsarev](https://github.com/ytsarev))
- Fix \#317, depresolver load new values when GSLB recreated [\#322](https://github.com/k8gb-io/k8gb/pull/322) ([kuritka](https://github.com/kuritka))
- Bump github.com/miekg/dns from 1.1.39 to 1.1.40 [\#321](https://github.com/k8gb-io/k8gb/pull/321) ([dependabot[bot]](https://github.com/apps/dependabot))
- Offline v0.7.5 release notes [\#320](https://github.com/k8gb-io/k8gb/pull/320) ([somaritane](https://github.com/somaritane))
- disable CoreDNS cache [\#315](https://github.com/k8gb-io/k8gb/pull/315) ([k0da](https://github.com/k0da))
- Validate spec.ingress.http.path [\#313](https://github.com/k8gb-io/k8gb/pull/313) ([k0da](https://github.com/k0da))

## [v0.7.5](https://github.com/k8gb-io/k8gb/tree/v0.7.5) (2021-02-24)

[Full Changelog](https://github.com/k8gb-io/k8gb/compare/v0.7.4...v0.7.5)

**Implemented enhancements:**

- coredns CRD plugin [\#249](https://github.com/k8gb-io/k8gb/issues/249)
- Move away from coreos etcd-operator [\#220](https://github.com/k8gb-io/k8gb/issues/220)

**Fixed bugs:**

- k8gb crashes on malformed spec section in `gslb` custom resource [\#296](https://github.com/k8gb-io/k8gb/issues/296)

**Merged pull requests:**

- Release v0.7.5 preparation [\#318](https://github.com/k8gb-io/k8gb/pull/318) ([somaritane](https://github.com/somaritane))
- Use SetAnnotation helper [\#314](https://github.com/k8gb-io/k8gb/pull/314) ([k0da](https://github.com/k0da))
- Infoblox, heavy load fixed [\#312](https://github.com/k8gb-io/k8gb/pull/312) ([kuritka](https://github.com/kuritka))
- Sort externalTargets queried from DNS [\#311](https://github.com/k8gb-io/k8gb/pull/311) ([k0da](https://github.com/k0da))
- Bump k8s.io/client-go group from 0.20.3 to 0.20.4 [\#310](https://github.com/k8gb-io/k8gb/pull/310) ([kuritka](https://github.com/kuritka))
- group version bump [\#306](https://github.com/k8gb-io/k8gb/pull/306) ([kuritka](https://github.com/kuritka))
- Fail on config error [\#302](https://github.com/k8gb-io/k8gb/pull/302) ([kuritka](https://github.com/kuritka))
- bump k3d-action to v 1.2.0 [\#295](https://github.com/k8gb-io/k8gb/pull/295) ([kuritka](https://github.com/kuritka))
- Switch to coredns with DNSendpoint plugin [\#292](https://github.com/k8gb-io/k8gb/pull/292) ([k0da](https://github.com/k0da))
- Additional chart tweaks for ArtifactHub [\#291](https://github.com/k8gb-io/k8gb/pull/291) ([somaritane](https://github.com/somaritane))
- Trying to please ArtifactHub markdown render [\#290](https://github.com/k8gb-io/k8gb/pull/290) ([somaritane](https://github.com/somaritane))
- Add artifact hub badge [\#288](https://github.com/k8gb-io/k8gb/pull/288) ([ytsarev](https://github.com/ytsarev))
- README: Replaced screenshot with code excerpt [\#287](https://github.com/k8gb-io/k8gb/pull/287) ([somaritane](https://github.com/somaritane))
- Offline v0.7.4 release notes [\#285](https://github.com/k8gb-io/k8gb/pull/285) ([ytsarev](https://github.com/ytsarev))

## [v0.7.4](https://github.com/k8gb-io/k8gb/tree/v0.7.4) (2021-02-05)

[Full Changelog](https://github.com/k8gb-io/k8gb/compare/v0.7.2...v0.7.4)

**Implemented enhancements:**

- Document Struct fields in CRD spec [\#273](https://github.com/k8gb-io/k8gb/issues/273)
- Upgrade to operator-sdk v1.3.0 [\#266](https://github.com/k8gb-io/k8gb/issues/266)
- Missing finalizer for NS1 [\#262](https://github.com/k8gb-io/k8gb/issues/262)
- Include kube-linter into k8gb pipelines [\#254](https://github.com/k8gb-io/k8gb/issues/254)

**Fixed bugs:**

- terratests - Possible race condition [\#211](https://github.com/k8gb-io/k8gb/issues/211)
- Infoblox Zone Delegation not created under correct Auth Zone [\#99](https://github.com/k8gb-io/k8gb/issues/99)

**Closed issues:**

- Split dnsupdate into provider pattern [\#255](https://github.com/k8gb-io/k8gb/issues/255)
- Move Infoblox provider logic to ExternalDNS [\#222](https://github.com/k8gb-io/k8gb/issues/222)
- Feature Request to Possibly Host more than one DNS Zones on K8gb [\#151](https://github.com/k8gb-io/k8gb/issues/151)

**Merged pull requests:**

- Consolidate `v` part of version tag in the Chart metadata [\#284](https://github.com/k8gb-io/k8gb/pull/284) ([ytsarev](https://github.com/ytsarev))
- Enable docker experimental features in GHA [\#283](https://github.com/k8gb-io/k8gb/pull/283) ([k0da](https://github.com/k0da))
- Update CRD yaml metadata [\#282](https://github.com/k8gb-io/k8gb/pull/282) ([ytsarev](https://github.com/ytsarev))
- Prepare for 0.7.4 release [\#281](https://github.com/k8gb-io/k8gb/pull/281) ([ytsarev](https://github.com/ytsarev))
- fix dependabot version upgrade [\#279](https://github.com/k8gb-io/k8gb/pull/279) ([kuritka](https://github.com/kuritka))
- Upgrade to operator-sdk v1.3.0 [\#276](https://github.com/k8gb-io/k8gb/pull/276) ([kuritka](https://github.com/kuritka))
- Update embedded doc strings in CRD spec [\#275](https://github.com/k8gb-io/k8gb/pull/275) ([ytsarev](https://github.com/ytsarev))
- Bump github.com/miekg/dns from 1.1.37 to 1.1.38 [\#274](https://github.com/k8gb-io/k8gb/pull/274) ([dependabot[bot]](https://github.com/apps/dependabot))
- infoblox, extracting HTTPPoolConnections,HTTPRequestTimeout [\#272](https://github.com/k8gb-io/k8gb/pull/272) ([kuritka](https://github.com/kuritka))
- Bump github.com/miekg/dns from 1.1.35 to 1.1.37 [\#271](https://github.com/k8gb-io/k8gb/pull/271) ([dependabot[bot]](https://github.com/apps/dependabot))
- Refactor to providers [\#270](https://github.com/k8gb-io/k8gb/pull/270) ([kuritka](https://github.com/kuritka))
- Extend pipelines with KubeLinter [\#269](https://github.com/k8gb-io/k8gb/pull/269) ([ytsarev](https://github.com/ytsarev))
- Enable docker multiarch build [\#267](https://github.com/k8gb-io/k8gb/pull/267) ([k0da](https://github.com/k0da))
- Upgrade external-dns to v0.7.6 [\#265](https://github.com/k8gb-io/k8gb/pull/265) ([ytsarev](https://github.com/ytsarev))
- Bump github.com/stretchr/testify from 1.6.1 to 1.7.0 [\#264](https://github.com/k8gb-io/k8gb/pull/264) ([dependabot[bot]](https://github.com/apps/dependabot))
- bump to k3d-action@v1.1.0 [\#263](https://github.com/k8gb-io/k8gb/pull/263) ([kuritka](https://github.com/kuritka))
- Fix badge link to Terratest action executions [\#261](https://github.com/k8gb-io/k8gb/pull/261) ([ytsarev](https://github.com/ytsarev))
- Update k8gb config samples with new exposeCoreDNS param [\#260](https://github.com/k8gb-io/k8gb/pull/260) ([ytsarev](https://github.com/ytsarev))
- Make CoreDNS exposure controllable [\#259](https://github.com/k8gb-io/k8gb/pull/259) ([ytsarev](https://github.com/ytsarev))
- refactor prettyPrint [\#258](https://github.com/k8gb-io/k8gb/pull/258) ([kuritka](https://github.com/kuritka))
- Refactor \#2, Dig  [\#257](https://github.com/k8gb-io/k8gb/pull/257) ([kuritka](https://github.com/kuritka))
- Extract prometheus metrics \#1 [\#256](https://github.com/k8gb-io/k8gb/pull/256) ([kuritka](https://github.com/kuritka))
- Bump sigs.k8s.io/external-dns from 0.7.5 to 0.7.6 [\#251](https://github.com/k8gb-io/k8gb/pull/251) ([dependabot[bot]](https://github.com/apps/dependabot))
- Publish CodeQL status tag [\#248](https://github.com/k8gb-io/k8gb/pull/248) ([ytsarev](https://github.com/ytsarev))
- Explicit fqdns in roundrobin sample CR [\#247](https://github.com/k8gb-io/k8gb/pull/247) ([ytsarev](https://github.com/ytsarev))
- Makefile help [\#246](https://github.com/k8gb-io/k8gb/pull/246) ([ytsarev](https://github.com/ytsarev))
- fix terratests [\#245](https://github.com/k8gb-io/k8gb/pull/245) ([kuritka](https://github.com/kuritka))
- README support table update [\#242](https://github.com/k8gb-io/k8gb/pull/242) ([ytsarev](https://github.com/ytsarev))
- Bump github.com/stretchr/testify from 1.5.1 to 1.6.1 [\#241](https://github.com/k8gb-io/k8gb/pull/241) ([dependabot[bot]](https://github.com/apps/dependabot))
- Preparation for artifacthub [\#240](https://github.com/k8gb-io/k8gb/pull/240) ([ytsarev](https://github.com/ytsarev))
- Bump github.com/prometheus/client\_golang from 1.7.1 to 1.9.0 [\#239](https://github.com/k8gb-io/k8gb/pull/239) ([dependabot[bot]](https://github.com/apps/dependabot))
- Switch external-dns to upstream v0.7.5 image release [\#237](https://github.com/k8gb-io/k8gb/pull/237) ([ytsarev](https://github.com/ytsarev))
- Bump sigs.k8s.io/external-dns from 0.7.4 to 0.7.5 [\#235](https://github.com/k8gb-io/k8gb/pull/235) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/miekg/dns from 1.1.30 to 1.1.35 [\#232](https://github.com/k8gb-io/k8gb/pull/232) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/onsi/ginkgo from 1.12.1 to 1.14.2 [\#230](https://github.com/k8gb-io/k8gb/pull/230) ([dependabot[bot]](https://github.com/apps/dependabot))
- Add GitHub code scanning [\#228](https://github.com/k8gb-io/k8gb/pull/228) ([donovanmuller](https://github.com/donovanmuller))
- Add dependabot [\#227](https://github.com/k8gb-io/k8gb/pull/227) ([donovanmuller](https://github.com/donovanmuller))
- bump AbsaOSS/k3d-action to version v1.0.0 [\#226](https://github.com/k8gb-io/k8gb/pull/226) ([kuritka](https://github.com/kuritka))
- Changelog for v0.7.2 [\#225](https://github.com/k8gb-io/k8gb/pull/225) ([ytsarev](https://github.com/ytsarev))

## [v0.7.2](https://github.com/k8gb-io/k8gb/tree/v0.7.2) (2020-12-16)

[Full Changelog](https://github.com/k8gb-io/k8gb/compare/v0.7.1...v0.7.2)

**Implemented enhancements:**

- Support for NS1 Managed DNS [\#202](https://github.com/k8gb-io/k8gb/issues/202)
- Add ability to reuse existing Ingress [\#200](https://github.com/k8gb-io/k8gb/issues/200)
- Move crds from templates/ to crd/ [\#144](https://github.com/k8gb-io/k8gb/issues/144)
- Relax dependency on specific 'k8gb' namespace name [\#129](https://github.com/k8gb-io/k8gb/issues/129)

**Fixed bugs:**

- JSON unmarshall error in ohmyglb logs/ohmyglb status update [\#108](https://github.com/k8gb-io/k8gb/issues/108)

**Closed issues:**

- Replace k8gb.io/primarygeotag annotation with k8gb.io/primary-geotag [\#210](https://github.com/k8gb-io/k8gb/issues/210)
- Consider the switch from kind to k3d [\#141](https://github.com/k8gb-io/k8gb/issues/141)

**Merged pull requests:**

- Switch to new chart repos for dependency charts [\#224](https://github.com/k8gb-io/k8gb/pull/224) ([ytsarev](https://github.com/ytsarev))
- Add missing NS1 api key propagation to the doc [\#223](https://github.com/k8gb-io/k8gb/pull/223) ([ytsarev](https://github.com/ytsarev))
- k3d migration [\#218](https://github.com/k8gb-io/k8gb/pull/218) ([kuritka](https://github.com/kuritka))
- NS1 support [\#217](https://github.com/k8gb-io/k8gb/pull/217) ([ytsarev](https://github.com/ytsarev))
- Fix cluster communication in full local setup [\#216](https://github.com/k8gb-io/k8gb/pull/216) ([ytsarev](https://github.com/ytsarev))
- Relax requirement on k8gb namespace name [\#215](https://github.com/k8gb-io/k8gb/pull/215) ([ytsarev](https://github.com/ytsarev))
- Fix makefile regressions [\#214](https://github.com/k8gb-io/k8gb/pull/214) ([ytsarev](https://github.com/ytsarev))
- Admiralty integration tutorial [\#213](https://github.com/k8gb-io/k8gb/pull/213) ([ytsarev](https://github.com/ytsarev))
- Primary geotag annotation fix [\#212](https://github.com/k8gb-io/k8gb/pull/212) ([somaritane](https://github.com/somaritane))
- regarding Helm Best Practices move CRDs from /templates/crds to /crds [\#209](https://github.com/k8gb-io/k8gb/pull/209) ([kuritka](https://github.com/kuritka))
- Makefile changes on demand [\#208](https://github.com/k8gb-io/k8gb/pull/208) ([kuritka](https://github.com/kuritka))
- Changelog for v0.7.1 [\#206](https://github.com/k8gb-io/k8gb/pull/206) ([ytsarev](https://github.com/ytsarev))

## [v0.7.1](https://github.com/k8gb-io/k8gb/tree/v0.7.1) (2020-11-23)

[Full Changelog](https://github.com/k8gb-io/k8gb/compare/v0.7.0...v0.7.1)

**Implemented enhancements:**

- split of pkg/internal/env into reusable pkg [\#185](https://github.com/k8gb-io/k8gb/issues/185)
- Add support for Route 53 as another edge DNS implementation [\#49](https://github.com/k8gb-io/k8gb/issues/49)

**Closed issues:**

- switch linters to golang-ci [\#197](https://github.com/k8gb-io/k8gb/issues/197)
- Move input environment variables into depresolver [\#170](https://github.com/k8gb-io/k8gb/issues/170)
- Investigate the best place for initializing depresolver and consider it as internal [\#168](https://github.com/k8gb-io/k8gb/issues/168)
- Makefile refactoring [\#109](https://github.com/k8gb-io/k8gb/issues/109)

**Merged pull requests:**

- Enable Gslb with Ingress Annotation [\#205](https://github.com/k8gb-io/k8gb/pull/205) ([ytsarev](https://github.com/ytsarev))
- Contexts complient with kube-builder [\#204](https://github.com/k8gb-io/k8gb/pull/204) ([kuritka](https://github.com/kuritka))
- Change context initialisation, fix helm upgrade [\#203](https://github.com/k8gb-io/k8gb/pull/203) ([kuritka](https://github.com/kuritka))
- Simplify Makefile [\#201](https://github.com/k8gb-io/k8gb/pull/201) ([kuritka](https://github.com/kuritka))
- use AbsaOSS/gopkg  [\#199](https://github.com/k8gb-io/k8gb/pull/199) ([kuritka](https://github.com/kuritka))
- switch to golangci-lint [\#198](https://github.com/k8gb-io/k8gb/pull/198) ([kuritka](https://github.com/kuritka))
- Make diagram image clickable for enlargement [\#196](https://github.com/k8gb-io/k8gb/pull/196) ([ytsarev](https://github.com/ytsarev))
- Fix last 404 [\#195](https://github.com/k8gb-io/k8gb/pull/195) ([ytsarev](https://github.com/ytsarev))
- Use absolute URLs in case of file reference [\#194](https://github.com/k8gb-io/k8gb/pull/194) ([ytsarev](https://github.com/ytsarev))
- Publish CHANGELOG.md to Github Pages [\#193](https://github.com/k8gb-io/k8gb/pull/193) ([ytsarev](https://github.com/ytsarev))
- Integrate depresolver [\#192](https://github.com/k8gb-io/k8gb/pull/192) ([kuritka](https://github.com/kuritka))
- Include CONTRIBUTING.md into gh-pages publishing [\#190](https://github.com/k8gb-io/k8gb/pull/190) ([ytsarev](https://github.com/ytsarev))
- Github Workflow to publish documentation [\#189](https://github.com/k8gb-io/k8gb/pull/189) ([ytsarev](https://github.com/ytsarev))

## [v0.7.0](https://github.com/k8gb-io/k8gb/tree/v0.7.0) (2020-10-28)

[Full Changelog](https://github.com/k8gb-io/k8gb/compare/v0.6.6...v0.7.0)

**Implemented enhancements:**

- Upgrade to operator-sdk 1.0 [\#166](https://github.com/k8gb-io/k8gb/issues/166)
- Route53 support [\#162](https://github.com/k8gb-io/k8gb/issues/162)
- Move the rest of configuration into depresolver [\#122](https://github.com/k8gb-io/k8gb/issues/122)
- Recent gosec fails on generated deep copy code [\#115](https://github.com/k8gb-io/k8gb/issues/115)

**Closed issues:**

- refactor controller\_tests [\#136](https://github.com/k8gb-io/k8gb/issues/136)
- Document internal components of k8gb [\#89](https://github.com/k8gb-io/k8gb/issues/89)

**Merged pull requests:**

- Fix Helm release pipeline [\#188](https://github.com/k8gb-io/k8gb/pull/188) ([ytsarev](https://github.com/ytsarev))
- Commit 'offline' Changelog [\#187](https://github.com/k8gb-io/k8gb/pull/187) ([ytsarev](https://github.com/ytsarev))
- Provide diagram of k8gb internal components [\#186](https://github.com/k8gb-io/k8gb/pull/186) ([ytsarev](https://github.com/ytsarev))
- Finalize Gslb if no route53 DNSEndpoint found [\#184](https://github.com/k8gb-io/k8gb/pull/184) ([ytsarev](https://github.com/ytsarev))
- Include GSLB dns zone into NS server names [\#183](https://github.com/k8gb-io/k8gb/pull/183) ([ytsarev](https://github.com/ytsarev))
- Zone delegation garbage collection for Route53 [\#182](https://github.com/k8gb-io/k8gb/pull/182) ([ytsarev](https://github.com/ytsarev))
- Extend with fake environment variables [\#181](https://github.com/k8gb-io/k8gb/pull/181) ([kuritka](https://github.com/kuritka))
- Post revamp readme fixes [\#180](https://github.com/k8gb-io/k8gb/pull/180) ([ytsarev](https://github.com/ytsarev))
- Readme revamp and Route53 tutorial [\#179](https://github.com/k8gb-io/k8gb/pull/179) ([ytsarev](https://github.com/ytsarev))
- Remove redundant route53.domain from values [\#178](https://github.com/k8gb-io/k8gb/pull/178) ([ytsarev](https://github.com/ytsarev))
- Simplify values.yaml [\#177](https://github.com/k8gb-io/k8gb/pull/177) ([ytsarev](https://github.com/ytsarev))
- Isolate controller tests [\#176](https://github.com/k8gb-io/k8gb/pull/176) ([kuritka](https://github.com/kuritka))
- gosec; ignore generated code [\#174](https://github.com/k8gb-io/k8gb/pull/174) ([kuritka](https://github.com/kuritka))
- Extending DepResolver [\#173](https://github.com/k8gb-io/k8gb/pull/173) ([kuritka](https://github.com/kuritka))
- Route53 support [\#172](https://github.com/k8gb-io/k8gb/pull/172) ([ytsarev](https://github.com/ytsarev))
- Fix external-dns SA definition [\#171](https://github.com/k8gb-io/k8gb/pull/171) ([ytsarev](https://github.com/ytsarev))
- Initial configuration layout for Route53 support [\#169](https://github.com/k8gb-io/k8gb/pull/169) ([ytsarev](https://github.com/ytsarev))

## [v0.6.6](https://github.com/k8gb-io/k8gb/tree/v0.6.6) (2020-10-05)

[Full Changelog](https://github.com/k8gb-io/k8gb/compare/v0.6.5...v0.6.6)

**Closed issues:**

- Rework README to focus on first time users [\#101](https://github.com/k8gb-io/k8gb/issues/101)

**Merged pull requests:**

- Upgrade to operator-sdk 1.0 [\#167](https://github.com/k8gb-io/k8gb/pull/167) ([ytsarev](https://github.com/ytsarev))
- Switch back to upstream etcd-operator chart [\#163](https://github.com/k8gb-io/k8gb/pull/163) ([ytsarev](https://github.com/ytsarev))

## [v0.6.5](https://github.com/k8gb-io/k8gb/tree/v0.6.5) (2020-08-03)

[Full Changelog](https://github.com/k8gb-io/k8gb/compare/v0.6.3...v0.6.5)

**Implemented enhancements:**

- Report on dnsZone and Gslb Ingress host mismatch [\#149](https://github.com/k8gb-io/k8gb/issues/149)

**Merged pull requests:**

- Fix log message about gslb failover strategy execution [\#161](https://github.com/k8gb-io/k8gb/pull/161) ([somaritane](https://github.com/somaritane))
- Add ability to override k8gb image tag [\#160](https://github.com/k8gb-io/k8gb/pull/160) ([somaritane](https://github.com/somaritane))
- Detect mismatch of Ingress hostname and EdgeDNSZone [\#159](https://github.com/k8gb-io/k8gb/pull/159) ([ytsarev](https://github.com/ytsarev))
- Mitigate coredns etcd plugin bug [\#158](https://github.com/k8gb-io/k8gb/pull/158) ([ytsarev](https://github.com/ytsarev))
- Hopefully very last rebranding bit - diagrams [\#157](https://github.com/k8gb-io/k8gb/pull/157) ([ytsarev](https://github.com/ytsarev))
- Last missing rebranding due to the spaces [\#156](https://github.com/k8gb-io/k8gb/pull/156) ([ytsarev](https://github.com/ytsarev))
- Fix local failover example deploy, demo image and demo targets [\#155](https://github.com/k8gb-io/k8gb/pull/155) ([ytsarev](https://github.com/ytsarev))
- fixed wapi credientials and namespace creation [\#153](https://github.com/k8gb-io/k8gb/pull/153) ([jeffhelps](https://github.com/jeffhelps))
- Fix ingress nginx failure in local env and pipelines [\#152](https://github.com/k8gb-io/k8gb/pull/152) ([ytsarev](https://github.com/ytsarev))
- Fix code markup in the readme [\#150](https://github.com/k8gb-io/k8gb/pull/150) ([ytsarev](https://github.com/ytsarev))
- Remove unnecessary infoblox variables from the guide [\#148](https://github.com/k8gb-io/k8gb/pull/148) ([ytsarev](https://github.com/ytsarev))
- An attempt to create step-by-step howto [\#146](https://github.com/k8gb-io/k8gb/pull/146) ([ytsarev](https://github.com/ytsarev))
- Update demo application version [\#145](https://github.com/k8gb-io/k8gb/pull/145) ([ytsarev](https://github.com/ytsarev))
- Increase test app installation timeout [\#143](https://github.com/k8gb-io/k8gb/pull/143) ([ytsarev](https://github.com/ytsarev))
- Switch back to upstream releases [\#142](https://github.com/k8gb-io/k8gb/pull/142) ([ytsarev](https://github.com/ytsarev))

## [v0.6.3](https://github.com/k8gb-io/k8gb/tree/v0.6.3) (2020-06-11)

[Full Changelog](https://github.com/k8gb-io/k8gb/compare/v0.6.2...v0.6.3)

**Implemented enhancements:**

- Make project lintable from project root [\#131](https://github.com/k8gb-io/k8gb/issues/131)

**Merged pull requests:**

- Document currently tested configuration [\#140](https://github.com/k8gb-io/k8gb/pull/140) ([ytsarev](https://github.com/ytsarev))
- Mass rebranding to K8GB [\#139](https://github.com/k8gb-io/k8gb/pull/139) ([ytsarev](https://github.com/ytsarev))
- Mass rebranding to KGB [\#137](https://github.com/k8gb-io/k8gb/pull/137) ([ytsarev](https://github.com/ytsarev))
- Switch to safe geotag propagation with depresolver [\#135](https://github.com/k8gb-io/k8gb/pull/135) ([ytsarev](https://github.com/ytsarev))
- Ability to override registry image [\#133](https://github.com/k8gb-io/k8gb/pull/133) ([ytsarev](https://github.com/ytsarev))
- Make project lintable from project's root [\#132](https://github.com/k8gb-io/k8gb/pull/132) ([kuritka](https://github.com/kuritka))

## [v0.6.2](https://github.com/k8gb-io/k8gb/tree/v0.6.2) (2020-05-20)

[Full Changelog](https://github.com/k8gb-io/k8gb/compare/v0.6.1...v0.6.2)

**Merged pull requests:**

- Fix helm installation smoke test [\#130](https://github.com/k8gb-io/k8gb/pull/130) ([ytsarev](https://github.com/ytsarev))
- Fix issues with public release [\#128](https://github.com/k8gb-io/k8gb/pull/128) ([ytsarev](https://github.com/ytsarev))

## [v0.6.1](https://github.com/k8gb-io/k8gb/tree/v0.6.1) (2020-05-20)

[Full Changelog](https://github.com/k8gb-io/k8gb/compare/v0.6.0...v0.6.1)

**Merged pull requests:**

- Release 0.6.1 [\#127](https://github.com/k8gb-io/k8gb/pull/127) ([ytsarev](https://github.com/ytsarev))
- Simplify versioning process [\#126](https://github.com/k8gb-io/k8gb/pull/126) ([ytsarev](https://github.com/ytsarev))

## [v0.6.0](https://github.com/k8gb-io/k8gb/tree/v0.6.0) (2020-05-16)

[Full Changelog](https://github.com/k8gb-io/k8gb/compare/v0.5.6...v0.6.0)

**Implemented enhancements:**

- Streamline Gslb Status [\#116](https://github.com/k8gb-io/k8gb/issues/116)
- Propagate Gslb CR annotations down to Gslb ingress [\#113](https://github.com/k8gb-io/k8gb/issues/113)
- Make Gslb timeouts and synchronisation intervals configurable [\#82](https://github.com/k8gb-io/k8gb/issues/82)
- Prepare Helm chart for uploading various repositories [\#75](https://github.com/k8gb-io/k8gb/issues/75)
- Extend documentation with end-to-end application deployment scenario [\#69](https://github.com/k8gb-io/k8gb/issues/69)
- Add full end to end integration tests to build pipeline [\#48](https://github.com/k8gb-io/k8gb/issues/48)
- Expose metrics and tracing [\#47](https://github.com/k8gb-io/k8gb/issues/47)

**Fixed bugs:**

- Non-deterministic failure of EtcdCluster deployment in air-gapped on-prem environments [\#107](https://github.com/k8gb-io/k8gb/issues/107)
- Flaky terrarest `TestOhmyglbBasicAppExample` [\#105](https://github.com/k8gb-io/k8gb/issues/105)

**Closed issues:**

- Can't install chart successfully [\#104](https://github.com/k8gb-io/k8gb/issues/104)

**Merged pull requests:**

- Extend release pipeline with docker build and push [\#125](https://github.com/k8gb-io/k8gb/pull/125) ([ytsarev](https://github.com/ytsarev))
- Streamline Gslb Status [\#121](https://github.com/k8gb-io/k8gb/pull/121) ([ytsarev](https://github.com/ytsarev))
- Extend `deploy-gslb-cr` target with failover strategy [\#118](https://github.com/k8gb-io/k8gb/pull/118) ([ytsarev](https://github.com/ytsarev))
- Configurable timeouts and synchronisation intervals [\#117](https://github.com/k8gb-io/k8gb/pull/117) ([kuritka](https://github.com/kuritka))
- Propagate Gslb CR annotations down to Gslb ingress [\#114](https://github.com/k8gb-io/k8gb/pull/114) ([ytsarev](https://github.com/ytsarev))
- Properly propagate etcd version in EtcdCluster CR [\#112](https://github.com/k8gb-io/k8gb/pull/112) ([ytsarev](https://github.com/ytsarev))
- Make basic app terratest reliable [\#111](https://github.com/k8gb-io/k8gb/pull/111) ([ytsarev](https://github.com/ytsarev))
- Optimize and cleanup test-apps target and samples [\#110](https://github.com/k8gb-io/k8gb/pull/110) ([ytsarev](https://github.com/ytsarev))
- Optimize CI status badges [\#106](https://github.com/k8gb-io/k8gb/pull/106) ([ytsarev](https://github.com/ytsarev))
- Failover demo [\#103](https://github.com/k8gb-io/k8gb/pull/103) ([kuritka](https://github.com/kuritka))
- Non deterministic round robin demo [\#98](https://github.com/k8gb-io/k8gb/pull/98) ([kuritka](https://github.com/kuritka))
- Initial operator metrics [\#97](https://github.com/k8gb-io/k8gb/pull/97) ([somaritane](https://github.com/somaritane))
- Add capability to end-to-end test HEAD of the branch [\#96](https://github.com/k8gb-io/k8gb/pull/96) ([ytsarev](https://github.com/ytsarev))
- Enhance terratest pipeline [\#95](https://github.com/k8gb-io/k8gb/pull/95) ([ytsarev](https://github.com/ytsarev))
- Etcd-operator as own subchart [\#94](https://github.com/k8gb-io/k8gb/pull/94) ([ytsarev](https://github.com/ytsarev))
- Include gosec into pipeline [\#93](https://github.com/k8gb-io/k8gb/pull/93) ([ytsarev](https://github.com/ytsarev))
- Terratest based end-to-end pipeline  [\#91](https://github.com/k8gb-io/k8gb/pull/91) ([ytsarev](https://github.com/ytsarev))
- Document Helm repo and installation [\#88](https://github.com/k8gb-io/k8gb/pull/88) ([ytsarev](https://github.com/ytsarev))
- How to run Oh My GLB locally [\#87](https://github.com/k8gb-io/k8gb/pull/87) ([kuritka](https://github.com/kuritka))

## [v0.5.6](https://github.com/k8gb-io/k8gb/tree/v0.5.6) (2020-04-14)

[Full Changelog](https://github.com/k8gb-io/k8gb/compare/v0.5.1...v0.5.6)

**Implemented enhancements:**

- When using the failover load balancing strategy, investigate and validate how resolution will be handled effectively when clusters are configured for mutual failover [\#67](https://github.com/k8gb-io/k8gb/issues/67)
- TTL control for splitbrain TXT record [\#61](https://github.com/k8gb-io/k8gb/issues/61)
- Implement failover load balancing strategy [\#46](https://github.com/k8gb-io/k8gb/issues/46)
- Posssible Routing Peering Capabilities BGP protocols [\#33](https://github.com/k8gb-io/k8gb/issues/33)

**Fixed bugs:**

- Missing endpoints in `localtargets.\*` A records [\#62](https://github.com/k8gb-io/k8gb/issues/62)
- Non-deterministic issue with `localtargets.\*` DNSEntrypoint population [\#38](https://github.com/k8gb-io/k8gb/issues/38)

**Closed issues:**

- Upgrade underlying operator-sdk version from v0.12.0 to latest upstream [\#71](https://github.com/k8gb-io/k8gb/issues/71)
- High Five [\#41](https://github.com/k8gb-io/k8gb/issues/41)

**Merged pull requests:**

- Helm package and publish on release event [\#86](https://github.com/k8gb-io/k8gb/pull/86) ([ytsarev](https://github.com/ytsarev))
- test upgraded build pipe [\#85](https://github.com/k8gb-io/k8gb/pull/85) ([kuritka](https://github.com/kuritka))
- Test mutual failover setup [\#84](https://github.com/k8gb-io/k8gb/pull/84) ([ytsarev](https://github.com/ytsarev))
- Upgrade operator sdk to v0.16.0 [\#83](https://github.com/k8gb-io/k8gb/pull/83) ([somaritane](https://github.com/somaritane))
- Reduce external-dns sync interval to 20s [\#81](https://github.com/k8gb-io/k8gb/pull/81) ([ytsarev](https://github.com/ytsarev))
- Time measure failover process [\#80](https://github.com/k8gb-io/k8gb/pull/80) ([ytsarev](https://github.com/ytsarev))
- Terratest e2e for Failover strategy [\#79](https://github.com/k8gb-io/k8gb/pull/79) ([ytsarev](https://github.com/ytsarev))
- Fix cluster namespaces permission for ohmyglb [\#77](https://github.com/k8gb-io/k8gb/pull/77) ([somaritane](https://github.com/somaritane))
- Terratest for standard ohmyglb deployment with app [\#76](https://github.com/k8gb-io/k8gb/pull/76) ([ytsarev](https://github.com/ytsarev))
- Terratest e2e testing proposal [\#74](https://github.com/k8gb-io/k8gb/pull/74) ([ytsarev](https://github.com/ytsarev))
- Expose all namespaces in ServeCRMetrics [\#73](https://github.com/k8gb-io/k8gb/pull/73) ([ytsarev](https://github.com/ytsarev))
- Fix docker repo link for external-dns [\#72](https://github.com/k8gb-io/k8gb/pull/72) ([ytsarev](https://github.com/ytsarev))
- Bump to include external-dns image with the bugfix [\#70](https://github.com/k8gb-io/k8gb/pull/70) ([ytsarev](https://github.com/ytsarev))
- Use custom build of external-dns with multi A fixes [\#68](https://github.com/k8gb-io/k8gb/pull/68) ([ytsarev](https://github.com/ytsarev))
- Failover strategy post e2e stabilization [\#66](https://github.com/k8gb-io/k8gb/pull/66) ([ytsarev](https://github.com/ytsarev))
- Failover strategy implementation [\#65](https://github.com/k8gb-io/k8gb/pull/65) ([ytsarev](https://github.com/ytsarev))
- Set low TTL on split brain TXT record via infoblox API [\#64](https://github.com/k8gb-io/k8gb/pull/64) ([ytsarev](https://github.com/ytsarev))
- Fully automated multicluster ohmyglb local deployment [\#63](https://github.com/k8gb-io/k8gb/pull/63) ([ytsarev](https://github.com/ytsarev))
- Splitbrain enhancements and fixes [\#60](https://github.com/k8gb-io/k8gb/pull/60) ([ytsarev](https://github.com/ytsarev))
- Bump to 5.3 to stabilize split brain handling [\#59](https://github.com/k8gb-io/k8gb/pull/59) ([ytsarev](https://github.com/ytsarev))
- Infoblox update [\#58](https://github.com/k8gb-io/k8gb/pull/58) ([ytsarev](https://github.com/ytsarev))
- Splitbrain fixes [\#57](https://github.com/k8gb-io/k8gb/pull/57) ([ytsarev](https://github.com/ytsarev))
- Config and helpers for local multicluster setup [\#56](https://github.com/k8gb-io/k8gb/pull/56) ([ytsarev](https://github.com/ytsarev))
- Move to `absaoss` in dockerhub and version bump [\#55](https://github.com/k8gb-io/k8gb/pull/55) ([ytsarev](https://github.com/ytsarev))
- Split brain handling [\#44](https://github.com/k8gb-io/k8gb/pull/44) ([ytsarev](https://github.com/ytsarev))
- Disable `external-dns` ownership for local coredns [\#43](https://github.com/k8gb-io/k8gb/pull/43) ([ytsarev](https://github.com/ytsarev))
- Quote geo tag declaration [\#42](https://github.com/k8gb-io/k8gb/pull/42) ([ytsarev](https://github.com/ytsarev))

## [v0.5.1](https://github.com/k8gb-io/k8gb/tree/v0.5.1) (2020-02-02)

[Full Changelog](https://github.com/k8gb-io/k8gb/compare/d834431a8236e7bbe2769df41bc0e02ceb5afeb3...v0.5.1)

**Merged pull requests:**

- CRUD gslb zone delegation in infoblox [\#39](https://github.com/k8gb-io/k8gb/pull/39) ([ytsarev](https://github.com/ytsarev))
- Multi node local kind cluster [\#37](https://github.com/k8gb-io/k8gb/pull/37) ([ytsarev](https://github.com/ytsarev))
- Initial Edge DNS support  [\#36](https://github.com/k8gb-io/k8gb/pull/36) ([ytsarev](https://github.com/ytsarev))
- Use `podinfo` as example test app [\#35](https://github.com/k8gb-io/k8gb/pull/35) ([ytsarev](https://github.com/ytsarev))
- Enable periodic reconciliation [\#34](https://github.com/k8gb-io/k8gb/pull/34) ([ytsarev](https://github.com/ytsarev))
- External dns ownership fix [\#32](https://github.com/k8gb-io/k8gb/pull/32) ([ytsarev](https://github.com/ytsarev))
- Tolerate external Gslb downtime [\#31](https://github.com/k8gb-io/k8gb/pull/31) ([ytsarev](https://github.com/ytsarev))
- DNS based cross Gslb communication [\#30](https://github.com/k8gb-io/k8gb/pull/30) ([ytsarev](https://github.com/ytsarev))
- BUGFIX: populate record status only when it's ready [\#29](https://github.com/k8gb-io/k8gb/pull/29) ([ytsarev](https://github.com/ytsarev))
- Expose DNS records for heatlhy hosts in Gslb Status [\#28](https://github.com/k8gb-io/k8gb/pull/28) ([ytsarev](https://github.com/ytsarev))
- Change example domain to `example.com` [\#27](https://github.com/k8gb-io/k8gb/pull/27) ([ytsarev](https://github.com/ytsarev))
- Ohmyglb operator chart [\#26](https://github.com/k8gb-io/k8gb/pull/26) ([ytsarev](https://github.com/ytsarev))
- Simple push/build helpers [\#25](https://github.com/k8gb-io/k8gb/pull/25) ([ytsarev](https://github.com/ytsarev))
- Expose coredns\(53 udp\) with nginx ingress controller [\#24](https://github.com/k8gb-io/k8gb/pull/24) ([ytsarev](https://github.com/ytsarev))
- Enhancements to local test configuration [\#23](https://github.com/k8gb-io/k8gb/pull/23) ([ytsarev](https://github.com/ytsarev))
- E2e test suite extension and optimization [\#22](https://github.com/k8gb-io/k8gb/pull/22) ([ytsarev](https://github.com/ytsarev))
- e2e tests for Gslb creation [\#21](https://github.com/k8gb-io/k8gb/pull/21) ([ytsarev](https://github.com/ytsarev))
- Foundation for e2e tests [\#20](https://github.com/k8gb-io/k8gb/pull/20) ([ytsarev](https://github.com/ytsarev))
- Deprecate coreDNS hosts config and worker health checks [\#19](https://github.com/k8gb-io/k8gb/pull/19) ([ytsarev](https://github.com/ytsarev))
- Switch source of addresses for A records to Ingress [\#18](https://github.com/k8gb-io/k8gb/pull/18) ([ytsarev](https://github.com/ytsarev))
- Dynamically populate DNSEndpoints according to health status  [\#17](https://github.com/k8gb-io/k8gb/pull/17) ([ytsarev](https://github.com/ytsarev))
- Register and watch for DNSEndpoints [\#16](https://github.com/k8gb-io/k8gb/pull/16) ([ytsarev](https://github.com/ytsarev))
- Foundation for external-dns DNSEndpoint creation [\#15](https://github.com/k8gb-io/k8gb/pull/15) ([ytsarev](https://github.com/ytsarev))
- Prototype of external-dns + coredns based configuration [\#14](https://github.com/k8gb-io/k8gb/pull/14) ([ytsarev](https://github.com/ytsarev))
- Make OhMyGlb operator watch all namespaces for Gslb CRs [\#13](https://github.com/k8gb-io/k8gb/pull/13) ([ytsarev](https://github.com/ytsarev))
- Add some badges [\#12](https://github.com/k8gb-io/k8gb/pull/12) ([ytsarev](https://github.com/ytsarev))
- Reconcile Gslb when relevant Endpoint is updated [\#11](https://github.com/k8gb-io/k8gb/pull/11) ([ytsarev](https://github.com/ytsarev))
- Enable golint in the pipeline, fix code accordingly [\#10](https://github.com/k8gb-io/k8gb/pull/10) ([ytsarev](https://github.com/ytsarev))
- Control coredns hosts config map [\#9](https://github.com/k8gb-io/k8gb/pull/9) ([ytsarev](https://github.com/ytsarev))
- Expose healthy workers and their ip addresses [\#8](https://github.com/k8gb-io/k8gb/pull/8) ([ytsarev](https://github.com/ytsarev))
- Install CoreDNS from stable chart with custom values [\#7](https://github.com/k8gb-io/k8gb/pull/7) ([ytsarev](https://github.com/ytsarev))
- Gslb Controller Unit Tests [\#6](https://github.com/k8gb-io/k8gb/pull/6) ([ytsarev](https://github.com/ytsarev))
- Gslb Ingress management and associated health checks [\#5](https://github.com/k8gb-io/k8gb/pull/5) ([ytsarev](https://github.com/ytsarev))
- \[WIP\] First iteration of ohmyglb operator [\#3](https://github.com/k8gb-io/k8gb/pull/3) ([ytsarev](https://github.com/ytsarev))
- Additional doc links [\#2](https://github.com/k8gb-io/k8gb/pull/2) ([ytsarev](https://github.com/ytsarev))
- Take readiness probes into account [\#1](https://github.com/k8gb-io/k8gb/pull/1) ([ytsarev](https://github.com/ytsarev))



\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
