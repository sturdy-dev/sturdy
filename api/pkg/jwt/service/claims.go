package service

import (
	"mash/pkg/jwt"
)

type jwtClaims struct {
	Type jwt.TokenType `json:"type,omitempty"`
}

type deprecatedClaims struct {
	UserID string `json:"user_id,omitempty"`
}
