package unit_test

import (
	"context"
	"testing"

	"github.com/vmware-tanzu/cartographer-catalog/tests/unit"
)

func TestAppClusterTemplate(t *testing.T) {
	ctx := context.Background()

	template, err := unit.NewRepository().GetTemplate(ctx,
		"ClusterSourceTemplate", "git-repository",
	)
	if err != nil {
		t.Fatal(err)
	}

	(&TestSuite{
		Template: template,

		Cases: []TestCase{{
			Name: "no source in spec",

			GivenOwner: MustUnstructured(YAML(`
				kind: Workload
				apiVersion: carto.run/v1alpha1
				metadata:
				  name: foo
				spec: {}
			`)),
			ExpectedErr: "struct has no .source field or method",
		}, {
			Name: "source in spec",

			GivenOwner: MustUnstructured(YAML(`
				kind: Workload
				apiVersion: carto.run/v1alpha1
				metadata:
				  name: foo
				spec:
				  source:
				    git:
				      url: https://github.com/foo/bar
				      ref: {branch: main}
			`)),
			Expected: MustUnstructured(YAML(`
				apiVersion: source.toolkit.fluxcd.io/v1beta1
				kind: GitRepository
				metadata:
				  name: foo
				  labels:
				    app.kubernetes.io/component: source
				  ownerReferences:
				  - apiVersion: carto.run/v1alpha1
				    blockOwnerDeletion: true
				    controller: true
				    kind: Workload
				    name: foo
				    uid: ""
				spec:
				  gitImplementation: go-git
				  ignore: '!.git'
				  interval: 1m0s
				  ref:
				    branch: main
				  url: https://github.com/foo/bar
			`)),
		}},
	}).Run(ctx, t)
}
