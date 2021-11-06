package tracker

import (
	discoverConfig "github.com/hultan/gitdiscover/internal/config"
)

// Tracker helps to create the TrackedFolders.
type Tracker struct {
	Config  *discoverConfig.Config
	Folders TrackedFolders
}

// NewTracker creates a new Tracker.
func NewTracker(config *discoverConfig.Config) *Tracker {
	g := &Tracker{Config: config}
	g.Refresh()

	return g
}

// Refresh refreshes the tracker.
func (g *Tracker) Refresh() {
	// Get the git statuses of the paths in the config
	var trackedFolders TrackedFolders

	for _, configRepo := range g.Config.Repositories {
		folder := newFolder(configRepo.Path)
		folder.setImagePath(configRepo.ImagePath)

		trackedFolders = append(trackedFolders, folder)
	}

	g.Folders = trackedFolders
}
