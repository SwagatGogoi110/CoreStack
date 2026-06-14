package mq

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/aws/mq/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type MqService struct {
	brokerStore storage.Backend[string, *model.BrokerSummary]
}

func NewMqService(factory *storage.Factory) (*MqService, error) {
	brokerStore, _ := storage.CreateAccountAware[*model.BrokerSummary](factory, "mq", "mq-brokers.json", "wal")

	return &MqService{
		brokerStore: brokerStore,
	}, nil
}

func (s *MqService) CreateBroker(ctx context.Context, name, engineType, deploymentMode string) (*model.BrokerSummary, error) {
	id := "b-" + uuid.New().String()[:10]
	broker := &model.BrokerSummary{
		BrokerId:       id,
		BrokerName:     name,
		BrokerArn:      fmt.Sprintf("arn:aws:mq:us-east-1:000000000000:broker:%s", name),
		BrokerState:    "RUNNING",
		EngineType:     engineType,
		DeploymentMode: deploymentMode,
	}

	if err := s.brokerStore.Put(ctx, id, broker); err != nil {
		return nil, err
	}
	return broker, nil
}

func (s *MqService) DescribeBroker(ctx context.Context, id string) (*model.BrokerSummary, error) {
	broker, ok, err := s.brokerStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("NotFoundException: Broker %s not found", id)
	}
	return broker, nil
}

func (s *MqService) ListBrokers(ctx context.Context) ([]*model.BrokerSummary, error) {
	return s.brokerStore.Scan(ctx, func(k string) bool { return true })
}

func (s *MqService) DeleteBroker(ctx context.Context, id string) error {
	return s.brokerStore.Delete(ctx, id)
}
