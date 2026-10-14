// Package types holds the shared value types used across every namespace:
// Symbol, AssetCode, ClOrdID, TxID, BookDepth, ChannelName, ConnectionState.
//
// Symbols are modern form only ("BTC/USD"); legacy codes (XBT, XXBT, ZUSD)
// are rejected with a typed error.
//
// Rust counterpart: rust/src/types.
package types
