// Note: Ignore this (but don't delete it) if you are using CRDs.  If using
// CRDs this file is necessary to generate docs.

package main

import (
	// Make sure dep tools picks up these dependencies
	_ "github.com/go-openapi/loads"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubernetes-sigs/kubebuilder/pkg/cmd/server"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // Enable cloud provider auth

	"github.com/openshift/origin-idler/pkg/apis"
	"github.com/openshift/origin-idler/pkg/openapi"
)

// Extension (aggregated) apiserver main.
func main() {
	version := "v0"
	server.StartApiServer("/registry/openshift.io", apis.APIMeta.GetAllApiBuilders(), openapi.GetOpenAPIDefinitions, "Api", version)
}
