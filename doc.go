// Package defaultdict implements Python's defaultdict, in a way that's both
// thread-safe and memory-efficient.
//
// There are two example use cases for it:
//
// 1. To implement a row lock, that every row (key) has its own lock.
//
// 2. To implement a concurrent counter, that every key has its own atomic int
// as the counter.
package defaultdict // import "go.yhsif.com/defaultdict"
