package stack

import (
	"fmt"
	"sync"
)

type Stack[T any] struct {
	mu    sync.Mutex
	items []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Push(v T) {
	s.mu.Lock()
	s.items = append(s.items, v)
	s.mu.Unlock()
}

func (s *Stack[T]) Pop() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var zero T
	l := s.sizeLocked()

	if l == 0 {
		return zero, false
	}

	result := s.items[l-1]
	s.items = s.items[:l-1]
	return result, true
}

func (s *Stack[T]) Peek() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var zero T
	l := s.sizeLocked()

	if l == 0 {
		return zero, false
	}
	return s.items[l-1], true
}

func (s *Stack[T]) Clone() *Stack[T] {
	s.mu.Lock()
	defer s.mu.Unlock()

	items := make([]T, s.sizeLocked())
	copy(items, s.items)
	return &Stack[T]{items: items}
}

func (s *Stack[T]) DeepClone(cloneFn func(T) T) *Stack[T] {
	s.mu.Lock()
	defer s.mu.Unlock()

	newItems := make([]T, s.sizeLocked())
	for i, v := range s.items {
		newItems[i] = cloneFn(v)
	}
	return &Stack[T]{items: newItems}
}

func (s *Stack[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = nil
}

func (s *Stack[T]) Size() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.sizeLocked()
}

func (s *Stack[T]) sizeLocked() int {
	return len(s.items)
}

func (s *Stack[T]) String() string {
	return fmt.Sprintf("%v", s.items)
}
