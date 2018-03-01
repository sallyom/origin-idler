package idler

import (
	"log"
	"fmt"

	"github.com/kubernetes-sigs/kubebuilder/pkg/builders"

	"github.com/openshift/origin-idler/pkg/apis/idling/v1alpha2"
	listers "github.com/openshift/origin-idler/pkg/client/listers_generated/idling/v1alpha2"
	idlingclient "github.com/openshift/origin-idler/pkg/client/clientset_generated/clientset/typed/idling/v1alpha2"
	corelisters "k8s.io/client-go/listers/core/v1"
	"github.com/openshift/origin-idler/pkg/controller/sharedinformers"
	"k8s.io/client-go/scale"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/cache"
)

const (
	triggerServicesIndex = "triggerServices"
)

func triggerServicesIndexFunc(obj interface{}) ([]string, error) {
	idler, wasIdler := obj.(*v1alpha2.Idler)
	if !wasIdler {
		return nil, fmt.Errorf("trigger services indexer received object %v that wasn't an Idler", obj)
	}

	return idler.Spec.TriggerServiceNames, nil
}

// TODO: emit events for errors!
// TODO: set status conditions?

func (c *IdlerControllerImpl) ensureIdle(u *v1alpha2.Idler) (*v1alpha2.Idler, error) {
	var updatedIdler *v1alpha2.Idler
	var delayedErrors []error

	// put our previous scale records in a form slightly more conducive to looking up on the fly
	prevScales := make(map[v1alpha2.CrossGroupObjectReference]int32, len(u.Status.UnidledScales))
	for _, record := range u.Status.UnidledScales {
		prevScales[record.CrossGroupObjectReference] = record.PreviousScale
	}

	// ensure that all scalables in TargetScalables
	// are scaled to zero, and that their previous scale is recorded
	currScales := map[v1alpha2.CrossGroupObjectReference]int32{}
	for _, target := range u.Spec.TargetScalables {
		groupRes := schema.GroupResource{
			Group: target.Group,
			Resource: target.Resource,
		}
		currScaleObj, err := c.scaleClient.Scales(u.Namespace).Get(groupRes, target.Name)
		if err != nil {
			fullErr := fmt.Errorf("unable to fetch scale for target scalable %s %s: %v", groupRes.String(), target.Name, err)
			delayedErrors = append(delayedErrors, fullErr)
			// TODO(directxman12): emit event
			// continue on, we'll try again later
			continue
		}
		currScale := currScaleObj.Spec.Replicas
		if currScale == 0 {
			// mitigate data loss from update failure (see below) by declaring that
			// users can't scale to zero themselves without removing this from the list
			// of target scalables.  If we see something that we own with scale zero,
			// and we don't have a scale record for it, assume we lost data and
			// consider its previous scale to be one.
			if _, knownScale := prevScales[target]; !knownScale {
				currScales[target] = 1
			}

			// skip the actual scaling -- we don't need to do it
			continue
		}

		currScaleObj.Spec.Replicas = 0
		_, err = c.scaleClient.Scales(u.Namespace).Update(groupRes, currScaleObj)
		if err != nil {
			fullErr := fmt.Errorf("unable to fetch scale for target scalable %s %s: %v", groupRes.String(), target.Name, err)
			delayedErrors = append(delayedErrors, fullErr)
			// TODO(directxman12): emit event
			// continue on, we'll try again later
			continue
		}
		currScales[target] = currScale
	}

	// NB: It's possible for us to lose information in this step --
	// if we fail to update the idler at the end, we'll fail to record
	// previous scales.  Since we have no way of telling where the
	// scales came from, things might be permantently stuck in idled state.
	// Currently we mitigate this with the following strategies:
	// - assume users never manually scale things to zero without removing them from the list,
	//   and assume lost scales were 1 (partial mitigation, works across controller restarts)
	// Other possible strategies are:
	// - put updates in a retrying queue (mitigates, except in case of controller restart)
	// - get a helpful change into Kubernetes that causes Kubernetes to save the previous scale for us
	//   (could mitigate fully)

	// we only actually need to bother updating recorded scales if we made changes
	if len(currScales) > 0 {
		for _, prevRecord := range u.Status.UnidledScales {
			newVal, hasNewVal := currScales[prevRecord.CrossGroupObjectReference]
			if hasNewVal {
				// warn, but use the new val
				groupRes := schema.GroupResource{
					Group: prevRecord.Group,
					Resource: prevRecord.Resource,
				}
				log.Printf("found a new non-zero scale %v for target scalable %s %s with previously recorded scale %v on idler %s/%s, using the new scale", newVal, groupRes.String(), prevRecord.Name, prevRecord.PreviousScale, u.Namespace, u.Name)
				continue
			}
			currScales[prevRecord.CrossGroupObjectReference] = prevRecord.PreviousScale
		}

		// actually copy our old value
		updatedIdler = u.DeepCopy()
		updatedIdler.Status.UnidledScales = make([]v1alpha2.UnidleInfo, 0, len(currScales))
		for ref, scale := range currScales {
			updatedIdler.Status.UnidledScales = append(updatedIdler.Status.UnidledScales, v1alpha2.UnidleInfo{
				CrossGroupObjectReference: ref,
				PreviousScale: scale,
			})
		}
	}

	// if we've made a change and at least one scalable is idled,
	// we've started idling, so indicate that by setting Idled to true.
	if updatedIdler != nil && len(updatedIdler.Status.UnidledScales) > 0 {
		updatedIdler.Status.Idled = true
	}

	var err error
	if len(delayedErrors) > 0 {
		err = fmt.Errorf("unable to fully idle idler %s/%s: %v", u.Namespace, u.Name, utilerrors.NewAggregate(delayedErrors))
	}

	return updatedIdler, err
}

