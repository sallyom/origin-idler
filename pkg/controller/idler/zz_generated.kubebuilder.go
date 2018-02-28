package idler

import (
	"github.com/golang/glog"
	"github.com/kubernetes-sigs/kubebuilder/pkg/controller"
	"github.com/openshift/origin-idler/pkg/controller/sharedinformers"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

// IdlerController implements the controller.IdlerController interface
type IdlerController struct {
	queue *controller.QueueWorker

	// Handles messages
	controller *IdlerControllerImpl

	Name string

	BeforeReconcile func(key string)
	AfterReconcile  func(key string, err error)

	Informers *sharedinformers.SharedInformers
}

// NewController returns a new IdlerController for responding to Idler events
func NewIdlerController(config *rest.Config, si *sharedinformers.SharedInformers) *IdlerController {
	q := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Idler")

	queue := &controller.QueueWorker{q, 10, "Idler", nil}
	c := &IdlerController{queue, nil, "Idler", nil, nil, si}

	// For non-generated code to add events
	uc := &IdlerControllerImpl{}
	var ci sharedinformers.Controller = uc

	if i, ok := ci.(sharedinformers.ControllerInit); ok {
		i.Init(&sharedinformers.ControllerInitArgumentsImpl{si, config, c.LookupAndReconcile})
	}

	c.controller = uc

	queue.Reconcile = c.LookupAndReconcile
	if c.Informers.WorkerQueues == nil {
		c.Informers.WorkerQueues = map[string]*controller.QueueWorker{}
	}
	c.Informers.WorkerQueues["Idler"] = queue
	si.Factory.Idling().V1alpha2().Idlers().Informer().
		AddEventHandler(&controller.QueueingEventHandler{q, nil, false})
	return c
}

func (c *IdlerController) GetName() string {
	return c.Name
}

func (c *IdlerController) LookupAndReconcile(key string) (err error) {
	var namespace, name string

	if c.BeforeReconcile != nil {
		c.BeforeReconcile(key)
	}
	if c.AfterReconcile != nil {
		// Wrap in a function so err is evaluated after it is set
		defer func() { c.AfterReconcile(key, err) }()
	}

	namespace, name, err = cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return
	}

	u, err := c.controller.Get(namespace, name)
	if errors.IsNotFound(err) {
		glog.Infof("Not doing work for Idler %v because it has been deleted", key)
		// Set error so it is picked up by AfterReconcile and the return function
		err = nil
		return
	}
	if err != nil {
		glog.Errorf("Unable to retrieve Idler %v from store: %v", key, err)
		return
	}

	// Set error so it is picked up by AfterReconcile and the return function
	err = c.controller.Reconcile(u)

	return
}

func (c *IdlerController) Run(stopCh <-chan struct{}) {
	for _, q := range c.Informers.WorkerQueues {
		q.Run(stopCh)
	}
	controller.GetDefaults(c.controller).Run(stopCh)
}
