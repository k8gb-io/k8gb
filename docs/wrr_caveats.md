# Weighted Round Robin caveats
k8gb successfully implements **weighted round robin (WRR)** at the level of authoritative CoreDNS servers. The correct weighted records are published by each cluster’s CoreDNS instance as the authoritative DNS server.

However, in most production setups, DNS queries are first answered by a parent DNS (sometimes referred to as a delegation DNS), which acts as an upper-level authoritative server delegating queries to cluster CoreDNS servers. The actual order and frequency in which DNS answers are returned to clients depends on this delegation DNS (e.g., BIND, Infoblox, Unbound, etc.), which may also act as a caching or forwarding layer.

Many delegation (parent) DNS servers will randomize or reorder A/NS records, or apply their own internal load-balancing or optimizations (such as sorting by RTT). As a result, the expected traffic split according to weights may not be reflected on the client side, even though k8gb and CoreDNS publish the correct weighted records. **This is out of the scope of K8gb**.

If you need strict control over record ordering or weight distribution, you have these options:

 - Review and configure your parent DNS server to ensure it does not shuffle or reorder answers, or that it supports weighted policies appropriately (see your DNS server documentation for options).

 - Bypass parent DNS and integrate directly with your custom application logic to implement advanced load-balancing based on your needs.

For more information and real-world examples, see the following k8gb GitHub issues:

 - [Weighted round robin DNS limitations #1950](https://github.com/k8gb-io/k8gb/issues/1950)
 - [Further discussion on WRR behavior #1943](https://github.com/k8gb-io/k8gb/issues/1943)
