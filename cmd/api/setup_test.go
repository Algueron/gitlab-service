package api

import (
	"gitlab-service/pkg/repository/gitlabrepo"
	"os"
	"testing"
)

var app GitlabService

func TestMain(m *testing.M) {
	app.GitlabRepo = &gitlabrepo.UnitTestRepo{}
	app.GitlabRepo.Connect("", "")
	os.Exit(m.Run())
}
