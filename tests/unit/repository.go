package unit

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/vmware-tanzu/cartographer/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type repository struct {
	cl client.Client
}

type RepositoryOption func(r *repository)

func NewRepository(opts ...RepositoryOption) *repository {
	r := &repository{}
	for _, opt := range opts {
		opt(r)
	}

	if r.cl == nil {
		r.cl = mustNewDefaultClient()
	}

	return r
}

func mustNewDefaultClient() client.Client {
	cl, err := newDefaultClient()
	if err != nil {
		panic(fmt.Errorf("new default client: %w", err))
	}

	return cl
}

func newDefaultClient() (client.Client, error) {
	scheme := k8sruntime.NewScheme()
	if err := v1alpha1.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("add to scheme: %w", err)
	}

	cl, err := client.New(config.GetConfigOrDie(), client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, fmt.Errorf("new client: %w", err)
	}

	return cl, nil
}

// GetTemplate retrieves from the Kubernetes cluster a Cartographer template
// of a given `kind` and `name` always encoded as a ClusterTemplate as that's
// the mininum resource necessary for what we care about in testing templates.
//
func (r *repository) GetTemplate(ctx context.Context, kind, name string) (*v1alpha1.ClusterTemplate, error) {
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(schema.GroupVersionKind{
		Version: v1alpha1.SchemeGroupVersion.Version,
		Group:   v1alpha1.SchemeGroupVersion.Group,
		Kind:    kind,
	})

	err := r.cl.Get(ctx, client.ObjectKey{Name: name}, obj)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	b, err := json.Marshal(obj.Object)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	target := &v1alpha1.ClusterTemplate{}
	if err := json.Unmarshal(b, target); err != nil {
		return nil, fmt.Errorf("unmarshal to clustertemplate: %w", err)
	}

	return target, nil
}
