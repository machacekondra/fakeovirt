package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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

func main() {
	port, available := os.LookupEnv("PORT")
	if !available {
		port = defaultPort
	}

	http.HandleFunc("/ovirt-engine/sso/oauth/token", SsoToken)
	http.HandleFunc(apiEndpoint("disks/"), OvirtDisks)
	http.HandleFunc(apiEndpoint("vms"), OvirtVms)
	http.HandleFunc(apiEndpoint("vms/"), OvirtVM)
	http.HandleFunc("/namespace", GetNamespace)
	http.HandleFunc(apiEndpoint("imagetransfers/"), OvirtImageTransfers)
	err := http.ListenAndServeTLS(":"+port, "server.crt", "server.key", nil)
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
	setContentType(w, xmlContentType)
	w.Write([]byte("<vms><vm id=\"123\"><name>cirrosvm</name><status>down</status></vm></vms>"))
}

// OvirtVM host Vms endpotint
func OvirtVM(w http.ResponseWriter, r *http.Request) {
	vmID := r.URL.Path[len(apiEndpoint("vms/")):]
	setContentType(w, xmlContentType)
	w.Write([]byte("<vm id=\"" + vmID + "\"><name>cirrosvm</name><status>down</status></vm>"))
}

// OvirtDisks host disks endpotint
func OvirtDisks(w http.ResponseWriter, r *http.Request) {
	diskSize, available := os.LookupEnv("DISKSIZE")
	if !available {
		diskSize = defaultSize
	}

	diskID := r.URL.Path[len(apiEndpoint("disks/")):]
	setContentType(w, xmlContentType)
	w.Write([]byte("<disk id=\"" + diskID + "\"><total_size>" + diskSize + "</total_size></disk>"))
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
	imageName, available := os.LookupEnv("IMAGENAME")
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
