package stream

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/ctxlog"
	"getsturdy.com/api/pkg/events"
	service_suggestions "getsturdy.com/api/pkg/suggestions/service"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/pkg/view"
	"getsturdy.com/api/pkg/workspaces"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"

	"go.uber.org/zap"
)

type Event struct {
	Name    EventName
	Message interface{}
}

type EventName string

const (
	Ping            EventName = "Ping"
	Diffs           EventName = "Diffs"
	CodebaseUpdated EventName = "CodebaseUpdated"
	ConflictDiffs   EventName = "ConflictDiffs"
)

type ViewDiffEvent struct {
	Diffs []unidiff.FileDiff `json:"diffs"`
}

type getAllowerer interface {
	GetAllower(ctx context.Context, obj interface{}) (*unidiff.Allower, error)
}

// Stream
//
// The expectedWorkspaceID makes sure that the view is using the expected workspace
func Stream(
	ctx context.Context,
	logger *zap.Logger,
	ws *workspaces.Workspace,
	vw *view.View,
	done chan bool,
	viewEventsReader events.EventReader,
	allowerProvider getAllowerer,
	workspaceService service_workspace.Service,
	suggestionsServcie *service_suggestions.Service,
) (chan Event, error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("recovered in stream", zap.Any("err", r))
		}
	}()

	chanStream := make(chan Event)
	ticker := time.Tick(time.Second * 10)

	// mx protects disconnected
	var disconnected bool
	var mx sync.RWMutex

	send := func(msg Event) {
		// To read disconnected flag and use chanStream
		mx.Lock()
		if disconnected {
			mx.Unlock()
			return
		}
		mx.Unlock()

		// It's possible for this to panic (send on closed channel)
		// TODO: Can mx protect this send?
		chanStream <- msg
	}

	sendCodebaseUpdated := func() {
		send(Event{
			Name: CodebaseUpdated,
		})
	}

	sendDiff := func() {
		evt, err := getAllDiffs(ctx, workspaceService, ws, allowerProvider, suggestionsServcie)
		if err != nil {
			ctxlog.ErrorOrWarn(logger, "could not get diffs", err)
			return
		}

		send(evt)
	}

	// Subscribe to real-time updates from the view
	var viewID *string
	if vw != nil {
		viewID = &vw.ID
	}
	sendDiffsChan := make(chan time.Time, 4)

	callbackFunc := func(eventType events.EventType, reference string) error {
		// This can be racy, but that's OK.
		if disconnected {
			return events.ErrClientDisconnected
		}

		workspaceUpdated := eventType == events.WorkspaceUpdated && reference == ws.ID
		viewUpdated := eventType == events.ViewUpdated && viewID != nil && reference == *viewID
		workspaceSnapshotUpdated := eventType == events.WorkspaceUpdatedSnapshot && reference == ws.ID
		diffsUpdated := workspaceUpdated || viewUpdated || workspaceSnapshotUpdated

		if diffsUpdated {
			select {
			case sendDiffsChan <- time.Now():
			default:
				logger.Info("dropping event to sendDiffsChan")
			}
		}

		if eventType == events.CodebaseEvent {
			sendCodebaseUpdated()
		}

		return nil
	}

	cancelWSSubscription := viewEventsReader.SubscribeWorkspace(ws.ID, callbackFunc)
	cancelSubscription := func() {
		cancelWSSubscription()
	}
	if userID, err := auth.UserID(ctx); err == nil {
		cancelUserSubscription := viewEventsReader.SubscribeUser(userID, callbackFunc)
		cancelSubscription = func() {
			cancelUserSubscription()
			cancelWSSubscription()
		}
	}

	var lastDiffAt time.Time

	go func() {
		for enqueuedAt := range sendDiffsChan {
			// a diff has been generated after this event, no need to create a new one
			if !lastDiffAt.IsZero() && enqueuedAt.Before(lastDiffAt) {
				logger.Info("dropping event in sendDiffsChan, we already have newer data")
				continue
			}
			sendDiff()
			lastDiffAt = time.Now()
		}
	}()

	// Early send diff. Needs to be done in the background, so that a message can be sent to done
	// if the client is disconnected.
	sendDiffsChan <- time.Now()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("recovered in stream loop", zap.Any("err", r))
			}
		}()

		for {
			select {
			case <-done:
				mx.Lock()
				disconnected = true
				close(chanStream)
				close(sendDiffsChan)
				cancelSubscription()
				mx.Unlock()

				// Exit this goroutine
				return
			case <-ticker:
				// Send ping to see if user is still connected
				chanStream <- Event{
					Name:    Ping,
					Message: struct{}{},
				}
			}
		}
	}()

	return chanStream, nil
}

func getAllDiffs(
	ctx context.Context,
	workspaceService service_workspace.Service,
	ws *workspaces.Workspace,
	allowerProvider getAllowerer,
	suggestionsService *service_suggestions.Service,
) (Event, error) {
	allower, err := allowerProvider.GetAllower(ctx, ws)
	if err != nil {
		return Event{}, fmt.Errorf("failed to get allowed patterns: %w", err)
	}

	suggestion, err := suggestionsService.GetByWorkspaceID(ctx, ws.ID)
	switch {
	case err == nil:
		diffs, err := suggestionsService.Diffs(ctx, suggestion, unidiff.WithAllower(allower))
		if err != nil {
			return Event{}, fmt.Errorf("failed to get diffs: %w", err)
		}
		return Event{
			Name:    Diffs,
			Message: ViewDiffEvent{diffs},
		}, nil
	case errors.Is(err, sql.ErrNoRows):
		diffs, isConflicting, err := workspaceService.Diffs(ctx, ws.ID, service_workspace.WithAllower(allower))
		if err != nil {
			return Event{}, fmt.Errorf("failed to get diffs: %w", err)
		}
		if isConflicting {
			return Event{Name: ConflictDiffs, Message: ViewDiffEvent{diffs}}, nil
		}
		return Event{Name: Diffs, Message: ViewDiffEvent{diffs}}, nil
	default:
		return Event{}, fmt.Errorf("failed to get suggestion: %w", err)
	}
}
