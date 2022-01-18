package db

import (
	"context"

	"mash/pkg/installations"
)

type Memory struct {
	list []*installations.Installation
}

func NewMemory() Repository {
	return &Memory{}
}

func (m *Memory) Create(_ context.Context, installation *installations.Installation) error {
	m.list = append(m.list, installation)
	return nil
}

func (m *Memory) ListAll(_ context.Context) ([]*installations.Installation, error) {
	return m.list, nil
}
