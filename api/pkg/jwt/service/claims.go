package service

import (
	"getsturdy.com/api/pkg/jwt"
)

type jwtClaims struct {
	Type jwt.TokenType `json:"type,omitempty"`
}

type deprecatedClaims struct {
	UserID string `json:"user_id,omitempty"`
}
