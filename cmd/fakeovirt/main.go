package main

import (
	"log"
	"net/http"
	"os"

	"github.com/machacekondra/fakeovirt/internal/app/routes"
)

const (
	defaultPort = "12346"
)

func main() {
	port, available := os.LookupEnv("PORT")
	if !available {
		port = defaultPort
	}
	err := http.ListenAndServeTLS(":"+port, "imageio.crt", "server.key", routes.CreateRouter())
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
