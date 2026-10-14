// Package trade is the order namespace. Every method returns a PendingTrade
// that can be routed with Via(transport) before Do(ctx) sends it. A UUID v4
// client order id is allocated up front. AddOrder is never auto-retried.
//
// Rust counterpart: rust/src/api/trade.
package trade
