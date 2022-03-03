package service_test

import (
	"context"
	"testing"

	"go.uber.org/zap"

	"getsturdy.com/api/pkg/analytics/disabled"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/auth"
	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/codebase"
	service_codebase "getsturdy.com/api/pkg/codebase/service"
	"getsturdy.com/api/pkg/internal/inmemory"
	"getsturdy.com/api/pkg/organization"
	service_organization "getsturdy.com/api/pkg/organization/service"
	"getsturdy.com/api/pkg/users"

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
	codebaseService := service_codebase.New(codebaseRepo, codebaseUserRepo, nil, nil, nil, nil, nil, nil, nil)

	authService := service_auth.New(
		codebaseService,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cb := codebase.Codebase{ID: uuid.NewString(), IsPublic: tc.codebaseIsPublic}
			assert.NoError(t, codebaseRepo.Create(cb))

			userID := users.ID(uuid.NewString())

			if tc.isMember {
				cbu := codebase.CodebaseUser{ID: uuid.NewString(), CodebaseID: cb.ID, UserID: userID}
				assert.NoError(t, codebaseUserRepo.Create(cbu))
			}

			ctx := context.Background()

			if tc.isAuthenticated {
				ctx = auth.NewContext(ctx, &auth.Subject{ID: userID.String(), Type: auth.SubjectUser})
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

		isMemberOfOrganization bool

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

		{
			name:            "user-can-read-private-codebase-member-of-organization",
			isAuthenticated: true, isMember: false, isMemberOfOrganization: true, codebaseIsPublic: false,
			expected: true,
		},
	}

	codebaseRepo := inmemory.NewInMemoryCodebaseRepo()
	codebaseUserRepo := inmemory.NewInMemoryCodebaseUserRepo()
	analyticsService := service_analytics.New(zap.NewNop(), disabled.NewClient(zap.NewNop()))
	codebaseService := service_codebase.New(codebaseRepo, codebaseUserRepo, nil, nil, nil, nil, nil, analyticsService, nil)

	organizationRepo := inmemory.NewInMemoryOrganizationRepo()
	organizationMemberRepo := inmemory.NewInMemoryOrganizationMemberRepository()
	organizationService := service_organization.New(nil, organizationRepo, organizationMemberRepo, analyticsService)

	authService := service_auth.New(
		codebaseService,
		nil,
		nil,
		nil,
		nil,
		organizationService,
	)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			orgID := uuid.NewString()

			cb := codebase.Codebase{ID: uuid.NewString(), IsPublic: tc.codebaseIsPublic}

			if tc.isMemberOfOrganization {
				cb.OrganizationID = &orgID
			}

			assert.NoError(t, codebaseRepo.Create(cb))

			userID := users.ID(uuid.NewString())

			if tc.isMember {
				cbu := codebase.CodebaseUser{ID: uuid.NewString(), CodebaseID: cb.ID, UserID: userID}
				assert.NoError(t, codebaseUserRepo.Create(cbu))
			}

			if tc.isMemberOfOrganization {
				org := organization.Organization{ID: orgID}
				assert.NoError(t, organizationRepo.Create(context.Background(), org))
				orgmember := organization.Member{ID: uuid.NewString(), OrganizationID: org.ID, UserID: userID}
				assert.NoError(t, organizationMemberRepo.Create(context.Background(), orgmember))
			}

			ctx := context.Background()

			if tc.isAuthenticated {
				ctx = auth.NewContext(ctx, &auth.Subject{ID: userID.String(), Type: auth.SubjectUser})
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

func TestCanReadWrite_organization(t *testing.T) {
	cases := []struct {
		name string

		isAuthenticated    bool
		isMember           bool
		isMemberOfCodebase bool

		expectedCanRead  bool
		expectedCanWrite bool
	}{
		{
			name:            "member-can-access",
			isAuthenticated: true, isMember: true, isMemberOfCodebase: false,

			expectedCanRead:  true,
			expectedCanWrite: true,
		},
		{
			name:            "non-member-can-not-access",
			isAuthenticated: true, isMember: false, isMemberOfCodebase: false,

			expectedCanRead:  false,
			expectedCanWrite: false,
		},
		{
			name:            "member-can-access-if-member-of-codebase",
			isAuthenticated: true, isMember: true, isMemberOfCodebase: true,

			expectedCanRead:  true,
			expectedCanWrite: true,
		},
		{
			name:            "non-member-can-read-only-access-if-member-of-codebase",
			isAuthenticated: true, isMember: false, isMemberOfCodebase: true,

			expectedCanRead:  true,
			expectedCanWrite: false,
		},
	}

	codebaseRepo := inmemory.NewInMemoryCodebaseRepo()
	codebaseUserRepo := inmemory.NewInMemoryCodebaseUserRepo()
	analyticsService := service_analytics.New(zap.NewNop(), disabled.NewClient(zap.NewNop()))
	codebaseService := service_codebase.New(codebaseRepo, codebaseUserRepo, nil, nil, nil, nil, nil, analyticsService, nil)

	organizationRepo := inmemory.NewInMemoryOrganizationRepo()
	organizationMemberRepo := inmemory.NewInMemoryOrganizationMemberRepository()

	organizationService := service_organization.New(nil, organizationRepo, organizationMemberRepo, analyticsService)

	authService := service_auth.New(
		codebaseService,
		nil,
		nil,
		nil,
		nil,
		organizationService,
	)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bgCtx := context.Background()

			org := organization.Organization{ID: uuid.NewString()}
			assert.NoError(t, organizationRepo.Create(bgCtx, org))

			userID := users.ID(uuid.NewString())

			if tc.isMember {
				orgmember := organization.Member{ID: uuid.NewString(), OrganizationID: org.ID, UserID: userID}
				assert.NoError(t, organizationMemberRepo.Create(bgCtx, orgmember))
			}

			if tc.isMemberOfCodebase {
				cb := codebase.Codebase{ID: uuid.NewString(), OrganizationID: &org.ID}
				assert.NoError(t, codebaseRepo.Create(cb))

				cbu := codebase.CodebaseUser{ID: uuid.NewString(), CodebaseID: cb.ID, UserID: userID}
				assert.NoError(t, codebaseUserRepo.Create(cbu))
			}

			ctx := context.Background()

			if tc.isAuthenticated {
				ctx = auth.NewContext(ctx, &auth.Subject{ID: userID.String(), Type: auth.SubjectUser})
			} else {
				ctx = auth.NewContext(ctx, &auth.Subject{Type: auth.SubjectAnonymous})
			}

			canReadAccessErr := authService.CanRead(ctx, org)
			if tc.expectedCanRead {
				assert.NoError(t, canReadAccessErr)
			} else {
				assert.Error(t, canReadAccessErr)
			}

			canWriteAccessErr := authService.CanWrite(ctx, org)
			if tc.expectedCanWrite {
				assert.NoError(t, canWriteAccessErr)
			} else {
				assert.Error(t, canWriteAccessErr)
			}
		})
	}
}
