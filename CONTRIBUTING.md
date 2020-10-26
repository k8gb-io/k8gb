# How to contribute

k8gb is MIT licensed and accepts contributions via GitHub pull requests. This document outlines some of the conventions on commit message formatting, contact points for developers, and other resources to help get contributions into k8gb project.

# Slack chat

- Find us at `#sig-multicluster` on kubernetes.slack.com

## Getting started

- Fork the repository on GitHub
- See the [local playground guide][docs/local.md] for testing environment setup

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

### Testing

* Unit tests should be updated for any functional code change at [test suite location](/controllers/gslb_controller_test.go).
* Acceptance terratest suite is located [here](/terratest) and executable by `make terratest` target. These tests are changed only if the
 change is substantial enough to affect the main end-to-end flow.

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

If the contribution changes the existing APIs or user interface it must include sufficient documentation to explain the use of the new or updated feature. Likewise the [CHANGELOG][changelog] should be updated with a summary of the change and link to the pull request.


[golang-style-doc]: https://github.com/golang/go/wiki/CodeReviewComments
