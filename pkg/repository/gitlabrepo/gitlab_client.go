package gitlabrepo

import (
	"errors"
	"gitlab-service/cmd/api"

	"github.com/xanzy/go-gitlab"
)

type GitlabClientRepo struct {
	client *gitlab.Client
	url    string
	token  string
}

// Connect to Gitlab
func (g *GitlabClientRepo) Connect(url string, token string) error {
	return errors.New("Unimplemented")
}

// Retrieve the list of available groups
func (g *GitlabClientRepo) GetAllGroups() ([]*api.Group, error) {
	return nil, errors.New("Unimplemented")
}

// Retrieve the list of subgroups of a group
func (g *GitlabClientRepo) GetGroupSubgroups(groupId int32) ([]*api.Group, error) {
	return nil, errors.New("Unimplemented")
}

// Retrieve the list of projects of a group
func (g *GitlabClientRepo) GetGroupProjects(groupId int32) ([]*api.Project, error) {
	return nil, errors.New("Unimplemented")
}

// Retrieve the list of projects
func (g *GitlabClientRepo) GetProjects() ([]*api.Project, error) {
	return nil, errors.New("Unimplemented")
}

// Create a new project
func (g *GitlabClientRepo) CreateProject(*api.Project) error {
	return errors.New("Unimplemented")
}

// Delete a single project
func (g *GitlabClientRepo) DeleteProject(projectId int32) error {
	return errors.New("Unimplemented")
}

// Retrieve a single project
func (g *GitlabClientRepo) GetProject(projectId int32) (*api.Project, error) {
	return nil, errors.New("Unimplemented")
}
