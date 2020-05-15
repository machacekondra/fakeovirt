package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/pat"
)

const (
	apiPrefix             = "/ovirt-engine/api/"
	defaultPort           = "12346"
	defaultNamespace      = "cdi"
	defaultAuthToken      = "thessotoken"
	jsonContentType       = "application/json"
	xmlContentType        = "application/xml"
	defaultSize           = "46137344"
	defaultImageioService = "imageio"
	defaultImageioPort    = "12345"
	defaultImageioImage   = "cirros"
)

var images = map[string]string{
	"invalid": "invalid",
}

func main() {
	port, available := os.LookupEnv("PORT")
	if !available {
		port = defaultPort
	}

	mux := pat.New()
	mux.HandleFunc("/ovirt-engine/sso/oauth/token", SsoToken)
	mux.HandleFunc(apiEndpoint("vms"), OvirtVms)

	mux.HandleFunc(apiEndpoint("vms/{id}"), OvirtResoruceHandler("vms"))
	mux.HandleFunc(apiEndpoint("storagedomains/{id}"), OvirtResoruceHandler("storagedomains"))
	mux.HandleFunc(apiEndpoint("vnicprofiles/{id}"), OvirtResoruceHandler("vnicprofiles"))
	mux.HandleFunc(apiEndpoint("networks/{id}"), OvirtResoruceHandler("networks"))
	mux.HandleFunc(apiEndpoint("disks/{id}"), OvirtDisks)

	mux.HandleFunc(apiEndpoint("vms/{id}/diskattachments"), OvirtVMSubresource)
	mux.HandleFunc(apiEndpoint("vms/{id}/graphicsconsoles"), OvirtVMSubresource)
	mux.HandleFunc(apiEndpoint("vms/{id}/nics"), OvirtVMSubresource)

	mux.HandleFunc("/namespace", GetNamespace)
	mux.HandleFunc(apiEndpoint("imagetransfers"), OvirtImageTransfers)

	err := http.ListenAndServeTLS(":"+port, "imageio.crt", "server.key", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// GetNamespace endpoint return the namespace which will be used by fake ovirt
func GetNamespace(w http.ResponseWriter, r *http.Request) {
	namespace, available := os.LookupEnv("NAMESPACE")
	if !available {
		namespace = defaultNamespace
	}

	setContentType(w, jsonContentType)
	w.Write([]byte("{\"namespace\": \"" + namespace + "\"}"))
}

// SsoToken endpoint fake the answer to SSO token request
func SsoToken(w http.ResponseWriter, r *http.Request) {
	token, available := os.LookupEnv("TOKEN")
	if !available {
		token = defaultAuthToken
	}

	setContentType(w, jsonContentType)
	w.Write([]byte("{\"access_token\":\"" + token + "\",\"scope\":\"\",\"exp\":\"9223372036854775807\",\"token_type\":\"bearer\"}"))
}

// OvirtVms host Vms endpotint
func OvirtVms(w http.ResponseWriter, r *http.Request) {
	//TODO: add support for searching with name and cluster name
	setContentType(w, xmlContentType)
	content, err := ioutil.ReadFile("vms/123/content")
	if err != nil {
		w.Write([]byte("<error/>"))
	}
	w.Write([]byte("<vms>" + string(content) + "</vms>"))
}

func OvirtVMSubresource(w http.ResponseWriter, r *http.Request) {
	vmID := r.URL.Path[len(apiEndpoint("vms/")):]
	setContentType(w, xmlContentType)
	content, err := ioutil.ReadFile("vms/" + vmID)
	if err != nil {
		w.Write([]byte("<error/>"))
	}
	w.Write(content)
}

func OvirtResoruceHandler(resource string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get(":id")
		setContentType(w, xmlContentType)
		content, err := ioutil.ReadFile(resource + "/" + id + "/content")
		if err != nil {
			w.Write([]byte("<error/>"))
		}
		w.Write(content)
	}
}

// OvirtDisks host disks endpoint
func OvirtDisks(w http.ResponseWriter, r *http.Request) {
	diskSize, available := os.LookupEnv("DISKSIZE")
	if !available {
		diskSize = defaultSize
	}

	id := r.URL.Query().Get(":id")
	setContentType(w, xmlContentType)
	content, err := ioutil.ReadFile("disks/" + id + "/content")
	if err != nil {
		w.Write([]byte("<error/>"))
	}
	w.Write([]byte(strings.ReplaceAll(string(content), "@DISKSIZE", diskSize)))
}

// OvirtImageTransfers host imagetransfer endpoint
func OvirtImageTransfers(w http.ResponseWriter, r *http.Request) {
	port, available := os.LookupEnv("PORT")
	if !available {
		port = defaultImageioPort
	}
	service, available := os.LookupEnv("SERVICE")
	if !available {
		service = defaultImageioService

	}
	imageName, available := images[GetImageId(r)]
	if !available {
		imageName = defaultImageioImage

	}
	namespace, available := os.LookupEnv("NAMESPACE")
	if !available {
		namespace = defaultNamespace
	}

	setContentType(w, xmlContentType)
	w.Write([]byte("<image_transfer id=\"64302e7f-3f08-4d32-9fe1-59b6c383acb5\"><signed_ticket>abc123</signed_ticket><phase>transferring</phase><transfer_url>https://" + service + "." + namespace + ":" + port + "/images/" + imageName + "</transfer_url></image_transfer>"))
}

func setContentType(w http.ResponseWriter, contentType string) {
	w.Header().Set("Content-Type", contentType)
}

func apiEndpoint(path string) string {
	return fmt.Sprintf("%s%s", apiPrefix, path)
}
