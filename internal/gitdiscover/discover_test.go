package gitdiscover

import (
	"github.com/hultan/gitdiscover/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDiscover_GetRepositories(t *testing.T) {
	config := config.NewConfig()
	config.Load()
	discover := GitNew(config)

	repos, err := discover.GetRepositories()
	if err != nil {
		assert.Empty(t, err)
	}
	assert.Equal(t, repos[4].Path, "/home/per/code/gitdiscover" )
	assert.NotEmpty(t, repos[4].Status)
	assert.NotNil(t, repos[4].Date)
}
