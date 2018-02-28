package apis

import (
	"github.com/kubernetes-sigs/kubebuilder/pkg/builders"
	"github.com/openshift/origin-idler/pkg/apis/idling"
	idlingv1alpha2 "github.com/openshift/origin-idler/pkg/apis/idling/v1alpha2"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type MetaData struct{}

var APIMeta = MetaData{}

// GetAllApiBuilders returns all known APIGroupBuilders
// so they can be registered with the apiserver
func (MetaData) GetAllApiBuilders() []*builders.APIGroupBuilder {
	return []*builders.APIGroupBuilder{
		GetIdlingAPIBuilder(),
	}
}

// GetCRDs returns all the CRDs for known resource types
func (MetaData) GetCRDs() []v1beta1.CustomResourceDefinition {
	return []v1beta1.CustomResourceDefinition{
		idlingv1alpha2.IdlerCRD,
	}
}

func (MetaData) GetRules() []rbacv1.PolicyRule {
	return []rbacv1.PolicyRule{
		{
			APIGroups: []string{"idling.openshift.io"},
			Resources: []string{"*"},
			Verbs:     []string{"*"},
		},
	}
}

func (MetaData) GetGroupVersions() []schema.GroupVersion {
	return []schema.GroupVersion{
		{
			Group:   "idling.openshift.io",
			Version: "v1alpha2",
		},
	}
}

var idlingApiGroup = builders.NewApiGroupBuilder(
	"idling.openshift.io",
	"github.com/openshift/origin-idler/pkg/apis/idling").
	WithUnVersionedApi(idling.ApiVersion).
	WithVersionedApis(
		idlingv1alpha2.ApiVersion,
	).
	WithRootScopedKinds()

func GetIdlingAPIBuilder() *builders.APIGroupBuilder {
	return idlingApiGroup
}
