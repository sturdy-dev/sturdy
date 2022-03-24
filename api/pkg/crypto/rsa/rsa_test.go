package rsa

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRsaKeypair(t *testing.T) {
	public, private, err := GenerateRsaKeypair()
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(string(public), "ssh-rsa "))
	assert.True(t, strings.HasPrefix(string(private), "-----BEGIN RSA PRIVATE KEY-----\n"))
	assert.True(t, strings.HasSuffix(string(private), "\n-----END RSA PRIVATE KEY-----\n"))
}

func TestGenerateRsaSSHKeypair(t *testing.T) {
	public, private, err := GenerateRsaSSHKeypair()
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(string(public), "ssh-rsa "))
	assert.True(t, strings.HasPrefix(string(private), "-----BEGIN OPENSSH PRIVATE KEY-----\n"))
	assert.True(t, strings.HasSuffix(string(private), "\n-----END OPENSSH PRIVATE KEY-----\n"))
}
