// Package run provides panic isolation for the aggregation engines: one
// misbehaving source must not take down the whole run.
package run

import "log"

// Safe runs fn, recovering from any panic and logging it instead of
// crashing the engine.
func Safe(name string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ENGINE PANIC] %s recovered: %v", name, r)
		}
	}()
	fn()
}

// SafeIterate runs fn, recovering from any panic and logging it. It wraps
// per-item work inside an engine so one bad item does not abort the rest.
func SafeIterate(name string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ENGINE ITERATION] %s recovered: %v", name, r)
		}
	}()
	fn()
}
