package db

import (
	"fmt"

	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/users"

	"github.com/jmoiron/sqlx"
)

type GitHubUserRepo interface {
	Create(user github.GitHubUser) error
	GetByUsername(username string) (*github.GitHubUser, error)
	GetByUserID(users.ID) (*github.GitHubUser, error)
	Update(ouser *github.GitHubUser) error
}

type gitHubUserRepo struct {
	db *sqlx.DB
}

func NewGitHubUserRepo(db *sqlx.DB) GitHubUserRepo {
	return &gitHubUserRepo{db: db}
}

func (r *gitHubUserRepo) Create(ouser github.GitHubUser) error {
	_, err := r.db.NamedExec(`INSERT INTO github_users (id, user_id, username, access_token, created_at, access_token_last_validated_at)
		VALUES (:id, :user_id, :username, :access_token, :created_at, :access_token_last_validated_at)`, &ouser)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

func (r *gitHubUserRepo) GetByUsername(username string) (*github.GitHubUser, error) {
	var user github.GitHubUser
	err := r.db.Get(&user, "SELECT * FROM github_users WHERE username=$1", username)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &user, nil
}

func (r *gitHubUserRepo) GetByUserID(userID users.ID) (*github.GitHubUser, error) {
	var user github.GitHubUser
	err := r.db.Get(&user, "SELECT * FROM github_users WHERE user_id=$1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &user, nil
}

func (r *gitHubUserRepo) Update(ouser *github.GitHubUser) error {
	_, err := r.db.NamedExec(`UPDATE github_users
		SET username = :username,
		    access_token = :access_token,
		    access_token_last_validated_at = :access_token_last_validated_at
		WHERE id=:id`, ouser)
	if err != nil {
		return fmt.Errorf("failed to update %w", err)
	}
	return nil
}
