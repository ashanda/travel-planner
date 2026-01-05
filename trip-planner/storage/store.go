package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type JSONStore[T any] struct {
	Path string
	mu   sync.Mutex
}

func NewJSONStore[T any](path string) *JSONStore[T] {
	return &JSONStore[T]{Path: path}
}

func (s *JSONStore[T]) ReadAll() ([]T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_ = os.MkdirAll(filepath.Dir(s.Path), 0755)

	b, err := os.ReadFile(s.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return []T{}, nil
		}
		return nil, err
	}
	if len(b) == 0 {
		return []T{}, nil
	}

	var items []T
	if err := json.Unmarshal(b, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *JSONStore[T]) WriteAll(items []T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_ = os.MkdirAll(filepath.Dir(s.Path), 0755)

	b, _ := json.MarshalIndent(items, "", "  ")
	return os.WriteFile(s.Path, b, 0644)
}
