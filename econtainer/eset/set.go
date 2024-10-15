package eset

import (
	"cmp"
	"sync"
)

type Set[T cmp.Ordered] interface {
	Add(items ...T)
	Remove(items ...T)
	Has(items ...T) bool
}

type set[T cmp.Ordered] struct {
	data map[T]struct{}
	mu   sync.Mutex
}

func (s *set[T]) Add(items ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range items {
		s.data[item] = struct{}{}
	}
}

func (s *set[T]) Remove(items ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range items {
		delete(s.data, item)
	}
}

func (s *set[T]) Has(items ...T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range items {
		if _, ok := s.data[item]; ok {
			return true
		}
	}

	return false
}

func New[T cmp.Ordered](items ...T) Set[T] {
	res := set[T]{
		data: make(map[T]struct{}, len(items)),
	}
	res.Add(items...)
	return &res
}
