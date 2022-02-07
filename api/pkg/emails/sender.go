package emails

import (
	"context"
)

type Email struct {
	To      string
	Subject string
	Html    string
}

type Sender interface {
	Send(context.Context, *Email) error
}
