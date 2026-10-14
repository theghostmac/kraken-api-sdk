# Go SDK roadmap

Work bottom-up: each milestone only depends on the ones before it, and each
one ships with tests. Rust paths are relative to `../rust/src`.

## Decisions to make before milestone 1

| Topic | Recommendation | Why |
|---|---|---|
| Decimal type | `github.com/shopspring/decimal` | De facto standard; string round-trip; JSON support. |
| WebSocket library | `github.com/coder/websocket` | Context-aware, maintained, concurrent-write safe. |
| UUID | `github.com/google/uuid` | `ClOrdID` allocation (v4). |
| Callbacks vs channels | Callbacks (`func(T)`) run on the dispatch goroutine, matching Rust | Keeps the design one SDK across languages. |
| Async model | Blocking methods taking `context.Context`; goroutines internally | Idiomatic Go replacement for Rust `async`. |

## Milestone 0 — value types and errors (no I/O)

1. **`types/symbol.go`** ← `types/symbol.rs`
   - `type Symbol struct{ s string }`, `NewSymbol(string) (Symbol, error)`, `(Symbol) String() string`
   - Sentinel errors: `ErrSymbolEmpty`, `ErrLegacyPrefix`, `ErrNonAlphanumeric`, `ErrMissingSeparator`
   - Port every test in `symbol.rs` into `types/symbol_test.go` (table-driven).
2. **`types/asset.go`** ← `types/asset.rs`
   - `AssetCode`, `NewAssetCode`, `ErrNonAlphabetic`, unexported `normaliseAsset(raw string) string` (legacy → modern).
3. **`types/ids.go`** ← `types/ids.rs`
   - `TxID`, `ClOrdID` + `NewClOrdID(string) (ClOrdID, error)` + `AllocateClOrdID() ClOrdID`, `APIKey`, `APISecret` + `APISecretFromBase64(string) (APISecret, error)`.
   - `APIKey`/`APISecret` must redact in `String()`, `GoString()` and `fmt` verbs.
4. **`sdkerr/error.go`** ← `error.rs`
   - `type Category string` + constants `Config, Auth, Network, Exchange, Client, RateLimit`
   - `type Error interface { error; Code() string; Category() Category; Retryable() bool; RequestID() string; Message() string }`
   - Helpers ported from `error.rs`: `krakenEClassCode`, `isKrakenRateLimit`, `parseThrottleUntil` (these belong in `internal/rest` later if only REST uses them).

## Milestone 1 — first public REST call

5. **`internal/clock/clock.go`** ← `clock/mod.rs`: `type Clock interface { Now() time.Time }`, `System`, and a fake for tests.
6. **`internal/transport/http.go`** ← `transport/http.rs`: `type HTTP interface { Do(ctx, *http.Request) (*http.Response, error) }`, default with timeout + TLS 1.2 minimum; typed `TransportError`.
7. **`internal/rest/surface.go`** ← `rest/surface.rs`: `Surface.Public(ctx, method, path string, params url.Values, out any) error` — decodes the `{"error":[],"result":…}` envelope and maps errors.
8. **`internal/config/config.go`** ← `build/config_resolver.rs`: `Config` struct + `Resolve(opts) (Config, error)`; reject `http://`/`ws://` endpoints.
9. **Root `kraken/client.go` + `options.go`** ← `build/client`, `build/client/builder.rs`
   - `New(opts ...Option) (*Client, error)` (no network), `WithAPIKey(key, secret string)`, `WithRESTBaseURL`, `WithHTTPClient`.
   - `(*Client) Market() *market.Namespace`
10. **`market/rest.go`** ← `api/market/rest.rs`: `ServerTime(ctx)`, `Status(ctx)`, then `Ticker(ctx, ...Symbol)`.
11. **`examples/quickstart/main.go`** ← `examples/quickstart.rs`.

**Done when:** `go run ./examples/quickstart` prints server time and BTC/USD last price, and every step above has `httptest` tests.

## Milestone 2 — private REST

12. **`internal/auth/nonce.go`** ← `auth/nonce.rs`: monotonic nonce source.
13. **`internal/auth/signer.go`** ← `auth/signer.rs`: `computeAPISign(secret []byte, nonce uint64, path, body string) string`.
    Golden test from Rust: secret = 32×`0x01`, nonce `1700000000000000000`, path `/0/private/Balance`, body `nonce=1700000000000000000` → `6i7+cpXVLDDnDOekgsRcgFQhkYFyPDSSxhSVA8OrA6vjziXcy2K1DP3n7f9B/Kgni71WjBSQusZKxduRsdF20Q==`.
14. **`internal/rest/retry.go`** ← `rest/retry.rs`: retry only idempotent calls. Never `AddOrder`.
15. **`account/`** ← `api/account`: `Balance`, `TradeBalance`, `OpenOrders`, `ClosedOrders`, `Ledgers`, …
16. **`internal/ratelimit/`** ← `rate_limit/`: counters + warning thresholds (non-blocking).

## Milestone 3 — REST trading

17. **`trade/`** ← `api/trade`: `PendingTrade` with `Via(Transport)` and `Do(ctx)`; `LimitBuy/Sell`, `MarketBuy/Sell`, `Cancel`, `CancelAll`, `CancelAllOrdersAfter`, batch. REST transport only for now.

## Milestone 4 — order book (no I/O)

18. **`book/checksum.go`** ← `book/checksum.rs`: `Checksum(asks, bids []PriceLevel) uint32` (IEEE CRC32 over stripped price/qty strings).
    Golden test: Kraken's 281-char worked example hashes to `3310070434` (see `book/checksum_tests.rs`).
19. **`book/builder.go`** ← `book/builder.rs`: `Builder.ApplySnapshot`, `ApplyDelta`, gap detection, resync state.

## Milestone 5 — WebSocket

20. `internal/transport/ws.go`, `internal/conn` (managed connection FSM, reconnect + backoff with `internal/jitter`, subscription registry), `internal/dispatch` (I/O goroutine + dispatch goroutine, bounded drop-oldest queue), `internal/auth` token lifecycle.
21. `subscription/`, `events/`, `market` and `account` `On*` callbacks, WebSocket order transport in `trade`.

Read `../rust/docs/architecture.md` and `../rust/docs/dispatch-flow.md` before starting this milestone.
