#!/bin/bash
# Workaround of Make being to smart on up-to-date PHONY targets
# If we execute all of them normal way, then targets from `deploy-second-ohmyglb`
# will never be executed as they contain the same underlying target as `deploy-first-ohmyglb`
# but with different variables

make deploy-two-local-clusters use-first-context deploy-first-ohmyglb deploy-test-apps
make use-second-context deploy-second-ohmyglb deploy-test-apps
