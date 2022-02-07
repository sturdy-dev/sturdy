package emails

import (
	"context"
)

type disabledClient struct{}

func NewDisabled() Sender {
	return &disabledClient{}
}

func (*disabledClient) Send(context.Context, *Email) error {
	return nil
}
