package storage

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	opPut    byte = 0x01
	opDelete byte = 0x02
)

type WalStorage[K comparable, V any] struct {
	mu           sync.RWMutex
	store        map[K]V
	snapshotPath string
	walPath      string
	walWriter    *bufio.Writer
	walFile      *os.File
}

func NewWalStorage[K comparable, V any](snapshotPath, walPath string, compactionInterval time.Duration) (*WalStorage[K, V], error) {
	s := &WalStorage[K, V]{
		store:        make(map[K]V),
		snapshotPath: snapshotPath,
		walPath:      walPath,
	}

	if err := s.Load(context.Background()); err != nil {
		return nil, err
	}

	// Start periodic compaction
	go func() {
		ticker := time.NewTicker(compactionInterval)
		for range ticker.C {
			if err := s.Flush(context.Background()); err != nil {
				log.Printf("WAL compaction failed: %v", err)
			}
		}
	}()

	return s, nil
}

func (s *WalStorage[K, V]) Put(ctx context.Context, key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[key] = value
	return s.appendWal(opPut, key, value)
}

func (s *WalStorage[K, V]) Get(ctx context.Context, key K) (V, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.store[key]
	return val, ok, nil
}

func (s *WalStorage[K, V]) Delete(ctx context.Context, key K) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.store, key)
	return s.appendWal(opDelete, key, nil)
}

func (s *WalStorage[K, V]) Scan(ctx context.Context, filter func(K) bool) ([]V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []V
	for k, v := range s.store {
		if filter(k) {
			result = append(result, v)
		}
	}
	return result, nil
}

func (s *WalStorage[K, V]) Keys(ctx context.Context) ([]K, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]K, 0, len(s.store))
	for k := range s.store {
		keys = append(keys, k)
	}
	return keys, nil
}

func (s *WalStorage[K, V]) Flush(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. Write snapshot to tmp file
	tmpPath := s.snapshotPath + ".tmp"
	if err := os.MkdirAll(filepath.Dir(s.snapshotPath), 0755); err != nil {
		return err
	}
	f, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s.store); err != nil {
		return err
	}
	f.Close()

	// 2. Atomic move
	if err := os.Rename(tmpPath, s.snapshotPath); err != nil {
		return err
	}

	// 3. Truncate WAL
	if s.walFile != nil {
		s.walFile.Close()
	}
	if err := os.Truncate(s.walPath, 0); err != nil {
		// If truncate fails, we might still be okay but it's risky
	}
	return s.openWal()
}

func (s *WalStorage[K, V]) Load(ctx context.Context) error {
	// Load snapshot
	if _, err := os.Stat(s.snapshotPath); err == nil {
		f, err := os.Open(s.snapshotPath)
		if err == nil {
			defer f.Close()
			if err := json.NewDecoder(f).Decode(&s.store); err != nil {
				log.Printf("Failed to decode snapshot: %v", err)
			}
		}
	}

	// Replay WAL
	if _, err := os.Stat(s.walPath); err == nil {
		if err := s.replayWal(); err != nil {
			log.Printf("WAL replay error: %v", err)
		}
	}

	return s.openWal()
}

func (s *WalStorage[K, V]) Clear(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store = make(map[K]V)
	if s.walFile != nil {
		s.walFile.Close()
	}
	os.Remove(s.walPath)
	os.Remove(s.snapshotPath)
	return s.openWal()
}

func (s *WalStorage[K, V]) openWal() error {
	if err := os.MkdirAll(filepath.Dir(s.walPath), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(s.walPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	s.walFile = f
	s.walWriter = bufio.NewWriter(f)
	return nil
}

func (s *WalStorage[K, V]) appendWal(op byte, key K, value any) error {
	if s.walWriter == nil {
		return nil
	}

	keyBytes, _ := json.Marshal(key)
	var valBytes []byte
	if op == opPut {
		valBytes, _ = json.Marshal(value)
	}

	s.walWriter.WriteByte(op)
	
	binary.Write(s.walWriter, binary.BigEndian, int32(len(keyBytes)))
	s.walWriter.Write(keyBytes)

	if op == opPut {
		binary.Write(s.walWriter, binary.BigEndian, int32(len(valBytes)))
		s.walWriter.Write(valBytes)
	}

	return s.walWriter.Flush()
}

func (s *WalStorage[K, V]) replayWal() error {
	f, err := os.Open(s.walPath)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	for {
		op, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		var keyLen int32
		if err := binary.Read(reader, binary.BigEndian, &keyLen); err != nil {
			break
		}
		keyBytes := make([]byte, keyLen)
		if _, err := io.ReadFull(reader, keyBytes); err != nil {
			break
		}

		var key K
		json.Unmarshal(keyBytes, &key)

		if op == opPut {
			var valLen int32
			if err := binary.Read(reader, binary.BigEndian, &valLen); err != nil {
				break
			}
			valBytes := make([]byte, valLen)
			if _, err := io.ReadFull(reader, valBytes); err != nil {
				break
			}
			var val V
			json.Unmarshal(valBytes, &val)
			s.store[key] = val
		} else if op == opDelete {
			delete(s.store, key)
		}
	}
	return nil
}
