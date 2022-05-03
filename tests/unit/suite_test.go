package unit_test

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/vmware-tanzu/cartographer-catalog/tests/unit"
	"github.com/vmware-tanzu/cartographer/pkg/apis/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type TestCase struct {
	Name string

	GivenOwner   client.Object
	GivenContext unit.TemplatingContext

	ExpectedErr string
	Expected    client.Object
}

func (tc *TestCase) Run(ctx context.Context, t *testing.T, template *v1alpha1.ClusterTemplate) {
	tc.GivenContext.Workload = tc.GivenOwner

	actual, err := unit.Stamp(ctx, template, tc.GivenContext)
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

	diff := cmp.Diff(tc.Expected, actual)
	if diff != "" {
		t.Errorf("(-expected, +actual) = %v", diff)
		PrintYAML(actual)
	}
}

type TestSuite struct {
	Template *v1alpha1.ClusterTemplate

	Cases []TestCase
}

func (s *TestSuite) Run(ctx context.Context, t *testing.T) {
	for _, tc := range s.Cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			tc.Run(ctx, t, s.Template)
		})
	}
}
