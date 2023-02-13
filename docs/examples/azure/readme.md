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

## Run the sample

* To run the provided sample, please use the available Makefile
```
make deploy-infra #will deploy all the required infra
```

```
make deploy-k8gb #will deploy K8gb on both clusters
```

```
make deploy-demo #will deploy the sample Podinfo workload with failover GLSB
```

```
make destroy-infra #destroys all the created
```
* This lab requires a running AD Domain Controller with DNS and KDC services working
    * There are several tutorials available online, but this Microsoft Learn link will probably help you out 
    * https://learn.microsoft.com/en-us/windows-server/identity/ad-ds/deploy/install-active-directory-domain-services--level-100-

## Configure GSS-TSIG authentication for DNS updates

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
* ExternalDNS configuration
    * A sample configuration can be found at k8gb folder
```
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
        - kerberos-realm: mbcpk8gb.local
```
 