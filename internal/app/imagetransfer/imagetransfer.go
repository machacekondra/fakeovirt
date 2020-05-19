package imagetransfer

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Image struct {
	XMLName xml.Name `xml:"image"`
	Id      string   `xml:"id"`
}

type ImageTransfer struct {
	XMLName   xml.Name `xml:"image_transfer"`
	Image     Image    `xml:"image"`
	Direction string   `xml:"direction"`
	Format    string   `xml:"format"`
}

func GetImageId(r *http.Request) string {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return ""
	}

	v := ImageTransfer{}
	err = xml.Unmarshal([]byte(body), &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return ""
	}

	return v.Image.Id
}
