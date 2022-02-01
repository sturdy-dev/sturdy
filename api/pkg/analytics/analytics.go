package analytics

import (
	"github.com/posthog/posthog-go"
)

type Client interface {
	posthog.Client
	// this is to prevevnt the posthog.Client from being used as a Client
	_incompatible()
}

type client struct {
	posthog.Client
}

func New(pc posthog.Client) *client {
	return &client{pc}
}

func (c *client) _incompatible() {}

type Message = posthog.Message
type Capture = posthog.Capture
type Identify = posthog.Identify
type Properties = posthog.Properties

func NewProperties() Properties {
	return posthog.NewProperties()
}
