package graphql

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"getsturdy.com/api/pkg/analytics"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/auth"
	service_auth "getsturdy.com/api/pkg/auth/service"
	service_change "getsturdy.com/api/pkg/changes/service"
	"getsturdy.com/api/pkg/codebases"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	service_codebase "getsturdy.com/api/pkg/codebases/service"
	"getsturdy.com/api/pkg/events"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	service_organization "getsturdy.com/api/pkg/organization/service"
	service_remote "getsturdy.com/api/pkg/remote/service"
	"getsturdy.com/api/pkg/users"
	db_user "getsturdy.com/api/pkg/users/db"
	"getsturdy.com/api/pkg/view"
	db_view "getsturdy.com/api/pkg/view/db"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	"getsturdy.com/api/vcs/executor"

	"github.com/graph-gophers/graphql-go"
	"github.com/jxskiss/base62"
	"go.uber.org/zap"
)

type CodebaseRootResolver struct {
	codebaseRepo     db_codebases.CodebaseRepository
	codebaseUserRepo db_codebases.CodebaseUserRepository
	viewRepo         db_view.Repository
	workspaceReader  db_workspaces.WorkspaceReader
	userRepo         db_user.Repository

	workspaceResolver                 *resolvers.WorkspaceRootResolver
	authorResolver                    resolvers.AuthorRootResolver
	viewResolver                      *resolvers.ViewRootResolver
	aclResolver                       resolvers.ACLRootResolver
	changeRootResolver                resolvers.ChangeRootResolver
	fileRootResolver                  resolvers.FileRootResolver
	instantIntegrationRootResolver    resolvers.IntegrationRootResolver
	codebaseGitHubIntegrationResolver resolvers.CodebaseGitHubIntegrationRootResolver
	organizationRootResolver          *resolvers.OrganizationRootResolver
	remoteRootResolver                resolvers.RemoteRootResolver

	logger           *zap.Logger
	viewEvents       events.EventReader
	eventsSender     events.EventSender
	executorProvider executor.Provider
	analyticsService *service_analytics.Service

	authService         *service_auth.Service
	codebaseService     *service_codebase.Service
	organizationService *service_organization.Service
	changeService       *service_change.Service
	remoteService       service_remote.Service
}

func NewCodebaseRootResolver(
	codebaseRepo db_codebases.CodebaseRepository,
	codebaseUserRepo db_codebases.CodebaseUserRepository,
	viewRepo db_view.Repository,
	workspaceReader db_workspaces.WorkspaceReader,
	userRepo db_user.Repository,

	workspaceResolver *resolvers.WorkspaceRootResolver,
	authorResolver resolvers.AuthorRootResolver,
	viewResolver *resolvers.ViewRootResolver,
	aclResolver resolvers.ACLRootResolver,
	changeRootResolver resolvers.ChangeRootResolver,
	fileRootResolver resolvers.FileRootResolver,
	instantIntegrationRootResolver resolvers.IntegrationRootResolver,
	codebaseGitHubIntegrationResolver resolvers.CodebaseGitHubIntegrationRootResolver,
	organizationRootResolver *resolvers.OrganizationRootResolver,
	remoteRootResolver resolvers.RemoteRootResolver,

	logger *zap.Logger,
	viewEvents events.EventReader,
	eventsSender events.EventSender,
	analyticsService *service_analytics.Service,
	executorProvider executor.Provider,

	authService *service_auth.Service,
	codebaseService *service_codebase.Service,
	organizationService *service_organization.Service,
	changeService *service_change.Service,
	remoteService service_remote.Service,
) resolvers.CodebaseRootResolver {
	return &CodebaseRootResolver{
		codebaseRepo:     codebaseRepo,
		codebaseUserRepo: codebaseUserRepo,
		viewRepo:         viewRepo,
		workspaceReader:  workspaceReader,
		userRepo:         userRepo,

		workspaceResolver:                 workspaceResolver,
		authorResolver:                    authorResolver,
		viewResolver:                      viewResolver,
		aclResolver:                       aclResolver,
		changeRootResolver:                changeRootResolver,
		fileRootResolver:                  fileRootResolver,
		instantIntegrationRootResolver:    instantIntegrationRootResolver,
		codebaseGitHubIntegrationResolver: codebaseGitHubIntegrationResolver,
		organizationRootResolver:          organizationRootResolver,
		remoteRootResolver:                remoteRootResolver,

		logger:           logger.Named("CodebaseRootResolver"),
		viewEvents:       viewEvents,
		eventsSender:     eventsSender,
		executorProvider: executorProvider,
		analyticsService: analyticsService,

		authService:         authService,
		codebaseService:     codebaseService,
		organizationService: organizationService,
		changeService:       changeService,
		remoteService:       remoteService,
	}
}

