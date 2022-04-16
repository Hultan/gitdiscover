package gitdiscover

import (
	"github.com/hultan/gitdiscover/internal/config"
)

// Discover keeps track of repos, external applications etc.
type Discover struct {
	Config *config.Config

	Folders              Repositories
	ExternalApplications []ExternalApplication
}

// NewDiscover creates a new Discover object.
func NewDiscover(config *config.Config) *Discover {
	g := &Discover{Config: config}
	g.Refresh()

	return g
}

// Refresh refreshes the list of repositories.
func (t *Discover) Refresh() {
	// Get the git statuses of the paths in the config
	var repositories Repositories

	for _, configRepo := range t.Config.Repositories {
		folder := newFolder(configRepo.Path)
		folder.setImagePath(configRepo.ImagePath)

		repositories = append(repositories, folder)
	}

	t.Folders = repositories
}
