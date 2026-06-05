package eventbridge

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/eventbridge/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type EventBridgeService struct {
	busStore    storage.Backend[string, *model.EventBus]
	ruleStore   storage.Backend[string, *model.Rule]
	targetStore storage.Backend[string, []*model.Target]
}

func NewEventBridgeService(factory *storage.Factory) (*EventBridgeService, error) {
	busStore, _ := storage.CreateAccountAware[*model.EventBus](factory, "eventbridge", "eventbridge-buses.json", "wal")
	ruleStore, _ := storage.CreateAccountAware[*model.Rule](factory, "eventbridge", "eventbridge-rules.json", "wal")
	targetStore, _ := storage.CreateAccountAware[[]*model.Target](factory, "eventbridge", "eventbridge-targets.json", "wal")

	s := &EventBridgeService{
		busStore:    busStore,
		ruleStore:   ruleStore,
		targetStore: targetStore,
	}

	return s, nil
}

func (s *EventBridgeService) CreateEventBus(ctx context.Context, name string) (*model.EventBus, error) {
	arn := fmt.Sprintf("arn:aws:events:us-east-1:000000000000:event-bus/%s", name)
	bus := &model.EventBus{
		Name:        name,
		Arn:         arn,
		CreatedTime: time.Now(),
	}

	if err := s.busStore.Put(ctx, name, bus); err != nil {
		return nil, err
	}

	log.Printf("Created EventBus: %s", name)
	return bus, nil
}

func (s *EventBridgeService) PutRule(ctx context.Context, name, busName, pattern, schedule string) (*model.Rule, error) {
	arn := fmt.Sprintf("arn:aws:events:us-east-1:000000000000:rule/%s/%s", busName, name)
	rule := &model.Rule{
		Name:               name,
		Arn:                arn,
		EventBusName:       busName,
		EventPattern:       pattern,
		ScheduleExpression: schedule,
		State:              "ENABLED",
		CreatedAt:          time.Now(),
	}

	if err := s.ruleStore.Put(ctx, s.ruleKey(busName, name), rule); err != nil {
		return nil, err
	}

	return rule, nil
}

func (s *EventBridgeService) PutTargets(ctx context.Context, ruleName, busName string, targets []*model.Target) error {
	key := s.ruleKey(busName, ruleName)
	existing, _, _ := s.targetStore.Get(ctx, key)
	
	targetMap := make(map[string]*model.Target)
	for _, t := range existing {
		targetMap[t.ID] = t
	}
	for _, t := range targets {
		targetMap[t.ID] = t
	}

	var updated []*model.Target
	for _, t := range targetMap {
		updated = append(updated, t)
	}

	return s.targetStore.Put(ctx, key, updated)
}

func (s *EventBridgeService) ruleKey(bus, rule string) string {
	if bus == "" {
		bus = "default"
	}
	return fmt.Sprintf("%s:%s", bus, rule)
}
