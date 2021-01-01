package defaultdict

import (
	"sync"
)

// Generator defines the function used to generate the default value of the map.
type Generator func() Pointer

// ToPool creates a *sync.Pool from this Generator.
func (g Generator) ToPool() *sync.Pool {
	return &sync.Pool{
		New: g,
	}
}

// SharedPoolMapGenerator creates a Generator that returns a Map using g as the
// Generator.
//
// It's different from just:
//
//     func() defaultdict.Pointer {
//       return defaultdict.New(g)
//     }
//
// That the Map returned by the same SharedPoolMapGenerator shares the same
// underlying pool,
// so it's more memory efficient when used as the second (or third, etc.) layer
// of a Map.
//
// This is an example of running the benchmark test with go1.15.6:
//
//     $ go test -bench . -benchmem
//     goos: linux
//     goarch: amd64
//     pkg: go.yhsif.com/defaultdict
//     BenchmarkSharedPoolMapGenerator/shared-4                    9459            121219 ns/op            1093 B/op          5 allocs/op
//     BenchmarkSharedPoolMapGenerator/naive-4                     9246            123866 ns/op            3289 B/op         14 allocs/op
//     PASS
//     ok      go.yhsif.com/defaultdict        2.322s
func SharedPoolMapGenerator(g Generator) Generator {
	p := g.ToPool()
	return func() Pointer {
		return &defaultdict{
			pool: p,
		}
	}
}
