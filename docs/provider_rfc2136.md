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
rfc2136:
  enabled: true
  rfc2136Opts:
    - host: yourAcDc.k8gb.local
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
        - kerberos-username: someServiceAccount
        - kerberos-password: insecurePlainTextPassword
        - kerberos-realm: yourKerberosRealm.domain
```
