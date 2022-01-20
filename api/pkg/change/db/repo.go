package db

import (
	"fmt"
	"getsturdy.com/api/pkg/change"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Get(id change.ID) (change.Change, error)
	Insert(change.Change) error
	Update(change.Change) error
}

func NewRepo(db *sqlx.DB) Repository {
	return &repo{db: db}
}

type repo struct {
	db *sqlx.DB
}

func (r *repo) Get(id change.ID) (change.Change, error) {
	var res change.Change
	err := r.db.Get(&res, `SELECT id, codebase_id, title, updated_description, user_id, git_creator_name, git_creator_email, created_at, git_created_at FROM changes WHERE id = $1`, id)
	if err != nil {
		return change.Change{}, err
	}
	return res, nil
}

func (r *repo) Insert(ch change.Change) error {
	_, err := r.db.NamedExec(`INSERT INTO changes
		(id, codebase_id, title, updated_description, user_id, git_creator_name, git_creator_email, created_at, git_created_at)
		VALUES(:id, :codebase_id, :title, :updated_description, :user_id, :git_creator_name, :git_creator_email, :created_at, :git_created_at)
    	`, &ch)
	if err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (r *repo) Update(ch change.Change) error {
	_, err := r.db.NamedExec(`UPDATE changes
    	SET updated_description = :updated_description,
    	    title = :title,
    	    user_id = :user_id,
    	    git_creator_name = :git_creator_name,
    	    git_creator_email = :git_creator_email,
    	    created_at = :created_at,
    	    git_created_at = :git_created_at
    	WHERE id = :id`, &ch)
	if err != nil {
		return fmt.Errorf("failed to update change: %w", err)
	}
	return nil
}
