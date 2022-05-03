package unit_test

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

func MustUnstructured(manifest string) *unstructured.Unstructured {
	u, err := Unstructured(manifest)
	if err != nil {
		panic(fmt.Errorf("unstructured: %w", err))
	}

	return u
}

func Unstructured(manifest string) (*unstructured.Unstructured, error) {
	obj := &unstructured.Unstructured{}
	if err := yaml.Unmarshal([]byte(manifest), obj); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return obj, nil
}

func PrintYAML(obj interface{}) {
	b, err := yaml.Marshal(obj)
	if err != nil {
		panic(fmt.Errorf("yaml marshal: %w", err))
	}

	fmt.Println(string(b))
	return
}

func YAML(y string) string {
	y = strings.ReplaceAll(y, "\t", "    ")
	return heredoc.Doc(y)
}
