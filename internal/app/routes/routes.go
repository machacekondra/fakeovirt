package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/pat"

	"github.com/machacekondra/fakeovirt/internal/app/handlers/dynamic"
	"github.com/machacekondra/fakeovirt/internal/app/router"

	"github.com/machacekondra/fakeovirt/internal/app/handlers/static"
)

const (
	apiPrefix = "/ovirt-engine/api/"
)

// CreateRouter creates and configures the root http router
func CreateRouter() *router.ReplacableDelegatingRouter {
	rootRouter := router.NewReplacableDelegatingRouter()
	configurators := map[string]dynamic.HandlersConfigurator{
		"static-vms":       ConfigureVms,
		"static-sso":       ConfigureSSO,
		"static-namespace": ConfigureNamespace,
		"static-transfers": ConfigureImageTransfers,
	}
	dynamic.NewStubbingHandler(rootRouter).Configure(configurators)
	return rootRouter
}

// ConfigureSSO configures the SSO endpoint
func ConfigureSSO(router *pat.Router) {
	router.HandleFunc("/ovirt-engine/sso/oauth/token", static.SsoToken)
}

// ConfigureImageTransfers configures the image transfers endpoint
func ConfigureImageTransfers(router *pat.Router) {
	router.HandleFunc(apiEndpoint("imagetransfers"), static.OvirtImageTransfers)
}

// ConfigureImageTransfers configures the SSO endpoint
func ConfigureNamespace(router *pat.Router) {
	router.HandleFunc("/namespace", static.GetNamespace)
}

// ConfigureVms defines the default VM-related routes
func ConfigureVms(router *pat.Router) {
	router.HandleFunc(apiEndpoint("disks/{id}"), static.OvirtDisks)

	// When the endpoint is not specified try to get stub from path
	router.NotFoundHandler = http.HandlerFunc(static.DynamicResource)
}

func apiEndpoint(path string) string {
	return fmt.Sprintf("%s%s", apiPrefix, path)
}
