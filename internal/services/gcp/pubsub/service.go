package pubsub

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/gcp/pubsub/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type PubSubService struct {
	topicStore    storage.Backend[string, *model.Topic]
	subStore      storage.Backend[string, *model.Subscription]
	queues        map[string][]*model.Message
	delivered     map[string]map[string]*model.Message
	messageIdCounter atomic.Uint64
	mu            sync.RWMutex
}

func NewPubSubService(factory *storage.Factory) (*PubSubService, error) {
	topicStore, _ := storage.CreateAccountAware[*model.Topic](factory, "pubsub", "pubsub-topics.json", "wal")
	subStore, _ := storage.CreateAccountAware[*model.Subscription](factory, "pubsub", "pubsub-subscriptions.json", "wal")

	return &PubSubService{
		topicStore: topicStore,
		subStore:   subStore,
		queues:     make(map[string][]*model.Message),
		delivered:  make(map[string]map[string]*model.Message),
	}, nil
}

// Topics

func (s *PubSubService) CreateTopic(ctx context.Context, name string) (*model.Topic, error) {
	if _, ok, _ := s.topicStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("AlreadyExists: Topic %s already exists", name)
	}

	topic := &model.Topic{
		Name: name,
	}

	if err := s.topicStore.Put(ctx, name, topic); err != nil {
		return nil, err
	}

	return topic, nil
}

func (s *PubSubService) ListTopics(ctx context.Context) ([]*model.Topic, error) {
	return s.topicStore.Scan(ctx, func(k string) bool { return true })
}

func (s *PubSubService) DeleteTopic(ctx context.Context, name string) error {
	return s.topicStore.Delete(ctx, name)
}

// Subscriptions

func (s *PubSubService) CreateSubscription(ctx context.Context, name, topic string) (*model.Subscription, error) {
	if _, ok, _ := s.subStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("AlreadyExists: Subscription %s already exists", name)
	}

	sub := &model.Subscription{
		Name:               name,
		Topic:              topic,
		AckDeadlineSeconds: 10,
	}

	if err := s.subStore.Put(ctx, name, sub); err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *PubSubService) ListSubscriptions(ctx context.Context) ([]*model.Subscription, error) {
	return s.subStore.Scan(ctx, func(k string) bool { return true })
}

func (s *PubSubService) DeleteSubscription(ctx context.Context, name string) error {
	return s.subStore.Delete(ctx, name)
}

// Messaging

func (s *PubSubService) Publish(ctx context.Context, topicName string, messages []*model.Message) ([]string, error) {
	if _, ok, _ := s.topicStore.Get(ctx, topicName); !ok {
		return nil, fmt.Errorf("NotFound: Topic %s not found", topicName)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	ids := make([]string, 0, len(messages))
	for _, m := range messages {
		id := fmt.Sprintf("%d", s.messageIdCounter.Add(1))
		m.MessageId = id
		m.PublishTime = time.Now()
		ids = append(ids, id)

		// Fan-out to all subscriptions for this topic
		subs, _ := s.subStore.Scan(ctx, func(k string) bool {
			sub, _, _ := s.subStore.Get(ctx, k)
			return sub.Topic == topicName
		})

		for _, sub := range subs {
			s.queues[sub.Name] = append(s.queues[sub.Name], m)
		}
	}

	return ids, nil
}

func (s *PubSubService) Pull(ctx context.Context, subName string, maxMessages int) ([]*model.ReceivedMessage, error) {
	if _, ok, _ := s.subStore.Get(ctx, subName); !ok {
		return nil, fmt.Errorf("NotFound: Subscription %s not found", subName)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	queue := s.queues[subName]
	if len(queue) == 0 {
		return []*model.ReceivedMessage{}, nil
	}

	count := maxMessages
	if len(queue) < count {
		count = len(queue)
	}

	received := make([]*model.ReceivedMessage, 0, count)
	for i := 0; i < count; i++ {
		msg := queue[i]
		ackId := uuid.New().String()
		received = append(received, &model.ReceivedMessage{
			AckId:   ackId,
			Message: msg,
		})

		// Mark as delivered
		if s.delivered[subName] == nil {
			s.delivered[subName] = make(map[string]*model.Message)
		}
		s.delivered[subName][ackId] = msg
	}

	// Remove from queue
	s.queues[subName] = queue[count:]

	return received, nil
}

func (s *PubSubService) Acknowledge(ctx context.Context, subName string, ackIds []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, id := range ackIds {
		delete(s.delivered[subName], id)
	}
	return nil
}
