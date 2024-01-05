package repository

import "gitlab-service/pkg/openapi"

type GitlabRepo interface {
	// Connect to Gitlab
	Connect(url string, token string) error

	// Retrieve the list of available groups
	GetAllGroups() ([]*openapi.Group, error)

	// Retrieve the list of subgroups of a group
	GetGroupSubgroups(groupId int32) ([]*openapi.Group, error)

	// Retrieve the list of projects of a group
	GetGroupProjects(groupId int32) ([]*openapi.Project, error)

	// Retrieve the list of projects
	GetProjects() ([]*openapi.Project, error)

	// Create a new project
	CreateProject(*openapi.Project) error

	// Delete a single project
	DeleteProject(projectId int32) error

	// Retrieve a single project
	GetProject(projectId int32) (*openapi.Project, error)
}
