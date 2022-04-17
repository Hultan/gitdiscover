package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testConfigPath = "code/gitdiscover/configs/config.json"

func TestConfig_NewConfig(t *testing.T) {
	c := NewConfig()
	assert.NotNil(t, c)
}

func TestConfig_Load(t *testing.T) {
	c := NewConfig()
	err := c.Load(testConfigPath)
	assert.Nil(t, err)

	assert.Equal(t, "2006-01-02, kl. 15:04", c.DateFormat)
	assert.Equal(t, 40, c.PathColumnWidth)
}

func TestConfig_Save(t *testing.T) {
	c := NewConfig()
	_ = c.Load(testConfigPath)
	width := c.PathColumnWidth
	c.PathColumnWidth = width + 1
	c.Save(testConfigPath)
	err := c.Load(testConfigPath)
	assert.Nil(t, err, "Load should not return an error in save")
	assert.Equal(t, width+1, c.PathColumnWidth)

	// Restore values
	c.PathColumnWidth = 40
	c.Save(testConfigPath)
}

func TestConfig_GetConfigPath(t *testing.T) {
	c := NewConfig()
	path := c.GetConfigPath("")
	assert.Equal(t, "/home/per/.config/softteam/gitdiscover/config.json", path)
	path = c.GetConfigPath(testConfigPath)
	assert.Equal(t, "/home/per/"+testConfigPath, path)
}

func TestConfig_getHomeDirectory(t *testing.T) {
	c := NewConfig()
	home := c.getHomeDirectory()
	assert.Equal(t, "/home/per", home)
}

func TestConfig_AddRepository(t *testing.T) {
	c := NewConfig()
	_ = c.Load(testConfigPath)

	count := len(c.Repositories)
	c.AddRepository(
		"/home/per/code/gitdiscover/assets/",
		"/home/per/code/gitdiscover/assets/application.png",
		false,
	)
	assert.Equal(t, count+1, len(c.Repositories))
}

func TestConfig_ClearRepositories(t *testing.T) {
	c := NewConfig()
	_ = c.Load(testConfigPath)

	count := len(c.Repositories)
	c.AddRepository(
		"/home/per/code/gitdiscover/assets/",
		"/home/per/code/gitdiscover/assets/application.png",
		false,
	)
	assert.Equal(t, count+1, len(c.Repositories))
	c.ClearRepositories()
	assert.Equal(t, 0, len(c.Repositories))
}

func TestConfig_RemoveRepository(t *testing.T) {
	c := NewConfig()
	_ = c.Load(testConfigPath)

	count := len(c.Repositories)
	c.AddRepository(
		"/home/per/code/gitdiscover/assets/",
		"/home/per/code/gitdiscover/assets/application.png",
		false,
	)
	assert.Equal(t, count+1, len(c.Repositories))
	c.RemoveRepository("/home/per/code/gitdiscover/assets/")
	assert.Equal(t, 1, len(c.Repositories))
}

func TestConfig_AddExternalApplication(t *testing.T) {
	c := NewConfig()
	_ = c.Load(testConfigPath)

	count := len(c.ExternalApplications)
	c.AddExternalApplication("test", "test", "test")
	assert.Equal(t, count+1, len(c.ExternalApplications))
}

func TestConfig_ClearExternalApplications(t *testing.T) {
	c := NewConfig()
	_ = c.Load(testConfigPath)

	count := len(c.ExternalApplications)
	c.AddExternalApplication("test", "test", "test")
	assert.Equal(t, count+1, len(c.ExternalApplications))
	c.ClearExternalApplications()
	assert.Equal(t, 0, len(c.ExternalApplications))
}

func TestConfig_RemoveExternalApplication(t *testing.T) {
	c := NewConfig()
	_ = c.Load(testConfigPath)

	count := len(c.ExternalApplications)
	c.AddExternalApplication("test", "test", "test")
	assert.Equal(t, count+1, len(c.ExternalApplications))
	c.RemoveExternalApplication("test")
	assert.Equal(t, 0, len(c.ExternalApplications))
}
