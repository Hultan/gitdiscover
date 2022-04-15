package gitdiscover

// Tracker helps to create the TrackedFolders.
type Tracker struct {
	Config  *Config
	Folders TrackedFolders
}

// NewTracker creates a new Tracker.
func NewTracker(config *Config) *Tracker {
	g := &Tracker{Config: config}
	g.Refresh()

	return g
}

// Refresh refreshes the tracker.
func (t *Tracker) Refresh() {
	// Get the git statuses of the paths in the config
	var trackedFolders TrackedFolders

	for _, configRepo := range t.Config.Repositories {
		folder := newFolder(configRepo.Path)
		folder.setImagePath(configRepo.ImagePath)

		trackedFolders = append(trackedFolders, folder)
	}

	t.Folders = trackedFolders
}
