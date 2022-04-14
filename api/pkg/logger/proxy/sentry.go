package proxy

import (
	"fmt"
	"net/http"
	"net/url"

	"getsturdy.com/api/pkg/installations"
	"getsturdy.com/api/pkg/version"

	"github.com/getsentry/sentry-go"
)

var proxyURL, _ = url.Parse("https://api.getsturdy.com/v3/sentry/store/v2")

func NewClient(installation *installations.Installation) (*sentry.Client, error) {
	return sentry.NewClient(sentry.ClientOptions{
		Dsn:         "https://doesntmatter@anything.ingest.sentry.io/whatever",
		ServerName:  fmt.Sprintf("installation-%s", installation.ID),
		Release:     version.Version,
		Environment: version.Type.String(),
		HTTPTransport: &proxyRoundTripper{
			ProxyURL: proxyURL,
		},
	})
}

// proxyRoundTripper is a RoundTripper that proxies requests through the Sturdy Cloud.
type proxyRoundTripper struct {
	ProxyURL   *url.URL
	HTTPClient http.Client
}

// RoundTrip routes the request to the proxy server.
func (t *proxyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL = t.ProxyURL
	return t.HTTPClient.Do(req)
}
