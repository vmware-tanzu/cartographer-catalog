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

	"github.com/google/go-cmp/cmp"
	"github.com/vmware-tanzu/cartographer/pkg/eval"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Expectation
//
type Expectation func(stampedObject *unstructured.Unstructured) error

type field struct {
	path   string
	target interface{}
}

func Field(jsonpath string) *field {
	return &field{
		path: jsonpath,
	}
}

func (f *field) SetTo(val interface{}) Expectation {
	f.target = val
	return f.checkSetTo
}

func (f *field) NotSet() Expectation {
	f.target = nil
	return f.checkNotSet
}

// checkNotSet validates that a given object `obj` DOES NOT properly evaluate a
// jsonpath query.
//
func (f *field) checkNotSet(obj *unstructured.Unstructured) error {
	_, err := f.eval(f.path, obj.Object)
	if err == nil {
		return fmt.Errorf("expected jsonpath evaluation to fail, but didn't")
	}

	if strings.Contains(err.Error(), "failed to find results") {
		return nil
	}

	return err
}

// checkSetSet validates whether a given object `obj` has under a particular
// location (dictated by a jsonpath query) some value.
//
func (f *field) checkSetTo(obj *unstructured.Unstructured) error {
	res, err := f.eval(f.path, obj.Object)
	if err != nil {
		return fmt.Errorf("evaluate: %w", err)
	}

	if d := cmp.Diff(res, f.target); d != "" {
		return fmt.Errorf("(-expected, +actual) = %v", d)
	}

	return nil
}

func (f *field) eval(path string, obj interface{}) (interface{}, error) {
	res, err := eval.EvaluatorBuilder().EvaluateJsonPath(f.path, obj)
	if err != nil {
		return nil, fmt.Errorf("evaluate: %w", err)
	}

	return res, nil
}
