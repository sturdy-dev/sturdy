package emails

import (
	"context"
	"fmt"
)

type disabledClient struct{}

func NewDisabled() Sender {
	return &disabledClient{}
}

func (*disabledClient) Send(context.Context, *Email) error {
	return fmt.Errorf("emails are disabled")
}
