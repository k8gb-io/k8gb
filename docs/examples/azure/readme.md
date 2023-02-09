<h1 align="center" style="margin-top: 0;">Using K8GB on Azure</h1>

## Sample solution

The provided lab sample solution will create a simple hub and spoke architecture.

## Technical decisions

* Azure Private DNS Zones was discarded since they don't allow the creation of NS records
* Azure DNS was discarded because our DNS is internal only
* Windows DNS was our choice, since we relly on Active Directory for historical reasons
    * In order to setup DNS dynamic updates, External DNS should be configured to use the GSS-TSIG protocol for Kerberos authentication
    * The Helm template of K8GB needed to be changed, since it only supported TSIG configuration

## Run the sample

## Configure GSS-TSIG authentication

* Ensure that the Network Security is configured only for AES256
* Ensure that the DNS Zone has only Secure updates option enabled
* Create a new Active Directory user
    * The user should be created Encryptions options for Kerberos AES256 encryption
    * The user needs to be added to the DNSAdmin group, or,
    * Select the zone that will have dynamic updates in DNS Manager, right click and select Properties. Under the Security tab, add the created user and add the permissions Write, Create all child objects and Delete all child objects
* ExternalDNS configuration
```
 rfc2136:
  enabled: true
  rfc2136Opts:
    - host: dc.domain.com
    - port: 53
    - gss-tsig-user: username
    - gss-tsig-password: password
    - gss-tsig-kerberos-realm: kerberosRealm
```
    * The hostname needs to be the FQDN of the DC

