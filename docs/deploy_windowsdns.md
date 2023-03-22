# Azure based deployment with Windows DNS integration

Here we provide an example of k8gb deployment in Azure context with Windows DNS as edgeDNS provider.

## Reference Setup

The reference setup includes two AKS clusters that can deployed on two different regions for load balancing or provide a failover solution.

The solution design can be found [here](https://github.com/k8gb-io/k8gb/tree/master/docs/examples/azure/).

Configurable resources:

* Resource groups
* VNet and subnets
* Peering
* Managed Identity
* Clusters

## Run the sample

* This lab requires a running AD Domain Controller with DNS and KDC services working
    * There are several tutorials available online, but this Microsoft Learn article will probably help you out 
    * [Microsoft Learn](https://learn.microsoft.com/en-us/windows-server/identity/ad-ds/deploy/install-active-directory-domain-services--level-100- "Install Active Directory")

* To run the provided sample, please use the provided Makefile [here](https://github.com/k8gb-io/k8gb/tree/master/docs/examples/azure/).
    * Deploys all the required infrastructure and configurations
    * Before execute, please fill all the local variables in the scripts with the correct naming for the resources in order to avoid having problems with your azure policies
    * Scripts will use AZ CLI, ensure that is installed and logged when trying to execute the command
        * [Microsoft Learn](https://learn.microsoft.com/en-us/cli/azure/install-azure-cli "Install AZ CLI")

### Deploy infrastructure 

This action will create resource groups, vnets, peering between vnets and AKS clusters to run all required workloads

```sh
make deploy-infra
```

### Setup clusters

Install required Ingress controllers in both clusters

```sh
make setup-clusters
```

### Configure GSS-TSIG authentication for DNS updates

Before deploying K8GB and the demo workload, ensure required configurations on Windows DNS 

#### Domain Controller config

* Ensure that the Network Security is configured only for AES256

![Network Policy - Kerberos auth](/docs/examples/azure/images/LocalSecuryPolicyNetworkKerberos.png?raw=true "Network Policy - Kerberos auth")
* Ensure that the DNS Zone has only Secure updates option enabled

![DNS Secure Updates](/docs/examples/azure/images/DNSSecureUpdates.png "DNS Secure Updates")
* Ensure that the DNS Zone has the option "Allow zone transfers" check with the option "To any server" under the tab Zone Transfers on the zone properties

![DNS Zone Transfers](/docs/examples/azure/images/DNSZoneTransfers.png "DNS Zone Transfers")

* Create a new Active Directory user
    * The user should be created with "Encryptions options" for Kerberos AES256 encryption
    * The user needs to be added to the DNSAdmin group, or,
    * Select the zone that will have dynamic updates in DNS Manager, right click and select Properties. Under the Security tab, add the created user and add the permissions Write, Create all child objects and Delete all child objects

#### K8GB / ExternalDNS configuration

* ExternalDNS configuration
    * For communication with WindowsDNS, ExternalDNS should be configured with the RFC2136 provider with GSS-TSIG option
    * [External DNS - RFC2126](https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/rfc2136.md "RFC2136 documentation")
    * A sample values.yaml for K8GB configuration can be found [here](https://github.com/k8gb-io/k8gb/tree/master/docs/examples/azure/k8gb/).
        * Ensure that the following properties are updated with your values:
            * dnsZone
            * edgeDNSZone
            * edgeDNSServers
            * host - always use FQDN with GSS-TSIG, not IP address
            * kerberos-username
            * kerberos-password
            * kerberos-realm
    * At this moment ExternalDNS doesn't provide a way to use secrets as the source for the kerberos-password setting, so if you store the manifest in a git repo, please ensure that only required persons can access it
```yaml
rfc2136:
  enabled: true
  rfc2136Opts:
    - host: AD-DC.mbcpk8gb.local #when using gssTsig, use the FQDN of the host, not an IP
    - port: 53
  rfc2136auth:
    insecure: 
      enabled: false
    tsig:
      enabled: false
      tsigCreds:
        - tsig-secret-alg: hmac-sha256
        - tsig-keyname: externaldns-key
    gssTsig:
      enabled: true
      gssTsigCreds:
        - kerberos-username: ad-user-account
        - kerberos-password: ad-user-account-password
        - kerberos-realm: mbcpcloud.lab
```

### Install K8gb

This action will install K8gb in both clusters using the provided [sample](https://github.com/k8gb-io/k8gb/tree/master/docs/examples/azure/k8gb/) values.yaml for each cluster. Please ensure that the are correctly updated before execution

```sh
make deploy-k8gb
```

### Install demo app

Deploys the sample Podinfo workload with failover GLSB configured using annotations in the Ingress resource [samples](https://github.com/k8gb-io/k8gb/tree/master/docs/examples/azure/demo/).
Ensure that the hosts on the samples are correctly updated.

```sh
make deploy-demo
```

### Destroy lab
* Destroys the lab environment created for this sample
```sh
make destroy-infra
```
