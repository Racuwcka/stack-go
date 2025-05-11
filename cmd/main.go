package main

import (
	"fmt"
	stackGo "github.com/Racuwcka/stack-go/internal/stack"
	"sync"
)

func main() {
	st := stackGo.NewStack[int]()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			st.Push(i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 11; i < 20; i++ {
			st.Push(i)
		}
	}()

	wg.Wait()

	it := st.NewIterator(stackGo.BottomToTop)
	for it.HasNext() {
		v, _ := it.Next()
		fmt.Printf("%v ", v)
	}
}
