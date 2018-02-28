package v1alpha2_test

import (
	"testing"

	"github.com/kubernetes-sigs/kubebuilder/pkg/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/rest"

	"github.com/openshift/origin-idler/pkg/apis"
	"github.com/openshift/origin-idler/pkg/client/clientset_generated/clientset"
)

var testenv *test.TestEnvironment
var config *rest.Config
var cs *clientset.Clientset

func TestV1alpha2(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "v1 Suite", []Reporter{test.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	testenv = &test.TestEnvironment{CRDs: apis.APIMeta.GetCRDs()}

	var err error
	config, err = testenv.Start()
	Expect(err).NotTo(HaveOccurred())

	cs = clientset.NewForConfigOrDie(config)
})

var _ = AfterSuite(func() {
	testenv.Stop()
})
