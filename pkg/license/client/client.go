package client

type client struct{}

func New() *client {
	return &client{}
}

func (c *client) Ping() error {
	// TODO
	return nil
}
