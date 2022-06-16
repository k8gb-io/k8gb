package assistant

import "sort"

/*
Copyright 2022 The k8gb Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/

type Target struct {
	IPs []string
}

type Targets map[string]*Target

func NewTargets() Targets {
	return make(map[string]*Target, 0)
}

func (t Targets) GetIPs() (ips []string) {
	// initializing targets to avoid possible nil reference errors (serialization etc.)
	ips = []string{}
	for _, v := range t {
		ips = append(ips, v.IPs...)
	}
	return ips
}

func (t Targets) Append(tag string, ips []string) {
	if target, found := t[tag]; found {
		target.IPs = append(target.IPs, ips...)
		return
	}
	t[tag] = &Target{IPs: ips}
}

func (t Targets) AppendTargets(targets Targets) {
	for k, v := range targets {
		t.Append(k, v.IPs)
	}
}

func (t Targets) Sort() {
	sort := func(targets []string) []string {
		sort.Slice(targets, func(i, j int) bool {
			return targets[i] < targets[j]
		})
		return targets
	}
	for _, v := range t {
		v.IPs = sort(v.IPs)
	}
}
