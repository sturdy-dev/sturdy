package graphql

import (
	"context"
	"database/sql"
	"errors"

	"getsturdy.com/api/pkg/auth"
	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/events"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/suggestions"
	service_suggestions "getsturdy.com/api/pkg/suggestions/service"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"

	"go.uber.org/zap"
)

type RootResolver struct {
	logger *zap.Logger

	authService        *service_auth.Service
	suggestionsService *service_suggestions.Service
	workspaceService   service_workspace.Service

	authorResolver    resolvers.AuthorRootResolver
	fileDiffResolver  resolvers.FileDiffRootResolver
	workspaceResolver *resolvers.WorkspaceRootResolver

	eventsReader events.EventReader
}

func New(
	logger *zap.Logger,

	authService *service_auth.Service,
	suggestionsService *service_suggestions.Service,
	workspaceService service_workspace.Service,

	authorResolver resolvers.AuthorRootResolver,
	fileDiffResolver resolvers.FileDiffRootResolver,
	workspaceResolver *resolvers.WorkspaceRootResolver,

	eventsReader events.EventReader,
) resolvers.SuggestionRootResolver {
	return &RootResolver{
		logger: logger,

		authService:        authService,
		suggestionsService: suggestionsService,
		workspaceService:   workspaceService,

		authorResolver:    authorResolver,
		fileDiffResolver:  fileDiffResolver,
		workspaceResolver: workspaceResolver,

		eventsReader: eventsReader,
	}
}

func (r *RootResolver) InternalSuggestion(ctx context.Context, suggestion *suggestions.Suggestion) (resolvers.SuggestionResolver, error) {
	return &Resolver{
		root:       r,
		suggestion: suggestion,
	}, nil
}

func (r *RootResolver) InternalSuggestionByID(ctx context.Context, id suggestions.ID) (resolvers.SuggestionResolver, error) {
	s, err := r.suggestionsService.GetByID(ctx, id)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return r.InternalSuggestion(ctx, s)
}

func (r *RootResolver) CreateSuggestion(ctx context.Context, args resolvers.CreateSuggestionArgs) (resolvers.SuggestionResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	ws, err := r.workspaceService.GetByID(ctx, string(args.Input.WorkspaceID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanRead(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	suggestion, err := r.suggestionsService.Create(ctx, userID, ws)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	return r.InternalSuggestion(ctx, suggestion)
}

func (r *RootResolver) ApplySuggestionHunks(ctx context.Context, args resolvers.ApplySuggestionHunksArgs) (resolvers.SuggestionResolver, error) {
	suggestion, err := r.suggestionsService.GetByID(ctx, suggestions.ID(args.Input.ID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, suggestion); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.suggestionsService.ApplyHunks(ctx, suggestion, args.Input.HunkIDs...); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return r.InternalSuggestion(ctx, suggestion)
}

func (r *RootResolver) DismissSuggestion(ctx context.Context, args resolvers.DismissSuggestionArgs) (resolvers.SuggestionResolver, error) {
	suggestion, err := r.suggestionsService.GetByID(ctx, suggestions.ID(args.Input.ID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, suggestion); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.suggestionsService.Dismiss(ctx, suggestion); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return r.InternalSuggestion(ctx, suggestion)
}

func (r *RootResolver) DismissSuggestionHunks(ctx context.Context, args resolvers.DismissSuggestionHunksArgs) (resolvers.SuggestionResolver, error) {
	suggestion, err := r.suggestionsService.GetByID(ctx, suggestions.ID(args.Input.ID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, suggestion); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.suggestionsService.DismissHunks(ctx, suggestion, args.Input.HunkIDs...); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return r.InternalSuggestion(ctx, suggestion)
}

func (r *RootResolver) UpdatedSuggestion(ctx context.Context, args resolvers.UpdatedSuggestionArgs) (chan resolvers.SuggestionResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	res := make(chan resolvers.SuggestionResolver, 100)
	didErrorOut := false

	listeningTo := map[events.EventType]bool{
		events.WorkspaceUpdated:           true,
		events.WorkspaceUpdatedSuggestion: true,
	}

	cancelFunc := r.eventsReader.SubscribeUser(userID, func(eventType events.EventType, reference string) error {
		select {
		case <-ctx.Done():
			return events.ErrClientDisconnected
		default:
		}

		if !listeningTo[eventType] {
			return nil
		}

		// suggestions opened for the workspace
		ss, err := r.suggestionsService.ListForWorkspaceID(ctx, reference)
		if err != nil {
			return gqlerrors.Error(err)
		}

		// sugestion that the workspace might contain
		suggestion, err := r.suggestionsService.GetByWorkspaceID(ctx, reference)
		switch {
		case err == nil:
			ss = append(ss, suggestion)
		case errors.Is(err, sql.ErrNoRows):
			// do nothing
		default:
			return gqlerrors.Error(err)
		}

		for _, suggestion := range ss {
			resolver, err := r.InternalSuggestion(ctx, suggestion)
			if err != nil {
				return err
			}

			select {
			case <-ctx.Done():
				return events.ErrClientDisconnected
			case res <- resolver:
				if didErrorOut {
					didErrorOut = false
				}
			default:
				r.logger.Error("dropped subscription event",
					zap.Stringer("user_id", userID),
					zap.Stringer("event_type", eventType),
					zap.Int("channel_size", len(res)),
				)
				didErrorOut = true
			}
		}

		return nil
	})

	go func() {
		<-ctx.Done()
		cancelFunc()
	}()

	return res, nil
}
