package events

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"

	"getsturdy.com/api/pkg/users"
)

type EventReader interface {
	// Introduce event type filtering in the call to SubscribeUser(uesrID string, cb CallbackFunc, eventTypes ...EventType) CancelFunc
	// Introduce reference filtering in the call to SubscribeUser(uesrID string, cb CallbackFunc, map[EventType][]string) CancelFunc
	SubscribeUser(userID users.ID, cb CallbackFunc) CancelFunc
}

type eventWriter interface {
	UserEvent(userID users.ID, eventType EventType, reference string)
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
	OrganizationUpdated
)

var eventTypeString = map[EventType]string{
	CodebaseEvent:                  "CodebaseEvent",
	CodebaseUpdated:                "CodebaseUpdated",
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
	OrganizationUpdated:            "OrganizationUpdated",
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
	q           chan payload

	logger *zap.Logger
}

func NewInMemory(logger *zap.Logger) EventReadWriter {
	m := &inMemory{
		subscribers: make(map[Topic]map[string]CallbackFunc),
		q:           make(chan payload, 1024),
		logger:      logger.Named("EventReadWriter"),
	}

	go m.work()

	return m
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

type payload struct {
	topic     Topic
	eventType EventType
	reference string
}

func (i *inMemory) event(topic Topic, eventType EventType, reference string) {
	i.q <- payload{topic, eventType, reference}
}

var ErrClientDisconnected = errors.New("client disconnected")

func (i *inMemory) work() {
	for event := range i.q {
		sentEventCounterMetric.WithLabelValues(eventTypeString[event.eventType]).Inc()

		subs := make(map[TopicSubscriber]CallbackFunc, len(i.subscribers[event.topic]))

		i.mx.RLock()
		// Copy over all subscribers
		for subId, cb := range i.subscribers[event.topic] {
			subs[TopicSubscriber{Topic: event.topic, SubscriberKey: subId}] = cb
		}
		i.mx.RUnlock()

		var unregKeys []TopicSubscriber
		for ts, cb := range subs {
			err := cb(event.eventType, event.reference)
			switch {
			case errors.Is(err, ErrClientDisconnected):
				receivedEventCounterMetric.WithLabelValues(eventTypeString[event.eventType], "no").Inc()
				unregKeys = append(unregKeys, ts)
			case err != nil:
				i.logger.Error("failed to send message", zap.Error(err))
			default:
				receivedEventCounterMetric.WithLabelValues(eventTypeString[event.eventType], "yes").Inc()
			}
		}

		if len(unregKeys) > 0 {
			i.unreg(unregKeys...)
		}
	}
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
