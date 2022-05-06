package graphql

import (
	"testing"

	"getsturdy.com/api/pkg/views"

	"github.com/stretchr/testify/assert"
)

func TestShortMountName(t *testing.T) {
	var absolutePath string
	resolver := Resolver{
		v: &views.View{
			MountPath: &absolutePath,
		},
	}

	absolutePath = "/Users/emilbroman/code/my-project"
	assert.Equal(t, resolver.ShortMountPath(), "~/code/my-project")

	absolutePath = "/home/emilbroman/my-project"
	assert.Equal(t, resolver.ShortMountPath(), "~/my-project")

	absolutePath = "/code/my-project"
	assert.Equal(t, resolver.ShortMountPath(), "/code/my-project")
}
