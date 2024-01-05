package api

import (
	"gitlab-service/pkg/openapi"
	"gitlab-service/pkg/repository"
	"net/http"
)

type GitlabService struct {
	GitlabRepo  repository.GitlabRepo
	GitlabUrl   string
	GitlabToken string
}

// Make sure we conform to ServerInterface
var _ openapi.ServerInterface = (*GitlabService)(nil)

// Retrieve the list of available groups
// (GET /group)
func (g *GitlabService) GetAllGroups(w http.ResponseWriter, r *http.Request) {

}

// Retrieve the list of subgroups of a group
// (GET /group/{groupId}/subgroups)
func (g *GitlabService) GetGroupSubgroups(w http.ResponseWriter, r *http.Request, groupId int32) {

}

// Retrieve the list of projects of a group
// (GET /groups/{groupId}/projects)
func (g *GitlabService) GetGroupProjects(w http.ResponseWriter, r *http.Request, groupId int32) {

}

// Retrieve the list of projects
// (GET /project)
func (g *GitlabService) GetProjects(w http.ResponseWriter, r *http.Request) {

}

// Create a new project
// (POST /project)
func (g *GitlabService) CreateProject(w http.ResponseWriter, r *http.Request) {

}

// Delete a single project
// (DELETE /project/{projectId})
func (g *GitlabService) DeleteProject(w http.ResponseWriter, r *http.Request, projectId int32) {

}

// Retrieve a single project
// (GET /project/{projectId})
func (g *GitlabService) GetProject(w http.ResponseWriter, r *http.Request, projectId int32) {

}
