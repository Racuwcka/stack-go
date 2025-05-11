package stack

import (
	"fmt"
	"reflect"
	"testing"
)

func TestIterator(t *testing.T) {
	var tests = []struct {
		stack     []int
		direction Direction
		want      []int
	}{
		{[]int{1, 2, 3}, TopToBottom, []int{3, 2, 1}},
		{[]int{1, 2, 3}, BottomToTop, []int{1, 2, 3}},
	}

	for _, test := range tests {
		name := fmt.Sprintf("case(%v)", test.stack)
		t.Run(name, func(t *testing.T) {
			stack := NewStack[int]()
			for _, v := range test.stack {
				stack.Push(v)
			}
			it := stack.NewIterator(test.direction)
			var got []int
			for it.HasNext() {
				val, _ := it.Next()
				got = append(got, val)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}

func TestPeek(t *testing.T) {
	type PeekResult struct {
		value int
		ok    bool
	}

	var tests = []struct {
		stack []int
		want  PeekResult
	}{
		{[]int{1, 2, 3}, PeekResult{value: 3, ok: true}},
		{[]int{1}, PeekResult{value: 1, ok: true}},
		{[]int{}, PeekResult{value: 0, ok: false}},
	}

	for _, test := range tests {
		name := fmt.Sprintf("case(%v)", test.stack)
		t.Run(name, func(t *testing.T) {
			stack := NewStack[int]()
			for _, v := range test.stack {
				stack.Push(v)
			}
			val, ok := stack.Peek()
			if val != test.want.value || ok != test.want.ok {
				t.Errorf("got %v %v, want %v %v", val, ok, test.want.value, test.want.ok)
			}
		})
	}
}

func TestClone_Slice(t *testing.T) {
	var tests = []struct {
		stack []int
		want  []int
	}{
		{[]int{1, 2, 3}, []int{1, 2, 3, 4}},
	}

	for _, test := range tests {
		name := fmt.Sprintf("case(%v)", test.stack)
		t.Run(name, func(t *testing.T) {
			stack := NewStack[int]()
			for _, v := range test.stack {
				stack.Push(v)
			}
			newStack := stack.Clone()
			newStack.Push(4)
			itStack := stack.NewIterator(BottomToTop)
			var got1 []int
			for itStack.HasNext() {
				val, _ := itStack.Next()
				got1 = append(got1, val)
			}
			if !reflect.DeepEqual(got1, test.stack) {
				t.Errorf("got %v, want %v", got1, test.stack)
			}

			itNewStack := newStack.NewIterator(BottomToTop)
			var got2 []int
			for itNewStack.HasNext() {
				val, _ := itNewStack.Next()
				got2 = append(got2, val)
			}
			if !reflect.DeepEqual(got2, test.want) {
				t.Errorf("got %v, want %v", got2, test.want)
			}
		})
	}
}

type UserPtr struct {
	Name *string
	Data []int
}

// cloneUser создаёт полностью независимую копию UserPtr:
// 1) новый string под капотом,
// 2) новый срез Data.
func cloneUser(u UserPtr) UserPtr {
	// скопируем строку
	nameCopy := *u.Name
	// скопируем срез
	dataCopy := make([]int, len(u.Data))
	copy(dataCopy, u.Data)

	return UserPtr{
		Name: &nameCopy,
		Data: dataCopy,
	}
}

func TestDeepClone(t *testing.T) {
	// исходный элемент
	origName := "Alice"
	orig := UserPtr{
		Name: &origName,
		Data: []int{1, 2, 3},
	}

	// положим в стек
	st := NewStack[UserPtr]()
	st.Push(orig)

	// сделаем deep-clone
	cloned := st.DeepClone(cloneUser)

	// модифицируем клон
	// (указываем явно клонированные поля — они должны быть независимы)
	*cloned.items[0].Name = "Bob"
	cloned.items[0].Data[0] = 99

	// теперь проверим оригинал
	it := st.NewIterator(BottomToTop)
	gotOrig, _ := it.Next() // (вернёт единственный элемент)

	// имя должно остаться "Alice"
	if *gotOrig.Name != "Alice" {
		t.Errorf("original Name was modified: got %q, want %q", *gotOrig.Name, "Alice")
	}
	// Data[0] должно остаться 1
	if gotOrig.Data[0] != 1 {
		t.Errorf("original Data was modified: got %d, want %d", gotOrig.Data[0], 1)
	}

	// и проверим клон на то, что изменения применились
	it2 := cloned.NewIterator(BottomToTop)
	gotClone, _ := it2.Next() // (вернёт единственный элемент)
	if *gotClone.Name != "Bob" {
		t.Errorf("clone Name mismatch: got %q, want %q", *gotClone.Name, "Bob")
	}
	if gotClone.Data[0] != 99 {
		t.Errorf("clone Data mismatch: got %d, want %d", gotClone.Data[0], 99)
	}
}
