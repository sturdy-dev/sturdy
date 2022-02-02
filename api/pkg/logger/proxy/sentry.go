package proxy

import (
	"fmt"

	"getsturdy.com/api/pkg/installations"

	"github.com/getsentry/raven-go"
)

func NewClient(installation installations.GetInstallationFunc) (*raven.Client, error) {
	client, err := raven.New("https://na@api.getsturdy.com/will/be/rewritten/123")
	if err != nil {
		return nil, err
	}

	ins, err := installation()
	if err != nil {
		return nil, fmt.Errorf("could not get current installation: %w", err)
	}

	// rewrite url
	client.Transport = &proxyTransport{
		Transport:    client.Transport,
		installation: ins,
	}
	return client, nil
}

type proxyTransport struct {
	raven.Transport
	installation *installations.Installation
}

func (t *proxyTransport) Send(url, authHeader string, packet *raven.Packet) error {
	packet.ServerName = fmt.Sprintf("installation-%s", t.installation.ID)
	return t.Transport.Send("https://api.getsturdy.com/v3/sentry/store/", authHeader, packet)
}
