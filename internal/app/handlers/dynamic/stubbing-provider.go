package dynamic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/machacekondra/fakeovirt/internal/app/router"

	"github.com/gorilla/pat"
	"github.com/machacekondra/fakeovirt/api/stubbing"
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

func addStubbing(router *pat.Router, stub stubbing.Stubbing) {
	code := stub.ResponseCode
	if code == 0 {
		code = http.StatusOK
	}
	handler := stubbingHandler{responseCode: code}
	if stub.ResponseBody != nil {
		handler.body = *stub.ResponseBody
	}
	router.Add(stub.Method, stub.Path, handler)
}
