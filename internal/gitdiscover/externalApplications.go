package gitdiscover

import (
	"github.com/hultan/gitdiscover/internal/config"
)

// ExternalApplication : An external application in the config
type ExternalApplication struct {
	Name     string
	Command  string
	Argument string
}

func NewExternalEpplicationsFromConfig(config *config.Config) []ExternalApplication {
	var apps []ExternalApplication
	for _, application := range config.ExternalApplications {
		apps = append(apps, ExternalApplication{
			Name:     application.Name,
			Command:  application.Command,
			Argument: application.Argument,
		})
	}

	return apps
}