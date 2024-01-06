package gitlabrepo

import (
	"errors"
	"gitlab-service/pkg/openapi"

	"github.com/xanzy/go-gitlab"
)

type GitlabClientRepo struct {
	client       *gitlab.Client
	gitlab_url   string
	gitlab_token string
}

// Connect to Gitlab
func (g *GitlabClientRepo) Connect(url string, token string) error {
	return errors.New("Unimplemented")
}

// Retrieve the list of available groups
func (g *GitlabClientRepo) GetAllGroups() ([]*openapi.Group, error) {
	return nil, errors.New("Unimplemented")
}

// Retrieve the list of subgroups of a group
func (g *GitlabClientRepo) GetGroupSubgroups(groupId int32) ([]*openapi.Group, error) {
	return nil, errors.New("Unimplemented")
}

// Retrieve the list of projects of a group
func (g *GitlabClientRepo) GetGroupProjects(groupId int32) ([]*openapi.Project, error) {
	return nil, errors.New("Unimplemented")
}

// Retrieve the list of projects
func (g *GitlabClientRepo) GetProjects() ([]*openapi.Project, error) {
	return nil, errors.New("Unimplemented")
}

// Create a new project
func (g *GitlabClientRepo) CreateProject(*openapi.Project) (int32, error) {
	return -1, errors.New("Unimplemented")
}

// Delete a single project
func (g *GitlabClientRepo) DeleteProject(projectId int32) error {
	return errors.New("Unimplemented")
}

// Retrieve a single project
func (g *GitlabClientRepo) GetProject(projectId int32) (*openapi.Project, error) {
	return nil, errors.New("Unimplemented")
}
