package auth

import (
	"context"

	"getsturdy.com/api/pkg/jwt"
	"getsturdy.com/api/pkg/users"
)

type SubjectType string

const (
	SubjectUndefined SubjectType = ""
	SubjectUser      SubjectType = "user"
	SubjectCI        SubjectType = "ci"
	SubjectMutagen   SubjectType = "mutagen"
	SubjectAnonymous SubjectType = "anonymous"
)

func (st SubjectType) String() string {
	return string(st)
}

type Subject struct {
	ID   string
	Type SubjectType
}

var (
	convertType = map[jwt.TokenType]SubjectType{
		jwt.TokenTypeAuth: SubjectUser,
		jwt.TokenTypeCI:   SubjectCI,
	}
)

func subjectFromToken(token *jwt.Token) *Subject {
	if token == nil {
		return &Subject{Type: SubjectAnonymous}
	}

	return &Subject{
		ID:   token.Subject,
		Type: convertType[token.Type],
	}
}

type subjectKeyType struct{}

var subjectKey = subjectKeyType{}

func NewContext(ctx context.Context, s *Subject) context.Context {
	return context.WithValue(ctx, subjectKey, s)
}

func FromContext(ctx context.Context) (*Subject, bool) {
	s, ok := ctx.Value(subjectKey).(*Subject)
	return s, ok
}

// UserID returns authenticated user's id from the context.
//
// If context if unauthenticated or not user, returns an ErrUnauthenticated.
func UserID(ctx context.Context) (users.ID, error) {
	s, ok := FromContext(ctx)
	if !ok {
		return "", ErrUnauthenticated
	}
	if s.Type != SubjectUser {
		return "", ErrUnauthenticated
	}
	return users.ID(s.ID), nil
}
