package db

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/users"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repository {
	return &repo{db: db}
}

func (r *repo) List(ctx context.Context, limit uint64) ([]*users.User, error) {
	users := []*users.User{}
	if err := r.db.SelectContext(ctx, &users, `
		SELECT
			id,
			name,
			email,
			email_verified,
			password,
			created_at,
			avatar_url,
			status,
			referer,
		    "is"
		FROM users
		ORDER BY created_at DESC
		LIMIT $1
	`, limit); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return users, nil
}

// The ID value is set inside this method
func (r *repo) Create(newUser *users.User) error {
	_, err := r.db.NamedExec(`INSERT INTO users (id, name, email, email_verified, password, created_at, status, referer, "is")
		VALUES (:id, :name, :email, :email_verified, :password, :created_at, :status, :referer, :is)`, &newUser)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

func (r *repo) GetByIDs(ctx context.Context, ids ...users.ID) ([]*users.User, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT 
		id,
		name,
		email,
		email_verified,
		password,
		created_at,
		avatar_url,
		status,
		referer,
		"is"
	FROM 
		users where id = ANY($1)`, pq.Array(ids))
	if err != nil {
		return nil, err
	}

	var uu []*users.User
	for rows.Next() {
		u := new(users.User)
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.EmailVerified, &u.PasswordHash, &u.CreatedAt, &u.AvatarURL, &u.Status, &u.Referer, &u.Is); err != nil {
			return nil, err
		}
		uu = append(uu, u)
	}

	return uu, nil
}

func (r *repo) Get(id users.ID) (*users.User, error) {
	var res users.User
	err := r.db.Get(&res, `SELECT
			id,
			name,
			email,
			email_verified,
			password,
			created_at,
			avatar_url,
			status,
			referer,
			"is"
		FROM users WHERE id=$1`, id)
	if err != nil {
		return nil, fmt.Errorf("failed tow query table: %w", err)
	}
	return &res, nil
}

func (r *repo) GetByEmail(email string) (*users.User, error) {
	var res users.User
	err := r.db.Get(&res, `SELECT 
			id,
			name,
			email,
			email_verified,
			password,
			created_at,
			avatar_url,
			status,
			referer,
			"is"
		FROM users WHERE email=$1
	`, email)
	if err != nil {
		return nil, fmt.Errorf("failed tow query table: %w", err)
	}
	return &res, nil
}

func (r *repo) Update(u *users.User) error {
	_, err := r.db.NamedExec(`UPDATE users
		SET name = :name,
		    email = :email,
		    email_verified = :email_verified,
		    avatar_url = :avatar_url,
			status = :status,
			"is" = :is
		WHERE id = :id`, u)
	if err != nil {
		return fmt.Errorf("failed to update %w", err)
	}
	return nil
}

func (r *repo) UpdatePassword(u *users.User) error {
	_, err := r.db.NamedExec(`UPDATE users
		SET password = :password
		WHERE id = :id`, u)
	if err != nil {
		return fmt.Errorf("failed to update %w", err)
	}
	return nil
}

func (r *repo) Count(ctx context.Context) (uint64, error) {
	var res struct {
		Count uint64
	}
	if err := r.db.GetContext(ctx, &res, "SELECT count(1) as Count FROM users"); err != nil {
		return 0, fmt.Errorf("failed to select: %w", err)
	}
	return res.Count, nil
}
