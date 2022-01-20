package db

import (
	"context"
	"database/sql"

	"getsturdy.com/api/pkg/onetime"
)

var _ Repository = &Memory{}

type Memory struct {
	byKeyUserID map[string]*onetime.Token
}

func NewMemory() *Memory {
	return &Memory{
		byKeyUserID: make(map[string]*onetime.Token),
	}
}

func (m *Memory) Get(_ context.Context, userID, key string) (*onetime.Token, error) {
	token, ok := m.byKeyUserID[key+userID]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return token, nil
}

func (m *Memory) Create(_ context.Context, token *onetime.Token) error {
	m.byKeyUserID[token.Key+token.UserID] = token
	return nil
}

func (m *Memory) Update(_ context.Context, token *onetime.Token) error {
	m.byKeyUserID[token.Key+token.UserID] = token
	return nil
}
