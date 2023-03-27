# Enabling RFC2136 for ExternalDNS

In order to enable the provider RFC2136 on ExternalDNS, the following values should be changed in the values.yaml for the K8GB helm chart:

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| rfc2136.enabled | bool | `false` |  |
| rfc2136.rfc2136Opts[0].host | string | `"host.k3d.internal"` |  |
| rfc2136.rfc2136Opts[1].port | int | `1053` |  |
| rfc2136.rfc2136auth.insecure.enabled | bool | `false` | Set to True if insecure updates to the DNS provided can be executed by ExternalDNS |
| rfc2136.rfc2136auth.tsig.enabled | bool | `false` | Set to True if the DNS server uses TSIG authentication for DNS updates by ExternalDNS |
| rfc2136.rfc2136auth.tsig.tsigCreds[0].tsig-secret-alg | string | `"hmac-sha256"` | Algorithm used to generate the token for TSIG |
| rfc2136.rfc2136auth.tsig.tsigCreds[1].tsig-keyname | string | `"externaldns-key"` |  |
| rfc2136.rfc2136auth.gssTsig.enabled | bool | `false` | Set to True if the DNS server uses GSS-TSIG (Kerberos) authentication for DNS updates by ExternalDNS |
| rfc2136.rfc2136auth.gssTsig.kerberosConfigMap | string | `"kerberos-configmap"` | When using GSS-TSIG, a ConfigMap with a valid krb5.conf configuration should be provided |
| rfc2136.rfc2136auth.gssTsig.gssTsigCreds[0].kerberos-username | string | `"ad-user-account"` | AD user account with permissions for DNS updates |
| rfc2136.rfc2136auth.gssTsig.gssTsigCreds[1].kerberos-password | string | `"ad-user-account-password"` | Passowrd of the AD user account |
| rfc2136.rfc2136auth.gssTsig.gssTsigCreds[2].kerberos-realm | string | `"REALM.DOMAIN"` | Kerberos REALM that should be used for authentication |

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