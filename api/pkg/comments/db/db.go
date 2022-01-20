package db

import (
	"fmt"
	"mash/pkg/change"
	"mash/pkg/comments"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(comments.Comment) error
	Get(id comments.ID) (comments.Comment, error)
	Update(comment comments.Comment) error
	GetByCodebaseAndChange(codebaseID string, changeID change.ID) ([]comments.Comment, error)
	GetByWorkspace(workspaceID string) ([]comments.Comment, error)
	GetByParent(id comments.ID) ([]comments.Comment, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repository {
	return &repo{db: db}
}

func (r *repo) Get(id comments.ID) (comments.Comment, error) {
	var res comments.Comment
	err := r.db.Get(&res, "SELECT * FROM comments WHERE id=$1", id)
	if err != nil {
		return comments.Comment{}, fmt.Errorf("failed to query table: %w", err)
	}
	return res, nil
}

func (r *repo) Create(comment comments.Comment) error {
	_, err := r.db.NamedExec(`INSERT INTO comments (id, codebase_id, change_id, user_id, created_at, message, path, line_start, line_end, line_is_new, workspace_id, context, context_starts_at_line, parent_comment_id)
		VALUES (:id, :codebase_id, :change_id, :user_id, :created_at, :message, :path, :line_start, :line_end, :line_is_new, :workspace_id, :context, :context_starts_at_line, :parent_comment_id)`, &comment)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

func (r *repo) Update(comment comments.Comment) error {
	_, err := r.db.NamedExec(`UPDATE comments
    	SET deleted_at = :deleted_at,
			message = :message,
    	    workspace_id = :workspace_id,
    	    change_id = :change_id
    	WHERE id = :id`, &comment)
	if err != nil {
		return fmt.Errorf("failed to update change: %w", err)
	}
	return nil
}

func (r *repo) GetByCodebaseAndChange(codebaseID string, changeID change.ID) ([]comments.Comment, error) {
	var res []comments.Comment
	err := r.db.Select(&res, `SELECT id, codebase_id, change_id, user_id, created_at, message, path, line_start, line_end, line_is_new, workspace_id, context, context_starts_at_line, parent_comment_id
		FROM comments
		WHERE codebase_id = $1
		  AND change_id = $2
	  	  AND deleted_at IS NULL
	  	  AND parent_comment_id IS NULL
	  	ORDER BY created_at DESC`, codebaseID, changeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return res, nil
}

func (r *repo) GetByWorkspace(workspaceID string) ([]comments.Comment, error) {
	var res []comments.Comment
	err := r.db.Select(&res, `SELECT id, codebase_id, change_id, user_id, created_at, message, path, line_start, line_end, line_is_new, workspace_id, context, context_starts_at_line, parent_comment_id
		FROM comments
		WHERE workspace_id = $1
		  AND deleted_at IS NULL
		  AND parent_comment_id IS NULL
	  ORDER BY created_at DESC`, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return res, nil
}

func (r *repo) GetByParent(id comments.ID) ([]comments.Comment, error) {
	var res []comments.Comment
	err := r.db.Select(&res, `SELECT id, codebase_id, change_id, user_id, created_at, message, path, line_start, line_end, line_is_new, workspace_id, context, context_starts_at_line, parent_comment_id
		FROM comments
		WHERE parent_comment_id = $1
		  AND deleted_at IS NULL
		ORDER BY created_at ASC`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return res, nil
}