func (c *IdlerControllerImpl) ensureUnidle(u *v1alpha2.Idler) (*v1alpha2.Idler, error) {
	// if for some reason we don't have any previous records, just stop quickly
	if len(u.Status.UnidledScales) == 0 {
		return nil, nil
	}

	var updatedIdler *v1alpha2.Idler
	var delayedErrors []error

	// arrange all scale records for easy access
	prevScales := make(map[v1alpha2.CrossGroupObjectReference]int32, len(u.Status.UnidledScales))
	for _, record := range u.Status.UnidledScales {
		prevScales[record.CrossGroupObjectReference] = record.PreviousScale
	}

	// scale all targets with know previous scales back up
	for _, target := range u.Spec.TargetScalables {
		prevScale, hasRecord := prevScales[target]
		if !hasRecord {
			// skip any target scalable that we don't know about having scaled...
			continue
		}
		groupRes := schema.GroupResource{
			Group: target.Group,
			Resource: target.Resource,
		}
		_, err := c.scaleClient.Scales(u.Namespace).Update(groupRes, &autoscalingv1.Scale{
			ObjectMeta: metav1.ObjectMeta{
				Name: u.Name,
				Namespace: u.Namespace,
			},
			Spec: autoscalingv1.ScaleSpec{
				Replicas: prevScale,
			},
		})
		if err != nil {
			delayedErrors = append(delayedErrors, fmt.Errorf("unable restore target scalable %s %s to previous scale %v: %v", groupRes.String(), target.Name, prevScale, err))
			continue
		}

		delete(prevScales, target)
	}

	// NB: unlike the idle function, there's no issue if we fail to save our updates here --
	// we'll just see that we have leftover previous scales, and try to reconcile them again.

	// clean up unknown recorded scales, if needed
	if len(prevScales) > 0 {
		knownTargets := make(map[v1alpha2.CrossGroupObjectReference]struct{}, len(u.Spec.TargetScalables))
		for _, target := range u.Spec.TargetScalables {
			knownTargets[target] = struct{}{}
		}

		for target := range prevScales {
			if _, known := knownTargets[target]; !known {
				// NB: this is actually ok in Go
				delete(prevScales, target)
			}
		}
	}

	// check if we need to update the list of scale records
	if len(prevScales) != len(u.Status.UnidledScales) {
		// actually copy our idler here
		updatedIdler = u.DeepCopy()
		updatedIdler.Status.UnidledScales = make([]v1alpha2.UnidleInfo, 0, len(prevScales))
		for target, scale := range prevScales {
			updatedIdler.Status.UnidledScales = append(updatedIdler.Status.UnidledScales, v1alpha2.UnidleInfo{
				CrossGroupObjectReference: target,
				PreviousScale: scale,
			})
		}
	}

	// skip checking for unidling completion if we still have scales do deal with
	if len(prevScales) > 0 {
		var err error
		if len(delayedErrors) > 0 {
			err = fmt.Errorf("unable to fully unidle idler %s/%s: %v", u.Namespace, u.Name, utilerrors.NewAggregate(delayedErrors))
		}
		return updatedIdler, err
	}

	// NB: if the user ever scales the trigger services backing scalables down manually
	// before idling is finished, we'll never fully mark as unidled.  Such is life.

	// check if we've *actually* finished idling, which is determined by whether or not
	// all listed trigger services have at least one endpoint subset.
	activeSvcCount := 0
	for _, svcName := range u.Spec.TriggerServiceNames {
		ep, err := c.endpointsLister.Endpoints(u.Namespace).Get(svcName)
		if err != nil {
			delayedErrors = append(delayedErrors, fmt.Errorf("unable to check endpoints for service %s: %v", svcName, err))
			continue
		}
		// TODO: report unready endpoints, somehow?
		if len(ep.Subsets) > 0 && len(ep.Subsets[0].Addresses) > 0 {
			activeSvcCount++
		}
	}

	if activeSvcCount == len(u.Spec.TriggerServiceNames) {
		// consider ourselves updated
		// populate updatedIdler if necessary
		if updatedIdler == nil {
			updatedIdler = u.DeepCopy()
		}

		updatedIdler.Status.Idled = false
	}

	var err error
	if len(delayedErrors) > 0 {
		err = fmt.Errorf("unable to fully unidle idler %s/%s: %v", u.Namespace, u.Name, utilerrors.NewAggregate(delayedErrors))
	}
	return updatedIdler, err
}

