# Kraken SDK (Go)

> **Status: work in progress.** Package layout only; no functionality yet.
> See [`docs/ROADMAP.md`](docs/ROADMAP.md).

Go implementation of the Kraken unified SDK design: one client over Spot REST
and Spot WebSocket v2, grouped into `Market`, `Account`, `Trade`,
`Subscription` and `Events`. The Rust SDK in [`../rust`](../rust) is the
reference implementation; behaviour should match it.

```go
import kraken "github.com/krakenfx/kraken-api-sdk/go"
```

## Layout

| Package | Purpose | Rust counterpart |
|---|---|---|
| `kraken` (root) | `Client`, `New`, options | `src/build` |
| `types` | `Symbol`, `AssetCode`, `ClOrdID`, `TxID`, … | `src/types` |
| `sdkerr` | error interface + categories | `src/error.rs` |
| `market` | public data, REST + WS | `src/api/market` |
| `account` | private reads, REST + WS | `src/api/account` |
| `trade` | orders, `PendingTrade` | `src/api/trade` |
| `subscription` | stream management | `src/api/subscription_*` |
| `events` | lifecycle event bus | `src/api/events.rs`, `src/dispatch/event_bus` |
| `book` | CRC32-checked order book | `src/book` |
| `internal/auth` | signing, nonces, WS token | `src/auth` |
| `internal/transport` | HTTP + WS adapters | `src/transport` |
| `internal/rest` | REST surface + retry | `src/rest` |
| `internal/ratelimit` | rate-limit tracking | `src/rate_limit` |
| `internal/conn` | connection FSM, subscription registry | `src/conn` |
| `internal/dispatch` | I/O loop + dispatch loop | `src/dispatch` |
| `internal/config` | config resolution, knobs | `src/build/config_resolver.rs`, `knobs.rs` |
| `internal/clock`, `internal/jitter` | testable time and randomness | `src/clock`, `src/jitter` |

## Development

```bash
make check   # fmt + vet + test (race)
```

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

Apache License 2.0 — see [../LICENSE](../LICENSE).
