# Enabling RFC2136 for ExternalDNS

In order to enable the provider RFC2136 on ExternalDNS, the following `rfc2136.*` [parameters](https://github.com/k8gb-io/k8gb/blob/master/chart/k8gb/README.md#values) should be changed in the values.yaml of the K8GB helm chart:

* One authentication method should be enabled on the values:
  * Insecure
    * This method doesn't use any authentication and anonymous updates to the DNS records can be executed
  * TSIG
    * This method uses TSIG authentication that relies on a token provided for the DNS records update.
  * GSS-TSIG
    * This method uses GSS-TSIG authentication, which is a variation of the TSIG method, but uses Kerberos for the generation of tokens for authentication and authorization
    * Method used by Active Directory Windows DNS

* GSS-TSIG
  * kerberos-username
    * this key should have the value of a Active Directory user account that has permissions for DNS updates
  * kerberos-password
    * password of the user account that will be used. Be aware that this isn't encrypted and so far ExternalDNS doesn't support adding a Secret reference for this value, so it will be stored in clear text
  * kerberos-realm
    * domain that will be used for authentication of the user

# Sample for GSS-TSIG authentication

```yaml
extdns:
  enabled: true
  fullnameOverride: "k8gb-external-dns"
  provider:
    name: rfc2136
  txtPrefix: "k8gb-<geotag>-"
  txtOwnerId: "k8gb-<loadBalancedZone>-<geotag>"
  domainFilters:
    - "<parentZone>"
  dnsPolicy: ClusterFirst
  env:
  - name: EXTERNAL_DNS_RFC2136_TSIG_SECRET
    valueFrom:
      secretKeyRef:
        name: rfc2136
        key: secret
  extraArgs:
    rfc2136-host: yourAcDc.k8gb.local
    rfc2136-port: 53
    rfc2136-gss-tsig:
    rfc2136-kerberos-username: someServiceAccount
    rfc2136-kerberos-password: insecurePlainTextPassword
    rfc2136-kerberos-realm: yourKerberosRealm.domain
```
