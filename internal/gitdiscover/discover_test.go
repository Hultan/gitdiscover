package gitdiscover

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hultan/gitdiscover/internal/config"
)

const testConfigPath = "code/gitdiscover/configs/config.json"

func Test_NewDiscover(t *testing.T) {
	c := getConfig()
	d := NewDiscover(c)
	assert.NotNil(t, d)
}

func Test_Refresh(t *testing.T) {
	c := getConfig()
	d := NewDiscover(c)
	d.Refresh()
	assert.NotEmpty(t, d.Repositories)
	assert.Equal(t, 1, len(d.Repositories))
	assert.Equal(t, "test", d.Repositories[0].Path())
	assert.Equal(t, "test", d.Repositories[0].ImagePath())
}

func Test_Save(t *testing.T) {
	c := getConfig()
	d := NewDiscover(c)
	d.Refresh()
	d.AddExternalApplication("name", "command", "argument")
	d.AddRepository("path", "image-path")
	d.Save()
	d.Refresh()
	assert.Equal(t, 1, len(d.ExternalApplications))
	assert.Equal(t, 2, len(d.Repositories))
	d.RemoveExternalApplication("name")
	d.RemoveRepository("path")
	d.Save()
}

func getConfig() *config.Config {
	c := config.NewConfig()
	_ = c.Load(testConfigPath)
	return c
}
