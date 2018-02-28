package idler

import (
	"log"

	"github.com/kubernetes-sigs/kubebuilder/pkg/builders"

	"github.com/openshift/origin-idler/pkg/apis/idling/v1alpha2"
	listers "github.com/openshift/origin-idler/pkg/client/listers_generated/idling/v1alpha2"
	"github.com/openshift/origin-idler/pkg/controller/sharedinformers"
)

// EDIT THIS FILE!
// Created by "kubebuilder create resource" for you to implement controller logic for the Idler resource API

// Reconcile handles enqueued messages
func (c *IdlerControllerImpl) Reconcile(u *v1alpha2.Idler) error {
	// INSERT YOUR CODE HERE - implement controller logic to reconcile observed and desired state of the object
	log.Printf("Running reconcile Idler for %s\n", u.Name)
	return nil
}

// +controller:group=idling,version=v1alpha2,kind=Idler,resource=idlers
type IdlerControllerImpl struct {
	builders.DefaultControllerFns

	// lister indexes properties about Idler
	lister listers.IdlerLister
}

// Init initializes the controller and is called by the generated code
// Register watches for additional resource types here.
func (c *IdlerControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) {
	// INSERT YOUR CODE HERE - add logic for initializing the controller as needed

	// Use the lister for indexing idlers labels
	c.lister = arguments.GetSharedInformers().Factory.Idling().V1alpha2().Idlers().Lister()

	// To watch other resource types, uncomment this function and replace Foo with the resource name to watch.
	// Must define the func FooToIdler(i interface{}) (string, error) {} that returns the Idler
	// "namespace/name"" to reconcile in response to the updated Foo
	// Note: To watch Kubernetes resources, you must also update the StartAdditionalInformers function in
	// pkg/controllers/sharedinformers/informers.go
	//
	// arguments.Watch("IdlerFoo",
	//     arguments.GetSharedInformers().Factory.Bar().V1beta1().Bars().Informer(),
	//     c.FooToIdler)
}

func (c *IdlerControllerImpl) Get(namespace, name string) (*v1alpha2.Idler, error) {
	return c.lister.Idlers(namespace).Get(name)
}
