package dynamodb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/dynamodb/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type DynamoDbService struct {
	tableStore storage.Backend[string, *model.TableDefinition]
	itemStore  storage.Backend[string, map[string]json.RawMessage]
	mu         sync.RWMutex
	items      map[string]map[string]json.RawMessage // tableKey -> itemKey -> item
}

func NewDynamoDbService(factory *storage.Factory) (*DynamoDbService, error) {
	tableStore, _ := storage.CreateAccountAware[*model.TableDefinition](factory, "dynamodb", "dynamodb-tables.json", "wal")
	itemStore, _ := storage.CreateAccountAware[map[string]json.RawMessage](factory, "dynamodb", "dynamodb-items.json", "wal")

	s := &DynamoDbService{
		tableStore: tableStore,
		itemStore:  itemStore,
		items:      make(map[string]map[string]json.RawMessage),
	}

	// Load items from store
	ctx := context.Background()
	keys, _ := itemStore.Keys(ctx)
	for _, k := range keys {
		if items, ok, _ := itemStore.Get(ctx, k); ok {
			s.items[k] = items
		}
	}

	return s, nil
}

func (s *DynamoDbService) CreateTable(ctx context.Context, req *model.TableDefinition) (*model.TableDefinition, error) {
	if _, ok, _ := s.tableStore.Get(ctx, req.TableName); ok {
		return nil, fmt.Errorf("ResourceInUseException: Table already exists")
	}

	req.TableStatus = "ACTIVE"
	req.CreationDateTime = time.Now()
	req.TableArn = fmt.Sprintf("arn:aws:dynamodb:us-east-1:000000000000:table/%s", req.TableName)
	
	if err := s.tableStore.Put(ctx, req.TableName, req); err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.items[req.TableName] = make(map[string]json.RawMessage)
	s.mu.Unlock()

	log.Printf("Created DynamoDB table: %s", req.TableName)
	return req, nil
}

func (s *DynamoDbService) DescribeTable(ctx context.Context, name string) (*model.TableDefinition, error) {
	table, ok, err := s.tableStore.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("ResourceNotFoundException: Table not found")
	}
	return table, nil
}

func (s *DynamoDbService) ListTables(ctx context.Context) ([]string, error) {
	return s.tableStore.Keys(ctx)
}

func (s *DynamoDbService) PutItem(ctx context.Context, tableName string, item json.RawMessage) error {
	table, err := s.DescribeTable(ctx, tableName)
	if err != nil {
		return err
	}

	itemKey := s.buildItemKey(table, item)
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	tableItems, ok := s.items[tableName]
	if !ok {
		tableItems = make(map[string]json.RawMessage)
		s.items[tableName] = tableItems
	}
	
	tableItems[itemKey] = item
	
	// Persist items for this table
	return s.itemStore.Put(ctx, tableName, tableItems)
}

func (s *DynamoDbService) GetItem(ctx context.Context, tableName string, key json.RawMessage) (json.RawMessage, error) {
	table, err := s.DescribeTable(ctx, tableName)
	if err != nil {
		return nil, err
	}

	itemKey := s.buildItemKey(table, key)
	
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if tableItems, ok := s.items[tableName]; ok {
		return tableItems[itemKey], nil
	}
	
	return nil, nil
}

func (s *DynamoDbService) buildItemKey(table *model.TableDefinition, item json.RawMessage) string {
	var m map[string]any
	json.Unmarshal(item, &m)
	
	pk := ""
	sk := ""
	
	for _, k := range table.KeySchema {
		val := m[k.AttributeName]
		// DynamoDB JSON format: {"AttributeName": {"S": "value"}}
		// We need to extract the actual value to build a consistent key string
		strVal := fmt.Sprintf("%v", val)
		if k.KeyType == "HASH" {
			pk = strVal
		} else {
			sk = strVal
		}
	}
	
	if sk != "" {
		return fmt.Sprintf("%s#%s", pk, sk)
	}
	return pk
}
