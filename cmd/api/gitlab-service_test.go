package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGitlabServiceHandlers(t *testing.T) {
	var tests = []struct {
		name           string
		method         string
		url            string
		jsonBody       string
		handler        http.HandlerFunc
		expectedStatus int
	}{
		{
			name:           "GetAllGroups",
			method:         "GET",
			url:            "/groups",
			jsonBody:       "",
			handler:        app.GetAllGroups,
			expectedStatus: http.StatusOK,
		},
	}

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

		// Create the handler
		handler := http.HandlerFunc(test.handler)

		// Server the request
		handler.ServeHTTP(rr, req)

		// Check the returned code
		if rr.Code != test.expectedStatus {
			t.Errorf("%s: wrong status code returned; expected %d but got %d", test.name, test.expectedStatus, rr.Code)
		}
	}
}
