package gitlabrepo

import (
	"context"
	"flag"
	"fmt"
	"gitlab-service/pkg/repository"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/xanzy/go-gitlab"
)

var (
	unitTestRepo      repository.GitlabRepo
	gitlabTestRepo    repository.GitlabRepo
	gitlabTestUrl     string
	gitlabTestToken   string
	fooGroupId        int
	barGroupId        int
	project1ProjectId int
	project2ProjectId int
)

type gitlabContainer struct {
	testcontainers.Container
	URI string
}

func SetupGitlabContainer(ctx context.Context) (*gitlabContainer, error) {
	waitStrategy := wait.HTTPStrategy{
		Port:              "80",
		Path:              "/",
		StatusCodeMatcher: func(status int) bool { return status == http.StatusOK },
		PollInterval:      5 * time.Second,
	}
	waitStrategy.WithStartupTimeout(3 * time.Minute)
	req := testcontainers.ContainerRequest{
		Image:        "gitlab/gitlab-ce:16.7.0-ce.0",
		ShmSize:      256 * (1 << 20),
		ExposedPorts: []string{"80"},
		Hostname:     "gitlab.example.com",
		Env: map[string]string{
			"GITLAB_OMNIBUS_CONFIG": "external_url 'http://my.domain.com/';",
		},
		WaitingFor: &waitStrategy,
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "80")
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())

	return &gitlabContainer{Container: container, URI: uri}, nil
}

// Retrieve the root password by greping the appropriate file
func GetGitlabRootPassword(ctx context.Context, container testcontainers.Container) (string, error) {
	c, reader, err := container.Exec(ctx, []string{"grep", "Password:", "/etc/gitlab/initial_root_password"})
	if err != nil {
		return "", fmt.Errorf("error while executing command to retrieve root password: %v", err)
	}
	if c != 0 {
		return "", fmt.Errorf("error while retrieving password, got return code %d", c)
	}

	// Parse the command output
	buf := new(strings.Builder)
	_, err = io.Copy(buf, reader)
	if err != nil {
		return "", fmt.Errorf("Error while retrieving output of password command: %v", err)
	}
	passwordLine := buf.String()

	// Compute the password
	chunks := strings.Split(passwordLine, " ")
	rootPassword := strings.TrimSpace(chunks[len(chunks)-1])

	return rootPassword, nil
}

func CreateRootToken(gitlabURL string, rootPassword string) (string, error) {
	// Create the Gitlab Client
	gitlabClient, err := gitlab.NewBasicAuthClient("root", rootPassword, gitlab.WithBaseURL(gitlabURL))
	if err != nil {
		return "", fmt.Errorf("Error while creating Gitlab client: %v", err)
	}

	// Get the root user ID
	u, _, err := gitlabClient.Users.CurrentUser()
	if err != nil {
		return "", fmt.Errorf("Error while getting current user: %v", err)
	}

	// Create the personal access token
	t, _, err := gitlabClient.Users.CreatePersonalAccessToken(u.ID, &gitlab.CreatePersonalAccessTokenOptions{
		Name:   gitlab.Ptr("test_token"),
		Scopes: gitlab.Ptr([]string{"api"}),
	})
	if err != nil {
		return "", fmt.Errorf("Error while creating personal access token: %v", err)
	}

	return t.Token, nil
}

func CreateGitlabTestDataset(gitlabURL string, gitlabToken string) error {
	// Create the Gitlab client
	gitlabClient, err := gitlab.NewClient(gitlabToken, gitlab.WithBaseURL(gitlabURL))
	if err != nil {
		return fmt.Errorf("Error while creating Gitlab client: %v", err)
	}

	// Create the first group
	foo, _, err := gitlabClient.Groups.CreateGroup(&gitlab.CreateGroupOptions{
		Name: gitlab.Ptr("foo"),
		Path: gitlab.Ptr("foo"),
	})
	if err != nil {
		return fmt.Errorf("Error while creating foo group: %v", err)
	}
	fooGroupId = foo.ID

	// Create the second group
	bar, _, err := gitlabClient.Groups.CreateGroup(&gitlab.CreateGroupOptions{
		ParentID: gitlab.Ptr(foo.ID),
		Name:     gitlab.Ptr("bar"),
		Path:     gitlab.Ptr("bar"),
	})
	if err != nil {
		return fmt.Errorf("Error while creating foo/bar group: %v", err)
	}
	barGroupId = bar.ID

	// Create the first project
	proj1, _, err := gitlabClient.Projects.CreateProject(&gitlab.CreateProjectOptions{
		Name:        gitlab.Ptr("project1"),
		NamespaceID: gitlab.Ptr(bar.ID),
	})
	if err != nil {
		return fmt.Errorf("Error while creating first project: %v", err)
	}
	project1ProjectId = proj1.ID

	// Create the second project
	proj2, _, err := gitlabClient.Projects.CreateProject(&gitlab.CreateProjectOptions{
		Name:        gitlab.Ptr("project2"),
		NamespaceID: gitlab.Ptr(bar.ID),
	})
	if err != nil {
		return fmt.Errorf("Error while creating second project: %v", err)
	}
	project2ProjectId = proj2.ID

	return nil
}

func TestMain(m *testing.M) {
	// Parse the flags
	flag.Parse()

	// Always initiate the unit test repository
	unitTestRepo = &UnitTestRepo{}
	unitTestRepo.Connect("", "")

	// If we're running the integration tests, also initialize the gitlab repository
	if testing.Short() == false {
		// Retrieve the background context
		ctx := context.Background()

		// Start a Gitlab CE container
		gitlabContainer, err := SetupGitlabContainer(ctx)
		if err != nil {
			log.Fatalf("Could not start gitlab container: %v", err)
		}

		// Clean up the container after the test is complete
		defer func() {
			if err := gitlabContainer.Terminate(ctx); err != nil {
				log.Fatalf("Failed to terminate container: %s", err)
			}

		}()

		// Try to get the Home page
		gitlabTestUrl = gitlabContainer.URI
		log.Println("Container available at:", gitlabTestUrl)
		resp, _ := http.Get(gitlabContainer.URI)
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Expected status code %d. Got %d.", http.StatusOK, resp.StatusCode)
		}

		// Retrieve the root password
		rootPassword, err := GetGitlabRootPassword(ctx, gitlabContainer)
		if err != nil {
			log.Fatalf("Error while getting root password: %v", err)
		}
		log.Println("Root password is ", rootPassword)

		// Create a personal access token for root
		gitlabTestToken, err = CreateRootToken(gitlabContainer.URI, rootPassword)
		if err != nil {
			log.Fatalf("Error while creating root token: %v", err)
		}

		// Create the fake groups
		err = CreateGitlabTestDataset(gitlabTestUrl, gitlabTestToken)
		if err != nil {
			log.Fatalf("Error while creating Gitlab test dataset token: %v", err)
		}

		// Initialize the Gitlab Repository
		gitlabTestRepo = &GitlabClientRepo{}
	}

	// Run the tests
	code := m.Run()

	// Exit properly
	os.Exit(code)
}
