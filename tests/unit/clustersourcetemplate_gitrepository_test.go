// Copyright 2022 VMware
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package unit_test

import (
	"context"
	"testing"

	"github.com/vmware-tanzu/cartographer-catalog/pkg/templates"
)

func TestGitRepositoryParams(t *testing.T) {
	(&templates.TestSuite{
		TemplateName: "git-repository",
		TemplateKind: "ClusterSourceTemplate",

		Cases: []templates.TestCase{
			{
				Name: "none set, template defaults used",

				Given: templates.Stamper().
					Workload(templates.MustWorkload(`
						kind: Workload
						metadata: {name: "foo"}
						spec: {source: {git: {url: "foo", ref: {}}}}
					`)),

				Expect: []templates.Expectation{
					templates.Field(`spec.gitImplementation`).SetTo("go-git"),
					templates.Field(`spec.secretRef`).NotSet(),
				},
			},

			{
				Name: "git_implementation set, overrides default",

				Given: templates.Stamper().
					Workload(templates.MustWorkload(`
						kind: Workload
						metadata: {name: "foo"}
						spec:
						  source: {git: {url: "foo", ref: {}}}
						  params:
						    - name: git_implementation
						      value: foo
					`)),

				Expect: []templates.Expectation{
					templates.Field(`spec.gitImplementation`).SetTo("foo"),
				},
			},

			{
				Name: "git_secret set, configures spec.secretRef",

				Given: templates.Stamper().
					Workload(templates.MustWorkload(`
						kind: Workload
						metadata: {name: "foo"}
						spec:
						  source: {git: {url: "foo", ref: {}}}
						  params:
						    - name: git_secret
						      value: foo
					`)),
				Expect: []templates.Expectation{
					templates.Field(`spec.secretRef.name`).SetTo("foo"),
				},
			},
		},
	}).Run(context.Background(), t)
}

func TestGitRepositoryLabels(t *testing.T) {
	(&templates.TestSuite{
		TemplateName: "git-repository",
		TemplateKind: "ClusterSourceTemplate",

		Cases: []templates.TestCase{
			{
				Name: "k8s component label set",

				Given: templates.Stamper().
					Workload(templates.MustWorkload(`
						kind: Workload
						metadata: {name: "foo"}
						spec: {source: {git: {url: "foo", ref: {}}}}
				`)),

				Expect: []templates.Expectation{
					templates.Field(`metadata.labels`).SetTo(map[string]interface{}{
						"app.kubernetes.io/component": "source",
					}),
				},
			},

			{
				Name: "workload labels propagete",
				Given: templates.Stamper().
					Workload(templates.MustWorkload(`
						kind: Workload
						metadata: {name: "foo", labels: {"foo": "bar"}}
						spec: {source: {git: {url: "foo", ref: {}}}}
				`)),

				Expect: []templates.Expectation{
					templates.Field(`metadata.labels`).SetTo(map[string]interface{}{
						"app.kubernetes.io/component": "source",
						"foo":                         "bar",
					}),
				},
			},
		},
	}).Run(context.Background(), t)
}
