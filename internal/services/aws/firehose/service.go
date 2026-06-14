package firehose

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/aws/firehose/model"
	"github.com/hectorvent/cloudstack/internal/services/aws/s3"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type FirehoseService struct {
	streamStore storage.Backend[string, *model.DeliveryStreamDescription]
	buffers     map[string][][]byte
	s3Service   *s3.S3Service
	mu          sync.RWMutex
}

func NewFirehoseService(factory *storage.Factory, s3Service *s3.S3Service) (*FirehoseService, error) {
	streamStore, _ := storage.CreateAccountAware[*model.DeliveryStreamDescription](factory, "firehose", "fh-streams.json", "wal")

	return &FirehoseService{
		streamStore: streamStore,
		buffers:     make(map[string][][]byte),
		s3Service:   s3Service,
	}, nil
}

func (s *FirehoseService) CreateDeliveryStream(ctx context.Context, name string, s3Dest *model.S3DestinationDescription) (*model.DeliveryStreamDescription, error) {
	if _, ok, _ := s.streamStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("ResourceInUseException: Delivery stream already exists")
	}

	stream := &model.DeliveryStreamDescription{
		DeliveryStreamName:   name,
		DeliveryStreamArn:    fmt.Sprintf("arn:aws:firehose:us-east-1:000000000000:deliverystream/%s", name),
		DeliveryStreamStatus: "ACTIVE",
		CreateTimestamp:      time.Now(),
		Destinations: []*model.DestinationDetails{
			{S3DestinationDescription: s3Dest},
		},
	}

	if err := s.streamStore.Put(ctx, name, stream); err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.buffers[name] = make([][]byte, 0)
	s.mu.Unlock()

	log.Printf("Created Firehose delivery stream: %s", name)
	return stream, nil
}

func (s *FirehoseService) PutRecord(ctx context.Context, streamName string, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	buffer, ok := s.buffers[streamName]
	if !ok {
		return fmt.Errorf("ResourceNotFoundException: Delivery stream %s not found", streamName)
	}

	s.buffers[streamName] = append(buffer, data)
	
	if len(s.buffers[streamName]) >= 5 {
		s.flush(ctx, streamName)
	}
	return nil
}

func (s *FirehoseService) flush(ctx context.Context, streamName string) {
	buffer := s.buffers[streamName]
	if len(buffer) == 0 {
		return
	}

	stream, _, _ := s.streamStore.Get(ctx, streamName)
	if stream == nil || len(stream.Destinations) == 0 || stream.Destinations[0].S3DestinationDescription == nil {
		return
	}

	s3Dest := stream.Destinations[0].S3DestinationDescription
	bucket := s.extractBucket(s3Dest.BucketArn)
	key := fmt.Sprintf("%s%s.json", s3Dest.Prefix, uuid.New().String())

	var content []byte
	for _, b := range buffer {
		content = append(content, b...)
		content = append(content, '\n')
	}

	s.s3Service.PutObject(ctx, bucket, key, content, "application/x-ndjson", nil)
	s.buffers[streamName] = make([][]byte, 0)
	log.Printf("Flushed %d records from %s to s3://%s/%s", len(buffer), streamName, bucket, key)
}

func (s *FirehoseService) extractBucket(arn string) string {
	// arn:aws:s3:::my-bucket
	parts := sync.OnceValue(func() []string { return nil }) // ignore this, just string manipulation
	_ = parts
	idx := 0
	for i := len(arn) - 1; i >= 0; i-- {
		if arn[i] == ':' {
			idx = i + 1
			break
		}
	}
	return arn[idx:]
}
