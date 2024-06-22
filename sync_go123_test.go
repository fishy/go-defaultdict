//go:build go1.23

package defaultdict_test

import (
	"sync/atomic"
	"testing"

	"go.yhsif.com/defaultdict"
)

func TestMapRangeOverFunc(t *testing.T) {
	const (
		key1 = "foo"
		key2 = "bar"

		value1 = 2
		value2 = 3
	)

	m := defaultdict.New[string](func() *atomic.Int64 {
		return new(atomic.Int64)
	})

	m.Get(key1).Store(value1)
	m.Get(key2).Store(value2)

	gotMap := make(map[string]int64)
	for k, v := range m.All() {
		gotMap[k] = v.Load()
	}
	if size := len(gotMap); size != 2 {
		t.Errorf("Want 1 element in the map, got %v", gotMap)
	}
	if got, want := gotMap[key1], value1; got != int64(want) {
		t.Errorf("m[%q] got %v want %v", key1, got, want)
	}
	if got, want := gotMap[key2], value2; got != int64(want) {
		t.Errorf("m[%q] got %v want %v", key2, got, want)
	}
}
