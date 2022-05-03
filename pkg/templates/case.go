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

package templates

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/vmware-tanzu/cartographer/pkg/apis/v1alpha1"
	cartotemplates "github.com/vmware-tanzu/cartographer/pkg/templates"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type StamperIface interface {
	Stamp(
		ctx context.Context,
		template *v1alpha1.ClusterTemplate,
	) (*unstructured.Unstructured, error)
}

// TestCase
//
type TestCase struct {
	Name string

	Given StamperIface

	GivenWorkload    *v1alpha1.Workload
	GivenSupplyChain *v1alpha1.ClusterSupplyChain
	GivenInputs      cartotemplates.Inputs

	Expect      []Expectation
	ExpectedErr string
	Expected    client.Object
}

// Run
//
func (tc *TestCase) Run(ctx context.Context, t *testing.T, template *v1alpha1.ClusterTemplate) {
	actual, err := tc.Given.Stamp(ctx, template)
	if err != nil {
		if tc.ExpectedErr == "" {
			t.Fatalf("unexpected err: %v", err)
		}

		if !strings.Contains(err.Error(), tc.ExpectedErr) {
			t.Fatalf("err '%v' doesn't contain expected '%v'",
				err, tc.ExpectedErr)
		}

		return
	}

	PrintYAML(actual)

	if tc.Expected != nil {
		diff := cmp.Diff(tc.Expected, actual)
		if diff != "" {
			t.Fatalf("(-expected, +actual) = %v", diff)
		}
	}

	for _, expectation := range tc.Expect {
		if err := expectation(actual); err != nil {
			t.Fatal(err)
		}
	}
}
