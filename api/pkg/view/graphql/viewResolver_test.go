package graphql

import (
	"github.com/stretchr/testify/assert"
	"mash/pkg/view"
	"testing"
)

func TestShortMountName(t *testing.T) {
	var absolutePath string
	resolver := Resolver{
		v: &view.View{
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
