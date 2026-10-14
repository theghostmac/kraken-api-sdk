package types

import "errors"

// Validation errors returned by NewAssetCode and NewSymbol. Match them with
// errors.Is.
var (
	// ErrEmpty: the input, or one BASE/QUOTE leg of a symbol, is empty.
	ErrEmpty = errors.New("kraken: symbol or asset code is empty")
	// ErrLegacyPrefix: a legacy code was passed (XBT, XXBT, ZUSD, XXBTZUSD).
	ErrLegacyPrefix = errors.New("kraken: legacy code not accepted; use the modern form (e.g. BTC/USD, not XXBTZUSD)")
	// ErrNonAlphanumeric: a symbol leg contains a character outside ASCII letters and digits.
	ErrNonAlphanumeric = errors.New("kraken: symbol leg must be alphanumeric")
	// ErrNonAlphabetic: an asset code contains a character outside ASCII letters.
	ErrNonAlphabetic = errors.New("kraken: asset code must be alphabetic")
	// ErrMissingSeparator: a symbol has no '/' separator (e.g. "BTCUSD").
	ErrMissingSeparator = errors.New("kraken: symbol must use the BASE/QUOTE form (e.g. BTC/USD)")
)
