// Package kraken is the Go SDK for Kraken Spot REST and Spot WebSocket v2.
//
// One Client, built once with New, exposes five namespaces: Market, Account,
// Trade, Subscription and Events. New does no network I/O.
//
// Rust counterpart: rust/src/build (ClientBuilder, Client, knobs).
package kraken
