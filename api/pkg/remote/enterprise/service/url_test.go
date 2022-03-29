package service

import (
	"testing"
)

func TestRewriteSshUrl(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			"gitlab-no-proto",
			"git@gitlab.com:zegl/sturdy-push.git",
			"ssh://git@gitlab.com/zegl/sturdy-push.git",
		},
		{
			"gitlab-valid-no-change",
			"ssh://git@gitlab.com/zegl/sturdy-push.git",
			"ssh://git@gitlab.com/zegl/sturdy-push.git",
		},
		{
			"azure-no-proto",
			"git@ssh.dev.azure.com:v3/getsturdy/gustav-sturdy-haxx/sturdy-on-azure",
			"ssh://git@ssh.dev.azure.com/v3/getsturdy/gustav-sturdy-haxx/sturdy-on-azure",
		},
		{
			"azure-valid-no-change",
			"ssh://git@ssh.dev.azure.com/v3/getsturdy/gustav-sturdy-haxx/sturdy-on-azure",
			"ssh://git@ssh.dev.azure.com/v3/getsturdy/gustav-sturdy-haxx/sturdy-on-azure",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rewriteSshUrl(tt.arg); got != tt.want {
				t.Errorf("RewriteURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
