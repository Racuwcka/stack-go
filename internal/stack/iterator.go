package stack

type Direction int

const (
	TopToBottom Direction = -1
	BottomToTop Direction = 1
)

type Iterator[T any] struct {
	stack     *Stack[T]
	direction Direction
	index     int
	end       int
}

func (s *Stack[T]) NewIterator(d Direction) *Iterator[T] {
	s.mu.Lock()
	defer s.mu.Unlock()

	copied := make([]T, len(s.items))
	copy(copied, s.items)

	var index, end int
	if d == TopToBottom {
		index = s.Size() - 1
	} else {
		end = s.Size() - 1
	}
	return &Iterator[T]{
		stack:     s,
		direction: d,
		index:     index,
		end:       end,
	}
}

func (it *Iterator[T]) HasNext() bool {
	if it.direction == TopToBottom {
		return it.index >= it.end
	}
	return it.index <= it.end
}

func (it *Iterator[T]) Next() (T, bool) {
	var zero T
	if !it.HasNext() {
		return zero, false
	}

	val := it.stack.items[it.index]

	if it.direction == TopToBottom {
		it.index--
	} else {
		it.index++
	}

	return val, true
}
