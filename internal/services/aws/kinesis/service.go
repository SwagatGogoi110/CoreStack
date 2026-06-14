package kinesis

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/aws/kinesis/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type KinesisService struct {
	streamStore storage.Backend[string, *model.KinesisStream]
	records     map[string][]*model.KinesisRecord // streamName -> records
	mu          sync.RWMutex
}

func NewKinesisService(factory *storage.Factory) (*KinesisService, error) {
	streamStore, _ := storage.CreateAccountAware[*model.KinesisStream](factory, "kinesis", "kinesis-streams.json", "wal")

	return &KinesisService{
		streamStore: streamStore,
		records:     make(map[string][]*model.KinesisRecord),
	}, nil
}

func (s *KinesisService) CreateStream(ctx context.Context, name string, shardCount int) (*model.KinesisStream, error) {
	arn := fmt.Sprintf("arn:aws:kinesis:us-east-1:000000000000:stream/%s", name)
	
	stream := &model.KinesisStream{
		StreamName:              name,
		StreamArn:               arn,
		StreamStatus:            "ACTIVE",
		RetentionPeriodHours:    24,
		StreamCreationTimestamp: time.Now(),
		Shards:                  make([]*model.KinesisShard, 0),
	}

	for i := 0; i < shardCount; i++ {
		stream.Shards = append(stream.Shards, &model.KinesisShard{
			ShardID: fmt.Sprintf("shardId-%012d", i),
		})
	}

	if err := s.streamStore.Put(ctx, name, stream); err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.records[name] = make([]*model.KinesisRecord, 0)
	s.mu.Unlock()

	log.Printf("Created Kinesis stream: %s with %d shards", name, shardCount)
	return stream, nil
}

func (s *KinesisService) PutRecord(ctx context.Context, name, partitionKey string, data []byte) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	streamRecords, ok := s.records[name]
	if !ok {
		return "", fmt.Errorf("ResourceNotFoundException: Stream %s not found", name)
	}

	seq := fmt.Sprintf("%d", time.Now().UnixNano())
	record := &model.KinesisRecord{
		Data:                        data,
		PartitionKey:                partitionKey,
		SequenceNumber:              seq,
		ApproximateArrivalTimestamp: time.Now(),
	}

	s.records[name] = append(streamRecords, record)
	return seq, nil
}

func (s *KinesisService) ListStreams(ctx context.Context) ([]string, error) {
	streams, err := s.streamStore.Scan(ctx, func(k string) bool { return true })
	if err != nil {
		return nil, err
	}
	var names []string
	for _, st := range streams {
		names = append(names, st.StreamName)
	}
	return names, nil
}
