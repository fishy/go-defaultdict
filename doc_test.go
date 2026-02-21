package defaultdict_test

import (
	"fmt"
	"sync"
	"sync/atomic"

	"go.yhsif.com/defaultdict"
)

// This example demonstrates how to use defaultdict to implement a thread-safe
// counter.
func Example() {
	generator := func() *int64 {
		// Just create a new *int64 so it can be used as atomic int64.
		return new(int64)
	}
	m := defaultdict.New[string](generator)
	var wg sync.WaitGroup
	for i := range 10 {
		for j := 0; j < i; j++ {
			wg.Add(1)
			go func(key string) {
				defer wg.Done()
				atomic.AddInt64(m.Get(key), 1)
			}(fmt.Sprintf("key-%d", i))
		}
	}

	wg.Wait()
	m.Range(func(key string, value *int64) bool {
		fmt.Printf("Key %v was added %d times\n", key, atomic.LoadInt64(value))
		return true
	})

	// Unordered Output:
	//
	// Key key-1 was added 1 times
	// Key key-2 was added 2 times
	// Key key-3 was added 3 times
	// Key key-4 was added 4 times
	// Key key-5 was added 5 times
	// Key key-6 was added 6 times
	// Key key-7 was added 7 times
	// Key key-8 was added 8 times
	// Key key-9 was added 9 times
}

// This example demonstrates how to use SharedPoolGenerator to implement a
// thread-safe counter with 2 layers of keys.
func ExampleSharedPoolMapGenerator() {
	generator := defaultdict.SharedPoolMapGenerator[string](func() *int64 {
		// Just create a new *int64 so it can be used as atomic int64.
		return new(int64)
	})
	m := defaultdict.New[string](generator)
	var wg sync.WaitGroup
	for i := 1; i < 4; i++ {
		for j := 1; j < 4; j++ {
			for k := 0; k < i*j; k++ {
				wg.Add(1)
				go func(key1, key2 string) {
					defer wg.Done()
					atomic.AddInt64(m.Get(key1).Get(key2), 1)
				}(fmt.Sprintf("key1-%d", i), fmt.Sprintf("key2-%d", j))
			}
		}
	}

	wg.Wait()
	m.Range(func(key1 string, value1 defaultdict.Map[string, *int64]) bool {
		value1.Range(func(key2 string, value2 *int64) bool {
			fmt.Printf("%v/%v was added %d times\n", key1, key2, atomic.LoadInt64(value2))
			return true
		})
		fmt.Println()
		return true
	})

	// Unordered Output:
	//
	// key1-1/key2-1 was added 1 times
	// key1-1/key2-2 was added 2 times
	// key1-1/key2-3 was added 3 times
	//
	// key1-2/key2-1 was added 2 times
	// key1-2/key2-2 was added 4 times
	// key1-2/key2-3 was added 6 times
	//
	// key1-3/key2-1 was added 3 times
	// key1-3/key2-2 was added 6 times
	// key1-3/key2-3 was added 9 times
}
