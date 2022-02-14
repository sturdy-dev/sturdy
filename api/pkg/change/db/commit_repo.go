package db

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/change"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type CommitRepository interface {
	ListByChangeIDs(context.Context, ...change.ID) ([]*change.ChangeCommit, error)
	GetByCommitID(commitID, codebaseID string) (change.ChangeCommit, error)
	GetByChangeIDOnTrunk(id change.ID) (change.ChangeCommit, error)
	Insert(ch change.ChangeCommit) error
	Update(ch change.ChangeCommit) error
}

func NewCommitRepository(db *sqlx.DB) CommitRepository {
	return &commitRepo{db: db}
}

type commitRepo struct {
	db *sqlx.DB
}

func (r *commitRepo) GetByCommitID(commitID, codebaseID string) (change.ChangeCommit, error) {
	var res change.ChangeCommit
	err := r.db.Get(&res, `SELECT change_id, commit_id, codebase_id, trunk FROM change_commits WHERE commit_id = $1 AND codebase_id = $2`, commitID, codebaseID)
	if err != nil {
		return change.ChangeCommit{}, err
	}
	return res, nil
}

func (r *commitRepo) ListByChangeIDs(ctx context.Context, ids ...change.ID) ([]*change.ChangeCommit, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			change_id,
			commit_id,
			codebase_id,
			trunk
		FROM
			change_commits
		WHERE
			change_id = ANY($1)
	`, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	res := []*change.ChangeCommit{}
	for rows.Next() {
		c := new(change.ChangeCommit)
		if err := rows.Scan(&c.ChangeID, &c.CommitID, &c.CodebaseID, &c.Trunk); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		res = append(res, c)
	}

	return res, nil
}

func (r *commitRepo) GetByChangeIDOnTrunk(id change.ID) (change.ChangeCommit, error) {
	var res change.ChangeCommit
	err := r.db.Get(&res, `SELECT change_id, commit_id, codebase_id, trunk
		FROM change_commits
		WHERE change_id = $1 AND trunk = true LIMIT 1`, id)
	if err != nil {
		return change.ChangeCommit{}, err
	}
	return res, nil
}

func (r *commitRepo) Insert(ch change.ChangeCommit) error {
	_, err := r.db.NamedExec(`INSERT INTO change_commits
		(change_id, commit_id, codebase_id, trunk)
		VALUES(:change_id, :commit_id, :codebase_id, :trunk)
    	`, &ch)
	if err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (r *commitRepo) Update(ch change.ChangeCommit) error {
	_, err := r.db.NamedExec(`UPDATE change_commits
		SET trunk = :trunk
		WHERE change_id = :change_id AND commit_id = :commit_id AND codebase_id = :codebase_id
    	`, &ch)
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}
	return nil
}
