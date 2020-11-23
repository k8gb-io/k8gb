# Ingress Annotations

Instead of direct Gslb resource creation there is ability to enable global load balancing
by setting annotations on the standard Ingress objects.

| Annotation            | Possible Values        |
|-----------------------|------------------------|
| k8gb.io/strategy      | roundRobin \| failover |
| k8gb.io/primarygeotag | arbitrary tag, e.g. eu |
