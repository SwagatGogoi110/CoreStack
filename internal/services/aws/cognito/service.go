package cognito

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/aws/cognito/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type CognitoService struct {
	poolStore   storage.Backend[string, *model.UserPool]
	clientStore storage.Backend[string, *model.UserPoolClient]
	userStore   storage.Backend[string, *model.CognitoUser]
	groupStore  storage.Backend[string, *model.CognitoGroup]
}

func NewCognitoService(factory *storage.Factory) (*CognitoService, error) {
	poolStore, _ := storage.CreateAccountAware[*model.UserPool](factory, "cognito", "cognito-pools.json", "wal")
	clientStore, _ := storage.CreateAccountAware[*model.UserPoolClient](factory, "cognito", "cognito-clients.json", "wal")
	userStore, _ := storage.CreateAccountAware[*model.CognitoUser](factory, "cognito", "cognito-users.json", "wal")
	groupStore, _ := storage.CreateAccountAware[*model.CognitoGroup](factory, "cognito", "cognito-groups.json", "wal")

	return &CognitoService{
		poolStore:   poolStore,
		clientStore: clientStore,
		userStore:   userStore,
		groupStore:  groupStore,
	}, nil
}

// User Pools

func (s *CognitoService) CreateUserPool(ctx context.Context, name string) (*model.UserPool, error) {
	id := fmt.Sprintf("us-east-1_%s", strings.ReplaceAll(uuid.New().String(), "-", "")[:9])
	
	pool := &model.UserPool{
		ID:               id,
		Name:             name,
		Arn:              fmt.Sprintf("arn:aws:cognito-idp:us-east-1:000000000000:userpool/%s", id),
		Status:           "Enabled",
		CreationDate:     time.Now().Unix(),
		LastModifiedDate: time.Now().Unix(),
	}

	if err := s.poolStore.Put(ctx, id, pool); err != nil {
		return nil, err
	}

	log.Printf("Created User Pool: %s", id)
	return pool, nil
}

func (s *CognitoService) DescribeUserPool(ctx context.Context, id string) (*model.UserPool, error) {
	pool, ok, err := s.poolStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("ResourceNotFoundException: User pool not found")
	}
	return pool, nil
}

func (s *CognitoService) ListUserPools(ctx context.Context) ([]*model.UserPool, error) {
	return s.poolStore.Scan(ctx, func(k string) bool { return true })
}

// User Pool Clients

func (s *CognitoService) CreateUserPoolClient(ctx context.Context, userPoolID, clientName string) (*model.UserPoolClient, error) {
	clientID := strings.ReplaceAll(uuid.New().String(), "-", "")[:26]
	
	client := &model.UserPoolClient{
		ClientID:         clientID,
		UserPoolID:       userPoolID,
		ClientName:       clientName,
		CreationDate:     time.Now().Unix(),
		LastModifiedDate: time.Now().Unix(),
	}

	if err := s.clientStore.Put(ctx, clientID, client); err != nil {
		return nil, err
	}

	log.Printf("Created User Pool Client: %s", clientID)
	return client, nil
}

// Users

func (s *CognitoService) AdminCreateUser(ctx context.Context, userPoolID, username string, attributes map[string]string) (*model.CognitoUser, error) {
	key := s.userKey(userPoolID, username)
	
	user := &model.CognitoUser{
		Username:         username,
		UserPoolID:       userPoolID,
		UserStatus:       "CONFIRMED",
		Enabled:          true,
		Attributes:       attributes,
		CreationDate:     time.Now().Unix(),
		LastModifiedDate: time.Now().Unix(),
	}

	if user.Attributes == nil {
		user.Attributes = make(map[string]string)
	}
	if _, ok := user.Attributes["sub"]; !ok {
		user.Attributes["sub"] = uuid.New().String()
	}

	if err := s.userStore.Put(ctx, key, user); err != nil {
		return nil, err
	}

	log.Printf("Created user %s in pool %s", username, userPoolID)
	return user, nil
}

func (s *CognitoService) userKey(userPoolID, username string) string {
	return fmt.Sprintf("%s:%s", userPoolID, username)
}
