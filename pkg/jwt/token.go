package jwt

import "time"

type TokenType string

const (
	TokenTypeUndefined TokenType = ""
	// TokenTypeVerifyEmail is the token type for verifying emails. It must have user_id as a subject.
	TokenTypeVerifyEmail TokenType = "verify_email"
	// TokenTypeAuth is the token type for authenticating users. it must have user_id as a subject.
	TokenTypeAuth TokenType = "auth"
	// TokenTypeCI is the token type for CI authentication. It must have change_id as a subject.
	TokenTypeCI TokenType = "ci"
)

type Token struct {
	Token     string
	Type      TokenType
	Subject   string
	ExpiresAt time.Time
}
