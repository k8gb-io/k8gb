# Ingress Annotations

Instead of direct Gslb resource creation there is ability to enable global load balancing
by setting annotations on the standard Ingress objects.

| Annotation             | Description      | Type                           |
| ---------------------- | ---------------- | ------------------------------ |
| k8gb.io/strategy       | Glsb strategy    | "`roundRobin`" \| "`failover`" |
| k8gb.io/primary-geotag | Arbitrary geotag | string (e.g. "`eu`")           |
