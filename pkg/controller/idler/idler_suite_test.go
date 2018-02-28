


package idler_test

import (
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "k8s.io/client-go/rest"
    "github.com/kubernetes-sigs/kubebuilder/pkg/test"

    "github.com/openshift/origin-idler/pkg/apis"
    "github.com/openshift/origin-idler/pkg/client/clientset_generated/clientset"
    "github.com/openshift/origin-idler/pkg/controller/sharedinformers"
    "github.com/openshift/origin-idler/pkg/controller/idler"
)

var testenv *test.TestEnvironment
var config *rest.Config
var cs *clientset.Clientset
var shutdown chan struct{}
var controller *idler.IdlerController
var si *sharedinformers.SharedInformers

func TestIdler(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecsWithDefaultAndCustomReporters(t, "Idler Suite", []Reporter{test.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
    testenv = &test.TestEnvironment{CRDs: apis.APIMeta.GetCRDs()}
    var err error
    config, err = testenv.Start()
    Expect(err).NotTo(HaveOccurred())
    cs = clientset.NewForConfigOrDie(config)

    shutdown = make(chan struct{})
    si = sharedinformers.NewSharedInformers(config, shutdown)
    controller = idler.NewIdlerController(config, si)
    controller.Run(shutdown)
})

var _ = AfterSuite(func() {
    testenv.Stop()
})
