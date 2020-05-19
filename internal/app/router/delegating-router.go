package router

import (
	"net/http"
	"sync/atomic"

	"github.com/gorilla/pat"
)

// ReplacableDelegatingRouter supports replacing defined routes
type ReplacableDelegatingRouter struct {
	delegate atomic.Value
}

// NewReplacableDelegatingRouter creates new router that supports replacing defined routes
func NewReplacableDelegatingRouter() *ReplacableDelegatingRouter {
	router := ReplacableDelegatingRouter{}
	router.delegate.Store(pat.New())
	return &router
}

// Set replaces the delegate with given configuration
func (rdr *ReplacableDelegatingRouter) Set(newRouter *pat.Router) {
	rdr.delegate.Store(newRouter)
}

// ServeHTTP implements http.Handler.ServeHTTP
func (rdr *ReplacableDelegatingRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rdr.Delegate().ServeHTTP(w, r)
}

// Delegate returns delegate router
func (rdr *ReplacableDelegatingRouter) Delegate() *pat.Router {
	return rdr.delegate.Load().(*pat.Router)
}
