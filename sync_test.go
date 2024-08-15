package defaultdict_test

import (
	"sync/atomic"
	"testing"

	"go.yhsif.com/defaultdict"
)

func incAndCheckAtomicInt64(tb testing.TB, got *int64, want int64, msg string) {
	tb.Helper()

	if v := atomic.AddInt64(got, 1); v != want {
		tb.Errorf("%s: got %d, want %d", msg, v, want)
	}
}

func TestMap(t *testing.T) {
	const key = "foo"

	m := defaultdict.New[string](func() *int64 {
		return new(int64)
	})

	v, ok := m.LoadAndDelete(key)
	if ok {
		t.Error("Expected LoadAndDelete non-exist key to return false, got true")
	}
	incAndCheckAtomicInt64(t, v, 1, "LoadAndDelete non-exist key")

	v, ok = m.Load(key)
	if ok {
		t.Error("Expected Load new key to return false, got true")
	}
	incAndCheckAtomicInt64(t, v, 1, "Load new key")
	v, ok = m.Load(key)
	if !ok {
		t.Error("Expected Load same key to return true, got false")
	}
	incAndCheckAtomicInt64(t, v, 2, "Load same key")
	incAndCheckAtomicInt64(t, m.Get(key), 3, "Get same key")

	m.Delete(key)
	incAndCheckAtomicInt64(t, m.Get(key), 1, "Get deleted key")

	got := make(map[string]int64)
	m.Range(func(key string, value *int64) bool {
		v := atomic.LoadInt64(value)
		got[key] = v
		return true
	})
	if size := len(got); size != 1 {
		t.Errorf("Want 1 element in the map, got %v", got)
	}
}

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
