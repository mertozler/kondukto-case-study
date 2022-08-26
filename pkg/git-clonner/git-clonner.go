package git_clonner

import (
	"github.com/go-git/go-git/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
)

type GitClonner struct {
	TargetDirectory string
}

func NewGitClonner() *GitClonner {
	dir, _ := os.Getwd()
	gitCloner := &GitClonner{
		TargetDirectory: dir + "/tmp/src/",
	}
	return gitCloner
}

func (g *GitClonner) CloneRepo(url string) (string, error) {
	projectId := uuid.New()
	logrus.Info("Clonning repository for scan id ", projectId)
	_, err := git.PlainClone(g.TargetDirectory+projectId.String(), false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		logrus.Errorf("Error while cloning repository for scan id %v", projectId, err)
		return projectId.String(), err
	}
	return projectId.String(), nil
}
