package router

import (
	"fmt"
	"github.com/gorilla/pat"
	"github.com/machacekondra/fakeovirt/internal/app/handlers/static"
)

const (
	apiPrefix             = "/ovirt-engine/api/"
)

// CreateRouter creates and configures http router
func CreateRouter() *pat.Router {
	mux := pat.New()
	mux.HandleFunc("/ovirt-engine/sso/oauth/token", static.SsoToken)

	mux.HandleFunc(apiEndpoint("vms"), static.OvirtVms)

	mux.HandleFunc(apiEndpoint("vms/{id}"), static.OvirtResoruceHandler("vms"))
	mux.HandleFunc(apiEndpoint("storagedomains/{id}"), static.OvirtResoruceHandler("storagedomains"))
	mux.HandleFunc(apiEndpoint("vnicprofiles/{id}"), static.OvirtResoruceHandler("vnicprofiles"))
	mux.HandleFunc(apiEndpoint("networks/{id}"), static.OvirtResoruceHandler("networks"))
	mux.HandleFunc(apiEndpoint("disks/{id}"), static.OvirtDisks)

	vmSubresourceHandler := static.OvirtVMSubresource(apiEndpoint("/vms"))
	mux.HandleFunc(apiEndpoint("vms/{id}/diskattachments"), vmSubresourceHandler)
	mux.HandleFunc(apiEndpoint("vms/{id}/graphicsconsoles"), vmSubresourceHandler)
	mux.HandleFunc(apiEndpoint("vms/{id}/nics"), vmSubresourceHandler)

	mux.HandleFunc("/namespace", static.GetNamespace)
	mux.HandleFunc(apiEndpoint("imagetransfers"), static.OvirtImageTransfers)

	return mux
}

func apiEndpoint(path string) string {
	return fmt.Sprintf("%s%s", apiPrefix, path)
}
