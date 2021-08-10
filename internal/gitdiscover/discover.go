package gitdiscover

import (
	gitConfig "github.com/hultan/gitdiscover/internal/config"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

// TODO : Rename RepositoryStatus, do not need to be a repository anymore

type RepositoryStatus struct {
	Path      string
	ImagePath string
	Status    string
	Date      *time.Time
	IsGit     bool
}

type Git struct {
	Logger *logrus.Logger
	Config *gitConfig.Config
}

func NewGit(config *gitConfig.Config, logger *logrus.Logger) *Git {
	git := new(Git)
	git.Config = config
	git.Logger = logger
	return git
}

func (g *Git) GetRepositories() ([]RepositoryStatus, error) {
	// Get the git statuses of the paths in the config
	var gitStatuses []RepositoryStatus
	for _, repo := range g.Config.Repositories {
		basePath := repo.Path
		gitPath := path.Join(basePath, ".git")
		status := RepositoryStatus{Path: basePath, ImagePath: repo.ImagePath}

		if _, err := os.Stat(gitPath); os.IsNotExist(err) {
			status.Date = g.getModifiedDate(basePath)
			status.Status = ""
			status.IsGit = false
		} else {
			gs := g.getGitStatus(basePath)
			status.Date = g.getModifiedDate(basePath)
			status.Status = strings.Replace(gs, "\n", "", -1)
			status.IsGit = true
		}

		gitStatuses = append(gitStatuses, status)
	}

	return gitStatuses, nil
}

// Get the git status
func (g *Git) getGitStatus(path string) string {
	cmd := exec.Command("/home/per/bin/gitprompt-go")
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		g.Logger.Error("failed to check git status")
		return ""
	}
	return string(out)
}

// Get the modified date of a file
func (g *Git) getModifiedDate(path string) *time.Time {
	info, err := os.Stat(path)
	if err != nil {
		g.Logger.Error("Failed to find modified date for path : ", path)
		return nil
	}
	date := info.ModTime()
	return &date
}

// Create format string for failed git statuses
func (g *Git) createPathFormatString() string {
	return "%-" + strconv.Itoa(g.Config.PathColumnWidth) + "s"
}