// Reconcile handles enqueued messages
func (c *IdlerControllerImpl) Reconcile(u *v1alpha2.Idler) error {
	log.Printf("Running reconcile Idler for %s\n", u.Name)

	var updatedIdler *v1alpha2.Idler
	var err error

	// when WantIdle is true...
	if u.Spec.WantIdle {
		updatedIdler, err = c.ensureIdle(u)
	} else if !u.Spec.WantIdle && u.Status.Idled {
		updatedIdler, err = c.ensureUnidle(u)
	}

	if updatedIdler != nil {
		_, updateErr := c.idlerClient.Idlers(u.Namespace).Update(updatedIdler)
		if updateErr != nil {
			err = fmt.Errorf("unable to update idler %s/%s: %v (errors while idling: %v)", u.Namespace, u.Name, updateErr, err)
		}
	}

	// NB: this *must* come last -- we treat all errors as non-fatal, and want to
	// perform an update first if we have one...
	return err
}

// +controller:group=idling,version=v1alpha2,kind=Idler,resource=idlers
type IdlerControllerImpl struct {
	builders.DefaultControllerFns

	// lister indexes properties about Idler
	lister listers.IdlerLister

	// indexer allows us to look up Idlers by index
	indexer cache.Indexer

	// endpointsLister indexes Endpoints
	endpointsLister corelisters.EndpointsLister

	// idlerClient can update idler
	idlerClient idlingclient.IdlersGetter

	// scaleClient knows how to access scales
	scaleClient scale.ScalesGetter
}

// Init initializes the controller and is called by the generated code
// Register watches for additional resource types here.
func (c *IdlerControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) {
	// INSERT YOUR CODE HERE - add logic for initializing the controller as needed

	idlerInformer := arguments.GetSharedInformers().Factory.Idling().V1alpha2().Idlers()
	// make sure we index by service names, for easy lookup later
	idlerInformer.Informer().AddIndexers(cache.Indexers{
		triggerServicesIndex: triggerServicesIndexFunc,
	})
	// Use the lister for indexing idlers labels
	c.lister = idlerInformer.Lister()
	c.indexer = idlerInformer.Informer().GetIndexer()

	c.endpointsLister = arguments.GetSharedInformers().KubernetesFactory.Core().V1().Endpoints().Lister()

	// TODO: set scale client and idlers client

	// To watch other resource types, uncomment this function and replace Foo with the resource name to watch.
	// Must define the func FooToIdler(i interface{}) (string, error) {} that returns the Idler
	// "namespace/name"" to reconcile in response to the updated Foo
	// Note: To watch Kubernetes resources, you must also update the StartAdditionalInformers function in
	// pkg/controllers/sharedinformers/informers.go
	//
	// arguments.Watch("IdlerFoo",
	//     arguments.GetSharedInformers().Factory.Bar().V1beta1().Bars().Informer(),
	//     c.FooToIdler)

	arguments.Watch("IdlerEndpoints",
		arguments.GetSharedInformers().KubernetesFactory.Core().V1().Endpoints().Informer(),
		c.EndpointsToIdler)
}

func (c *IdlerControllerImpl) EndpointsToIdler(epRaw interface{}) ([]string, error) {
	ep, wasEp := epRaw.(*corev1.Endpoints)
	if !wasEp {
		return nil, fmt.Errorf("endpoints-to-idler lookup received non-Endpoints object %v", epRaw)
	}
	idlers, err := c.indexer.ByIndex(triggerServicesIndex, ep.Name)
	if err != nil {
		return nil, err
	}

	keys := make([]string, len(idlers))

	for i, idlerRaw := range idlers {
		idler, wasIdler := idlers[0].(*v1alpha2.Idler)
		if !wasIdler {
			return nil, fmt.Errorf("idler indexer returned non-idler object %v", idlerRaw)
		}
		keys[i] = idler.Namespace+"/"+idler.Name
	}
	return keys, nil
}

func (c *IdlerControllerImpl) Get(namespace, name string) (*v1alpha2.Idler, error) {
	return c.lister.Idlers(namespace).Get(name)
}
