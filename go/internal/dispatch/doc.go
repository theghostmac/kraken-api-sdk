// Package dispatch runs the two loops: the I/O loop (owns sockets, decodes
// frames, maintains books) and the dispatch loop (runs on* callbacks and event
// handlers), joined by a bounded drop-oldest queue.
//
// Rust counterpart: rust/src/dispatch.
package dispatch
