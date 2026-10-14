// Package types holds the value types shared across the SDK: asset codes,
// symbols and ids. Rust counterpart: rust/src/types.
package types

import "strings"

// AssetCode is a modern user-facing asset code (e.g., BTC, USD).
//
// Build one from caller input with NewAssetCode, which rejects legacy codes.
// Decoding (UnmarshalText, and therefore encoding/json) is the lenient wire
// path for exchange payloads: legacy codes are normalised, unknown codes pass
// through verbatim, and nothing is rejected. Don't decode untrusted user
// input into an AssetCode; call NewAssetCode instead.
type AssetCode struct {
	value string
}

// NewAssetCode constructs an AssetCode from caller-supplied input. It
// uppercases the input and rejects empty strings (ErrEmpty), non-letters
// (ErrNonAlphabetic) and legacy X/Z codes such as XXBT or ZUSD (ErrLegacyPrefix).
func NewAssetCode(s string) (AssetCode, error) {
	if s == "" {
		return AssetCode{}, ErrEmpty
	}
	for _, r := range s {
		if !isASCIILetter(r) {
			return AssetCode{}, ErrNonAlphabetic
		}
	}

	s = strings.ToUpper(s)

	// Table-driven reject, not a leading X/Z strip, so XRP/XLM/ZEC stay valid.
	if normaliseAsset(s) != s {
		return AssetCode{}, ErrLegacyPrefix
	}

	return AssetCode{value: s}, nil
}

// assetCodeFromWire constructs an AssetCode from a server-returned code,
// normalizing legacy forms. Infallible.
func assetCodeFromWire(raw string) AssetCode {
	return AssetCode{value: normaliseAsset(raw)}
}

// String returns the asset code.
func (a AssetCode) String() string {
	return a.value
}

// IsZero reports whether `a` is the zero AssetCode.
func (a AssetCode) IsZero() bool {
	return a.value == ""
}

// MarshalText implements encoding.TextMarshaler.
func (a AssetCode) MarshalText() ([]byte, error) {
	return []byte(a.value), nil
}

// UnmarshalText implements encoding.TextUnmarshaler using the lenient wire
// path (see AssetCode).
func (a *AssetCode) UnmarshalText(text []byte) error {
	*a = assetCodeFromWire(string(text))
	return nil
}

func isASCIILetter(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z'
}

// normaliseAsset normalises a raw Kraken asset code to its modern user-facing form.
// Unknown or already-modern codes pass through unchanged.
func normaliseAsset(raw string) string {
	switch raw {
	case "XXBT", "XBT":
		return "BTC"
	case "XETC":
		return "ETC"
	case "XETH":
		return "ETH"
	case "XLTC":
		return "LTC"
	case "XMLN":
		return "MLN"
	case "XREP":
		return "REP"
	case "XXDG":
		return "XDG"
	case "XXLM":
		return "XLM"
	case "XXMR":
		return "XMR"
	case "XXRP":
		return "XRP"
	case "XZEC":
		return "ZEC"
	case "ZARS":
		return "ARS"
	case "ZAUD":
		return "AUD"
	case "ZCAD":
		return "CAD"
	case "ZCLP":
		return "CLP"
	case "ZCOP":
		return "COP"
	case "ZDKK":
		return "DKK"
	case "ZEUR":
		return "EUR"
	case "ZGBP":
		return "GBP"
	case "ZGEL":
		return "GEL"
	case "ZGHS":
		return "GHS"
	case "ZJPY":
		return "JPY"
	case "ZLKR":
		return "LKR"
	case "ZMXN":
		return "MXN"
	case "ZPLN":
		return "PLN"
	case "ZSEK":
		return "SEK"
	case "ZUGX":
		return "UGX"
	case "ZUSD":
		return "USD"
	case "ZVND":
		return "VND"
	case "ZXOF":
		return "XOF"
	case "KFEE":
		return "FEE"
	default:
		return raw
	}
}
