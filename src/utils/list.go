package utils

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
