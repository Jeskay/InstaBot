package utils

import (
	"math/rand"
	"time"
)

type IterableList[T any] struct {
	list    []T
	current int
}

func NewIterableList[V any](items []V) *IterableList[V] {
	return &IterableList[V]{
		list:    items,
		current: 0,
	}
}

func (l *IterableList[T]) Next() bool {
	if l.current < len(l.list) {
		l.current++
		return true
	}
	return false
}

func (l *IterableList[T]) Finished() bool {
	return l.current == len(l.list)
}

func (l *IterableList[T]) Current(reverse bool) T {
	if reverse {
		return l.list[len(l.list)-1-l.current]
	}
	return l.list[l.current]
}

type Spread struct {
	Min    int
	Max    int
	source rand.Source
}

func NewSpread(min int, max int) *Spread {
	src := rand.NewSource(time.Now().UnixNano())
	return &Spread{
		Min:    min,
		Max:    max,
		source: src,
	}
}

func (s *Spread) Rand() int {
	return rand.New(s.source).Intn(s.Max-s.Min) + s.Min
}
