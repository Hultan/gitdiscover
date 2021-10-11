package gitdiscover

import (
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	gitConfig "github.com/hultan/gitdiscover/internal/config"
)

type Repository struct {
	Name         string
	Path         string
	ImagePath    string
	Status       string
	ModifiedDate *time.Time
	IsGit        bool
}

type Git struct {
	Logger      *logrus.Logger
	Config       *gitConfig.Config
	Repositories []*Repository
}

func NewGit(config *gitConfig.Config, logger *logrus.Logger) *Git {
	git := new(Git)
	git.Config = config
	git.Logger = logger
	git.GetRepositories()
	return git
}

func (g *Git) GetRepositories() ([]*Repository, error) {
	// Get the git statuses of the paths in the config
	var directories []*Repository

	for _, repo := range g.Config.Repositories {
		basePath := repo.Path
		gitPath := path.Join(basePath, ".git")

		dir := Repository{
			Path: basePath,
			ImagePath: repo.ImagePath,
			ModifiedDate: g.getModifiedDate(repo.Path),
			Name : path.Base(repo.Path),
		}

		if _, err := os.Stat(gitPath); os.IsNotExist(err) {
			dir.Status = ""
			dir.IsGit = false
		} else {
			gs := g.getGitStatus(basePath)
			dir.Status = strings.Replace(gs, "\n", "", -1)
			dir.IsGit = true
		}

		directories = append(directories, &dir)
	}

	g.Repositories = directories

	return directories, nil
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

func (g *Git) GetRepositoryByName(name string) []*Repository {
	repos := make([]*Repository,0)

	for _, repo := range g.Repositories {
		if repo.Name == name {
			repos = append(repos, repo)
		}
	}
	return repos
}

func (g *Git) GetRepositoryByPath(path string) *Repository {
	for _, repo := range g.Repositories {
		if repo.Path == path {
			return repo
		}
	}
	return nil
}
