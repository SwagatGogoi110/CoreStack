package storage

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type Factory struct {
	persistentPath      string
	defaultAccountID    string
	compactionInterval  time.Duration
}

func NewFactory(persistentPath, defaultAccountID string, compactionInterval time.Duration) *Factory {
	return &Factory{
		persistentPath:     persistentPath,
		defaultAccountID:   defaultAccountID,
		compactionInterval: compactionInterval,
	}
}

func (f *Factory) Root() string {
	return f.persistentPath
}

func CreateAccountAware[V any](f *Factory, serviceName, fileName, mode string) (Backend[string, V], error) {
	var inner Backend[string, V]
	var err error

	switch mode {
	case "memory":
		inner = NewInMemoryStorage[string, V]()
	case "wal":
		snapshotPath := filepath.Join(f.persistentPath, strings.Replace(fileName, ".json", "-snapshot.json", 1))
		walPath := filepath.Join(f.persistentPath, strings.Replace(fileName, ".json", ".wal", 1))
		inner, err = NewWalStorage[string, V](snapshotPath, walPath, f.compactionInterval)
	default:
		inner = NewInMemoryStorage[string, V]()
	}


	if err != nil {
		return nil, err
	}

	return NewAccountAwareBackend[V](inner, f.defaultAccountID, func(ctx context.Context) string {
		rc := common.FromContext(ctx)
		return rc.AccountID
	}), nil
}
