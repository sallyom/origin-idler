package fake

import (
	clientset "github.com/openshift/origin-idler/pkg/client/clientset_generated/clientset"
	idlingv1alpha2 "github.com/openshift/origin-idler/pkg/client/clientset_generated/clientset/typed/idling/v1alpha2"
	fakeidlingv1alpha2 "github.com/openshift/origin-idler/pkg/client/clientset_generated/clientset/typed/idling/v1alpha2/fake"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/testing"
)

// NewSimpleClientset returns a clientset that will respond with the provided objects.
// It's backed by a very simple object tracker that processes creates, updates and deletions as-is,
// without applying any validations and/or defaults. It shouldn't be considered a replacement
// for a real clientset and is mostly useful in simple unit tests.
func NewSimpleClientset(objects ...runtime.Object) *Clientset {
	o := testing.NewObjectTracker(scheme, codecs.UniversalDecoder())
	for _, obj := range objects {
		if err := o.Add(obj); err != nil {
			panic(err)
		}
	}

	fakePtr := testing.Fake{}
	fakePtr.AddReactor("*", "*", testing.ObjectReaction(o))
	fakePtr.AddWatchReactor("*", testing.DefaultWatchReactor(watch.NewFake(), nil))

	return &Clientset{fakePtr, &fakediscovery.FakeDiscovery{Fake: &fakePtr}}
}

// Clientset implements clientset.Interface. Meant to be embedded into a
// struct to get a default implementation. This makes faking out just the method
// you want to test easier.
type Clientset struct {
	testing.Fake
	discovery *fakediscovery.FakeDiscovery
}

func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	return c.discovery
}

var _ clientset.Interface = &Clientset{}

// IdlingV1alpha2 retrieves the IdlingV1alpha2Client
func (c *Clientset) IdlingV1alpha2() idlingv1alpha2.IdlingV1alpha2Interface {
	return &fakeidlingv1alpha2.FakeIdlingV1alpha2{Fake: &c.Fake}
}

// Idling retrieves the IdlingV1alpha2Client
func (c *Clientset) Idling() idlingv1alpha2.IdlingV1alpha2Interface {
	return &fakeidlingv1alpha2.FakeIdlingV1alpha2{Fake: &c.Fake}
}