func (r *CodebaseRootResolver) Codebase(ctx context.Context, args resolvers.CodebaseArgs) (resolvers.CodebaseResolver, error) {
	// Lookup single
	if args.ID != nil {
		cb, err := r.resolveCodebase(ctx, *args.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}
			return nil, gqlerrors.Error(fmt.Errorf("failed to resolve by id: %w", err))
		}
		return cb, nil
	}

	// Lookup single by short ID
	if args.ShortID != nil {
		cb, err := r.resolveCodebaseByShort(ctx, *args.ShortID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}
			return nil, gqlerrors.Error(fmt.Errorf("failed to resolve by short id: %w", err))
		}
		return cb, nil
	}

	return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "message", "one of 'id' or 'shortID' must be present")
}

func (r *CodebaseRootResolver) Codebases(ctx context.Context) ([]resolvers.CodebaseResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		// for unauthenticated users, we return 0 codebases
		return nil, nil
	}

	// Lookup all by user
	cbs, err := r.codebaseUserRepo.GetByUser(userID)
	if err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to get codebases by user: %w", err))
	}

	var res []resolvers.CodebaseResolver
	for _, cb := range cbs {
		cbr, err := r.resolveCodebase(ctx, graphql.ID(cb.CodebaseID))
		switch {
		case err == nil:
			res = append(res, cbr)
		case errors.Is(err, sql.ErrNoRows):
			continue
		default:
			return nil, gqlerrors.Error(fmt.Errorf("failed to resolve codebase: %w", err))
		}
	}

	return res, nil
}

