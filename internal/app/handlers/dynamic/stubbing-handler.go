package dynamic

import (
	"fmt"
	"net/http"
)

type stubbingHandler struct {
	body         string
	responseCode int
}

// ServeHTTP returns body with given response code that the stubbingHandler was configured with
func (h stubbingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(h.responseCode)
	_, err := w.Write([]byte(h.body))
	if err != nil {
		fmt.Printf("Can't write the response: %v", err)
	}
}
