package db

import (
	"context"
	"database/sql"

	"getsturdy.com/api/pkg/servicetokens"
)

var _ Repository = &memory{}

type memory struct {
	byID map[string]*servicetokens.Token
}

func NewMemory() *memory {
	return &memory{
		byID: map[string]*servicetokens.Token{},
	}
}

func (m *memory) Create(_ context.Context, token *servicetokens.Token) error {
	m.byID[token.ID] = token
	return nil
}

func (m *memory) Update(_ context.Context, token *servicetokens.Token) error {
	m.byID[token.ID] = token
	return nil
}

func (m *memory) GetByID(_ context.Context, id string) (*servicetokens.Token, error) {
	token, found := m.byID[id]
	if !found {
		return nil, sql.ErrNoRows
	}
	return token, nil
}

func (m *memory) DeleteByID(_ context.Context, id string) error {
	delete(m.byID, id)
	return nil
}
