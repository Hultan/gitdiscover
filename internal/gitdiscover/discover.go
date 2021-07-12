package gitdiscover

import (
	gitConfig "github.com/hultan/gitdiscover/internal/config"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

type RepositoryStatus struct {
	Path   string
	Status string
	Date   *time.Time
}

type Git struct {
	Config *gitConfig.Config
}

func GitNew(config *gitConfig.Config) *Git {
	git := new(Git)
	git.Config = config
	return git
}

func (g *Git) GetRepositories() ([]RepositoryStatus, error) {
	// Get the git statuses of the paths in the config
	var gitStatuses []RepositoryStatus
	for _, basePath := range g.Config.Paths {
		gitPath := path.Join(basePath, ".git")
		status := RepositoryStatus{Path: basePath}

		if _, err := os.Stat(gitPath); os.IsNotExist(err) {
			status.Date = nil
			status.Status = "Not a git directory!"
		} else {
			gs := g.getGitStatus(basePath)
			status.Date = g.getModifiedDate(basePath)
			status.Status = strings.Replace(gs,"\n","",-1)
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
		return "failed to check git status"
	}
	return string(out)
}

// Get the modified date of a file
func (g *Git) getModifiedDate(path string) *time.Time {
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}
	date := info.ModTime()
	return &date
}

// Create format string for failed git statuses
func (g *Git) createPathFormatString() string {
	return "%-" + strconv.Itoa(g.Config.PathColumnWidth) + "s"
}

