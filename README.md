[![Go Reference](https://pkg.go.dev/badge/go.yhsif.com/defaultdict.svg)](https://pkg.go.dev/go.yhsif.com/defaultdict)
[![Go Report Card](https://goreportcard.com/badge/go.yhsif.com/defaultdict)](https://goreportcard.com/report/go.yhsif.com/defaultdict)

# go-defaultdict

Go implementation of [Python's `defaultdict`][python-defaultdict],
in a way that's both thread-safe and memory efficient.

## Overview

Underneath it pairs a [`sync.Map`][sync-map] with a [`sync.Pool`][sync-pool],
and removed all direct store/write accesses to the map.
As a result, the only way to mutate the map is through `Load`/`Get`,
(which either create a new value for you if this is the first access to the key,
or return the value created by a previous `Load`/`Get`),
then mutate the value returned directly (in a thread-safe way).

Here are 2 example usages:

1. To implement a rowlock.
   See [my rowlock package][rowlock] for detailed example.

2. To implement a concurrent counter-by-key.
   See [package example][package-example] or below for details.

## Example Code

Here's a step-by-step example to create a concurrent counter-by-key.

First, create a generator,
which simply returns an `*int64` so it can be used by atomic int64 functions:

```go
generator := func() defaultdict.Pointer {
  return new(int64)
}
```

Then, create the map:

```go
m := defaultdict.New(generator)
```

When you need to add the counter, get by key then use `atomic.AddInt64`:

```go
atomic.AddInt64(m.Get(key).(*int64), 1)
```

When you need to get the counter value,
just get by key then use `atomic.LoadInt64`:

```go
fmt.Printf(
  "Key %v was added %d times\n",
  key,
  atomic.LoadInt64(
    m.Get(key).(*int64),
  ),
)
```

Or use `Range`:

```go
m.Range(func(key defaultdict.Comparable, value defaultdict.Pointer) bool {
  fmt.Printf("Key %v was added %d times\n", key, atomic.LoadInt64(value.(*int64)))
  return true
})
```

## License

[BSD License](LICENSE).

[python-defaultdict]: https://docs.python.org/3/library/collections.html#collections.defaultdict
[sync-map]: https://pkg.go.dev/sync#Map
[sync-pool]: https://pkg.go.dev/sync#Pool
[rowlock]: https://github.com/fishy/rowlock/blob/v0.4.0/rowlock.go
[package-example]: https://pkg.go.dev/go.yhsif.com/defaultdict#example-package
