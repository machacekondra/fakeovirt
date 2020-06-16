package dynamic

import (
	"fmt"
	"net/http"
	"sync"
)

type response struct {
	body         string
	responseCode int
	times        int
}

type responseHolder struct {
	responses []*response
}

type stubbingHandler struct {
	mutex  *sync.Mutex
	holder *responseHolder
}

func newStubbingHandler(responses []*response) stubbingHandler {
	return stubbingHandler{holder: &responseHolder{responses: responses}, mutex: &sync.Mutex{}}
}

// ServeHTTP returns body with given response code that the stubbingHandler was configured with
func (h stubbingHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	responses := h.holder.responses
	if len(responses) > 0 {
		resp := responses[0]
		w.WriteHeader(resp.responseCode)
		_, err := w.Write([]byte(resp.body))
		if err != nil {
			fmt.Printf("Can't write the response: %v", err)
		}

		switch resp.times {
		case 0:
			break
		case 1:
			shortened := responses[1:]
			h.holder.responses = shortened
		default:
			resp.times--
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
