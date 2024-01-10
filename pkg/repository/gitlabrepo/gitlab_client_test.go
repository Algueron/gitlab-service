package gitlabrepo

import "testing"

func TestGitlabClientRepoConnect(t *testing.T) {
	// Only during integration tests
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}

	// Try to connect to invalid URL (empty)
	emptyRepo := GitlabClientRepo{}
	err := emptyRepo.Connect("://", "")
	if err == nil {
		t.Errorf("no error raised while connecting to an empty URL")
	}

	// Try to connect to a non-existing URL
	wrongRepo := GitlabClientRepo{}
	err = wrongRepo.Connect("http://test:60000", "test")
	if err == nil {
		t.Errorf("no error raised while connecting to a non-existing Repository")
	}

	// Successfull connection is already tested in TestMain func
}
