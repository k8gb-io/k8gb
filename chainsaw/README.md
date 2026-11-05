# Useful commands

Run a single test:

`chainsaw test --config ./config.yaml --values ./values.yaml --test-dir ./tests/<testname>`

`make chainsaw` saves the console log and a JUnit step report in the repository's
`tmp/chainsaw/` directory.
On a failed step, the global catch handler uses Chainsaw's native `get`, `events`,
and `podLogs` collectors to write resource YAML, events, and controller logs from
both clusters to `chainsaw.log` before test cleanup. Resource YAML includes owner
references, deletion timestamps, and finalizers. Collection runs without enabling
runner debug logging. The CI workflow uploads this directory.

Repeat both legacy migration tests (from this directory) to check for timing issues:

```sh
chainsaw test --config ./config.yaml --values ./values.yaml \
  --test-dir ./tests/legacy-migration-dnsendpoint-handoff \
  --test-dir ./tests/legacy-migration-requested-pauses-runtime \
  --parallel 2 --repeat-count 20
```
