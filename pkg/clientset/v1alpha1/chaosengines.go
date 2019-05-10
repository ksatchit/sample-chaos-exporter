package v1alpha1

import (
	// "github.com/martin-helmich/kubernetes-crd-example/api/types/v1alpha1"
        "github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ChaosEngineInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.ChaosEngineList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.ChaosEngine, error)
	Create(*v1alpha1.ChaosEngine) (*v1alpha1.ChaosEngine, error)
	// Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type chaosEngineClient struct {
	restClient rest.Interface
	ns         string
}

func (c *chaosEngineClient) List(opts metav1.ListOptions) (*v1alpha1.ChaosEngineList, error) {
	result := v1alpha1.ChaosEngineList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("chaosengines").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *chaosEngineClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.ChaosEngine, error) {
	result := v1alpha1.ChaosEngine{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("chaosengines").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *chaosEngineClient) Create(chaosengine *v1alpha1.ChaosEngine) (*v1alpha1.ChaosEngine, error) {
	result := v1alpha1.ChaosEngine{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("chaosengines").
		Body(chaosengine).
		Do().
		Into(&result)

	return &result, err
}

/*
func (c *projectClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("projects").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
*/
