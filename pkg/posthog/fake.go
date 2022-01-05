package posthog

import "github.com/posthog/posthog-go"

func NewFakeClient() posthog.Client {
	return &fake{}
}

type fake struct{}

func (*fake) Close() error {
	return nil
}

func (*fake) Enqueue(posthog.Message) error {
	return nil
}
