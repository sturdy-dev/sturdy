package service_test

import (
	"context"
	"testing"

	"mash/pkg/auth"
	service_auth "mash/pkg/auth/service"
	"mash/pkg/codebase"
	service_codebase "mash/pkg/codebase/service"
	"mash/pkg/internal/inmemory"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCanWrite_codebase(t *testing.T) {
	cases := []struct {
		name string

		isAuthenticated  bool
		isMember         bool
		codebaseIsPublic bool

		expected bool
	}{
		{
			name:            "anon-private-codebase-no-access",
			isAuthenticated: false, isMember: false, codebaseIsPublic: false,
			expected: false,
		},
		{
			name:            "anon-public-codebase-has-no-access",
			isAuthenticated: false, isMember: false, codebaseIsPublic: true,
			expected: false,
		},
		{
			name:            "user-private-codebase-no-member-no-access",
			isAuthenticated: true, isMember: false, codebaseIsPublic: false,
			expected: false,
		},
		{
			name:            "user-private-codebase-is-member-has-access",
			isAuthenticated: true, isMember: true, codebaseIsPublic: false,
			expected: true,
		},
		{
			name:            "user-public-codebase-no-member-has-no-access",
			isAuthenticated: true, isMember: false, codebaseIsPublic: true,
			expected: false,
		},
		{
			name:            "user-public-codebase-is-member-has-access",
			isAuthenticated: true, isMember: true, codebaseIsPublic: true,
			expected: true,
		},
	}

	codebaseRepo := inmemory.NewInMemoryCodebaseRepo()
	codebaseUserRepo := inmemory.NewInMemoryCodebaseUserRepo()
	codebaseService := service_codebase.New(codebaseRepo, codebaseUserRepo, nil, nil, nil, nil, nil)

	authService := service_auth.New(
		codebaseService,
		nil,
		nil,
		nil,
		nil,
	)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cb := codebase.Codebase{ID: uuid.NewString(), IsPublic: tc.codebaseIsPublic}
			assert.NoError(t, codebaseRepo.Create(cb))

			userID := uuid.NewString()

			if tc.isMember {
				cbu := codebase.CodebaseUser{ID: uuid.NewString(), CodebaseID: cb.ID, UserID: userID}
				assert.NoError(t, codebaseUserRepo.Create(cbu))
			}

			ctx := context.Background()

			if tc.isAuthenticated {
				ctx = auth.NewContext(ctx, &auth.Subject{ID: userID, Type: auth.SubjectUser})
			} else {
				ctx = auth.NewContext(ctx, &auth.Subject{Type: auth.SubjectAnonymous})
			}

			hasAccessErr := authService.CanWrite(ctx, cb)

			if tc.expected {
				assert.NoError(t, hasAccessErr)
			} else {
				assert.Error(t, hasAccessErr)
			}

		})
	}
}

func TestCanRead_codebase(t *testing.T) {
	cases := []struct {
		name string

		isAuthenticated  bool
		isMember         bool
		codebaseIsPublic bool

		expected bool
	}{
		{
			name:            "anon-private-codebase-no-access",
			isAuthenticated: false, isMember: false, codebaseIsPublic: false,
			expected: false,
		},
		{
			name:            "anon-public-codebase-has-access",
			isAuthenticated: false, isMember: false, codebaseIsPublic: true,
			expected: true,
		},
		{
			name:            "user-private-codebase-no-member-no-access",
			isAuthenticated: true, isMember: false, codebaseIsPublic: false,
			expected: false,
		},
		{
			name:            "user-private-codebase-is-member-has-access",
			isAuthenticated: true, isMember: true, codebaseIsPublic: false,
			expected: true,
		},
		{
			name:            "user-public-codebase-no-member-has-access",
			isAuthenticated: true, isMember: false, codebaseIsPublic: true,
			expected: true,
		},
		{
			name:            "user-public-codebase-is-member-has-access",
			isAuthenticated: true, isMember: true, codebaseIsPublic: true,
			expected: true,
		},
	}

	codebaseRepo := inmemory.NewInMemoryCodebaseRepo()
	codebaseUserRepo := inmemory.NewInMemoryCodebaseUserRepo()
	codebaseService := service_codebase.New(codebaseRepo, codebaseUserRepo, nil, nil, nil, nil, nil)

	authService := service_auth.New(
		codebaseService,
		nil,
		nil,
		nil,
		nil,
	)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cb := codebase.Codebase{ID: uuid.NewString(), IsPublic: tc.codebaseIsPublic}
			assert.NoError(t, codebaseRepo.Create(cb))

			userID := uuid.NewString()

			if tc.isMember {
				cbu := codebase.CodebaseUser{ID: uuid.NewString(), CodebaseID: cb.ID, UserID: userID}
				assert.NoError(t, codebaseUserRepo.Create(cbu))
			}

			ctx := context.Background()

			if tc.isAuthenticated {
				ctx = auth.NewContext(ctx, &auth.Subject{ID: userID, Type: auth.SubjectUser})
			} else {
				ctx = auth.NewContext(ctx, &auth.Subject{Type: auth.SubjectAnonymous})
			}

			hasAccessErr := authService.CanRead(ctx, cb)
			if tc.expected {
				assert.NoError(t, hasAccessErr)
			} else {
				assert.Error(t, hasAccessErr)
			}
		})
	}
}
