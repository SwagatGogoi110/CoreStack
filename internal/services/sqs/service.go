package sqs

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/sqs/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type SqsService struct {
	queueStore storage.Backend[string, *model.Queue]
	messages   map[string][]*model.Message
	mu         sync.RWMutex
	baseUrl    string
}

func NewSqsService(factory *storage.Factory, baseUrl string) (*SqsService, error) {
	queueStore, _ := storage.CreateAccountAware[*model.Queue](factory, "sqs", "sqs-queues.json", "wal")

	return &SqsService{
		queueStore: queueStore,
		messages:   make(map[string][]*model.Message),
		baseUrl:    baseUrl,
	}, nil
}

func (s *SqsService) CreateQueue(ctx context.Context, name string, attributes map[string]string) (*model.Queue, error) {
	if _, ok, _ := s.queueStore.Get(ctx, name); ok {
		// In a real SQS, if it exists and attributes match, it's fine.
		// For simplicity, let's just return the existing one.
	}

	queueURL := fmt.Sprintf("%s/000000000000/%s", s.baseUrl, name)
	queue := &model.Queue{
		QueueName:             name,
		QueueURL:              queueURL,
		AccountID:             "000000000000",
		Attributes:            attributes,
		CreatedTimestamp:      time.Now(),
		LastModifiedTimestamp: time.Now(),
	}

	if err := s.queueStore.Put(ctx, name, queue); err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.messages[name] = make([]*model.Message, 0)
	s.mu.Unlock()

	log.Printf("Created SQS queue: %s", name)
	return queue, nil
}

func (s *SqsService) SendMessage(ctx context.Context, queueName, body string) (*model.Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	queueMsgs, ok := s.messages[queueName]
	if !ok {
		return nil, fmt.Errorf("AWS.SimpleQueueService.NonExistentQueue: The specified queue does not exist")
	}

	msg := &model.Message{
		MessageID:     uuid.New().String(),
		Body:          body,
		SentTimestamp: time.Now(),
		MD5OfBody:     fmt.Sprintf("%x", md5.Sum([]byte(body))),
		VisibleAt:     time.Now(),
	}

	s.messages[queueName] = append(queueMsgs, msg)
	return msg, nil
}

func (s *SqsService) ReceiveMessage(ctx context.Context, queueName string, maxMessages int) ([]*model.Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	queueMsgs, ok := s.messages[queueName]
	if !ok {
		return nil, fmt.Errorf("AWS.SimpleQueueService.NonExistentQueue: The specified queue does not exist")
	}

	var received []*model.Message
	now := time.Now()
	
	// Very basic visibility logic
	for _, m := range queueMsgs {
		if m.VisibleAt.Before(now) {
			m.ReceiveCount++
			m.ReceiptHandle = uuid.New().String()
			m.VisibleAt = now.Add(30 * time.Second) // Default visibility timeout
			received = append(received, m)
			if len(received) >= maxMessages {
				break
			}
		}
	}

	return received, nil
}

func (s *SqsService) GetQueueUrl(ctx context.Context, name string) (string, error) {
	q, ok, _ := s.queueStore.Get(ctx, name)
	if !ok {
		return "", fmt.Errorf("AWS.SimpleQueueService.NonExistentQueue: The specified queue does not exist")
	}
	return q.QueueURL, nil
}

func (s *SqsService) ListQueues(ctx context.Context, prefix string) ([]string, error) {
	queues, _ := s.queueStore.Scan(ctx, func(k string) bool {
		return strings.HasPrefix(k, prefix)
	})
	urls := make([]string, 0, len(queues))
	for _, q := range queues {
		urls = append(urls, q.QueueURL)
	}
	return urls, nil
}
