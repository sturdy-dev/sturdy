package events

import (
	"sync"

	"getsturdy.com/api/pkg/users"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type EventReader interface {
	// Introduce event type filtering in the call to SubscribeUser(uesrID string, cb CallbackFunc, eventTypes ...EventType) CancelFunc
	// Introduce reference filtering in the call to SubscribeUser(uesrID string, cb CallbackFunc, map[EventType][]string) CancelFunc
	SubscribeUser(userID users.ID, cb CallbackFunc) CancelFunc
	SubscribeWorkspace(workspaceID string, cb CallbackFunc) CancelFunc
}

type eventWriter interface {
	UserEvent(userID users.ID, eventType EventType, reference string)
	WorkspaceEvent(workspaceID string, eventType EventType, reference string)
}

type EventReadWriter interface {
	EventReader
	eventWriter
}

type EventType int

func (e EventType) String() string {
	return eventTypeString[e]
}

const (
	CodebaseEvent EventType = iota
	CodebaseUpdated
	ViewUpdated
	ViewStatusUpdated
	WorkspaceUpdated
	WorkspaceUpdatedComments
	WorkspaceUpdatedReviews
	ReviewUpdated
	WorkspaceUpdatedActivity
	WorkspaceUpdatedSnapshot
	WorkspaceUpdatedPresence
	WorkspaceUpdatedSuggestion
	GitHubPRUpdated
	NotificationEvent
	StatusUpdated
	CompletedOnboardingStep
	WorkspaceWatchingStatusUpdated
)

var eventTypeString = map[EventType]string{
	CodebaseEvent:                  "CodebaseEvent",
	CodebaseUpdated:                "CodebaseUpdated",
	ViewUpdated:                    "ViewUpdated",
	ViewStatusUpdated:              "ViewStatusUpdated",
	WorkspaceUpdated:               "WorkspaceUpdated",
	WorkspaceUpdatedComments:       "WorkspaceUpdatedComments",
	WorkspaceUpdatedReviews:        "WorkspaceUpdatedReviews",
	WorkspaceUpdatedActivity:       "WorkspaceUpdatedActivity",
	WorkspaceUpdatedPresence:       "WorkspaceUpdatedPresence",
	WorkspaceUpdatedSuggestion:     "WorkspaceUpdatedSuggestion",
	ReviewUpdated:                  "ReviewUpdated",
	GitHubPRUpdated:                "GitHubPRUpdated",
	NotificationEvent:              "NotificationEvent",
	StatusUpdated:                  "StatusUpdated",
	CompletedOnboardingStep:        "CompletedOnboardingStep",
	WorkspaceWatchingStatusUpdated: "WorkspaceWatchingStatusUpdated",
}

type CallbackFunc func(eventType EventType, reference string) error

type CancelFunc func()

var (
	sentEventCounterMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sturdy_sent_event_counter",
	}, []string{"eventType"})
	receivedEventCounterMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sturdy_received_event_counter",
	}, []string{"eventType", "success"})
)

type inMemory struct {
	mx          sync.RWMutex
	subscribers map[Topic]map[string]CallbackFunc
}

func NewInMemory() EventReadWriter {
	return &inMemory{
		subscribers: make(map[Topic]map[string]CallbackFunc),
	}
}

func (i *inMemory) UserEvent(userID users.ID, eventType EventType, reference string) {
	i.event(Topic(userID), eventType, reference)
}

func (i *inMemory) WorkspaceEvent(workspaceID string, eventType EventType, reference string) {
	i.event(Topic(workspaceID), eventType, reference)
}

type Topic string

type TopicSubscriber struct {
	Topic         Topic
	SubscriberKey string
}

func (i *inMemory) event(topic Topic, eventType EventType, reference string) {
	sentEventCounterMetric.WithLabelValues(eventTypeString[eventType]).Inc()

	subs := make(map[TopicSubscriber]CallbackFunc, len(i.subscribers[topic]))

	i.mx.RLock()
	// Copy over all subscribers
	for subId, cb := range i.subscribers[topic] {
		subs[TopicSubscriber{Topic: topic, SubscriberKey: subId}] = cb
	}
	i.mx.RUnlock()

	var unregKeys []TopicSubscriber
	for ts, cb := range subs {
		if err := cb(eventType, reference); err != nil {
			receivedEventCounterMetric.WithLabelValues(eventTypeString[eventType], "no").Inc()
			// TODO: Log errors here (if they are not of some specific ClientDisconnected-type)
			unregKeys = append(unregKeys, ts)
		} else {
			receivedEventCounterMetric.WithLabelValues(eventTypeString[eventType], "yes").Inc()
		}
	}

	if len(unregKeys) > 0 {
		i.unreg(unregKeys...)
	}
}

func (i *inMemory) SubscribeWorkspace(workspaceID string, cb CallbackFunc) CancelFunc {
	workspaceTopic := Topic(workspaceID)

	id := uuid.New().String()

	i.mx.Lock()
	_, ok := i.subscribers[workspaceTopic]
	if !ok {
		i.subscribers[workspaceTopic] = make(map[string]CallbackFunc)
	}
	i.subscribers[workspaceTopic][id] = cb

	unregKey := TopicSubscriber{
		Topic:         workspaceTopic,
		SubscriberKey: id,
	}
	i.mx.Unlock()

	return func() { i.unreg(unregKey) }
}

func (i *inMemory) SubscribeUser(userID users.ID, cb CallbackFunc) CancelFunc {
	userTopic := Topic(userID)

	id := uuid.New().String()

	i.mx.Lock()
	_, ok := i.subscribers[userTopic]
	if !ok {
		i.subscribers[userTopic] = make(map[string]CallbackFunc)
	}
	i.subscribers[userTopic][id] = cb

	unregKey := TopicSubscriber{
		Topic:         userTopic,
		SubscriberKey: id,
	}
	i.mx.Unlock()

	return func() { i.unreg(unregKey) }
}

func (i *inMemory) unreg(keys ...TopicSubscriber) {
	i.mx.Lock()
	for _, k := range keys {
		delete(i.subscribers[k.Topic], k.SubscriberKey)
		// delete map if empty
		if len(i.subscribers[k.Topic]) == 0 {
			delete(i.subscribers, k.Topic)
		}
	}
	i.mx.Unlock()
}