func (r *CodebaseRootResolver) CreateCodebase(ctx context.Context, args resolvers.CreateCodebaseArgs) (resolvers.CodebaseResolver, error) {
	var orgID *string
	if args.Input.OrganizationID != nil {
		o := string(*args.Input.OrganizationID)
		orgID = &o

		// Verify access to organization
		org, err := r.organizationService.GetByID(ctx, o)
		if err != nil {
			return nil, gqlerrors.Error(err)
		}
		if err := r.authService.CanWrite(ctx, org); err != nil {
			return nil, gqlerrors.Error(err)
		}
	}

	cb, err := r.codebaseService.Create(ctx, args.Input.Name, orgID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	return r.resolveCodebase(ctx, graphql.ID(cb.ID))
}

func (r *CodebaseRootResolver) UpdatedCodebase(ctx context.Context) (<-chan resolvers.CodebaseResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	c := make(chan resolvers.CodebaseResolver, 100)
	didErrorOut := false

	cancelFunc := r.viewEvents.SubscribeUser(userID, func(et events.EventType, reference string) error {
		select {
		case <-ctx.Done():
			return events.ErrClientDisconnected
		default:
		}

		if et == events.CodebaseUpdated {
			id := graphql.ID(reference)
			resolver, err := r.Codebase(ctx, resolvers.CodebaseArgs{ID: &id})
			if err != nil {
				return err
			}
			select {
			case <-ctx.Done():
				return events.ErrClientDisconnected
			case c <- resolver:
				if didErrorOut {
					didErrorOut = false
				}
				return nil
			default:
				r.logger.Error("dropped subscription event",
					zap.Stringer("user_id", userID),
					zap.Stringer("event_type", et),
					zap.Int("channel_size", len(c)),
				)
				didErrorOut = true
				return nil
			}
		}
		return nil
	})

	go func() {
		<-ctx.Done()
		cancelFunc()
	}()

	return c, nil
}

func (r *CodebaseRootResolver) UpdateCodebase(ctx context.Context, args resolvers.UpdateCodebaseArgs) (resolvers.CodebaseResolver, error) {
	cb, err := r.codebaseService.GetByID(ctx, codebases.ID(args.Input.ID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	// Auth
	if err := r.authService.CanWrite(ctx, cb); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if args.Input.Name != nil && len(*args.Input.Name) > 0 {
		cb.Name = *args.Input.Name
	}
	if args.Input.DisableInviteCode != nil {
		cb.InviteCode = nil
	}
	if args.Input.GenerateInviteCode != nil {
		// Generate new code
		token := make([]byte, 10)
		_, err = rand.Read(token)
		if err != nil {
			return nil, gqlerrors.Error(fmt.Errorf("failed to generate invite code: %w", err))
		}
		inviteCode := base62.EncodeToString(token)
		cb.InviteCode = &inviteCode
	}
	if args.Input.Archive != nil {
		t := time.Now()
		cb.ArchivedAt = &t
	}
	if args.Input.IsPublic != nil {
		cb.IsPublic = *args.Input.IsPublic
		// track, will be used to review malicious activity and the codebases that are made public
		r.analyticsService.Capture(ctx, "ser codebase is_public", analytics.CodebaseID(cb.ID))
	}

	if err := r.codebaseService.Update(ctx, cb); err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to update codebase: %w", err))
	}

	return &CodebaseResolver{c: cb, root: r}, nil
}
func (r *CodebaseRootResolver) AddUserToCodebase(ctx context.Context, args resolvers.AddUserToCodebaseArgs) (resolvers.CodebaseResolver, error) {
	cb, err := r.codebaseService.GetByID(ctx, codebases.ID(args.Input.CodebaseID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	// Auth
	if err := r.authService.CanWrite(ctx, cb); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if _, err := r.codebaseService.AddUserByEmail(ctx, codebases.ID(args.Input.CodebaseID), args.Input.Email); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &CodebaseResolver{c: cb, root: r}, nil
}

func (r *CodebaseRootResolver) RemoveUserFromCodebase(ctx context.Context, args resolvers.RemoveUserFromCodebaseArgs) (resolvers.CodebaseResolver, error) {
	cb, err := r.codebaseService.GetByID(ctx, codebases.ID(args.Input.CodebaseID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	// Auth
	if err := r.authService.CanWrite(ctx, cb); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.codebaseService.RemoveUser(ctx, codebases.ID(args.Input.CodebaseID), users.ID(args.Input.UserID)); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &CodebaseResolver{c: cb, root: r}, nil
}

type CodebaseResolver struct {
	c    *codebases.Codebase
	root *CodebaseRootResolver

	lastUpdatedAt     *int32
	lastUpdatedAtOnce sync.Once
}

func (r *CodebaseResolver) ID() graphql.ID {
	return graphql.ID(r.c.ID)
}

func (r *CodebaseResolver) Name() string {
	return r.c.Name
}

func (r *CodebaseResolver) Slug() string {
	return r.c.Slug()
}

func (r *CodebaseResolver) ShortID() graphql.ID {
	return graphql.ID(r.c.ShortCodebaseID)
}

func (r *CodebaseResolver) Description() string {
	return r.c.Description
}

func (r *CodebaseResolver) InviteCode() *string {
	if r.c.InviteCode == nil {
		return nil
	}
	return r.c.InviteCode
}

func (r *CodebaseResolver) CreatedAt() int32 {
	if r.c.CreatedAt == nil {
		return 0
	}
	return int32(r.c.CreatedAt.Unix())
}

func (r *CodebaseResolver) ArchivedAt() *int32 {
	if r.c.ArchivedAt == nil {
		return nil
	}
	t := int32(r.c.ArchivedAt.Unix())
	return &t
}

func (r *CodebaseResolver) calculateLastUpdatedAt(ctx context.Context) *int32 {
	var largestTime int32
	var zero int32 = 0

	headChange, err := r.root.changeService.HeadChange(ctx, r.c)
	if err != nil {
		return &zero
	}

	maybeTime := []*time.Time{
		headChange.CreatedAt,
		headChange.GitCreatedAt,
	}

	for _, t := range maybeTime {
		if t == nil {
			continue
		}
		t2 := int32(t.Unix())
		if t2 > largestTime {
			largestTime = t2
		}
	}

	if largestTime > 0 {
		return &largestTime
	}

	return nil
}

func (r *CodebaseResolver) LastUpdatedAt(ctx context.Context) *int32 {
	r.lastUpdatedAtOnce.Do(func() {
		r.lastUpdatedAt = r.calculateLastUpdatedAt(ctx)
	})
	return r.lastUpdatedAt
}

func (r *CodebaseResolver) Workspaces(ctx context.Context) ([]resolvers.WorkspaceResolver, error) {
	workspaces, err := r.root.workspaceReader.ListByCodebaseIDs([]codebases.ID{r.c.ID}, false)
	if err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to list workspaces by codebase id: %w", err))
	}

	var res []resolvers.WorkspaceResolver
	for _, ws := range workspaces {
		resolver, err := (*r.root.workspaceResolver).Workspace(ctx, resolvers.WorkspaceArgs{ID: graphql.ID(ws.ID)})
		switch {
		case err == nil:
			res = append(res, resolver)
		case errors.Is(err, sql.ErrNoRows):
			continue
		default:
			return nil, gqlerrors.Error(fmt.Errorf("failed to get workspace by id: %w", err))
		}
	}

	return res, nil
}

func (r *CodebaseResolver) Members(ctx context.Context, args resolvers.CodebaseMembersArgs) (resolvers []resolvers.AuthorResolver, err error) {
	userIDs := make(map[users.ID]struct{})

	// Get direct members (members of the codebase)
	if args.FilterDirectAccess == nil || *args.FilterDirectAccess {
		codebaseUsers, err := r.root.codebaseUserRepo.GetByCodebase(r.c.ID)
		if err != nil {
			return nil, gqlerrors.Error(fmt.Errorf("failed to get codebase members: %w", err))
		}
		for _, cu := range codebaseUsers {
			userIDs[cu.UserID] = struct{}{}
		}
	}

	// Get indirect members (members of the organization)
	if args.FilterDirectAccess == nil || !*args.FilterDirectAccess {
		if r.c.OrganizationID != nil {
			members, err := r.root.organizationService.Members(ctx, *r.c.OrganizationID)
			if err != nil {
				return nil, gqlerrors.Error(fmt.Errorf("failed to get organization members: %w", err))
			}
			for _, member := range members {
				userIDs[member.UserID] = struct{}{}
			}
		}
	}

	// stable order
	ids := make([]users.ID, 0, len(userIDs))
	for userID := range userIDs {
		ids = append(ids, userID)
	}
	sort.Slice(ids, func(i, j int) bool {
		return strings.Compare(string(ids[i]), string(ids[j])) < 0
	})

	for _, userID := range ids {
		author, err := r.root.authorResolver.Author(ctx, graphql.ID(userID))
		switch {
		case err == nil:
			resolvers = append(resolvers, author)
		case errors.Is(err, sql.ErrNoRows):
			continue
		default:
			return nil, gqlerrors.Error(fmt.Errorf("failed to get author by user id: %w", err))
		}
	}

	return
}

func (r *CodebaseResolver) Views(ctx context.Context, args resolvers.CodebaseViewsArgs) (res []resolvers.ViewResolver, err error) {
	var views []*view.View

	if args.IncludeOthers != nil && *args.IncludeOthers {
		views, err = r.root.viewRepo.ListByCodebase(r.c.ID)
	} else if subj, ok := auth.FromContext(ctx); ok && subj.Type == auth.SubjectUser {
		views, err = r.root.viewRepo.ListByCodebaseAndUser(r.c.ID, users.ID(subj.ID))
	}
	if err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to list views: %w", err))
	}

	for _, v := range views {
		viewResolver, err := (*r.root.viewResolver).View(ctx, resolvers.ViewArgs{ID: graphql.ID(v.ID)})
		switch {
		case err == nil:
			res = append(res, viewResolver)
		case errors.Is(err, sql.ErrNoRows):
			continue
		default:
			return nil, gqlerrors.Error(fmt.Errorf("failed to resolve view by id: %w", err))
		}
	}
	return
}

func (r *CodebaseResolver) LastUsedView(ctx context.Context) (resolvers.ViewResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		// for unauthenticated users, no view is considered the last used view
		return nil, nil
	}
	return (*r.root.viewResolver).InternalLastUsedViewByUser(ctx, r.c.ID, userID)
}

func (r *CodebaseResolver) GitHubIntegration(ctx context.Context) (resolvers.CodebaseGitHubIntegrationResolver, error) {
	resolver, err := r.root.codebaseGitHubIntegrationResolver.InternalCodebaseGitHubIntegration(ctx, r.ID())
	switch {
	case err == nil:
		return resolver, nil
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	default:
		return nil, gqlerrors.Error(err)
	}
}

func (r *CodebaseResolver) IsReady() bool {
	return r.c.IsReady
}

func (r *CodebaseResolver) ACL(ctx context.Context) (resolvers.ACLResolver, error) {
	resolver, err := r.root.aclResolver.InternalACLByCodebaseID(ctx, graphql.ID(r.c.ID))
	switch {
	case err == nil:
		return resolver, nil
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	default:
		return nil, gqlerrors.Error(err)
	}
}

func (r *CodebaseResolver) Changes(ctx context.Context, args *resolvers.CodebaseChangesArgs) ([]resolvers.ChangeResolver, error) {
	const defaultLimit int = 100
	var (
		limit  = defaultLimit
		before *graphql.ID
	)
	if args != nil && args.Input != nil {
		if args.Input.Limit != nil && *args.Input.Limit <= 100 {
			limit = int(*args.Input.Limit)
		}

		before = args.Input.Before
	}
	return r.root.changeRootResolver.IntenralListChanges(ctx, r.c.ID, limit, before)
}

func (r *CodebaseResolver) Readme(ctx context.Context) (resolvers.FileResolver, error) {
	// GitHub supported names:
	// https://github.com/github/markup/blob/master/README.md
	fileResolver, err := r.root.fileRootResolver.InternalFile(ctx, r.c, "README.md", "README.mkdn", "README.mdown", "README.markdown")
	switch {
	case err == nil && fileResolver != nil:
		if file, ok := fileResolver.ToFile(); ok {
			return file, nil
		} else {
			return nil, nil
		}
	case errors.Is(err, gqlerrors.ErrNotFound):
		return nil, nil
	default:
		return nil, gqlerrors.Error(err)
	}
}

func (r *CodebaseResolver) File(ctx context.Context, args resolvers.CodebaseFileArgs) (resolvers.FileOrDirectoryResolver, error) {
	fr, err := r.root.fileRootResolver.InternalFile(ctx, r.c, args.Path)
	switch {
	case err == nil:
		return fr, nil
	case errors.Is(err, gqlerrors.ErrNotFound):
		return nil, nil
	default:
		return nil, gqlerrors.Error(err)
	}
}

func (r *CodebaseResolver) Integrations(ctx context.Context, args resolvers.IntegrationsArgs) ([]resolvers.IntegrationResolver, error) {
	if args.ID != nil {
		single, err := r.root.instantIntegrationRootResolver.InternalIntegrationByID(ctx, string(*args.ID))
		if err != nil {
			return nil, err
		}
		return []resolvers.IntegrationResolver{single}, nil
	}

	return r.root.instantIntegrationRootResolver.InternalIntegrationsByCodebaseID(ctx, r.c.ID)
}

func (r *CodebaseResolver) IsPublic() bool {
	return r.c.IsPublic
}

func (r *CodebaseResolver) Organization(ctx context.Context) (resolvers.OrganizationResolver, error) {
	if r.c.OrganizationID == nil {
		return nil, nil
	}
	id := graphql.ID(*r.c.OrganizationID)
	res, err := (*r.root.organizationRootResolver).Organization(ctx, resolvers.OrganizationArgs{ID: &id})
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return res, nil
}

func (r *CodebaseResolver) Remote(ctx context.Context) (resolvers.RemoteResolver, error) {
	resolver, err := r.root.remoteRootResolver.InternalRemoteByCodebaseID(ctx, codebases.ID(r.ID()))
	switch {
	case err == nil:
		return resolver, nil
	case errors.Is(err, sql.ErrNoRows), errors.Is(err, gqlerrors.ErrNotFound):
		return nil, nil
	default:
		return nil, gqlerrors.Error(err)
	}
}

func (r *CodebaseResolver) Writeable(ctx context.Context) bool {
	if err := r.root.authService.CanWrite(ctx, r.c); err == nil {
		return true
	}
	return false
}

func (r *CodebaseRootResolver) resolveCodebase(ctx context.Context, id graphql.ID) (*CodebaseResolver, error) {
	c, err := r.codebaseRepo.Get(codebases.ID(id))
	if err != nil {
		return nil, err
	}

	if err := r.authService.CanRead(ctx, c); err != nil {
		return nil, err
	}

	return &CodebaseResolver{c: c, root: r}, nil
}

func (r *CodebaseRootResolver) resolveCodebaseByShort(ctx context.Context, shortID graphql.ID) (*CodebaseResolver, error) {
	s := string(shortID)
	if idx := strings.LastIndex(s, "-"); idx >= 0 {
		s = s[idx+1:]
	}

	c, err := r.codebaseRepo.GetByShortID(codebases.ShortCodebaseID(s))
	if err != nil {
		return nil, err
	}

	if err := r.authService.CanRead(ctx, c); err != nil {
		return nil, err
	}

	return &CodebaseResolver{c: c, root: r}, nil
}

func (r *CodebaseRootResolver) PullCodebase(ctx context.Context, args resolvers.PullCodebaseArgs) (resolvers.CodebaseResolver, error) {
	c, err := r.codebaseRepo.Get(codebases.ID(args.Input.CodebaseID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, c); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.remoteService.Pull(ctx, c.ID); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &CodebaseResolver{c: c, root: r}, nil
}

func (r *CodebaseRootResolver) PushCodebase(ctx context.Context, args resolvers.PushCodebaseArgs) (resolvers.CodebaseResolver, error) {
	c, err := r.codebaseRepo.Get(codebases.ID(args.Input.CodebaseID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, c); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.remoteService.PushTrunk(ctx, c.ID); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &CodebaseResolver{c: c, root: r}, nil
}
