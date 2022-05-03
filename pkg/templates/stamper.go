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
	"fmt"

	"github.com/vmware-tanzu/cartographer/pkg/apis/v1alpha1"
	"github.com/vmware-tanzu/cartographer/pkg/templates"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type stamper struct {
	workload    *v1alpha1.Workload
	supplyChain *v1alpha1.ClusterSupplyChain
	resource    *v1alpha1.SupplyChainResource

	inputs templates.Inputs
}

func Stamper() *stamper {
	return &stamper{
		resource: &v1alpha1.SupplyChainResource{
			Name: "<none>",
		},
		supplyChain: &v1alpha1.ClusterSupplyChain{
			ObjectMeta: metav1.ObjectMeta{
				Name: "<none>",
			},
		},
	}
}

func (s *stamper) Source(name, url, rev string) *stamper {
	if len(s.inputs.Sources) == 0 {
		s.inputs.Sources = map[string]templates.SourceInput{}
	}

	s.inputs.Sources[name] = templates.SourceInput{
		Name:     name,
		URL:      url,
		Revision: rev,
	}
	return s
}

func (s *stamper) ClusterSupplyChain(v *v1alpha1.ClusterSupplyChain) *stamper {
	s.supplyChain = v
	return s
}

func (s *stamper) Workload(v *v1alpha1.Workload) *stamper {
	s.workload = v
	return s
}

func (s *stamper) validate() error {
	if s.workload == nil {
		return fmt.Errorf("workload must be set")
	}

	if s.supplyChain == nil {
		return fmt.Errorf("supplychain must be set")
	}

	if s.resource == nil {
		return fmt.Errorf("resource must be set")
	}

	return nil
}

func (s *stamper) Stamp(
	ctx context.Context, template *v1alpha1.ClusterTemplate,
) (*unstructured.Unstructured, error) {
	if err := s.validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	templatingContext := map[string]interface{}{
		"workload": s.workload,
		"params": templates.ParamsBuilder(
			template.Spec.Params,
			s.supplyChain.Spec.Params,
			s.resource.Params,
			s.workload.Spec.Params,
		),
		"sources": s.inputs.Sources,
		"images":  s.inputs.Images,
		"configs": s.inputs.Configs,
	}

	if s.inputs.OnlyConfig() != nil {
		templatingContext["config"] = s.inputs.OnlyConfig()
	}
	if s.inputs.OnlyImage() != nil {
		templatingContext["image"] = s.inputs.OnlyImage()
	}
	if s.inputs.OnlySource() != nil {
		templatingContext["source"] = s.inputs.OnlySource()
	}

	stampContext := templates.StamperBuilder(
		s.workload,
		templatingContext,
		nil,
	)

	obj, err := stampContext.Stamp(ctx, template.Spec)
	if err != nil {
		return nil, fmt.Errorf("stamp: %w", err)
	}

	return obj, nil
}
