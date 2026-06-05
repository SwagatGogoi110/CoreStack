package storage

import (
	"context"
	"fmt"
)

// AccountAwareBackend wraps a Backend and prefixes keys with the account ID from the context.
type AccountAwareBackend[V any] struct {
	inner          Backend[string, V]
	defaultAccount string
	getAccount     func(ctx context.Context) string
}

func NewAccountAwareBackend[V any](inner Backend[string, V], defaultAccount string, getAccount func(ctx context.Context) string) *AccountAwareBackend[V] {
	return &AccountAwareBackend[V]{
		inner:          inner,
		defaultAccount: defaultAccount,
		getAccount:     getAccount,
	}
}

func (b *AccountAwareBackend[V]) account(ctx context.Context) string {
	acc := b.getAccount(ctx)
	if acc == "" {
		return b.defaultAccount
	}
	return acc
}

func (b *AccountAwareBackend[V]) prefixKey(ctx context.Context, key string) string {
	return fmt.Sprintf("%s:%s", b.account(ctx), key)
}

func (b *AccountAwareBackend[V]) Put(ctx context.Context, key string, value V) error {
	return b.inner.Put(ctx, b.prefixKey(ctx, key), value)
}

func (b *AccountAwareBackend[V]) Get(ctx context.Context, key string) (V, bool, error) {
	return b.inner.Get(ctx, b.prefixKey(ctx, key))
}

func (b *AccountAwareBackend[V]) Delete(ctx context.Context, key string) error {
	return b.inner.Delete(ctx, b.prefixKey(ctx, key))
}

func (b *AccountAwareBackend[V]) Scan(ctx context.Context, filter func(string) bool) ([]V, error) {
	acc := b.account(ctx)
	prefix := acc + ":"
	return b.inner.Scan(ctx, func(k string) bool {
		if len(k) < len(prefix) || k[:len(prefix)] != prefix {
			return false
		}
		return filter(k[len(prefix):])
	})
}

func (b *AccountAwareBackend[V]) Keys(ctx context.Context) ([]string, error) {
	acc := b.account(ctx)
	prefix := acc + ":"
	allKeys, err := b.inner.Keys(ctx)
	if err != nil {
		return nil, err
	}
	var filtered []string
	for _, k := range allKeys {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			filtered = append(filtered, k[len(prefix):])
		}
	}
	return filtered, nil
}

func (b *AccountAwareBackend[V]) Flush(ctx context.Context) error { return b.inner.Flush(ctx) }
func (b *AccountAwareBackend[V]) Load(ctx context.Context) error  { return b.inner.Load(ctx) }
func (b *AccountAwareBackend[V]) Clear(ctx context.Context) error { return b.inner.Clear(ctx) }
