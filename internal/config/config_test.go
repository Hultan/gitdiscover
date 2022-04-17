package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testConfigPath = "code/gitdiscover/configs/config.json"

func TestConfig_NewConfig(t *testing.T) {
	config := NewConfig()
	assert.NotNil(t, config)
}

func TestConfig_Load(t *testing.T) {
	config := NewConfig()
	err := config.Load(testConfigPath)
	assert.Nil(t, err)

	assert.Equal(t, "2006-01-02, kl. 15:04", config.DateFormat)
	assert.Equal(t, 40, config.PathColumnWidth)
}

func TestConfig_Save(t *testing.T) {
	config := NewConfig()
	_ = config.Load(testConfigPath)
	width := config.PathColumnWidth
	config.PathColumnWidth = width + 1
	config.Save(testConfigPath)
	err := config.Load(testConfigPath)
	assert.Nil(t, err, "Load should not return an error in save")
	assert.Equal(t, width+1, config.PathColumnWidth)

	// Restore values
	config.PathColumnWidth = 40
	config.Save(testConfigPath)
}

func TestConfig_GetConfigPath(t *testing.T) {
	config := NewConfig()
	path := config.GetConfigPath("")
	assert.Equal(t, "/home/per/.config/softteam/gitdiscover/config.json", path)
	path = config.GetConfigPath(testConfigPath)
	assert.Equal(t, "/home/per/"+testConfigPath, path)
}

func TestConfig_getHomeDirectory(t *testing.T) {
	config := NewConfig()
	home := config.getHomeDirectory()
	assert.Equal(t, "/home/per", home)
}

func TestConfig_AddRepository(t *testing.T) {
	config := NewConfig()
	_ = config.Load(testConfigPath)

	count := len(config.Repositories)
	config.AddRepository("/home/per/code/gitdiscover/assets/", "/home/per/code/gitdiscover/assets/application.png")
	assert.Equal(t, count+1, len(config.Repositories))
}

func TestConfig_ClearRepositories(t *testing.T) {
	config := NewConfig()
	_ = config.Load(testConfigPath)

	count := len(config.Repositories)
	config.AddRepository("/home/per/code/gitdiscover/assets/", "/home/per/code/gitdiscover/assets/application.png")
	assert.Equal(t, count+1, len(config.Repositories))
	config.ClearRepositories()
	assert.Equal(t, 0, len(config.Repositories))
}

func TestConfig_RemoveRepository(t *testing.T) {
	config := NewConfig()
	_ = config.Load(testConfigPath)

	count := len(config.Repositories)
	config.AddRepository("/home/per/code/gitdiscover/assets/", "/home/per/code/gitdiscover/assets/application.png")
	assert.Equal(t, count+1, len(config.Repositories))
	config.RemoveRepository("/home/per/code/gitdiscover/assets/")
	assert.Equal(t, 1, len(config.Repositories))
}

func TestConfig_AddExternalApplication(t *testing.T) {
	config := NewConfig()
	_ = config.Load(testConfigPath)

	count := len(config.ExternalApplications)
	config.AddExternalApplication("test", "test", "test")
	assert.Equal(t, count+1, len(config.ExternalApplications))
}

func TestConfig_ClearExternalApplications(t *testing.T) {
	config := NewConfig()
	_ = config.Load(testConfigPath)

	count := len(config.ExternalApplications)
	config.AddExternalApplication("test", "test", "test")
	assert.Equal(t, count+1, len(config.ExternalApplications))
	config.ClearExternalApplications()
	assert.Equal(t, 0, len(config.ExternalApplications))
}

func TestConfig_RemoveExternalApplication(t *testing.T) {
	config := NewConfig()
	_ = config.Load(testConfigPath)

	count := len(config.ExternalApplications)
	config.AddExternalApplication("test", "test", "test")
	assert.Equal(t, count+1, len(config.ExternalApplications))
	config.RemoveExternalApplication("test")
	assert.Equal(t, 0, len(config.ExternalApplications))
}
