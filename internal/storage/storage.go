// Package storage provides thread-safe storage for pack sizes.
package storage

import (
	"sync"
)

// Storage interface for pack size persistence
type Storage interface {
	GetPackSizes() []int
	SetPackSizes(sizes []int) error
	AddPackSize(size int) bool
	RemovePackSize(size int) bool
}

// MemoryStorage is a thread-safe in-memory implementation
type MemoryStorage struct {
	mu        sync.RWMutex
	packSizes []int
}

func DefaultPackSizes() []int {
	return []int{250, 500, 1000, 2000, 5000}
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		packSizes: DefaultPackSizes(),
	}
}

func NewMemoryStorageWithSizes(sizes []int) *MemoryStorage {
	s := &MemoryStorage{
		packSizes: make([]int, len(sizes)),
	}
	copy(s.packSizes, sizes)
	return s
}

// GetPackSizes returns a copy to prevent external modifications
func (s *MemoryStorage) GetPackSizes() []int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]int, len(s.packSizes))
	copy(result, s.packSizes)
	return result
}

func (s *MemoryStorage) SetPackSizes(sizes []int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.packSizes = make([]int, len(sizes))
	copy(s.packSizes, sizes)
	return nil
}

func (s *MemoryStorage) AddPackSize(size int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, existing := range s.packSizes {
		if existing == size {
			return false
		}
	}

	s.packSizes = append(s.packSizes, size)
	return true
}

func (s *MemoryStorage) RemovePackSize(size int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, existing := range s.packSizes {
		if existing == size {
			s.packSizes = append(s.packSizes[:i], s.packSizes[i+1:]...)
			return true
		}
	}

	return false
}
