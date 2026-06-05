package iam

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/iam/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

const (
	chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type IamService struct {
	users            storage.Backend[string, *model.User]
	groups           storage.Backend[string, *model.Group]
	roles            storage.Backend[string, *model.Role]
	policies         storage.Backend[string, *model.Policy]
	accessKeys       storage.Backend[string, *model.AccessKey]
	instanceProfiles storage.Backend[string, any] // TODO: Add model
	sessions         storage.Backend[string, any] // TODO: Add model
	accountID        string
}

func NewIamService(factory *storage.Factory) (*IamService, error) {
	users, _ := storage.CreateAccountAware[*model.User](factory, "iam", "iam-users.json", "wal")
	groups, _ := storage.CreateAccountAware[*model.Group](factory, "iam", "iam-groups.json", "wal")
	roles, _ := storage.CreateAccountAware[*model.Role](factory, "iam", "iam-roles.json", "wal")
	policies, _ := storage.CreateAccountAware[*model.Policy](factory, "iam", "iam-policies.json", "wal")
	accessKeys, _ := storage.CreateAccountAware[*model.AccessKey](factory, "iam", "iam-access-keys.json", "wal")
	instanceProfiles, _ := storage.CreateAccountAware[any](factory, "iam", "iam-instance-profiles.json", "wal")
	sessions, _ := storage.CreateAccountAware[any](factory, "iam", "iam-sessions.json", "wal")

	s := &IamService{
		users:            users,
		groups:           groups,
		roles:            roles,
		policies:         policies,
		accessKeys:       accessKeys,
		instanceProfiles: instanceProfiles,
		sessions:         sessions,
	}

	return s, nil
}

func (s *IamService) CreateUser(ctx context.Context, userName, path string) (*model.User, error) {
	if _, ok, _ := s.users.Get(ctx, userName); ok {
		return nil, fmt.Errorf("EntityAlreadyExists: User with name %s already exists", userName)
	}

	userID := "AIDA" + s.randomID(16)
	normalizedPath := s.normalizePath(path)
	arn := s.iamArn("user", normalizedPath, userName)

	user := &model.User{
		UserID:     userID,
		UserName:   userName,
		Path:       normalizedPath,
		Arn:        arn,
		CreateDate: time.Now(),
		Tags:       make(map[string]string),
	}

	if err := s.users.Put(ctx, userName, user); err != nil {
		return nil, err
	}

	log.Printf("Created IAM user: %s", userName)
	return user, nil
}

func (s *IamService) GetUser(ctx context.Context, userName string) (*model.User, error) {
	user, ok, err := s.users.Get(ctx, userName)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("NoSuchEntity: The user with name %s cannot be found", userName)
	}
	return user, nil
}

func (s *IamService) DeleteUser(ctx context.Context, userName string) error {
	user, err := s.GetUser(ctx, userName)
	if err != nil {
		return err
	}
	if len(user.AttachedPolicyArns) > 0 {
		return fmt.Errorf("DeleteConflict: Cannot delete entity, must detach all policies first")
	}
	if len(user.GroupNames) > 0 {
		return fmt.Errorf("DeleteConflict: Cannot delete entity, must remove from all groups first")
	}
	return s.users.Delete(ctx, userName)
}

func (s *IamService) CreateAccessKey(ctx context.Context, userName string) (*model.AccessKey, error) {
	_, err := s.GetUser(ctx, userName)
	if err != nil {
		return nil, err
	}

	// TODO: check quota

	keyID := "AKIA" + s.randomID(16)
	secretKey := s.randomSecret(40)
	key := &model.AccessKey{
		AccessKeyID:     keyID,
		SecretAccessKey: secretKey,
		UserName:        userName,
		Status:          "Active",
		CreateDate:      time.Now(),
	}

	if err := s.accessKeys.Put(ctx, keyID, key); err != nil {
		return nil, err
	}

	return key, nil
}

// Helpers

func (s *IamService) randomID(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func (s *IamService) randomSecret(length int) string {
	secretChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	b := make([]byte, length)
	for i := range b {
		b[i] = secretChars[rand.Intn(len(secretChars))]
	}
	return string(b)
}

func (s *IamService) normalizePath(path string) string {
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	return path
}

func (s *IamService) iamArn(resourceType, path, name string) string {
	// TODO: Use account ID from context or config
	return fmt.Sprintf("arn:aws:iam::%s:%s%s%s", "000000000000", resourceType, path, name)
}
