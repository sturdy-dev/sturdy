package graphql

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"getsturdy.com/api/pkg/analytics"
	"getsturdy.com/api/pkg/auth"
	service_auth "getsturdy.com/api/pkg/auth/service"
	db_change "getsturdy.com/api/pkg/change/db"
	"getsturdy.com/api/pkg/change/decorate"
	"getsturdy.com/api/pkg/codebase"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	service_codebase "getsturdy.com/api/pkg/codebase/service"
	"getsturdy.com/api/pkg/codebase/vcs"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	service_organization "getsturdy.com/api/pkg/organization/service"
	db_user "getsturdy.com/api/pkg/user/db"
	"getsturdy.com/api/pkg/view"
	db_view "getsturdy.com/api/pkg/view/db"
	"getsturdy.com/api/pkg/view/events"
	db_workspace "getsturdy.com/api/pkg/workspace/db"
	vcsvcs "getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	"github.com/graph-gophers/graphql-go"
	"github.com/jxskiss/base62"
	"go.uber.org/zap"
)

type CodebaseRootResolver struct {
	codebaseRepo     db_codebase.CodebaseRepository
	codebaseUserRepo db_codebase.CodebaseUserRepository
	viewRepo         db_view.Repository
	workspaceReader  db_workspace.WorkspaceReader
	userRepo         db_user.Repository
	changeRepo       db_change.Repository
	changeCommitRepo db_change.CommitRepository

	workspaceResolver                 *resolvers.WorkspaceRootResolver
	authorResolver                    resolvers.AuthorRootResolver
	viewResolver                      *resolvers.ViewRootResolver
	aclResolver                       resolvers.ACLRootResolver
	changeRootResolver                resolvers.ChangeRootResolver
	fileRootResolver                  resolvers.FileRootResolver
	instantIntegrationRootResolver    resolvers.IntegrationRootResolver
	codebaseGitHubIntegrationResolver resolvers.CodebaseGitHubIntegrationRootResolver
	organizationRootResolver          *resolvers.OrganizationRootResolver

	logger           *zap.Logger
	viewEvents       events.EventReader
	eventsSender     events.EventSender
	analyticsClient  analytics.Client
	executorProvider executor.Provider

	authService         *service_auth.Service
	codebaseService     *service_codebase.Service
	organizationService *service_organization.Service
}

func NewCodebaseRootResolver(
	codebaseRepo db_codebase.CodebaseRepository,
	codebaseUserRepo db_codebase.CodebaseUserRepository,
	viewRepo db_view.Repository,
	workspaceReader db_workspace.WorkspaceReader,
	userRepo db_user.Repository,
	changeRepo db_change.Repository,
	changeCommitRepo db_change.CommitRepository,

	workspaceResolver *resolvers.WorkspaceRootResolver,
	authorResolver resolvers.AuthorRootResolver,
	viewResolver *resolvers.ViewRootResolver,
	aclResolver resolvers.ACLRootResolver,
	changeRootResolver resolvers.ChangeRootResolver,
	fileRootResolver resolvers.FileRootResolver,
	instantIntegrationRootResolver resolvers.IntegrationRootResolver,
	codebaseGitHubIntegrationResolver resolvers.CodebaseGitHubIntegrationRootResolver,
	organizationRootResolver *resolvers.OrganizationRootResolver,

	logger *zap.Logger,
	viewEvents events.EventReader,
	eventsSender events.EventSender,
	analyticsClient analytics.Client,
	executorProvider executor.Provider,

	authService *service_auth.Service,
	codebaseService *service_codebase.Service,
	organizationService *service_organization.Service,
) resolvers.CodebaseRootResolver {
	return &CodebaseRootResolver{
		codebaseRepo:     codebaseRepo,
		codebaseUserRepo: codebaseUserRepo,
		viewRepo:         viewRepo,
		workspaceReader:  workspaceReader,
		userRepo:         userRepo,
		changeRepo:       changeRepo,
		changeCommitRepo: changeCommitRepo,

		workspaceResolver:                 workspaceResolver,
		authorResolver:                    authorResolver,
		viewResolver:                      viewResolver,
		aclResolver:                       aclResolver,
		changeRootResolver:                changeRootResolver,
		fileRootResolver:                  fileRootResolver,
		instantIntegrationRootResolver:    instantIntegrationRootResolver,
		codebaseGitHubIntegrationResolver: codebaseGitHubIntegrationResolver,
		organizationRootResolver:          organizationRootResolver,

		logger:           logger.Named("CodebaseRootResolver"),
		viewEvents:       viewEvents,
		eventsSender:     eventsSender,
		analyticsClient:  analyticsClient,
		executorProvider: executorProvider,

		authService:         authService,
		codebaseService:     codebaseService,
		organizationService: organizationService,
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
		if et == events.CodebaseUpdated {
			id := graphql.ID(reference)
			resolver, err := r.Codebase(ctx, resolvers.CodebaseArgs{ID: &id})
			if err != nil {
				return err
			}
			select {
			case <-ctx.Done():
				return errors.New("disconnected")
			case c <- resolver:
				if didErrorOut {
					didErrorOut = false
				}
				return nil
			default:
				r.logger.Error("dropped subscription event",
					zap.String("user_id", userID),
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
		close(c)
	}()

	return c, nil
}

func (r *CodebaseRootResolver) UpdateCodebase(ctx context.Context, args resolvers.UpdateCodebaseArgs) (resolvers.CodebaseResolver, error) {
	// Auth
	if err := r.authService.CanWrite(ctx, &codebase.Codebase{ID: string(args.Input.ID)}); err != nil {
		return nil, gqlerrors.Error(err)
	}

	authSubject, ok := auth.FromContext(ctx)
	if !ok {
		return nil, gqlerrors.Error(fmt.Errorf("could not get auth"))
	}

	cb, err := r.codebaseRepo.Get(string(args.Input.ID))
	if err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to get codebase by id: %w", err))
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
		_ = r.analyticsClient.Enqueue(analytics.Capture{
			Event:      "set codebase is_public",
			DistinctId: authSubject.ID,
			Properties: analytics.NewProperties().Set("codebase_id", cb.ID),
		})
	}

	if err := r.codebaseRepo.Update(cb); err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to update codebase repo: %w", err))
	}

	// Send events
	if err := r.eventsSender.Codebase(cb.ID, events.CodebaseUpdated, cb.ID); err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to send codebase event: %w", err))
	}

	return &CodebaseResolver{c: cb, root: r}, nil
}

