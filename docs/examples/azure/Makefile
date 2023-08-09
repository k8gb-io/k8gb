# Copyright 2023 The k8gb Contributors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.PHONY: deploy-infra
deploy-infra:
	./infra/build.sh

.PHONY: setup-clusters
setup-clusters:
	./infra/setup.sh

.PHONY: destroy-infra
destroy-infra:
	./infra/destroy.sh

.PHONY: deploy-k8gb
deploy-k8gb:
	./k8gb/deploy.sh

.PHONY: deploy-demo
deploy-demo:
	./demo/deploy.sh

