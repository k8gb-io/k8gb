# How to contribute

k8gb is Apache 2 licensed and accepts contributions via GitHub pull requests. This document outlines some of the conventions on commit message formatting, contact points for developers, and other resources to help get contributions into k8gb project.

# Slack chat

- Find us at `#sig-multicluster` on kubernetes.slack.com

## Getting started

- Fork the repository on GitHub
- See the [local playground guide](/docs/local.md) for testing environment setup

## Reporting bugs and creating issues

Reporting bugs is one of the best ways to contribute.
Feel free to open an issue describing your problem or any question.

## Contribution flow

This is a rough outline of what a contributor's workflow looks like:

- Create a topic branch from where to base the contribution.
- Make commits of logical units.
- Make sure commit messages are in the proper format (see below).
- Push changes in a topic branch to a personal fork of the repository.
- Submit a pull request to AbsaOSS/k8gb.
- The PR must receive a LGTM from two maintainers found in the `CODEOWNERS` file.

Thanks for contributing!

### Local setup

```sh
make deploy-full-local-setup
```
Deploys k8gb from scracth including:

* 2 local clusters
* Stable k8gb helm chart
* Gslb Custom Resources examples
* Test applications

```sh
make upgrade-candidate
```
Performs upgrade of k8gb helm chart and controller
to the testing version that is built from your current
development tree.

### Testing

* Unit tests should be updated for any functional code change at [test suite location](https://github.com/AbsaOSS/k8gb/tree/master/controllers/gslb_controller_test.go).
* Acceptance terratest suite is located [here](https://github.com/AbsaOSS/k8gb/tree/master/terratest) and executable by `make terratest` target.
  These tests are updated only if the change is substantial enough to affect the main end-to-end flow.

* There is possibility to execute terratest suite over the real clusters.
For this you need to override the set of test settings as in example below.
```sh
PRIMARY_GEO_TAG=af-south-1 \
SECONDARY_GEO_TAG=eu-west-1 \
DNS_SERVER1=a377095726f1845fb85b95c2afef8ac0-9a1a10f24e634e28.elb.af-south-1.amazonaws.com \
DNS_SERVER1_PORT=53 \
DNS_SERVER2=a873f5c83be624a0a84c05a743d598a8-443f7e0285e4a28f.elb.eu-west-1.amazonaws.com \
DNS_SERVER2_PORT=53 \
GSLB_DOMAIN=test.k8gb.io \
K8GB_CLUSTER1=arn:aws:eks:af-south-1:<aws-account-id>:cluster/k8gb-cluster-af-south-1 \
K8GB_CLUSTER2=arn:aws:eks:eu-west-1:<aws-account-id>:cluster/k8gb-cluster-eu-west-1 \
make terratest
```

- [Debugging](#debugging)

### Debugging

Delve debugger needs to be installed first. Follow the [installation instructions](https://github.com/go-delve/delve/tree/master/Documentation/installation) for specific platforms from Delve's website.

1. Run the following script

```shell script
> make debug-local
```

2. Attach debugger of your favourite IDE to port `2345`.

### Code style

The coding style suggested by the Golang community is used in k8gb project. See the [style doc][golang-style-doc] for details.

Please follow this style to make k8gb easy to review, maintain and develop.

Run `make lint` to automatically check if your code is compliant.

### Format of the commit message

We follow a rough convention for commit messages that are designed to answer two
questions: what changed and why. The subject line should feature the what and
the body of the commit should describe the why.

```
scripts: add the test-cluster command

this uses tmux to setup a test cluster that can easily be killed and started for debugging.

Fixes #38
```

The format can be described more formally as follows:

```
<subsystem>: <what changed>
<BLANK LINE>
<why this change was made>
<BLANK LINE>
<footer>
```

The first line is the subject and should be no longer than 70 characters, the second line is always blank, and other lines should be wrapped at 80 characters. This allows the message to be easier to read on GitHub as well as in various git tools.

## Documentation

If the contribution changes the existing APIs or user interface it must include sufficient documentation to explain the use of the new or updated feature.
[CHANGELOG](CHANGELOG.md) is automatically generated from Github PRs and Issues. Please use special [keywords](https://docs.github.com/en/github/managing-your-work-on-github/linking-a-pull-request-to-an-issue#linking-a-pull-request-to-an-issue-using-a-keyword) to link PR to and Issue for a clean changelog generation.


[golang-style-doc]: https://github.com/golang/go/wiki/CodeReviewComments