type CodebaseResolver struct {
	c    *codebase.Codebase
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

func (r *CodebaseResolver) calculateLastUpdatedAt() *int32 {
	var largestTime int32

	var gitTime time.Time
	if err := r.root.executorProvider.New().Git(func(repo vcsvcs.Repo) error {
		changes, err := vcs.ListChanges(repo, 1)
		if err != nil || len(changes) == 0 {
			return fmt.Errorf("failed to list changes: %w", err)
		}
		gitTime = changes[0].Time
		return nil
	}).ExecTrunk(r.c.ID, "codebase.LastUpdatedAt"); err != nil {
		var zero int32 = 0
		return &zero
	}

	maybeTime := []*time.Time{
		&gitTime,
		r.c.CreatedAt,
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

func (r *CodebaseResolver) LastUpdatedAt() *int32 {
	r.lastUpdatedAtOnce.Do(func() {
		r.lastUpdatedAt = r.calculateLastUpdatedAt()
	})
	return r.lastUpdatedAt
}

func (r *CodebaseResolver) Workspaces(ctx context.Context) ([]resolvers.WorkspaceResolver, error) {
	workspaces, err := r.root.workspaceReader.ListByCodebaseIDs([]string{r.c.ID}, false)
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

func (r *CodebaseResolver) Members(ctx context.Context) (resolvers []resolvers.AuthorResolver, err error) {
	codebaseUsers, err := r.root.codebaseUserRepo.GetByCodebase(r.c.ID)
	if err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to get codebase members: %w", err))
	}

	userIDs := make(map[string]struct{})

	for _, cu := range codebaseUsers {
		userIDs[cu.UserID] = struct{}{}
	}

	// also list members of the organization
	if r.c.OrganizationID != nil {
		members, err := r.root.organizationService.Members(ctx, *r.c.OrganizationID)
		if err != nil {
			return nil, gqlerrors.Error(fmt.Errorf("failed to get organization members: %w", err))
		}
		for _, member := range members {
			userIDs[member.UserID] = struct{}{}
		}
	}

	for userID := range userIDs {
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
		views, err = r.root.viewRepo.ListByCodebaseAndUser(r.c.ID, subj.ID)
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
	var limit = 100
	if args != nil && args.Input != nil && args.Input.Limit != nil && *args.Input.Limit <= 100 {
		limit = int(*args.Input.Limit)
	}

	// vcs.ListChanges and decorate.DecorateChanges will import all commits to Sturdy.
	// This is not ideal. If we could make sure that the database is already is up to date with the Git state,
	// we would not have to read from disk here.
	var log []*vcsvcs.LogEntry
	if err := r.root.executorProvider.New().Git(func(repo vcsvcs.Repo) error {
		var err error
		log, err = vcs.ListChanges(repo, limit)
		if err != nil {
			return fmt.Errorf("failed to list changes: %w", err)
		}
		return nil
	}).ExecTrunk(r.c.ID, "codebase.Changes"); err != nil {
		return nil, gqlerrors.Error(err)
	}

	decoratedLog, err := decorate.DecorateChanges(log, r.root.userRepo, r.root.logger, r.root.changeRepo, r.root.changeCommitRepo, r.root.codebaseUserRepo, r.c.ID, true)
	if err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to decorate changes: %w", err))
	}

	var res []resolvers.ChangeResolver
	for _, dc := range decoratedLog {
		id := graphql.ID(dc.ChangeID)
		r, err := r.root.changeRootResolver.Change(ctx, resolvers.ChangeArgs{ID: &id})
		switch {
		case err == nil:
			res = append(res, r)
		case errors.Is(err, sql.ErrNoRows):
			continue
		default:
			return nil, gqlerrors.Error(fmt.Errorf("failed to get change by id: %w", err))
		}
	}

	return res, nil
}

func (r *CodebaseResolver) Readme(ctx context.Context) (resolvers.FileResolver, error) {
	// GitHub supported names:
	// https://github.com/github/markup/blob/master/README.md
	fileResolver, err := r.root.fileRootResolver.InternalFile(ctx, r.c.ID, "README.md", "README.mkdn", "README.mdown", "README.markdown")
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
	fr, err := r.root.fileRootResolver.InternalFile(ctx, r.c.ID, args.Path)
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

func (r *CodebaseRootResolver) resolveCodebase(ctx context.Context, id graphql.ID) (*CodebaseResolver, error) {
	c, err := r.codebaseRepo.Get(string(id))
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

	c, err := r.codebaseRepo.GetByShortID(s)
	if err != nil {
		return nil, err
	}

	if err := r.authService.CanRead(ctx, c); err != nil {
		return nil, err
	}

	return &CodebaseResolver{c: c, root: r}, nil
}
