package gitlabrepo

import (
	"gitlab-service/pkg/repository"
	"os"
	"testing"
)

var testRepo repository.GitlabRepo

func TestMain(m *testing.M) {
	// Initiate the repository
	testRepo = &UnitTestRepo{}
	testRepo.Connect("", "")

	// Run the tests
	code := m.Run()

	// Exit properly
	os.Exit(code)
}

func TestUnitTestRepoConnect(t *testing.T) {
	// Call test repo
	testRepo.Connect("", "")

	// Check we have groups
	groups, _ := testRepo.GetAllGroups()
	if len(groups) == 0 {
		t.Error("Groups is empty in repository")
	}

	// Check we have projects
	projects, _ := testRepo.GetProjects()
	if len(projects) == 0 {
		t.Error("Projects is empty in repository")
	}
}

func TestUnitTestRepoGetAllGroups(t *testing.T) {
	// Get the groups
	groups, err := testRepo.GetAllGroups()

	// The error should always be returned nil
	if err != nil {
		t.Errorf("error is not nil: %v", err)
	}

	// We shoud have two groups
	if len(groups) != 2 {
		t.Errorf("not having 2 groups: %d", len(groups))
	}
}

func TestUnitTestRepoGetGroupSubgroups(t *testing.T) {
	// Get for an unexisting group
	_, err := testRepo.GetGroupSubgroups(3)
	if err == nil {
		t.Error("no error when getting a non-existing group")
	}

	// Get for a group with no child group
	groups, err := testRepo.GetGroupSubgroups(2)
	if err != nil {
		t.Error("error when getting an empty group")
	}
	if len(groups) != 0 {
		t.Errorf("getting subgroups for an empty groups: %d", len(groups))
	}

	// Get for a group with child group
	groups, err = testRepo.GetGroupSubgroups(1)
	if err != nil {
		t.Error("error when getting an non-empty group")
	}
	if len(groups) != 1 {
		t.Errorf("this group should have exactly one subgroup: %d", len(groups))
	}
}
