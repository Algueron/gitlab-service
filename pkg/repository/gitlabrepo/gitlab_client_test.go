package gitlabrepo

import (
	"gitlab-service/pkg/openapi"
	"testing"
)

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

	// Try to connect to the integration test generated instance
	gitlabTestRepo.Connect(gitlabTestUrl, gitlabTestToken)
}

func TestGitlabClientRepoGetAllGroups(t *testing.T) {
	// Only during integration tests
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}

	// Connect to Gitlab
	gitlabTestRepo.Connect(gitlabTestUrl, gitlabTestToken)

	// Get the groups
	groups, err := gitlabTestRepo.GetAllGroups()

	// The error should always be returned nil
	if err != nil {
		t.Errorf("error is not nil: %v", err)
	}

	// We shoud have two groups
	if len(groups) != 2 {
		t.Errorf("not having 2 groups: %d", len(groups))
	}
}

func TestGitlabClientRepoGetGroupSubgroups(t *testing.T) {
	// Only during integration tests
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}

	// Connect to Gitlab
	gitlabTestRepo.Connect(gitlabTestUrl, gitlabTestToken)

	// Get for a non-existent group
	_, err := gitlabTestRepo.GetGroupSubgroups(7777777)
	if err == nil {
		t.Error("no error when getting a non-existing group")
	}

	// Get for a group with no child group
	groups, err := gitlabTestRepo.GetGroupSubgroups(int32(barGroupId))
	if err != nil {
		t.Error("error when getting an empty group")
	}
	if len(groups) != 0 {
		t.Errorf("getting subgroups for an empty group: %d", len(groups))
	}

	// Get for a group with child group
	groups, err = gitlabTestRepo.GetGroupSubgroups(int32(fooGroupId))
	if err != nil {
		t.Error("error when getting an non-empty group")
	}
	if len(groups) != 1 {
		t.Errorf("this group should have exactly one subgroup: %d", len(groups))
	}
}

func TestGitlabClientRepoGetGroupProjects(t *testing.T) {
	// Only during integration tests
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}

	// Connect to Gitlab
	gitlabTestRepo.Connect(gitlabTestUrl, gitlabTestToken)

	// Get for a non-existent group
	_, err := gitlabTestRepo.GetGroupProjects(7777777)
	if err == nil {
		t.Error("no error when getting a non-existing group")
	}

	// Get for a group without child project
	projects, err := gitlabTestRepo.GetGroupProjects(int32(fooGroupId))
	if err != nil {
		t.Error("error when getting an empty group")
	}
	if len(projects) != 0 {
		t.Errorf("getting subgroups for an empty group: %d", len(projects))
	}

	// Get for a group with child projects
	projects, err = gitlabTestRepo.GetGroupProjects(int32(barGroupId))
	if err != nil {
		t.Error("error when getting an non-empty group")
	}
	if len(projects) == 0 {
		t.Error("this group should have projects")
	}

}

func TestGitlabClientRepoGetProjects(t *testing.T) {
	// Only during integration tests
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}

	// Connect to Gitlab
	gitlabTestRepo.Connect(gitlabTestUrl, gitlabTestToken)

	// Get the projects
	projects, err := gitlabTestRepo.GetProjects()

	// The error should always be returned nil
	if err != nil {
		t.Errorf("error is not nil: %v", err)
	}

	// We shoud have two projects
	if len(projects) == 0 {
		t.Errorf("no project")
	}
}

func TestGitlabClientRepoCreateProject(t *testing.T) {
	// Only during integration tests
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}

	// Connect to Gitlab
	gitlabTestRepo.Connect(gitlabTestUrl, gitlabTestToken)

	// Try to create in an non-existent group
	_, err := gitlabTestRepo.CreateProject(&openapi.Project{
		GroupId: int32adr(7777777),
		Name:    stradr("Non existing group"),
	})
	if err == nil {
		t.Error("no error when creating in a non-existing group")
	}

	// Get the number of projects
	projects, _ := gitlabTestRepo.GetProjects()
	prevCount := len(projects)

	// Create a project
	id, err := gitlabTestRepo.CreateProject(&openapi.Project{
		Name:    stradr("New Project"),
		GroupId: int32adr(int32(barGroupId)),
	})
	if err != nil {
		t.Errorf("error while creating a valid project: %v", err)
	}

	// Compare the number of projects
	projects, _ = gitlabTestRepo.GetProjects()
	if len(projects) != prevCount+1 {
		t.Errorf("invalid number of projects, shoud be %d but is %d", prevCount+1, len(projects))
	}

	// Check if the project exist
	p, err := gitlabTestRepo.GetProject(id)
	if err != nil {
		t.Errorf("error while getting newly created project: %v", err)
	}
	if *p.Name != "New Project" {
		t.Errorf("project retrieved is not matching by name: %v", *p.Name)
	}
}

func TestGitlabClientRepoDeleteProject(t *testing.T) {
	// Only during integration tests
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}

	// Connect to Gitlab
	gitlabTestRepo.Connect(gitlabTestUrl, gitlabTestToken)

	// Create a project
	id, err := gitlabTestRepo.CreateProject(&openapi.Project{
		Name:    stradr("Project to be deleted"),
		GroupId: int32adr(int32(barGroupId)),
	})
	if err != nil {
		t.Errorf("error while creating a valid project: %v", err)
	}

	// Get the number of projects
	projects, _ := gitlabTestRepo.GetProjects()
	prevCount := len(projects)

	// Delete the project
	err = gitlabTestRepo.DeleteProject(id)
	if err != nil {
		t.Errorf("error while deleting a valid project: %v", err)
	}

	// Compare the number of projects
	projects, _ = gitlabTestRepo.GetProjects()
	if len(projects) != prevCount-1 {
		t.Errorf("invalid number of projects, shoud be %d but is %d", prevCount-1, len(projects))
	}

	// Check the project does not exist
	_, err = gitlabTestRepo.GetProject(id)
	if err == nil {
		t.Errorf("no error while getting a deleted project")
	}

	// Try to delete a non-existing project
	err = gitlabTestRepo.DeleteProject(999999)
	if err == nil {
		t.Error("no error when deleting a non-existing project")
	}
}

func TestGitlabClientRepoGetProject(t *testing.T) {
	// Only during integration tests
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}

	// Connect to Gitlab
	gitlabTestRepo.Connect(gitlabTestUrl, gitlabTestToken)

	// Get for a non-existent project
	_, err := gitlabTestRepo.GetProject(999999)
	if err == nil {
		t.Error("no error when getting a non-existing project")
	}

	// Get an existing project
	_, err = gitlabTestRepo.GetProject(int32(project1ProjectId))
	if err != nil {
		t.Errorf("error while getting a valid project: %v", err)
	}
}
