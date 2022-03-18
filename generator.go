package defaultdict

import (
	"sync"
)

// Generator defines the function used to generate the default value of the map.
type Generator[T Pointer] func() T

// ToPool creates a *sync.Pool from this Generator.
func (g Generator[T]) ToPool() *sync.Pool {
	return &sync.Pool{
		New: func() any {
			return g()
		},
	}
}

// SharedPoolMapGenerator creates a Generator that returns a Map using g as the
// Generator.
//
// It's different from just:
//
//     func generator[K comparable, V defaultdict.Pointer](g defaultdict.Generator[V]) Map[K, V] {
//       return defaultdict.New[K](g)
//     }
//
// That the Map returned by the same SharedPoolMapGenerator shares the same
// underlying pool,
// so it's more memory efficient when used as the second (or third, etc.) layer
// of a Map.
//
// This is an example of running the benchmark test with go1.15.6:
//
//     $ go test -bench .
//     goos: linux
//     goarch: amd64
//     pkg: go.yhsif.com/defaultdict
//     cpu: Intel(R) Core(TM) i5-7260U CPU @ 2.20GHz
//     BenchmarkSharedPoolMapGenerator/shared-4                    9691            117636 ns/op            1093 B/op          5 allocs/op
//     BenchmarkSharedPoolMapGenerator/naive-4                     9882            121344 ns/op            3305 B/op         15 allocs/op
//     PASS
//     ok      go.yhsif.com/defaultdict        2.368s
func SharedPoolMapGenerator[K comparable, V Pointer](g Generator[V]) Generator[Map[K, V]] {
	p := g.ToPool()
	return func() Map[K, V] {
		return &defaultdict[K, V]{
			pool: p,
		}
	}
}
