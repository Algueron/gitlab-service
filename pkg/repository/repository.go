package repository

import "gitlab-service/cmd/api"

type GitlabRepo interface {
	// Connect to Gitlab
	Connect(url string, token string) error

	// Retrieve the list of available groups
	GetAllGroups() ([]*api.Group, error)

	// Retrieve the list of subgroups of a group
	GetGroupSubgroups(groupId int32) ([]*api.Group, error)

	// Retrieve the list of projects of a group
	GetGroupProjects(groupId int32) ([]*api.Project, error)

	// Retrieve the list of projects
	GetProjects() ([]*api.Project, error)

	// Create a new project
	CreateProject(*api.Project) error

	// Delete a single project
	DeleteProject(projectId int32) error

	// Retrieve a single project
	GetProject(projectId int32) (*api.Project, error)
}
