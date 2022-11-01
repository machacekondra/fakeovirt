package dynamic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/machacekondra/fakeovirt/internal/app/router"

	"github.com/gorilla/pat"
	"github.com/machacekondra/fakeovirt/pkg/api/stubbing"
)

// StubbingProvider allows for dynamic response stubbing
type StubbingProvider struct {
	handlerConfigurators map[string]HandlersConfigurator
	router               *router.ReplacableDelegatingRouter
}

// HandlersConfigurator creates initial router configuration
type HandlersConfigurator func(*pat.Router)

// NewStubbingHandler creates new stubbing handler
func NewStubbingHandler(router *router.ReplacableDelegatingRouter) *StubbingProvider {
	handler := StubbingProvider{
		router: router,
	}
	return &handler
}

// Configure adds paths to handlers serving stubbing requests and ones made by the handlersConfigurators to the router held by the StubbingProvider
func (h *StubbingProvider) Configure(handlersConfigurators map[string]HandlersConfigurator) {
	h.handlerConfigurators = handlersConfigurators
	h.configure()
	h.configureProvidedHandlers()
}

func (h *StubbingProvider) configureProvidedHandlers() {
	for key := range h.handlerConfigurators {
		handlersConfigurator := h.handlerConfigurators[key]
		handlersConfigurator(h.router.Delegate())
	}
}

func (h *StubbingProvider) configure() {
	delegate := h.router.Delegate()
	delegate.Post("/stub", h.Stub())
	delegate.Post("/reset", h.Reset())
}

// Reset restores original paths in the router held by the StubbingProvider
func (h *StubbingProvider) Reset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newRouter := h.purge()

		h.configure()
		c := r.URL.Query().Get("configurators")
		if c == "" {
			h.configureProvidedHandlers()
		} else {
			configurations := strings.Split(c, ",")
			for _, key := range configurations {
				handlersConfigurator := h.handlerConfigurators[key]
				if handlersConfigurator != nil {
					handlersConfigurator(newRouter)
				}
			}
		}
	}
}

func (h *StubbingProvider) purge() *pat.Router {
	newRouter := pat.New()
	h.router.Set(newRouter)
	return newRouter
}

// Stub adds new routes to the router held by the StubbingProvider
func (h *StubbingProvider) Stub() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stubbings := stubbing.Stubbings{}
		err := json.NewDecoder(r.Body).Decode(&stubbings)
		if err != nil {
			fmt.Printf("Unable to decode: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for _, stub := range stubbings {
			addStubbing(h.router.Delegate(), stub)
		}

		w.WriteHeader(200)
	}
}

// Stub adds new routes to the router held by the StubbingProvider
func (h *StubbingProvider) AddStaticStubs() {
	stubbings := stubbing.Stubbings{}
	// err := json.NewDecoder(r.Body).Decode(&stubbings)
	err := filepath.Walk("stubs",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Base(path) == "stub.json" {
				new := stubbing.Stubbings{}
				data, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				reader := bytes.NewReader(data)
				err = json.NewDecoder(reader).Decode(&stubbings)
				if err != nil {
					return err
				}
				stubbings = append(stubbings, new...)
			}
			return nil
		})

	if err != nil {
		fmt.Printf("Unable to decode: %v\n", err)
	}

	for _, stub := range stubbings {
		addStubbing(h.router.Delegate(), stub)
	}
}

func addStubbing(router *pat.Router, stub stubbing.Stubbing) {
	var responses []*response
	for _, rs := range stub.Responses {
		resp := response{}
		code := rs.ResponseCode
		if code == 0 {
			code = http.StatusOK
		}
		resp.responseCode = code
		if rs.ResponseBody != nil {
			resp.body = *rs.ResponseBody
		}
		resp.times = rs.Times
		responses = append(responses, &resp)
	}

	router.Add(stub.Method, stub.Path, newStubbingHandler(responses))
}
