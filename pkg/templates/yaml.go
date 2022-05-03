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
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/vmware-tanzu/cartographer/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

// MustUnstructured
//
func MustUnstructured(manifest string) *unstructured.Unstructured {
	manifest = YAML(manifest)

	obj, err := Unstructured(manifest)
	if err != nil {
		panic(fmt.Errorf("unstructured: %w", err))
	}

	return obj
}

func MustClusterSupplyChain(manifest string) *v1alpha1.ClusterSupplyChain {
	manifest = YAML(manifest)

	obj := &v1alpha1.ClusterSupplyChain{}
	if err := yaml.Unmarshal([]byte(manifest), obj); err != nil {
		panic(fmt.Errorf("unmarshal: %w", err))
	}

	return obj
}

func MustWorkload(manifest string) *v1alpha1.Workload {
	manifest = YAML(manifest)

	obj := &v1alpha1.Workload{}
	if err := yaml.Unmarshal([]byte(manifest), obj); err != nil {
		panic(fmt.Errorf("unmarshal: %w", err))
	}

	return obj
}

// Unstructured
//
func Unstructured(manifest string) (*unstructured.Unstructured, error) {
	obj := &unstructured.Unstructured{}
	if err := yaml.Unmarshal([]byte(manifest), obj); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return obj, nil
}

// PrintYAML
//
func PrintYAML(obj interface{}) {
	b, err := yaml.Marshal(obj)
	if err != nil {
		panic(fmt.Errorf("yaml marshal: %w", err))
	}

	fmt.Println(string(b))
	return
}

// YAML
//
func YAML(y string) string {
	y = strings.ReplaceAll(y, "\t", "    ")
	return heredoc.Doc(y)
}
