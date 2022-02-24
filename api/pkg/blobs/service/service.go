package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"

	"getsturdy.com/api/pkg/blobs"

	"github.com/jmoiron/sqlx"
)

type Service struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Service {
	return &Service{db: db}
}

var ErrNotFound = fmt.Errorf("not found: %w", sql.ErrNoRows)

func (s *Service) Fetch(ctx context.Context, id blobs.ID) (*blobs.Blob, error) {
	var blob blobs.Blob
	if err := s.db.GetContext(ctx, &blob, "SELECT * FROM blobs WHERE id = $1", id); errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to fetch blob: %w", err)
	} else {
		return &blob, nil
	}
}

func (s *Service) Store(ctx context.Context, id blobs.ID, reader io.Reader) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	blob := &blobs.Blob{
		ID:   id,
		Data: data,
	}
	if _, err := s.db.NamedExecContext(ctx, `
		INSERT INTO blobs (id, data)
		VALUES (:id, :data)
	`, blob); err != nil {
		return fmt.Errorf("failed to store blob: %w", err)
	}
	return nil
}
