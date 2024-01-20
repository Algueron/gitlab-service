package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGitlabServiceGetAllGroups(t *testing.T) {
	// Create the API request
	req, _ := http.NewRequest("GET", "/groups", nil)

	// Create a Recorder
	rr := httptest.NewRecorder()

	// Create a test server and call it
	handler := http.HandlerFunc(app.GetAllGroups)
	handler.ServeHTTP(rr, req)

	// Check return status
	if rr.Code != http.StatusOK {
		t.Errorf("GET /groups: wront status returned; expected %d but got %d", http.StatusOK, rr.Code)
	}
}
