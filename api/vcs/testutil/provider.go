package testutil

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"getsturdy.com/api/vcs/provider"
)

func TestingRepoProvider(t *testing.T) provider.RepoProvider {
	reposBasePath, err := ioutil.TempDir(os.TempDir(), "sturdy")
	assert.NoError(t, err)

	lsfHostname := "localhost:8888"
	if n := os.Getenv("E2E_LFS_HOSTNAME"); n != "" {
		lsfHostname = n
	}

	return provider.New(reposBasePath, lsfHostname)
}
