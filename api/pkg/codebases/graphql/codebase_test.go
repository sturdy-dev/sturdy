package graphql

import (
	"context"
	"testing"

	"getsturdy.com/api/pkg/auth"
	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/codebases"
	service_codebase "getsturdy.com/api/pkg/codebases/service"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/internal/inmemory"
	"getsturdy.com/api/pkg/users"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCodebaseAccess(t *testing.T) {
	codebaseRepo := inmemory.NewInMemoryCodebaseRepo()
	codebaseUserRepo := inmemory.NewInMemoryCodebaseUserRepo()
	codebaseService := service_codebase.New(codebaseRepo, codebaseUserRepo, nil, nil, nil, nil, nil, nil, nil)
	authService := service_auth.New(codebaseService, nil, nil, nil, nil, nil)
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
		nil,
		nil,
		authService,
		codebaseService,
		nil,
		nil,
	)

	privateCodebase := codebases.Codebase{ID: codebases.ID(uuid.NewString()), ShortCodebaseID: "short-private"}
	assert.NoError(t, codebaseRepo.Create(privateCodebase))

	publicCodebase := codebases.Codebase{ID: codebases.ID(uuid.NewString()), ShortCodebaseID: "short-public", IsPublic: true}
	assert.NoError(t, codebaseRepo.Create(publicCodebase))

	userID := users.ID(uuid.NewString())

	// Add member to both codebases
	assert.NoError(t, codebaseUserRepo.Create(codebases.CodebaseUser{ID: uuid.NewString(), CodebaseID: privateCodebase.ID, UserID: userID}))
	assert.NoError(t, codebaseUserRepo.Create(codebases.CodebaseUser{ID: uuid.NewString(), CodebaseID: publicCodebase.ID, UserID: userID}))

	cases := []struct {
		name         string
		ctx          context.Context
		codebaseID   codebases.ID
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
			ctx:          auth.NewContext(context.Background(), &auth.Subject{Type: auth.SubjectUser, ID: userID.String()}),
			codebaseID:   privateCodebase.ID,
			expectAccess: true,
		},
		{
			name:         "member-public-codebase",
			ctx:          auth.NewContext(context.Background(), &auth.Subject{Type: auth.SubjectUser, ID: userID.String()}),
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
