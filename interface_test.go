package defaultdict_test

import (
	"fmt"
	"sync"

	"go.yhsif.com/defaultdict"
)

type MySliceValue struct {
	lock  sync.Mutex
	slice []int
}

func (m *MySliceValue) Append(i int) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.slice = append(m.slice, i)
}

func (m *MySliceValue) Len() int {
	m.lock.Lock()
	defer m.lock.Unlock()
	return len(m.slice)
}

func MySliceValueGenerator() *MySliceValue {
	return new(MySliceValue)
}

func ExamplePointer_slice() {
	const key = "foo"
	m := defaultdict.New[string](MySliceValueGenerator)
	m.Get(key).Append(1)
	m.Get(key).Append(2)
	fmt.Println(m.Get(key).Len())   // 2
	fmt.Println(m.Get("bar").Len()) // 0

	// Output:
	// 2
	// 0
}
