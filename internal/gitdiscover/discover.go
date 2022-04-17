package gitdiscover

import (
	"github.com/hultan/gitdiscover/internal/config"
)

// Discover keeps track of repos, external applications etc.
type Discover struct {
	Config *config.Config

	Repositories         Repositories
	ExternalApplications []*ExternalApplication
}

// ExternalApplication : An external application in the config
type ExternalApplication struct {
	Name     string
	Command  string
	Argument string
}

// NewDiscover creates a new Discover object.
func NewDiscover(config *config.Config) *Discover {
	g := &Discover{Config: config}
	g.Refresh()

	return g
}

// Refresh refreshes the list of repositories.
func (d *Discover) Refresh() {
	// Git Repositories
	var repositories Repositories
	for _, configRepo := range d.Config.Repositories {
		folder := newFolder(configRepo.Path)
		folder.setImagePath(configRepo.ImagePath)

		repositories = append(repositories, folder)
	}
	d.Repositories = repositories

	// External applications
	var apps []*ExternalApplication
	for _, application := range d.Config.ExternalApplications {
		apps = append(apps, &ExternalApplication{
			Name:     application.Name,
			Command:  application.Command,
			Argument: application.Argument,
		})
	}
	d.ExternalApplications = apps
}

// Save saves the Discover object to the config file
func (d *Discover) Save() {
	d.Config.ClearRepositories()
	for _, repository := range d.Repositories {
		d.Config.AddRepository(repository.path, repository.imagePath)
	}

	d.Config.ClearExternalApplications()
	for _, application := range d.ExternalApplications {
		d.Config.AddExternalApplication(
			application.Name,
			application.Command,
			application.Argument,
		)
	}
	d.Config.Save("")
}

// ClearRepositories clears the slice of repositories
func (d *Discover) ClearRepositories() {
	d.Config.ClearRepositories()
	d.Refresh()
}

// AddRepository adds a new repository
func (d *Discover) AddRepository(path, imagePath string) {
	d.Config.AddRepository(path, imagePath)
	d.Refresh()
}

// RemoveRepository adds a new repository
func (d *Discover) RemoveRepository(path string) {
	d.Config.RemoveRepository(path)
	d.Refresh()
}

// GetRepositoryByIndex gets an external application by index
func (d *Discover) GetRepositoryByIndex(i int) *Repository {
	return d.Repositories[i]
}

// GetExternalApplicationByName gets an external application by name
func (d *Discover) GetExternalApplicationByName(name string) *ExternalApplication {
	ea := d.Config.GetExternalApplicationByName(name)
	return &ExternalApplication{
		Name:     ea.Name,
		Command:  ea.Command,
		Argument: ea.Argument,
	}
}

// GetExternalApplicationByName gets an external application by name
func (d *Discover) GetExternalApplicationByIndex(i int) *ExternalApplication {
	return d.ExternalApplications[i]
}

// ClearExternalApplications clears the slice of external applications
func (d *Discover) ClearExternalApplications() {
	d.Config.ClearExternalApplications()
	d.Refresh()
}

// AddExternalApplication adds an external application
func (d *Discover) AddExternalApplication(name, command, argument string) {
	d.Config.AddExternalApplication(name, command, argument)
	d.Refresh()
}

// RemoveExternalApplication adds a new extenal application
func (d *Discover) RemoveExternalApplication(name string) {
	d.Config.RemoveExternalApplication(name)
	d.Refresh()
}

// GetDateFormat returns the date format
func (d *Discover) GetDateFormat() string {
	return d.Config.DateFormat
}

// GetConfigPath returns the config path
func (d *Discover) GetConfigPath() string {
	return d.Config.GetConfigPath("")
}
