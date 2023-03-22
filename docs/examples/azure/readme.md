<h1 align="center" style="margin-top: 0;">Using K8GB on Azure</h1>

## Sample solution

The provided lab sample solution will create a simple hub and spoke architecture with two AKS clusters in different regions

![GLSB with K8gb on Azure](/docs/examples/azure/images/k8gb_solution.png?raw=true "GLSB with K8gb on Azure")

## Technical decisions

* Azure Private DNS Zones was discarded since they don't allow the creation of NS records
    * Wouldn't be possible to delegate zone into CoreDNS
* Azure DNS was discarded because our DNS is internal only
* Windows DNS was our choice, since we relly on Active Directory for historical reasons in our environments
    * In order to setup DNS dynamic updates, External DNS should be configured to use the GSS-TSIG protocol for Kerberos authentication
    * The Helm template of K8GB needed to be changed, since it only supported TSIG configuration


 