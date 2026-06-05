package configservice

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/configservice/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type ConfigService struct {
	ruleStore storage.Backend[string, *model.ConfigRule]
	mu        sync.RWMutex
}

func NewConfigService(factory *storage.Factory) (*ConfigService, error) {
	ruleStore, _ := storage.CreateAccountAware[*model.ConfigRule](factory, "configservice", "config-rules.json", "wal")

	return &ConfigService{
		ruleStore: ruleStore,
	}, nil
}

func (s *ConfigService) PutConfigRule(ctx context.Context, rule *model.ConfigRule) (*model.ConfigRule, error) {
	id := uuid.New().String()
	rule.ConfigRuleID = id
	rule.ConfigRuleArn = fmt.Sprintf("arn:aws:config:us-east-1:000000000000:config-rule/%s", id)
	rule.ConfigRuleState = "ACTIVE"

	if err := s.ruleStore.Put(ctx, rule.ConfigRuleName, rule); err != nil {
		return nil, err
	}

	log.Printf("PutConfigRule: %s", rule.ConfigRuleName)
	return rule, nil
}

func (s *ConfigService) DescribeConfigRules(ctx context.Context, names []string) ([]*model.ConfigRule, error) {
	return s.ruleStore.Scan(ctx, func(k string) bool {
		if len(names) == 0 {
			return true
		}
		for _, n := range names {
			if n == k {
				return true
			}
		}
		return false
	})
}
