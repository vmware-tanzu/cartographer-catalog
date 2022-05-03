package unit

import (
	"context"
	"fmt"

	"github.com/vmware-tanzu/cartographer/pkg/apis/v1alpha1"
	"github.com/vmware-tanzu/cartographer/pkg/templates"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TemplatingContext
//
type TemplatingContext struct {
	Workload    client.Object          `json:"workload"`
	Deliverable client.Object          `json:"deliverable"`
	Params      map[string]interface{} `json:"params"`
}

// Stamp stamps out a Kubernetes object based on the templatespec of a template
// and the data provided for its stamping context.
//
func Stamp(
	ctx context.Context,
	template *v1alpha1.ClusterTemplate,
	data TemplatingContext,
) (*unstructured.Unstructured, error) {
	params := map[string]interface{}{}
	params = mergeInto(params, templateDefaultParams(template))
	params = mergeInto(params, data.Params)

	owner := data.Workload
	if owner == nil {

		return nil, fmt.Errorf("workload or deliverable " +
			"must be set in templating ctx")
	}

	stamper := templates.StamperBuilder(owner, map[string]interface{}{
		"workload": owner,
		"params":   params,
	}, nil)

	stampedObj, err := stamper.Stamp(ctx, template.Spec)
	if err != nil {
		return nil, fmt.Errorf("stamp: %w", err)
	}

	return stampedObj, nil
}

// templateDefaultParams retrieves the default set of params from a
// Cartographer template.
//
func templateDefaultParams(template *v1alpha1.ClusterTemplate) map[string]interface{} {
	res := map[string]interface{}{}

	for _, templateParam := range template.Spec.Params {
		res[templateParam.Name] = templateParam.DefaultValue
	}

	return res
}

// mergeInto merges a given map `src` into a destination `dst` without
// modifying either.
//
func mergeInto(dst, src map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range dst {
		res[k] = v
	}

	for k, v := range src {
		res[k] = v
	}

	return res
}
