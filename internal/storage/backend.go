package storage

import (
	"context"
)

// Backend defines the generic storage interface for AWS services.
type Backend[K comparable, V any] interface {
	Put(ctx context.Context, key K, value V) error
	Get(ctx context.Context, key K) (V, bool, error)
	Delete(ctx context.Context, key K) error
	Scan(ctx context.Context, filter func(K) bool) ([]V, error)
	Keys(ctx context.Context) ([]K, error)
	Flush(ctx context.Context) error
	Load(ctx context.Context) error
	Clear(ctx context.Context) error
}
