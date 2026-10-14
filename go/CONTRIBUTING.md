# Contributing to the Go SDK

Repo-wide policy is in [../CONTRIBUTING.md](../CONTRIBUTING.md). This file
covers Go-specific rules.

## Checks

Run before every push (CI runs the same):

```bash
make check
```

That runs `gofmt`, `go vet`, and `go test -race ./...`.

## Rules

- **Match the Rust SDK's behaviour.** When porting a component, read the matching
  Rust module (see the layout table in `README.md`) and port its tests too.
- **No `float64` for money.** Prices, volumes, balances and fees are decimals.
- **`context.Context` first** on every call that can block or reach the network.
- **Errors implement `sdkerr.Error`** (`Code`, `Category`, `Retryable`,
  `RequestID`, `Message`). Wrap with `%w`; callers use `errors.As`.
- **Never log or format secrets**: API key, secret, WS token, nonce.
- **Keep internals in `internal/`.** Only the root and the namespace packages
  are public API.
- **No network in unit tests.** Use `httptest` and fake transports. Live tests
  go behind the `live` build tag.
- Every package that has code has tests. Table-driven where it fits.
