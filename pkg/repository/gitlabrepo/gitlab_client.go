package gitlabrepo

import (
	"fmt"
	"gitlab-service/pkg/openapi"
	"log"
	"strings"

	"github.com/xanzy/go-gitlab"
)

type GitlabClientRepo struct {
	client       *gitlab.Client
	gitlab_url   string
	gitlab_token string
}

// Connect to Gitlab
func (g *GitlabClientRepo) Connect(url string, token string) error {
	// Transform Gitlab URL to API endpoint
	gitlabUrl := url
	if !strings.HasSuffix(gitlabUrl, "/") {
		gitlabUrl = gitlabUrl + "/"
	}
	gitlabUrl = gitlabUrl + "api/v4"
	log.Println("creating client to Gitlab API endpoint at ", gitlabUrl)

	// Initialize members
	g.gitlab_url = gitlabUrl
	g.gitlab_token = token

	// Create the client
	git, err := gitlab.NewClient(g.gitlab_token, gitlab.WithBaseURL(g.gitlab_url))
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}

	// Perform a sanity check by calling the metadata API
	meta, _, err := git.Metadata.GetMetadata()
	if err != nil {
		return fmt.Errorf("failed to retrieve Gitlab metadata: %v", err)
	}
	log.Println("successfully connect to Gitlab, version is ", meta.Version)
	g.client = git
	return nil
}

// Retrieve the list of available groups
func (g *GitlabClientRepo) GetAllGroups() ([]*openapi.Group, error) {
	// Make a first call to calculate the number of pages
	_, resp, err := g.client.Groups.ListGroups(&gitlab.ListGroupsOptions{
		TopLevelOnly: gitlab.Ptr(false),
	})
	if err != nil {
		return nil, fmt.Errorf("error while retrieving groups: %v", err)
	}
	nbTotalPages := resp.TotalPages

	var retrievedGroups []*openapi.Group
	for i := 0; i < nbTotalPages; i++ {
		groups, _, err := g.client.Groups.ListGroups(&gitlab.ListGroupsOptions{
			TopLevelOnly: gitlab.Ptr(true),
			ListOptions: gitlab.ListOptions{
				Page: i,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("error while retrieving groups: %v", err)
		}
		for _, e := range groups {
			retrievedGroups = append(retrievedGroups, &openapi.Group{
				Id:       int32adr(int32(e.ID)),
				Name:     &e.Name,
				Path:     &e.Path,
				FullPath: &e.FullPath,
			})
		}
	}

	return retrievedGroups, nil
}

// Retrieve the list of subgroups of a group
func (g *GitlabClientRepo) GetGroupSubgroups(groupId int32) ([]*openapi.Group, error) {
	// Make a first call to calculate the number of pages
	_, resp, err := g.client.Groups.ListSubGroups(int(groupId), &gitlab.ListSubGroupsOptions{})
	if err != nil {
		return nil, fmt.Errorf("error while retrieving subgroups: %v", err)
	}
	nbTotalPages := resp.TotalPages

	var retrievedGroups []*openapi.Group
	for i := 0; i < nbTotalPages; i++ {
		groups, _, err := g.client.Groups.ListSubGroups(int(groupId), &gitlab.ListSubGroupsOptions{
			ListOptions: gitlab.ListOptions{
				Page: i,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("error while retrieving subgroups: %v", err)
		}
		for _, e := range groups {
			retrievedGroups = append(retrievedGroups, &openapi.Group{
				Id:       int32adr(int32(e.ID)),
				Name:     &e.Name,
				Path:     &e.Path,
				FullPath: &e.FullPath,
			})
		}
	}
	return retrievedGroups, nil
}

// Retrieve the list of projects of a group
func (g *GitlabClientRepo) GetGroupProjects(groupId int32) ([]*openapi.Project, error) {
	// Make a first call to calculate the number of pages
	_, resp, err := g.client.Groups.ListGroupProjects(int(groupId), &gitlab.ListGroupProjectsOptions{})
	if err != nil {
		return nil, fmt.Errorf("error while retrieving group's projects: %v", err)
	}
	nbTotalPages := resp.TotalPages

	var retrievedProjects []*openapi.Project
	for i := 0; i < nbTotalPages; i++ {
		projects, _, err := g.client.Groups.ListGroupProjects(int(groupId), &gitlab.ListGroupProjectsOptions{
			ListOptions: gitlab.ListOptions{
				Page: i,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("error while retrieving group's projects: %v", err)
		}
		for _, e := range projects {
			retrievedProjects = append(retrievedProjects, &openapi.Project{
				Id:            int32adr(int32(e.ID)),
				GroupId:       int32adr(int32(e.Namespace.ID)),
				Name:          &e.Name,
				DefaultBranch: &e.DefaultBranch,
				HttpUrlToRepo: &e.HTTPURLToRepo,
			})
		}
	}
	return retrievedProjects, nil
}

// Retrieve the list of projects
func (g *GitlabClientRepo) GetProjects() ([]*openapi.Project, error) {
	// Make a first call to calculate the number of pages
	_, resp, err := g.client.Projects.ListProjects(&gitlab.ListProjectsOptions{})
	if err != nil {
		return nil, fmt.Errorf("error while retrieving projects: %v", err)
	}
	nbTotalPages := resp.TotalPages

	var retrievedProjects []*openapi.Project
	for i := 0; i < nbTotalPages; i++ {
		projects, _, err := g.client.Projects.ListProjects(&gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				Page: i,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("error while retrieving projects: %v", err)
		}
		for _, e := range projects {
			retrievedProjects = append(retrievedProjects, &openapi.Project{
				Id:            int32adr(int32(e.ID)),
				GroupId:       int32adr(int32(e.Namespace.ID)),
				Name:          &e.Name,
				DefaultBranch: &e.DefaultBranch,
				HttpUrlToRepo: &e.HTTPURLToRepo,
			})
		}
	}
	return retrievedProjects, nil
}

// Create a new project
func (g *GitlabClientRepo) CreateProject(p *openapi.Project) (int32, error) {
	gp, _, err := g.client.Projects.CreateProject(&gitlab.CreateProjectOptions{
		NamespaceID: gitlab.Ptr(int(*p.GroupId)),
		Name:        p.Name,
	})
	if err != nil {
		return -1, fmt.Errorf("error while creationg project %s: %v", *p.Name, err)
	}

	return int32(gp.ID), nil
}

// Delete a single project
func (g *GitlabClientRepo) DeleteProject(projectId int32) error {
	_, err := g.client.Projects.DeleteProject(int(projectId))
	if err != nil {
		return fmt.Errorf("error while deleting project %d: %v", projectId, err)
	}

	return nil
}

// Retrieve a single project
func (g *GitlabClientRepo) GetProject(projectId int32) (*openapi.Project, error) {
	p, _, err := g.client.Projects.GetProject(int(projectId), &gitlab.GetProjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("error while getting project %d: %v", projectId, err)
	}

	project := openapi.Project{
		Id:            int32adr(int32(p.ID)),
		GroupId:       int32adr(int32(p.Namespace.ID)),
		Name:          &p.Name,
		DefaultBranch: &p.DefaultBranch,
		HttpUrlToRepo: &p.HTTPURLToRepo,
	}
	return &project, nil
}
