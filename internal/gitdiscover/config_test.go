package gitdiscover

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_NewConfig(t *testing.T) {
	config := NewConfig()
	assert.NotNil(t, config)
}

func TestConfig_GetConfigPath(t *testing.T) {
	config := NewConfig()
	path := config.GetConfigPath()
	assert.Equal(t, "/home/per/.config/softteam/gitdiscover/config.json", path)
}

func TestConfig_Load(t *testing.T) {
	config := NewConfig()
	err := config.Load()
	assert.Nil(t, err)

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
