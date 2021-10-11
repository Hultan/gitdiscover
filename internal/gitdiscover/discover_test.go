package gitdiscover

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/hultan/gitdiscover/internal/config"
)

func TestDiscover_GetRepositories(t *testing.T) {
	config := config.NewConfig()
	config.Load()

	logger := logrus.New()
	discover := NewGit(config, logger)
	assert.NotNil(t, discover)

	repos, err := discover.GetRepositories()
	if err != nil {
		assert.Empty(t, err)
	}
	assert.NotEmpty(t, repos[4].Status)
	assert.NotNil(t, repos[4].ModifiedDate)
}

func TestDiscover_GetRepositoriesByNameAndPath(t *testing.T) {
	config := config.NewConfig()
	config.Load()

	logger := logrus.New()
	discover := NewGit(config, logger)
	assert.NotNil(t, discover)

	repos := discover.GetRepositoryByName("gitdiscover")
	assert.Equal(t, len(repos),
		1)
	assert.Equal(t, repos[0].Path,
		"/home/per/code/gitdiscover")
	assert.Equal(t, discover.GetRepositoryByPath("/home/per/code/gitdiscover").Path,
		"/home/per/code/gitdiscover")
}

func TestDiscover_GetRepositoryByPath(t *testing.T) {
	config := config.NewConfig()
	config.Load()

	logger := logrus.New()
	discover := NewGit(config, logger)
	assert.NotNil(t, discover)

	assert.Equal(t, "/home/per/code/gitdiscover",
		discover.GetRepositoryByPath("/home/per/code/gitdiscover").Path)
	assert.Equal(t, "gitdiscover",
		discover.GetRepositoryByPath("/home/per/code/gitdiscover").Name)
	assert.Equal(t, true,
		discover.GetRepositoryByPath("/home/per/code/gitdiscover").IsGit)
	assert.Equal(t, "/home/per/code/gitdiscover/assets/application.png",
		discover.GetRepositoryByPath("/home/per/code/gitdiscover").ImagePath)
}
