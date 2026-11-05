# Useful commands

Run a single test:

`chainsaw test --config ./config.yaml --values ./values.yaml --test-dir ./tests/<testname>`

# When a test fails

Before cleanup, `config.yaml` saves the test namespace's Gslbs, DNSEndpoints, Ingresses and events, the ZoneDelegations, and the k8gb controller logs from both clusters. In CI these files are written to `tmp/chainsaw/<namespace>/<cluster>/` and uploaded as the `chainsaw-logs-*` artifact. Locally they are printed to the test output, unless you set `CHAINSAW_ARTIFACTS_DIR` to a directory.
