package defaultdict

import (
	"iter"
	"sync"
)

type defaultdict[K comparable, V Pointer] struct {
	m    sync.Map
	pool *sync.Pool
}

// New creates a new DefaultDict.
//
// It pairs a sync.Map with a sync.Pool under the hood.
func New[K comparable, V Pointer](g Generator[V]) Map[K, V] {
	return &defaultdict[K, V]{
		pool: g.ToPool(),
	}
}

func (d *defaultdict[K, V]) Delete(key K) {
	d.m.Delete(key)
}

func (d *defaultdict[K, V]) Load(key K) (V, bool) {
	newValue := d.pool.Get()
	value, loaded := d.m.LoadOrStore(key, newValue)
	if loaded {
		d.pool.Put(newValue)
	}
	return value.(V), loaded
}

func (d *defaultdict[K, V]) Get(key K) V {
	v, _ := d.Load(key)
	return v
}

func (d *defaultdict[K, V]) LoadAndDelete(key K) (V, bool) {
	value, loaded := d.m.LoadAndDelete(key)
	if !loaded {
		value = d.pool.Get()
	}
	return value.(V), loaded
}

func (d *defaultdict[K, V]) Range(f func(key K, value V) bool) {
	d.m.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}

func (d *defaultdict[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		d.m.Range(func(k, v any) bool {
			return yield(k.(K), v.(V))
		})
	}
}
