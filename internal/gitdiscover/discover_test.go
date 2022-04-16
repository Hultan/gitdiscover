package gitdiscover

import (
	"testing"

	"github.com/stretchr/testify/assert"

	config2 "github.com/hultan/gitdiscover/internal/config"
)

func TestDiscover_GetRepositories(t *testing.T) {
	config := config2.NewConfig()
	config.Load()

	discover := NewDiscover(config)
	assert.NotNil(t, discover)

	assert.NotEmpty(t, discover.Folders[4].GitStatus)
	assert.NotNil(t, discover.Folders[4].ModifiedDate)
}
