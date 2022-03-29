package service

import (
	"net/url"
	"strings"

	"getsturdy.com/api/pkg/remote"
)

func rewriteSshUrl(in string) string {
	if !strings.HasPrefix(in, "ssh://") {
		rewritten := "ssh://" + strings.ReplaceAll(in, ":", "/")

		// validate the new url
		if _, err := url.Parse(rewritten); err == nil {
			return rewritten
		}
	}

	// could not rewrite, return original
	return in
}

func RewriteRemoteURL(rem remote.Remote) remote.Remote {
	if rem.KeyPairID != nil {
		rem.URL = rewriteSshUrl(rem.URL)
	}
	return rem
}
