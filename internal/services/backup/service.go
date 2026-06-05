package backup

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/backup/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type BackupService struct {
	vaultStore storage.Backend[string, *model.BackupVault]
}

func NewBackupService(factory *storage.Factory) (*BackupService, error) {
	vaultStore, _ := storage.CreateAccountAware[*model.BackupVault](factory, "backup", "backup-vaults.json", "wal")

	return &BackupService{
		vaultStore: vaultStore,
	}, nil
}

func (s *BackupService) CreateBackupVault(ctx context.Context, name string) (*model.BackupVault, error) {
	if _, ok, _ := s.vaultStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("AlreadyExistsException: Backup vault already exists")
	}

	vault := &model.BackupVault{
		BackupVaultName: name,
		BackupVaultArn:  fmt.Sprintf("arn:aws:backup:us-east-1:000000000000:backup-vault:%s", name),
		CreationDate:    time.Now().Unix(),
	}

	if err := s.vaultStore.Put(ctx, name, vault); err != nil {
		return nil, err
	}

	log.Printf("Created Backup vault: %s", name)
	return vault, nil
}

func (s *BackupService) DescribeBackupVault(ctx context.Context, name string) (*model.BackupVault, error) {
	vault, ok, err := s.vaultStore.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("ResourceNotFoundException: Backup vault not found")
	}
	return vault, nil
}

func (s *BackupService) ListBackupVaults(ctx context.Context) ([]*model.BackupVault, error) {
	return s.vaultStore.Scan(ctx, func(k string) bool { return true })
}
