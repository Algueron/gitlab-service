package main

import (
	"flag"
	"fmt"
	"gitlab-service/cmd/api"
	"gitlab-service/pkg/openapi"
	"gitlab-service/pkg/repository/gitlabrepo"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	middleware "github.com/oapi-codegen/nethttp-middleware"
)

const port = 80

func main() {
	var app api.GitlabService
	port := flag.String("port", "8080", "Port for test HTTP server")
	flag.StringVar(&app.GitlabUrl, "url", "https://gitlab.example.com", "URL of the Gitlab instance")
	flag.StringVar(&app.GitlabToken, "token", "62f6aa1b-e0b7-4a6b-ad8c-78b180b4c606", "Token for Gitlab API")
	flag.Parse()

	// Create an instance of our handler which satisfies the generated interface
	app.GitlabRepo = &gitlabrepo.GitlabClientRepo{}
	err := app.GitlabRepo.Connect(app.GitlabUrl, app.GitlabToken)
	if err != nil {
		log.Fatalf("Failed to connect gitlab client: %v", err)
	}

	swagger, err := openapi.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// This is how you set up a basic Gorilla router
	r := mux.NewRouter()

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	r.Use(middleware.OapiRequestValidator(swagger))

	// We now register our petStore above as the handler for the interface
	openapi.HandlerFromMux(&app, r)

	s := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("0.0.0.0", *port),
	}

	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}
