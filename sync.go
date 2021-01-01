package defaultdict

import (
	"sync"
)

type defaultdict struct {
	m    sync.Map
	pool *sync.Pool
}

// New creates a new DefaultDict.
//
// It pairs a sync.Map with a sync.Pool under the hood.
func New(g Generator) Map {
	return &defaultdict{
		pool: g.ToPool(),
	}
}

func (d *defaultdict) Delete(key Comparable) {
	d.m.Delete(key)
}

func (d *defaultdict) Load(key Comparable) (Pointer, bool) {
	newValue := d.pool.Get()
	value, loaded := d.m.LoadOrStore(key, newValue)
	if loaded {
		d.pool.Put(newValue)
	}
	return value, loaded
}

func (d *defaultdict) Get(key Comparable) Pointer {
	v, _ := d.Load(key)
	return v
}

func (d *defaultdict) LoadAndDelete(key Comparable) (Pointer, bool) {
	value, loaded := d.m.LoadAndDelete(key)
	if !loaded {
		value = d.pool.Get()
	}
	return value, loaded
}

func (d *defaultdict) Range(f func(key Comparable, value Pointer) bool) {
	d.m.Range(f)
}
