package defaultdict_test

import (
	"runtime"
	"testing"

	"go.yhsif.com/defaultdict"
)

func BenchmarkSharedPoolMapGenerator(b *testing.B) {
	const n = 5
	g := func() *int64 {
		return new(int64)
	}
	for _, c := range []struct {
		label string
		g     defaultdict.Generator[defaultdict.Map[int, *int64]]
	}{
		{
			label: "shared",
			g:     defaultdict.SharedPoolMapGenerator[int](g),
		},
		{
			label: "naive",
			g: func() defaultdict.Map[int, *int64] {
				return defaultdict.New[int](g)
			},
		},
	} {
		b.Run(c.label, func(b *testing.B) {
			b.ReportAllocs()

			m := defaultdict.New[int](c.g)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				for j := range n {
					for k := range n {
						m.Get(j).Get(k)
					}
				}
				runtime.GC()
			}
		})
	}
}
