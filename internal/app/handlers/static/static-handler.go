package static

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/machacekondra/fakeovirt/internal/app/imagetransfer"
)

const (
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
	"cirros2": "cirros2",
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

// SsoLogout endpoint fake the answer to SSO logout
func SsoLogout(w http.ResponseWriter, r *http.Request) {
	setContentType(w, jsonContentType)
	w.Write([]byte("{ }"))
}

// Dynamic resource for unspecifed url
func DynamicResource(w http.ResponseWriter, r *http.Request) {
	// Remove from url /ovirt-engine/api prefix
	path := strings.TrimPrefix(r.URL.Path, "/ovirt-engine/api")
	fmt.Println(r.URL.Path)
	if r.Header.Get("Accept") == jsonContentType {
		path = fmt.Sprintf("stubs/%s/content.json", path)
		setContentType(w, jsonContentType)
	} else {
		path = fmt.Sprintf("stubs/%s/content.xml", path)
		setContentType(w, xmlContentType)
	}

	// Use the rest of url as path in to get the content
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("<error/>"))
	}
	w.Write(content)
}

// OvirtDisks host disks endpoint
func OvirtDisks(w http.ResponseWriter, r *http.Request) {
	diskSize, available := os.LookupEnv("DISKSIZE")
	if !available {
		diskSize = defaultSize
	}

	id := r.URL.Query().Get(":id")
	setContentType(w, xmlContentType)
	content, err := ioutil.ReadFile("stubs/disks/" + id + "/content")
	if err != nil {
		fmt.Println(err)
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
	imageName, available := images[imagetransfer.GetImageId(r)]
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
