// Package ratelimit tracks API and trading rate-limit counters per tier. It
// never blocks the caller; it reports usage so warnings can be emitted.
//
// Rust counterpart: rust/src/rate_limit.
package ratelimit
