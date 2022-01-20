package service_test

import (
	"context"
	"testing"
	"time"

	"getsturdy.com/api/pkg/jwt"
	db_keys "getsturdy.com/api/pkg/jwt/keys/db"
	"getsturdy.com/api/pkg/jwt/service"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestVerify_shouldVerifyIssuedKey(t *testing.T) {
	svc := service.NewService(zap.NewNop(), db_keys.NewInMemory())

	token, err := svc.IssueToken(context.Background(), "user-id", time.Hour, jwt.TokenTypeAuth)
	assert.NoError(t, err)

	verifiedToken, err := svc.Verify(context.Background(), token.Token, jwt.TokenTypeAuth)
	if assert.NoError(t, err) {
		assert.Equal(t, token, verifiedToken)
	}
}
