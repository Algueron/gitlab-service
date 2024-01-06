package gitlabrepo

import (
	"gitlab-service/pkg/openapi"
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
	// Get for a non-existent group
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
		t.Errorf("getting subgroups for an empty group: %d", len(groups))
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

func TestUnitTestRepoGetGroupProjects(t *testing.T) {
	// Get for a non-existent group
	_, err := testRepo.GetGroupProjects(3)
	if err == nil {
		t.Error("no error when getting a non-existing group")
	}

	// Get for a group without child project
	projects, err := testRepo.GetGroupProjects(1)
	if err != nil {
		t.Error("error when getting an empty group")
	}
	if len(projects) != 0 {
		t.Errorf("getting subgroups for an empty group: %d", len(projects))
	}

	// Get for a group with child projects
	projects, err = testRepo.GetGroupProjects(2)
	if err != nil {
		t.Error("error when getting an non-empty group")
	}
	if len(projects) == 0 {
		t.Error("this group should have projects")
	}
}

func TestUnitTestRepoGetProjects(t *testing.T) {
	// Get the groups
	projects, err := testRepo.GetProjects()

	// The error should always be returned nil
	if err != nil {
		t.Errorf("error is not nil: %v", err)
	}

	// We shoud have two groups
	if len(projects) == 0 {
		t.Errorf("no project")
	}
}

func TestUnitRepoCreateProject(t *testing.T) {
	// Try to create in an non-existent group
	_, err := testRepo.CreateProject(&openapi.Project{
		GroupId: int32adr(3),
		Name:    stradr("Non existing group"),
	})
	if err == nil {
		t.Error("no error when creating in a non-existing group")
	}

	// Get the number of projects
	projects, _ := testRepo.GetProjects()
	prevCount := len(projects)

	// Create a project
	id, err := testRepo.CreateProject(&openapi.Project{
		Name:          stradr("New Project"),
		GroupId:       int32adr(2),
		DefaultBranch: stradr("main"),
		HttpUrlToRepo: stradr("https://gitlab.example.com/foo/bar/new-project.git"),
	})
	if err != nil {
		t.Errorf("error while creating a valid project: %v", err)
	}

	// Compare the number of projects
	projects, _ = testRepo.GetProjects()
	if len(projects) != prevCount+1 {
		t.Errorf("invalid number of projects, shoud be %d but is %d", prevCount+1, len(projects))
	}

	// Check if the project exist
	p, err := testRepo.GetProject(id)
	if err != nil {
		t.Errorf("error while getting newly created project: %v", err)
	}
	if *p.Name != "New Project" {
		t.Errorf("project retrieved is not matching by name: %v", *p.Name)
	}
}

func TestUnitRepoDeleteProject(t *testing.T) {
	// Create a project
	id, err := testRepo.CreateProject(&openapi.Project{
		Name:          stradr("Project to be deleted"),
		GroupId:       int32adr(2),
		DefaultBranch: stradr("main"),
		HttpUrlToRepo: stradr("https://gitlab.example.com/foo/bar/new-to-be-deleted.git"),
	})
	if err != nil {
		t.Errorf("error while creating a valid project: %v", err)
	}

	// Get the number of projects
	projects, _ := testRepo.GetProjects()
	prevCount := len(projects)

	// Delete the project
	err = testRepo.DeleteProject(id)
	if err != nil {
		t.Errorf("error while deleting a valid project: %v", err)
	}

	// Compare the number of projects
	projects, _ = testRepo.GetProjects()
	if len(projects) != prevCount-1 {
		t.Errorf("invalid number of projects, shoud be %d but is %d", prevCount-1, len(projects))
	}

	// Check the project does not exist
	_, err = testRepo.GetProject(id)
	if err == nil {
		t.Errorf("no error while getting a deleted project")
	}

	// Try to delete a non-existing project
	err = testRepo.DeleteProject(999)
	if err == nil {
		t.Error("no error when deleting a non-existing project")
	}
}

func TestUnitRepoGetProject(t *testing.T) {
	// Get for a non-existent project
	_, err := testRepo.GetProject(999)
	if err == nil {
		t.Error("no error when getting a non-existing project")
	}

	// Get an existing project
	_, err = testRepo.GetProject(1)
	if err != nil {
		t.Errorf("error while getting a valid project: %v", err)
	}
}
