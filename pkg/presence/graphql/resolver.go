package graphql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"mash/pkg/auth"
	service_auth "mash/pkg/auth/service"
	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"
	"mash/pkg/presence"
	service_presence "mash/pkg/presence/service"
	"mash/pkg/view/events"

	"github.com/graph-gophers/graphql-go"
	"go.uber.org/zap"
)

type presenceRootResolver struct {
	presenceService service_presence.Service

	authorRootResolver    *resolvers.AuthorRootResolver
	workspaceRootResolver *resolvers.WorkspaceRootResolver
	authService           *service_auth.Service

	logger           *zap.Logger
	eventsReadWriter events.EventReadWriter
}

func NewRootResolver(
	presenceService service_presence.Service,
	authorRootResolver *resolvers.AuthorRootResolver,
	workspaceRootResolver *resolvers.WorkspaceRootResolver,
	logger *zap.Logger,
	eventsReadWriter events.EventReadWriter,
) resolvers.PresenceRootResolver {
	return &presenceRootResolver{
		presenceService:       presenceService,
		authorRootResolver:    authorRootResolver,
		workspaceRootResolver: workspaceRootResolver,
		logger:                logger.Named("PresenceRootResolver"),
		eventsReadWriter:      eventsReadWriter,
	}
}

func (r *presenceRootResolver) InternalWorkspacePresence(ctx context.Context, workspaceID string) ([]resolvers.PresenceResolver, error) {
	presences, err := r.presenceService.ListByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	res := make([]resolvers.PresenceResolver, 0, len(presences))
	for _, p := range presences {
		res = append(res, &presenceResolver{root: r, pre: p})
	}

	return res, nil
}

func (r *presenceRootResolver) ReportWorkspacePresence(ctx context.Context, args resolvers.ReportWorkspacePresenceArgs) (resolvers.PresenceResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	var state presence.State
	switch args.Input.State {
	case resolvers.PresenceStateCoding:
		state = presence.StateCoding
	case resolvers.PresenceStateViewing:
		state = presence.StateViewing
	case resolvers.PresenceStateIdle:
		state = presence.StateIdle
	default:
		return nil, gqlerrors.Error(fmt.Errorf("failed to set state, invalid value: %s", args.Input.State))
	}

	pre, err := r.presenceService.Record(ctx, userID, string(args.Input.WorkspaceID), state)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &presenceResolver{root: r, pre: pre}, nil
}

func (r *presenceRootResolver) UpdatedWorkspacePresence(ctx context.Context, args resolvers.UpdatedWorkspacePresenceArgs) (chan resolvers.PresenceResolver, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	res := make(chan resolvers.PresenceResolver, 100)
	didErrorOut := false

	cancelFunc := r.eventsReadWriter.SubscribeUser(userID, func(eventType events.EventType, reference string) error {
		if eventType == events.WorkspaceUpdatedPresence &&
			(args.WorkspaceID == nil || // Subscribed to all workspaces
				reference == string(*args.WorkspaceID)) { // Subscribed to a specific workspace

			// TODO: Send only the presence that updated
			resolvers, err := r.InternalWorkspacePresence(ctx, reference)
			if err != nil {
				return err
			}

			// Send all
			for _, resolver := range resolvers {
				r2 := resolver
				select {
				case <-ctx.Done():
					return errors.New("disconnected")
				case res <- r2:
					if didErrorOut {
						didErrorOut = false
					}
				default:
					r.logger.Error("dropped subscription event",
						zap.String("user_id", userID),
						zap.Stringer("event_type", eventType),
						zap.Int("channel_size", len(res)),
					)
					didErrorOut = true
				}
			}
		}
		return nil
	})

	go func() {
		<-ctx.Done()
		cancelFunc()
		close(res)
	}()

	return res, nil
}

type presenceResolver struct {
	root *presenceRootResolver
	pre  *presence.Presence
}

func (r *presenceResolver) ID() graphql.ID {
	return graphql.ID(r.pre.ID)
}

func (r *presenceResolver) Author(ctx context.Context) (resolvers.AuthorResolver, error) {
	author, err := (*r.root.authorRootResolver).Author(ctx, graphql.ID(r.pre.UserID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return author, nil
}

func (r *presenceResolver) State() (resolvers.PresenceState, error) {
	switch r.pre.State {
	case presence.StateCoding:
		return resolvers.PresenceStateCoding, nil
	case presence.StateIdle:
		return resolvers.PresenceStateIdle, nil
	case presence.StateViewing:
		return resolvers.PresenceStateViewing, nil
	default:
		return resolvers.PresenceStateIdle, fmt.Errorf("unexpected presence state")
	}
}

func (r *presenceResolver) LastActiveAt() int32 {
	return int32(r.pre.LastActiveAt.Unix())
}

func (r *presenceResolver) Workspace(ctx context.Context) (resolvers.WorkspaceResolver, error) {
	return (*r.root.workspaceRootResolver).Workspace(ctx, resolvers.WorkspaceArgs{ID: graphql.ID(r.pre.WorkspaceID)})
}
