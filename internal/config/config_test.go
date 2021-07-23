package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_GetConfigPath(t *testing.T) {
	config := NewConfig()
	path := config.GetConfigPath()
	assert.Equal(t, "/home/per/.config/softteam/gitdiscover/config.json", path)
}

func TestConfig_Load(t *testing.T) {
	config := NewConfig()
	err := config.Load()
	assert.Nil(t, err)
	assert.Equal(t, "/home/per/code/gotk3-more-examples", config.Repositories[2].Path)
	assert.Equal(t, "assets/application.png", config.Repositories[2].ImagePath)
	assert.Equal(t, "Nemo", config.OpenIns[1].Name)
	assert.Equal(t, "nemo", config.OpenIns[1].Command)
	assert.Equal(t, "%PATH%", config.OpenIns[1].Argument)
	assert.Equal(t, "2006-01-02, kl. 15:04", config.DateFormat)
	assert.Equal(t, 40, config.PathColumnWidth)
}

func TestConfig_Save(t *testing.T) {
}

func TestConfig_getHomeDirectory(t *testing.T) {
	config := NewConfig()
	home := config.getHomeDirectory()
	assert.Equal(t, "/home/per", home)
}

func TestConfig_NewConfig(t *testing.T) {
	config := NewConfig()
	assert.NotNil(t, config)
}

func TestConfig_ConfigExists(t *testing.T) {
	config := NewConfig()
	assert.True(t, config.ConfigExists())
}
