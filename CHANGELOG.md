# Changelog

## [v0.7.0](https://github.com/absaoss/k8gb/tree/v0.7.0) (2020-10-28)

[Full Changelog](https://github.com/absaoss/k8gb/compare/v0.6.6...v0.7.0)

**Implemented enhancements:**

- Upgrade to operator-sdk 1.0 [\#166](https://github.com/AbsaOSS/k8gb/issues/166)
- Route53 support [\#162](https://github.com/AbsaOSS/k8gb/issues/162)
- Move the rest of configuration into depresolver [\#122](https://github.com/AbsaOSS/k8gb/issues/122)
- Recent gosec fails on generated deep copy code [\#115](https://github.com/AbsaOSS/k8gb/issues/115)

**Closed issues:**

- refactor controller\_tests [\#136](https://github.com/AbsaOSS/k8gb/issues/136)
- Document internal components of k8gb [\#89](https://github.com/AbsaOSS/k8gb/issues/89)

**Merged pull requests:**

- Provide diagram of k8gb internal components [\#186](https://github.com/AbsaOSS/k8gb/pull/186) ([ytsarev](https://github.com/ytsarev))
- Finalize Gslb if no route53 DNSEndpoint found [\#184](https://github.com/AbsaOSS/k8gb/pull/184) ([ytsarev](https://github.com/ytsarev))
- Include GSLB dns zone into NS server names [\#183](https://github.com/AbsaOSS/k8gb/pull/183) ([ytsarev](https://github.com/ytsarev))
- Zone delegation garbage collection for Route53 [\#182](https://github.com/AbsaOSS/k8gb/pull/182) ([ytsarev](https://github.com/ytsarev))
- Extend with fake environment variables [\#181](https://github.com/AbsaOSS/k8gb/pull/181) ([kuritka](https://github.com/kuritka))
- Post revamp readme fixes [\#180](https://github.com/AbsaOSS/k8gb/pull/180) ([ytsarev](https://github.com/ytsarev))
- Readme revamp and Route53 tutorial [\#179](https://github.com/AbsaOSS/k8gb/pull/179) ([ytsarev](https://github.com/ytsarev))
- Remove redundant route53.domain from values [\#178](https://github.com/AbsaOSS/k8gb/pull/178) ([ytsarev](https://github.com/ytsarev))
- Simplify values.yaml [\#177](https://github.com/AbsaOSS/k8gb/pull/177) ([ytsarev](https://github.com/ytsarev))
- Isolate controller tests [\#176](https://github.com/AbsaOSS/k8gb/pull/176) ([kuritka](https://github.com/kuritka))
- gosec; ignore generated code [\#174](https://github.com/AbsaOSS/k8gb/pull/174) ([kuritka](https://github.com/kuritka))
- Extending DepResolver [\#173](https://github.com/AbsaOSS/k8gb/pull/173) ([kuritka](https://github.com/kuritka))
- Route53 support [\#172](https://github.com/AbsaOSS/k8gb/pull/172) ([ytsarev](https://github.com/ytsarev))
- Fix external-dns SA definition [\#171](https://github.com/AbsaOSS/k8gb/pull/171) ([ytsarev](https://github.com/ytsarev))
- Initial configuration layout for Route53 support [\#169](https://github.com/AbsaOSS/k8gb/pull/169) ([ytsarev](https://github.com/ytsarev))

## [v0.6.6](https://github.com/absaoss/k8gb/tree/v0.6.6) (2020-10-05)

[Full Changelog](https://github.com/absaoss/k8gb/compare/v0.6.5...v0.6.6)

**Closed issues:**

- Rework README to focus on first time users [\#101](https://github.com/AbsaOSS/k8gb/issues/101)

**Merged pull requests:**

- Upgrade to operator-sdk 1.0 [\#167](https://github.com/AbsaOSS/k8gb/pull/167) ([ytsarev](https://github.com/ytsarev))
- Switch back to upstream etcd-operator chart [\#163](https://github.com/AbsaOSS/k8gb/pull/163) ([ytsarev](https://github.com/ytsarev))

## [v0.6.5](https://github.com/absaoss/k8gb/tree/v0.6.5) (2020-08-03)

[Full Changelog](https://github.com/absaoss/k8gb/compare/v0.6.3...v0.6.5)

**Implemented enhancements:**

- Report on dnsZone and Gslb Ingress host mismatch [\#149](https://github.com/AbsaOSS/k8gb/issues/149)

**Merged pull requests:**

- Fix log message about gslb failover strategy execution [\#161](https://github.com/AbsaOSS/k8gb/pull/161) ([somaritane](https://github.com/somaritane))
- Add ability to override k8gb image tag [\#160](https://github.com/AbsaOSS/k8gb/pull/160) ([somaritane](https://github.com/somaritane))
- Detect mismatch of Ingress hostname and EdgeDNSZone [\#159](https://github.com/AbsaOSS/k8gb/pull/159) ([ytsarev](https://github.com/ytsarev))
- Mitigate coredns etcd plugin bug [\#158](https://github.com/AbsaOSS/k8gb/pull/158) ([ytsarev](https://github.com/ytsarev))
- Hopefully very last rebranding bit - diagrams [\#157](https://github.com/AbsaOSS/k8gb/pull/157) ([ytsarev](https://github.com/ytsarev))
- Last missing rebranding due to the spaces [\#156](https://github.com/AbsaOSS/k8gb/pull/156) ([ytsarev](https://github.com/ytsarev))
- Fix local failover example deploy, demo image and demo targets [\#155](https://github.com/AbsaOSS/k8gb/pull/155) ([ytsarev](https://github.com/ytsarev))
- fixed wapi credientials and namespace creation [\#153](https://github.com/AbsaOSS/k8gb/pull/153) ([jeffhelps](https://github.com/jeffhelps))
- Fix ingress nginx failure in local env and pipelines [\#152](https://github.com/AbsaOSS/k8gb/pull/152) ([ytsarev](https://github.com/ytsarev))
- Fix code markup in the readme [\#150](https://github.com/AbsaOSS/k8gb/pull/150) ([ytsarev](https://github.com/ytsarev))
- Remove unnecessary infoblox variables from the guide [\#148](https://github.com/AbsaOSS/k8gb/pull/148) ([ytsarev](https://github.com/ytsarev))
- An attempt to create step-by-step howto [\#146](https://github.com/AbsaOSS/k8gb/pull/146) ([ytsarev](https://github.com/ytsarev))
- Update demo application version [\#145](https://github.com/AbsaOSS/k8gb/pull/145) ([ytsarev](https://github.com/ytsarev))
- Increase test app installation timeout [\#143](https://github.com/AbsaOSS/k8gb/pull/143) ([ytsarev](https://github.com/ytsarev))
- Switch back to upstream releases [\#142](https://github.com/AbsaOSS/k8gb/pull/142) ([ytsarev](https://github.com/ytsarev))

## [v0.6.3](https://github.com/absaoss/k8gb/tree/v0.6.3) (2020-06-11)

[Full Changelog](https://github.com/absaoss/k8gb/compare/v0.6.2...v0.6.3)

**Implemented enhancements:**

- Make project lintable from project root [\#131](https://github.com/AbsaOSS/k8gb/issues/131)

**Merged pull requests:**

- Document currently tested configuration [\#140](https://github.com/AbsaOSS/k8gb/pull/140) ([ytsarev](https://github.com/ytsarev))
- Mass rebranding to K8GB [\#139](https://github.com/AbsaOSS/k8gb/pull/139) ([ytsarev](https://github.com/ytsarev))
- Mass rebranding to KGB [\#137](https://github.com/AbsaOSS/k8gb/pull/137) ([ytsarev](https://github.com/ytsarev))
- Switch to safe geotag propagation with depresolver [\#135](https://github.com/AbsaOSS/k8gb/pull/135) ([ytsarev](https://github.com/ytsarev))
- Ability to override registry image [\#133](https://github.com/AbsaOSS/k8gb/pull/133) ([ytsarev](https://github.com/ytsarev))
- Make project lintable from project's root [\#132](https://github.com/AbsaOSS/k8gb/pull/132) ([kuritka](https://github.com/kuritka))

## [v0.6.2](https://github.com/absaoss/k8gb/tree/v0.6.2) (2020-05-20)

[Full Changelog](https://github.com/absaoss/k8gb/compare/v0.6.0...v0.6.2)

**Merged pull requests:**

- Fix helm installation smoke test [\#130](https://github.com/AbsaOSS/k8gb/pull/130) ([ytsarev](https://github.com/ytsarev))
- Fix issues with public release [\#128](https://github.com/AbsaOSS/k8gb/pull/128) ([ytsarev](https://github.com/ytsarev))
- Release 0.6.1 [\#127](https://github.com/AbsaOSS/k8gb/pull/127) ([ytsarev](https://github.com/ytsarev))
- Simplify versioning process [\#126](https://github.com/AbsaOSS/k8gb/pull/126) ([ytsarev](https://github.com/ytsarev))

## [v0.6.0](https://github.com/absaoss/k8gb/tree/v0.6.0) (2020-05-16)

[Full Changelog](https://github.com/absaoss/k8gb/compare/v0.5.6...v0.6.0)

**Implemented enhancements:**

- Streamline Gslb Status [\#116](https://github.com/AbsaOSS/k8gb/issues/116)
- Propagate Gslb CR annotations down to Gslb ingress [\#113](https://github.com/AbsaOSS/k8gb/issues/113)
- Make Gslb timeouts and synchronisation intervals configurable [\#82](https://github.com/AbsaOSS/k8gb/issues/82)
- Prepare Helm chart for uploading various repositories [\#75](https://github.com/AbsaOSS/k8gb/issues/75)
- Extend documentation with end-to-end application deployment scenario [\#69](https://github.com/AbsaOSS/k8gb/issues/69)
- Add full end to end integration tests to build pipeline [\#48](https://github.com/AbsaOSS/k8gb/issues/48)
- Expose metrics and tracing [\#47](https://github.com/AbsaOSS/k8gb/issues/47)

**Fixed bugs:**

- Non-deterministic failure of EtcdCluster deployment in air-gapped on-prem environments [\#107](https://github.com/AbsaOSS/k8gb/issues/107)
- Flaky terrarest `TestOhmyglbBasicAppExample` [\#105](https://github.com/AbsaOSS/k8gb/issues/105)

**Closed issues:**

- Can't install chart successfully [\#104](https://github.com/AbsaOSS/k8gb/issues/104)

**Merged pull requests:**

- Extend release pipeline with docker build and push [\#125](https://github.com/AbsaOSS/k8gb/pull/125) ([ytsarev](https://github.com/ytsarev))
- Streamline Gslb Status [\#121](https://github.com/AbsaOSS/k8gb/pull/121) ([ytsarev](https://github.com/ytsarev))
- Extend `deploy-gslb-cr` target with failover strategy [\#118](https://github.com/AbsaOSS/k8gb/pull/118) ([ytsarev](https://github.com/ytsarev))
- Configurable timeouts and synchronisation intervals [\#117](https://github.com/AbsaOSS/k8gb/pull/117) ([kuritka](https://github.com/kuritka))
- Propagate Gslb CR annotations down to Gslb ingress [\#114](https://github.com/AbsaOSS/k8gb/pull/114) ([ytsarev](https://github.com/ytsarev))
- Properly propagate etcd version in EtcdCluster CR [\#112](https://github.com/AbsaOSS/k8gb/pull/112) ([ytsarev](https://github.com/ytsarev))
- Make basic app terratest reliable [\#111](https://github.com/AbsaOSS/k8gb/pull/111) ([ytsarev](https://github.com/ytsarev))
- Optimize and cleanup test-apps target and samples [\#110](https://github.com/AbsaOSS/k8gb/pull/110) ([ytsarev](https://github.com/ytsarev))
- Optimize CI status badges [\#106](https://github.com/AbsaOSS/k8gb/pull/106) ([ytsarev](https://github.com/ytsarev))
- Failover demo [\#103](https://github.com/AbsaOSS/k8gb/pull/103) ([kuritka](https://github.com/kuritka))
- Non deterministic round robin demo [\#98](https://github.com/AbsaOSS/k8gb/pull/98) ([kuritka](https://github.com/kuritka))
- Initial operator metrics [\#97](https://github.com/AbsaOSS/k8gb/pull/97) ([somaritane](https://github.com/somaritane))
- Add capability to end-to-end test HEAD of the branch [\#96](https://github.com/AbsaOSS/k8gb/pull/96) ([ytsarev](https://github.com/ytsarev))
- Enhance terratest pipeline [\#95](https://github.com/AbsaOSS/k8gb/pull/95) ([ytsarev](https://github.com/ytsarev))
- Etcd-operator as own subchart [\#94](https://github.com/AbsaOSS/k8gb/pull/94) ([ytsarev](https://github.com/ytsarev))
- Include gosec into pipeline [\#93](https://github.com/AbsaOSS/k8gb/pull/93) ([ytsarev](https://github.com/ytsarev))
- Terratest based end-to-end pipeline  [\#91](https://github.com/AbsaOSS/k8gb/pull/91) ([ytsarev](https://github.com/ytsarev))
- Document Helm repo and installation [\#88](https://github.com/AbsaOSS/k8gb/pull/88) ([ytsarev](https://github.com/ytsarev))
- How to run Oh My GLB locally [\#87](https://github.com/AbsaOSS/k8gb/pull/87) ([kuritka](https://github.com/kuritka))

## [v0.5.6](https://github.com/absaoss/k8gb/tree/v0.5.6) (2020-04-14)

[Full Changelog](https://github.com/absaoss/k8gb/compare/v0.5.1...v0.5.6)

**Implemented enhancements:**

- When using the failover load balancing strategy, investigate and validate how resolution will be handled effectively when clusters are configured for mutual failover [\#67](https://github.com/AbsaOSS/k8gb/issues/67)
- TTL control for splitbrain TXT record [\#61](https://github.com/AbsaOSS/k8gb/issues/61)
- Implement failover load balancing strategy [\#46](https://github.com/AbsaOSS/k8gb/issues/46)
- Posssible Routing Peering Capabilities BGP protocols [\#33](https://github.com/AbsaOSS/k8gb/issues/33)

**Fixed bugs:**

- Missing endpoints in `localtargets.\*` A records [\#62](https://github.com/AbsaOSS/k8gb/issues/62)
- Non-deterministic issue with `localtargets.\*` DNSEntrypoint population [\#38](https://github.com/AbsaOSS/k8gb/issues/38)

**Closed issues:**

- Upgrade underlying operator-sdk version from v0.12.0 to latest upstream [\#71](https://github.com/AbsaOSS/k8gb/issues/71)
- High Five [\#41](https://github.com/AbsaOSS/k8gb/issues/41)

**Merged pull requests:**

- Helm package and publish on release event [\#86](https://github.com/AbsaOSS/k8gb/pull/86) ([ytsarev](https://github.com/ytsarev))
- test upgraded build pipe [\#85](https://github.com/AbsaOSS/k8gb/pull/85) ([kuritka](https://github.com/kuritka))
- Test mutual failover setup [\#84](https://github.com/AbsaOSS/k8gb/pull/84) ([ytsarev](https://github.com/ytsarev))
- Upgrade operator sdk to v0.16.0 [\#83](https://github.com/AbsaOSS/k8gb/pull/83) ([somaritane](https://github.com/somaritane))
- Reduce external-dns sync interval to 20s [\#81](https://github.com/AbsaOSS/k8gb/pull/81) ([ytsarev](https://github.com/ytsarev))
- Time measure failover process [\#80](https://github.com/AbsaOSS/k8gb/pull/80) ([ytsarev](https://github.com/ytsarev))
- Terratest e2e for Failover strategy [\#79](https://github.com/AbsaOSS/k8gb/pull/79) ([ytsarev](https://github.com/ytsarev))
- Fix cluster namespaces permission for ohmyglb [\#77](https://github.com/AbsaOSS/k8gb/pull/77) ([somaritane](https://github.com/somaritane))
- Terratest for standard ohmyglb deployment with app [\#76](https://github.com/AbsaOSS/k8gb/pull/76) ([ytsarev](https://github.com/ytsarev))
- Terratest e2e testing proposal [\#74](https://github.com/AbsaOSS/k8gb/pull/74) ([ytsarev](https://github.com/ytsarev))
- Expose all namespaces in ServeCRMetrics [\#73](https://github.com/AbsaOSS/k8gb/pull/73) ([ytsarev](https://github.com/ytsarev))
- Fix docker repo link for external-dns [\#72](https://github.com/AbsaOSS/k8gb/pull/72) ([ytsarev](https://github.com/ytsarev))
- Bump to include external-dns image with the bugfix [\#70](https://github.com/AbsaOSS/k8gb/pull/70) ([ytsarev](https://github.com/ytsarev))
- Use custom build of external-dns with multi A fixes [\#68](https://github.com/AbsaOSS/k8gb/pull/68) ([ytsarev](https://github.com/ytsarev))
- Failover strategy post e2e stabilization [\#66](https://github.com/AbsaOSS/k8gb/pull/66) ([ytsarev](https://github.com/ytsarev))
- Failover strategy implementation [\#65](https://github.com/AbsaOSS/k8gb/pull/65) ([ytsarev](https://github.com/ytsarev))
- Set low TTL on split brain TXT record via infoblox API [\#64](https://github.com/AbsaOSS/k8gb/pull/64) ([ytsarev](https://github.com/ytsarev))
- Fully automated multicluster ohmyglb local deployment [\#63](https://github.com/AbsaOSS/k8gb/pull/63) ([ytsarev](https://github.com/ytsarev))
- Splitbrain enhancements and fixes [\#60](https://github.com/AbsaOSS/k8gb/pull/60) ([ytsarev](https://github.com/ytsarev))
- Bump to 5.3 to stabilize split brain handling [\#59](https://github.com/AbsaOSS/k8gb/pull/59) ([ytsarev](https://github.com/ytsarev))
- Infoblox update [\#58](https://github.com/AbsaOSS/k8gb/pull/58) ([ytsarev](https://github.com/ytsarev))
- Splitbrain fixes [\#57](https://github.com/AbsaOSS/k8gb/pull/57) ([ytsarev](https://github.com/ytsarev))
- Config and helpers for local multicluster setup [\#56](https://github.com/AbsaOSS/k8gb/pull/56) ([ytsarev](https://github.com/ytsarev))
- Move to `absaoss` in dockerhub and version bump [\#55](https://github.com/AbsaOSS/k8gb/pull/55) ([ytsarev](https://github.com/ytsarev))
- Split brain handling [\#44](https://github.com/AbsaOSS/k8gb/pull/44) ([ytsarev](https://github.com/ytsarev))
- Disable `external-dns` ownership for local coredns [\#43](https://github.com/AbsaOSS/k8gb/pull/43) ([ytsarev](https://github.com/ytsarev))
- Quote geo tag declaration [\#42](https://github.com/AbsaOSS/k8gb/pull/42) ([ytsarev](https://github.com/ytsarev))

## [v0.5.1](https://github.com/absaoss/k8gb/tree/v0.5.1) (2020-02-02)

[Full Changelog](https://github.com/absaoss/k8gb/compare/d834431a8236e7bbe2769df41bc0e02ceb5afeb3...v0.5.1)

**Merged pull requests:**

- CRUD gslb zone delegation in infoblox [\#39](https://github.com/AbsaOSS/k8gb/pull/39) ([ytsarev](https://github.com/ytsarev))
- Multi node local kind cluster [\#37](https://github.com/AbsaOSS/k8gb/pull/37) ([ytsarev](https://github.com/ytsarev))
- Initial Edge DNS support  [\#36](https://github.com/AbsaOSS/k8gb/pull/36) ([ytsarev](https://github.com/ytsarev))
- Use `podinfo` as example test app [\#35](https://github.com/AbsaOSS/k8gb/pull/35) ([ytsarev](https://github.com/ytsarev))
- Enable periodic reconciliation [\#34](https://github.com/AbsaOSS/k8gb/pull/34) ([ytsarev](https://github.com/ytsarev))
- External dns ownership fix [\#32](https://github.com/AbsaOSS/k8gb/pull/32) ([ytsarev](https://github.com/ytsarev))
- Tolerate external Gslb downtime [\#31](https://github.com/AbsaOSS/k8gb/pull/31) ([ytsarev](https://github.com/ytsarev))
- DNS based cross Gslb communication [\#30](https://github.com/AbsaOSS/k8gb/pull/30) ([ytsarev](https://github.com/ytsarev))
- BUGFIX: populate record status only when it's ready [\#29](https://github.com/AbsaOSS/k8gb/pull/29) ([ytsarev](https://github.com/ytsarev))
- Expose DNS records for heatlhy hosts in Gslb Status [\#28](https://github.com/AbsaOSS/k8gb/pull/28) ([ytsarev](https://github.com/ytsarev))
- Change example domain to `example.com` [\#27](https://github.com/AbsaOSS/k8gb/pull/27) ([ytsarev](https://github.com/ytsarev))
- Ohmyglb operator chart [\#26](https://github.com/AbsaOSS/k8gb/pull/26) ([ytsarev](https://github.com/ytsarev))
- Simple push/build helpers [\#25](https://github.com/AbsaOSS/k8gb/pull/25) ([ytsarev](https://github.com/ytsarev))
- Expose coredns\(53 udp\) with nginx ingress controller [\#24](https://github.com/AbsaOSS/k8gb/pull/24) ([ytsarev](https://github.com/ytsarev))
- Enhancements to local test configuration [\#23](https://github.com/AbsaOSS/k8gb/pull/23) ([ytsarev](https://github.com/ytsarev))
- E2e test suite extension and optimization [\#22](https://github.com/AbsaOSS/k8gb/pull/22) ([ytsarev](https://github.com/ytsarev))
- e2e tests for Gslb creation [\#21](https://github.com/AbsaOSS/k8gb/pull/21) ([ytsarev](https://github.com/ytsarev))
- Foundation for e2e tests [\#20](https://github.com/AbsaOSS/k8gb/pull/20) ([ytsarev](https://github.com/ytsarev))
- Deprecate coreDNS hosts config and worker health checks [\#19](https://github.com/AbsaOSS/k8gb/pull/19) ([ytsarev](https://github.com/ytsarev))
- Switch source of addresses for A records to Ingress [\#18](https://github.com/AbsaOSS/k8gb/pull/18) ([ytsarev](https://github.com/ytsarev))
- Dynamically populate DNSEndpoints according to health status  [\#17](https://github.com/AbsaOSS/k8gb/pull/17) ([ytsarev](https://github.com/ytsarev))
- Register and watch for DNSEndpoints [\#16](https://github.com/AbsaOSS/k8gb/pull/16) ([ytsarev](https://github.com/ytsarev))
- Foundation for external-dns DNSEndpoint creation [\#15](https://github.com/AbsaOSS/k8gb/pull/15) ([ytsarev](https://github.com/ytsarev))
- Prototype of external-dns + coredns based configuration [\#14](https://github.com/AbsaOSS/k8gb/pull/14) ([ytsarev](https://github.com/ytsarev))
- Make OhMyGlb operator watch all namespaces for Gslb CRs [\#13](https://github.com/AbsaOSS/k8gb/pull/13) ([ytsarev](https://github.com/ytsarev))
- Add some badges [\#12](https://github.com/AbsaOSS/k8gb/pull/12) ([ytsarev](https://github.com/ytsarev))
- Reconcile Gslb when relevant Endpoint is updated [\#11](https://github.com/AbsaOSS/k8gb/pull/11) ([ytsarev](https://github.com/ytsarev))
- Enable golint in the pipeline, fix code accordingly [\#10](https://github.com/AbsaOSS/k8gb/pull/10) ([ytsarev](https://github.com/ytsarev))
- Control coredns hosts config map [\#9](https://github.com/AbsaOSS/k8gb/pull/9) ([ytsarev](https://github.com/ytsarev))
- Expose healthy workers and their ip addresses [\#8](https://github.com/AbsaOSS/k8gb/pull/8) ([ytsarev](https://github.com/ytsarev))
- Install CoreDNS from stable chart with custom values [\#7](https://github.com/AbsaOSS/k8gb/pull/7) ([ytsarev](https://github.com/ytsarev))
- Gslb Controller Unit Tests [\#6](https://github.com/AbsaOSS/k8gb/pull/6) ([ytsarev](https://github.com/ytsarev))
- Gslb Ingress management and associated health checks [\#5](https://github.com/AbsaOSS/k8gb/pull/5) ([ytsarev](https://github.com/ytsarev))
- \[WIP\] First iteration of ohmyglb operator [\#3](https://github.com/AbsaOSS/k8gb/pull/3) ([ytsarev](https://github.com/ytsarev))
- Additional doc links [\#2](https://github.com/AbsaOSS/k8gb/pull/2) ([ytsarev](https://github.com/ytsarev))
- Take readiness probes into account [\#1](https://github.com/AbsaOSS/k8gb/pull/1) ([ytsarev](https://github.com/ytsarev))



\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
