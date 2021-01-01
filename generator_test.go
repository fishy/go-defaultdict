package defaultdict_test

import (
	"runtime"
	"testing"

	"go.yhsif.com/defaultdict"
)

func BenchmarkSharedPoolMapGenerator(b *testing.B) {
	const n = 5
	g := func() defaultdict.Pointer {
		return new(int64)
	}
	for _, c := range []struct {
		label string
		m     defaultdict.Map
	}{
		{
			label: "shared",
			m:     defaultdict.New(defaultdict.SharedPoolMapGenerator(g)),
		},
		{
			label: "naive",
			m: defaultdict.New(func() defaultdict.Pointer {
				return defaultdict.New(g)
			}),
		},
	} {
		b.Run(c.label, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := 0; j < n; j++ {
					for k := 0; k < n; k++ {
						c.m.Get(j).(defaultdict.Map).Get(k)
					}
				}
				runtime.GC()
			}
		})
	}
}
