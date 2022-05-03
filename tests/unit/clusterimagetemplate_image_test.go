package unit

import (
	"context"
	"testing"

	"github.com/vmware-tanzu/cartographer-catalog/pkg/templates"
)

func TestClusterImageTemplateImageLabels(t *testing.T) {
	// a supply chain with a default registry parameter configured
	//
	supplyChain := templates.MustClusterSupplyChain(`
		apiVersion: carto.run/v1alpha1
		kind: ClusterSupplyChain
		metadata:
		  name: foo
		spec:
		  params:
		    - name: registry
		      default: {server: "foo", repository: "bar"}
	`)

	(&templates.TestSuite{
		TemplateName: "image",
		TemplateKind: "ClusterImageTemplate",

		Cases: []templates.TestCase{
			{
				Name: "k8s component label",

				Given: templates.Stamper().
					ClusterSupplyChain(supplyChain).
					Source("source", "github.com/foo", "0xb33f").
					Workload(templates.MustWorkload(`
						metadata:
						  name: foo
						  namespace: default
						spec:
						  source: {}
					`)),

				Expect: []templates.Expectation{
					templates.Field(`metadata.labels`).SetTo(map[string]interface{}{
						"app.kubernetes.io/component": "build",
					}),
				},
			},
		},
	}).Run(context.Background(), t)
}
