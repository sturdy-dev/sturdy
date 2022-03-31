package users

import (
	"regexp"
	"strings"
	"time"
)

type ID string

func (id ID) String() string {
	return string(id)
}

type Status string

const (
	StatusUndefined Status = ""
	// StatusActive means the user has logged in and least once.
	StatusActive Status = "active"
	// StatusShadow means the user never logged in, see Referer field to know
	// why it's created.
	StatusShadow Status = "shadow"
)

type User struct {
	ID            ID         `db:"id" json:"id"`
	Name          string     `db:"name" json:"name"`
	Email         string     `db:"email" json:"email"`
	EmailVerified bool       `db:"email_verified" json:"email_verified"`
	PasswordHash  string     `db:"password" json:"-"`
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	AvatarURL     *string    `db:"avatar_url" json:"avatar_url"`
	Status        Status     `db:"status" json:"-"`
	Referer       string     `db:"referer" json:"-"`
}

var nonAlpha = regexp.MustCompile(`[^a-zA-Z]+`)

// EmailToName returns the name of the user based on the email address.
func EmailToName(email string) string {
	beforeAt, afterAt, _ := strings.Cut(email, "@")
	beforePlus, _, _ := strings.Cut(beforeAt, "+")
	spaceDivided := nonAlpha.ReplaceAllString(beforePlus, " ")
	lowerCased := strings.ToLower(spaceDivided)
	if lowerCased == "sturdy" {
		beforeDot, _, _ := strings.Cut(afterAt, ".")
		return EmailToName(beforeDot)
	}
	capitilized := strings.Title(lowerCased)
	return strings.TrimSpace(capitilized)
}
