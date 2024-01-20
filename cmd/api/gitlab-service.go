package api

import (
	"fmt"
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
	// Retrieve groups from Gitlab
	groups, err := g.GitlabRepo.GetAllGroups()
	if err != nil {
		g.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	g.writeJSON(w, http.StatusOK, groups)
}

// Retrieve the list of subgroups of a group
// (GET /group/{groupId}/subgroups)
func (g *GitlabService) GetGroupSubgroups(w http.ResponseWriter, r *http.Request, groupId int32) {
	groups, err := g.GitlabRepo.GetGroupSubgroups(groupId)
	if err != nil {
		g.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	g.writeJSON(w, http.StatusOK, groups)
}

// Retrieve the list of projects of a group
// (GET /groups/{groupId}/projects)
func (g *GitlabService) GetGroupProjects(w http.ResponseWriter, r *http.Request, groupId int32) {
	projects, err := g.GitlabRepo.GetGroupProjects(groupId)
	if err != nil {
		g.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	g.writeJSON(w, http.StatusOK, projects)
}

// Retrieve the list of projects
// (GET /project)
func (g *GitlabService) GetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := g.GitlabRepo.GetProjects()
	if err != nil {
		g.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	g.writeJSON(w, http.StatusOK, projects)
}

// Create a new project
// (POST /project)
func (g *GitlabService) CreateProject(w http.ResponseWriter, r *http.Request) {
	var project openapi.Project

	// Read the json payload
	err := g.readJSON(w, r, &project)
	if err != nil {
		g.errorJSON(w, fmt.Errorf("invalid project: %v", err), http.StatusBadRequest)
		return
	}

	projectId, err := g.GitlabRepo.CreateProject(&project)
	if err != nil {
		g.errorJSON(w, err, http.StatusBadRequest)
	}
	g.writeJSON(w, http.StatusCreated, projectId)
}

// Delete a single project
// (DELETE /project/{projectId})
func (g *GitlabService) DeleteProject(w http.ResponseWriter, r *http.Request, projectId int32) {
	err := g.GitlabRepo.DeleteProject(projectId)
	if err != nil {
		g.errorJSON(w, err, http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

// Retrieve a single project
// (GET /project/{projectId})
func (g *GitlabService) GetProject(w http.ResponseWriter, r *http.Request, projectId int32) {
	project, err := g.GitlabRepo.GetProject(projectId)
	if err != nil {
		g.errorJSON(w, err, http.StatusBadRequest)
	}
	g.writeJSON(w, http.StatusOK, project)
}
