package defaultdict

// Comparable is the key type of a Map.
//
// It's used for documentation purpose only.
//
// See https://golang.org/ref/spec#Comparison_operators for more info.
type Comparable = interface{}

// Pointer is the value type of a Map.
//
// It's used for documentation purpose only.
//
// In a Map, all values should be pointers, and those pointers should be safe to
// be mutated concurrently, for the following reasons:
//
// 1, A Map is for concurrent mutations (if you only need concurrent read-only
// access, a builtin map would suffice)
//
// 2. There's no Store function provided by Map interface. All mutations are
// done by Get/Load then mutate the returned value directly.
//
// As an example, you can use *int64 as Pointer, and do mutations via atomic
// int64 operations. But slices, even though they are pointers, should not be
// used here directly. You usually need to pair slice with a lock, for example:
//
//     type myValue struct{
//       lock sync.Lock
//       // Note that this must be the pointer to the slice,
//       // because append calls could change the slice itself.
//       slice *[]int
//     }
//
//     func (m *myValue) Append(i int) {
//       m.lock.Lock()
//       defer m.lock.Unlock()
//       *m.slice = append(*m.slice, i)
//     }
//
//     func myValueGenerator() defaultdict.Pointer {
//       return &myValue{
//         slice: new([]int),
//       }
//     }
type Pointer = interface{}

// Map defines a map.
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
type Map interface {
	// Same as in sync.Map, see above for notes about the bool returns.
	Delete(key Comparable)
	Load(key Comparable) (Pointer, bool)
	LoadAndDelete(key Comparable) (Pointer, bool)
	Range(f func(key Comparable, value Pointer) bool)

	// Same as Load, just without the bool return.
	Get(key Comparable) Pointer
}
