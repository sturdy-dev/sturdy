package oss

import (
	"context"
	"testing"

	"mash/pkg/auth"
	service_auth "mash/pkg/auth/service"
	"mash/pkg/codebase"
	service_codebase "mash/pkg/codebase/service"
	"mash/pkg/graphql/resolvers"
	"mash/pkg/internal/inmemory"
	"mash/pkg/posthog"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCodebaseAccess(t *testing.T) {
	codebaseRepo := inmemory.NewInMemoryCodebaseRepo()
	codebaseUserRepo := inmemory.NewInMemoryCodebaseUserRepo()
	codebaseService := service_codebase.New(codebaseRepo, codebaseUserRepo)
	authService := service_auth.New(codebaseService, nil, nil, nil, nil)
	resolver := NewCodebaseRootResolver(
		codebaseRepo,
		codebaseUserRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		zap.NewNop(),
		nil,
		nil,
		posthog.NewFakeClient(),
		nil,
		authService,
	)

	privateCodebase := codebase.Codebase{ID: uuid.NewString(), ShortCodebaseID: "short-private"}
	assert.NoError(t, codebaseRepo.Create(privateCodebase))

	publicCodebase := codebase.Codebase{ID: uuid.NewString(), ShortCodebaseID: "short-public", IsPublic: true}
	assert.NoError(t, codebaseRepo.Create(publicCodebase))

	userID := uuid.NewString()

	// Add member to both codebases
	assert.NoError(t, codebaseUserRepo.Create(codebase.CodebaseUser{ID: uuid.NewString(), CodebaseID: privateCodebase.ID, UserID: userID}))
	assert.NoError(t, codebaseUserRepo.Create(codebase.CodebaseUser{ID: uuid.NewString(), CodebaseID: publicCodebase.ID, UserID: userID}))

	cases := []struct {
		name         string
		ctx          context.Context
		codebaseID   string
		expectAccess bool
	}{
		{
			name:         "anon-private-codebase",
			ctx:          auth.NewContext(context.Background(), &auth.Subject{Type: auth.SubjectAnonymous}),
			codebaseID:   privateCodebase.ID,
			expectAccess: false,
		},
		{
			name:         "anon-public-codebase",
			ctx:          auth.NewContext(context.Background(), &auth.Subject{Type: auth.SubjectAnonymous}),
			codebaseID:   publicCodebase.ID,
			expectAccess: true,
		},
		{
			name:         "member-private-codebase",
			ctx:          auth.NewContext(context.Background(), &auth.Subject{Type: auth.SubjectUser, ID: userID}),
			codebaseID:   privateCodebase.ID,
			expectAccess: true,
		},
		{
			name:         "member-public-codebase",
			ctx:          auth.NewContext(context.Background(), &auth.Subject{Type: auth.SubjectUser, ID: userID}),
			codebaseID:   publicCodebase.ID,
			expectAccess: true,
		},
		{
			name:         "non-member-private-codebase",
			ctx:          auth.NewContext(context.Background(), &auth.Subject{Type: auth.SubjectUser, ID: uuid.NewString()}),
			codebaseID:   privateCodebase.ID,
			expectAccess: false,
		},
		{
			name:         "non-member-public-codebase",
			ctx:          auth.NewContext(context.Background(), &auth.Subject{Type: auth.SubjectUser, ID: uuid.NewString()}),
			codebaseID:   publicCodebase.ID,
			expectAccess: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gid := graphql.ID(tc.codebaseID)
			codebaseResolver, err := resolver.Codebase(tc.ctx, resolvers.CodebaseArgs{ID: &gid})

			if tc.expectAccess {
				assert.NoError(t, err)
				assert.NotNil(t, codebaseResolver)
			} else {
				assert.Error(t, err)
				assert.Nil(t, codebaseResolver)
			}
		})

	}
}
