package db

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/user"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository interface {
	Create(newUser *user.User) error
	Get(id string) (*user.User, error)
	GetByIDs(ctx context.Context, ids ...string) ([]*user.User, error)
	GetByEmail(email string) (*user.User, error)
	Update(*user.User) error
	UpdatePassword(u *user.User) error
	Count(context.Context) (int, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repository {
	return &repo{db: db}
}

// The ID value is set inside this method
func (r *repo) Create(newUser *user.User) error {
	_, err := r.db.NamedExec(`INSERT INTO users (id, name, email, email_verified, password, created_at)
		VALUES (:id, :name, :email, :email_verified, :password, :created_at)`, &newUser)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

func (r *repo) GetByIDs(ctx context.Context, ids ...string) ([]*user.User, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT 
		id,
		name,
		email,
		email_verified,
		password,
		created_at,
		avatar_url
	FROM 
		users where id = ANY($1)`, pq.Array(ids))
	if err != nil {
		return nil, err
	}

	var users []*user.User
	for rows.Next() {
		u := new(user.User)
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.EmailVerified, &u.PasswordHash, &u.CreatedAt, &u.AvatarURL); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *repo) Get(id string) (*user.User, error) {
	var res user.User
	err := r.db.Get(&res, `SELECT * FROM users WHERE id=$1`, id)
	if err != nil {
		return nil, fmt.Errorf("failed tow query table: %w", err)
	}
	return &res, nil
}

func (r *repo) GetByEmail(email string) (*user.User, error) {
	var res user.User
	err := r.db.Get(&res, `SELECT * FROM users WHERE email=$1`, email)
	if err != nil {
		return nil, fmt.Errorf("failed tow query table: %w", err)
	}
	return &res, nil
}

func (r *repo) Update(u *user.User) error {
	_, err := r.db.NamedExec(`UPDATE users
		SET name = :name,
		    email = :email,
		    email_verified = :email_verified,
		    avatar_url = :avatar_url
		WHERE id = :id`, u)
	if err != nil {
		return fmt.Errorf("failed to update %w", err)
	}
	return nil
}

func (r *repo) UpdatePassword(u *user.User) error {
	_, err := r.db.NamedExec(`UPDATE users
		SET password = :password
		WHERE id = :id`, u)
	if err != nil {
		return fmt.Errorf("failed to update %w", err)
	}
	return nil
}

func (r *repo) Count(ctx context.Context) (int, error) {
	var res struct {
		Count int
	}
	if err := r.db.GetContext(ctx, &res, "SELECT count(*) as Count FROM users"); err != nil {
		return 0, fmt.Errorf("failed to get user count: %w", err)
	}
	return res.Count, nil
}
