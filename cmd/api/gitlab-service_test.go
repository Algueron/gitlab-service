package api

import (
	"fmt"
	"gitlab-service/pkg/openapi"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	middleware "github.com/oapi-codegen/nethttp-middleware"
)

func TestGitlabServiceHandlers(t *testing.T) {
	var tests = []struct {
		name           string
		method         string
		url            string
		jsonBody       string
		expectedStatus int
	}{
		{
			name:           "GetAllGroups",
			method:         "GET",
			url:            "/group",
			jsonBody:       "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GetGroupSubgroups - invalid group",
			method:         "GET",
			url:            "/group/777777/subgroups",
			jsonBody:       "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "GetGroupSubgroups - empty group",
			method:         "GET",
			url:            "/group/2/subgroups",
			jsonBody:       "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GetGroupSubgroups",
			method:         "GET",
			url:            "/group/1/subgroups",
			jsonBody:       "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GetGroupProjects - invalid group",
			method:         "GET",
			url:            "/groups/777777/projects",
			jsonBody:       "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "GetGroupSubgroups - empty group",
			method:         "GET",
			url:            "/groups/2/projects",
			jsonBody:       "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GetGroupSubgroups",
			method:         "GET",
			url:            "/groups/1/projects",
			jsonBody:       "",
			expectedStatus: http.StatusOK,
		},
	}

	swagger, err := openapi.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// Create the handler
	router := mux.NewRouter()

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	router.Use(middleware.OapiRequestValidator(swagger))

	// We now register our petStore above as the handler for the interface
	openapi.HandlerFromMux(&app, router)

	for _, test := range tests {

		// Create the API request
		var req *http.Request
		if test.jsonBody == "" {
			req, _ = http.NewRequest(test.method, test.url, nil)
		} else {
			req, _ = http.NewRequest(test.method, test.url, strings.NewReader(test.jsonBody))
		}

		// Create the recorder
		rr := httptest.NewRecorder()

		// Server the request
		router.ServeHTTP(rr, req)

		// Check the returned code
		if rr.Code != test.expectedStatus {
			t.Errorf("%s: wrong status code returned; expected %d but got %d", test.name, test.expectedStatus, rr.Code)
		}
	}
}
