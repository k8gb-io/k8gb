<h1 align="center" style="margin-top: 0;">Using K8GB on Azure</h1>

## Sample solution

The provided lab sample solution will create a simple hub and spoke architecture with two AKS clusters in different regions

## Technical decisions

* Azure Private DNS Zones was discarded since they don't allow the creation of NS records
    * Wouldn't be possible to delegate zone into CoreDNS
* Azure DNS was discarded because our DNS is internal only
* Windows DNS was our choice, since we relly on Active Directory for historical reasons in our environments
    * In order to setup DNS dynamic updates, External DNS should be configured to use the GSS-TSIG protocol for Kerberos authentication
    * The Helm template of K8GB needed to be changed, since it only supported TSIG configuration

## Run the sample

* To run the provided sample, please use the available Makefile

## Configure GSS-TSIG authentication

* Ensure that the Network Security is configured only for AES256
* Ensure that the DNS Zone has only Secure updates option enabled
* Create a new Active Directory user
    * The user should be created with "Encryptions options" for Kerberos AES256 encryption
    * The user needs to be added to the DNSAdmin group, or,
    * Select the zone that will have dynamic updates in DNS Manager, right click and select Properties. Under the Security tab, add the created user and add the permissions Write, Create all child objects and Delete all child objects
* ExternalDNS configuration
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
