package gitlabrepo

import (
	"errors"
	"gitlab-service/pkg/openapi"
	"slices"
	"strings"
)

type UnitTestRepo struct {
	groups   []*openapi.Group
	projects []*openapi.Project
}

// Connect to Gitlab
func (u *UnitTestRepo) Connect(url string, token string) error {
	// Utilities functions
	stradr := func(s string) *string { return &s }
	int32adr := func(i int32) *int32 { return &i }

	// For unit testing, initialize fake data
	if len(u.groups) == 0 {
		// Create two groups
		u.groups = append(u.groups, &openapi.Group{
			Id:       int32adr(1),
			Name:     stradr("foo"),
			Path:     stradr("foo"),
			FullPath: stradr("foo"),
		})
		u.groups = append(u.groups, &openapi.Group{
			Id:       int32adr(2),
			Name:     stradr("bar"),
			Path:     stradr("bar"),
			FullPath: stradr("foo/bar"),
		})

		// Create two projects
		u.projects = append(u.projects, &openapi.Project{
			Id:            int32adr(1),
			GroupId:       int32adr(1),
			DefaultBranch: stradr("main"),
			Name:          stradr("project1"),
			HttpUrlToRepo: stradr("https://gitlab.example.com/foo/project1.git"),
		})
		u.projects = append(u.projects, &openapi.Project{
			Id:            int32adr(2),
			GroupId:       int32adr(2),
			DefaultBranch: stradr("main"),
			Name:          stradr("project2"),
			HttpUrlToRepo: stradr("https://gitlab.example.com/foo/bar/project2.git"),
		})
	}
	return nil
}

// Retrieve the list of available groups
func (u *UnitTestRepo) GetAllGroups() ([]*openapi.Group, error) {
	// Return the list of fake groups
	return u.groups, nil
}

// Retrieve the list of subgroups of a group
func (u *UnitTestRepo) GetGroupSubgroups(groupId int32) ([]*openapi.Group, error) {
	// Check if the group exists
	idx := slices.IndexFunc(u.groups, func(g *openapi.Group) bool { return *g.Id == groupId })
	if idx == -1 {
		return nil, errors.New("group does not exist")
	}

	// Return the groups with a starting full path
	var subgroups []*openapi.Group
	for _, e := range u.groups {
		if (*e.Id != groupId) && strings.HasPrefix(*e.FullPath, *u.groups[idx].FullPath) {
			subgroups = append(subgroups, e)
		}
	}
	return subgroups, nil
}

// Retrieve the list of projects of a group
func (u *UnitTestRepo) GetGroupProjects(groupId int32) ([]*openapi.Project, error) {
	// Check if the group exist
	idx := slices.IndexFunc(u.groups, func(g *openapi.Group) bool { return *g.Id == groupId })
	if idx == -1 {
		return nil, errors.New("group does not exist")
	}

	// Return the projects attached to the group
	var projects []*openapi.Project
	for _, p := range u.projects {
		if *p.GroupId == groupId {
			projects = append(projects, p)
		}
	}
	return projects, nil
}

// Retrieve the list of projects
func (u *UnitTestRepo) GetProjects() ([]*openapi.Project, error) {
	// return the list of fake projects
	return u.projects, nil
}

// Create a new project
func (u *UnitTestRepo) CreateProject(p *openapi.Project) (int, error) {
	// Check if the group exist
	idx := slices.IndexFunc(u.groups, func(g *openapi.Group) bool { return *g.Id == *p.GroupId })
	if idx == -1 {
		return -1, errors.New("group does not exist")
	}

	// Generate a new projectId
	projectWithMaxId := slices.MaxFunc(u.projects, func(a, b *openapi.Project) int { return max(int(*a.Id - *b.Id)) })
	newId := *projectWithMaxId.Id + 1
	p.Id = &newId

	// Adds the project to the fake list
	u.projects = append(u.projects, p)

	return int(*p.Id), nil
}

// Delete a single project
func (u *UnitTestRepo) DeleteProject(projectId int32) error {
	// Check if the project exist
	idx := slices.IndexFunc(u.projects, func(p *openapi.Project) bool { return *p.Id == projectId })
	if idx == -1 {
		return errors.New("project does not exist")
	}

	// Remove the project
	u.projects = append(u.projects[:idx], u.projects[idx+1:]...)

	return nil
}

// Retrieve a single project
func (u *UnitTestRepo) GetProject(projectId int32) (*openapi.Project, error) {
	// Check if the project exist
	idx := slices.IndexFunc(u.projects, func(p *openapi.Project) bool { return *p.Id == projectId })
	if idx == -1 {
		return nil, errors.New("project does not exist")
	}
	return u.projects[idx], nil
}
