package disabled

import "mash/pkg/analytics"

type Client struct{}

func NewClient() analytics.Client {
	return analytics.New(&Client{})
}

func (*Client) Enqueue(analytics.Message) error {
	return nil
}

func (*Client) Close() error {
	return nil
}
