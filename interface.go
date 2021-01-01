package defaultdict

import (
	"sync"
)

// Comparable is the key type of a DefaultDict.
//
// It's used for documentation purpose only.
//
// See https://golang.org/ref/spec#Comparison_operators for more info.
type Comparable = interface{}

// DefaultDict defines a map.
//
// There are a few slight differences in Load and LoadAndDelete comparing to
// sync.Map:
//
// 1. The value type is guaranteed to be the same as the type returned by the
// Generator used to create this DefaultDict, never nil,
// even if this is a new key.
//
// 2. The bool return being false means that the value is directly from the
// Generator.
type DefaultDict interface {
	// Same as in sync.Map.
	Delete(key Comparable)
	Load(key Comparable) (interface{}, bool)
	LoadAndDelete(key Comparable) (interface{}, bool)
	Range(f func(key Comparable, value interface{}) bool)

	// Same as Load, just without the bool return.
	Get(key Comparable) interface{}
}

// Generator defines the function used to generate the default value of the map.
type Generator func() interface{}

// ToPool creates a *sync.Pool from this Generator.
func (g Generator) ToPool() *sync.Pool {
	return &sync.Pool{
		New: g,
	}
}

// SharedPoolGenerator creates a Generator that returns a DefaultDict using g as
// the Generator.
//
// It's different from just `func() interface{} { return New(g) }` that the
// DefaultDicts returned by the same SharedPoolGenerator shares the same
// underlying pool, so it's more memory efficient when used as the second (or
// third, etc.) layer of DefaultDict.
func SharedPoolGenerator(g Generator) Generator {
	p := g.ToPool()
	return func() interface{} {
		return &defaultdict{
			pool: p,
		}
	}
}
