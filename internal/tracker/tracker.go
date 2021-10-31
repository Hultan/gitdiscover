package tracker

import (
	discoverConfig "github.com/hultan/gitdiscover/internal/config"
)

type Tracker struct {
	Config *discoverConfig.Config
	Folders TrackedFolders
}

func NewTracker(config *discoverConfig.Config) *Tracker {
	g := &Tracker{Config: config}
	g.Refresh()

	return g
}

func (g *Tracker) Refresh() {
	// Get the git statuses of the paths in the config
	var trackedFolders TrackedFolders

	for _, configRepo := range g.Config.Repositories {
		folder := NewFolder(configRepo.Path)
		folder.SetImagePath(configRepo.ImagePath)

		trackedFolders = append(trackedFolders, folder)
	}

	g.Folders = trackedFolders
}
