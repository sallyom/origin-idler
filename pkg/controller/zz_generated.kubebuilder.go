package controller

import (
	"github.com/kubernetes-sigs/kubebuilder/pkg/controller"
	"github.com/openshift/origin-idler/pkg/controller/idler"
	"github.com/openshift/origin-idler/pkg/controller/sharedinformers"
	"k8s.io/client-go/rest"
)

func GetAllControllers(config *rest.Config) ([]controller.Controller, chan struct{}) {
	shutdown := make(chan struct{})
	si := sharedinformers.NewSharedInformers(config, shutdown)
	return []controller.Controller{
		idler.NewIdlerController(config, si),
	}, shutdown
}
