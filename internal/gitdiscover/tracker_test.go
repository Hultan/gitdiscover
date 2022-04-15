package gitdiscover

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiscover_GetRepositories(t *testing.T) {
	config := NewConfig()
	config.Load()

	discover := NewTracker(config)
	assert.NotNil(t, discover)

	assert.NotEmpty(t, discover.Folders[4].GitStatus)
	assert.NotNil(t, discover.Folders[4].ModifiedDate)
}
