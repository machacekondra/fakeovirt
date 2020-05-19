package main

import (
	"github.com/machacekondra/fakeovirt/internal/app/router"
	"log"
	"net/http"
	"os"
)

const (
	defaultPort           = "12346"
)

func main() {
	port, available := os.LookupEnv("PORT")
	if !available {
		port = defaultPort
	}

	mux:= router.CreateRouter()

	err := http.ListenAndServeTLS(":"+port, "imageio.crt", "server.key", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

