// Package auth signs Spot REST requests (HMAC-SHA512), issues nonces, and
// manages the WebSocket auth-token lifecycle. Secrets are never logged.
//
// Rust counterpart: rust/src/auth.
package auth
