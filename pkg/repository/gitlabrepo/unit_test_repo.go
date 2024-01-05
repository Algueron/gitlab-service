package gitlabrepo

import (
	"errors"
	"gitlab-service/cmd/api"
)

type UnitTestRepo struct{}

// Connect to Gitlab
func (u *UnitTestRepo) Connect(url string, token string) error {
	return errors.New("Unimplemented")
}

// Retrieve the list of available groups
func (u *UnitTestRepo) GetAllGroups() ([]*api.Group, error) {
	return nil, errors.New("Unimplemented")
}

// Retrieve the list of subgroups of a group
func (u *UnitTestRepo) GetGroupSubgroups(groupId int32) ([]*api.Group, error) {
	return nil, errors.New("Unimplemented")
}

// Retrieve the list of projects of a group
func (u *UnitTestRepo) GetGroupProjects(groupId int32) ([]*api.Project, error) {
	return nil, errors.New("Unimplemented")
}

// Retrieve the list of projects
func (u *UnitTestRepo) GetProjects() ([]*api.Project, error) {
	return nil, errors.New("Unimplemented")
}

// Create a new project
func (u *UnitTestRepo) CreateProject(*api.Project) error {
	return errors.New("Unimplemented")
}

// Delete a single project
func (u *UnitTestRepo) DeleteProject(projectId int32) error {
	return errors.New("Unimplemented")
}

// Retrieve a single project
func (u *UnitTestRepo) GetProject(projectId int32) (*api.Project, error) {
	return nil, errors.New("Unimplemented")
}
