package types

import "strings"

// Symbol represents a wire-form trading pair, e.g. "BTC/USD" (modern BASE/QUOTE form).
// Legacy aliases (e.g. XBT, XXBTZUSD) are rejected at construction.
//
// Decoding (UnmarshalText, and therefore encoding/json) is strict: it calls
// NewSymbol. Exchange payloads that carry legacy pair names (e.g. AssetPairs
// wsname "XBT/USD") must be decoded as string and converted in the REST layer.
type Symbol struct {
	value string
}

// NewSymbol constructs a Symbol from a wire-form string, validating its shape:
// non-empty, modern BASE/QUOTE with a single '/' separator and valid non-legacy legs.
func NewSymbol(wire string) (Symbol, error) {
	if wire == "" {
		return Symbol{}, ErrEmpty
	}

	// Reject legacy concatenated X/Z forms (e.g., "XXBTZUSD") before the separator check.
	if upper := strings.ToUpper(wire); strings.HasPrefix(upper, "XX") || strings.HasPrefix(upper, "ZZ") {
		return Symbol{}, ErrLegacyPrefix
	}

	base, quote, found := strings.Cut(wire, "/")
	if !found {
		return Symbol{}, ErrMissingSeparator
	}

	if err := validateLeg(base); err != nil {
		return Symbol{}, err
	}
	if err := validateLeg(quote); err != nil {
		return Symbol{}, err
	}

	// Uppercase only after validation: both legs are ASCII by now, so Unicode
	// case mapping (e.g. 'ı' -> 'I') can't turn invalid input into a valid pair.
	return Symbol{value: strings.ToUpper(wire)}, nil
}

// validateLeg validates a single leg: non-empty, ASCII alphanumeric, and non-legacy.
// Digits are allowed (e.g., 1INCH, AI16Z).
func validateLeg(leg string) error {
	if leg == "" {
		return ErrEmpty
	}

	for _, r := range leg {
		if !isASCIIAlphanumeric(r) {
			return ErrNonAlphanumeric
		}
	}

	upper := strings.ToUpper(leg)
	if normaliseAsset(upper) != upper {
		return ErrLegacyPrefix
	}

	return nil
}

// String returns the bare wire-form pair string (e.g. "BTC/USD").
func (s Symbol) String() string {
	return s.value
}

// IsZero reports whether s is the uninitialized zero value.
func (s Symbol) IsZero() bool {
	return s.value == ""
}

// MarshalText implements encoding.TextMarshaler.
func (s Symbol) MarshalText() ([]byte, error) {
	return []byte(s.value), nil
}

// UnmarshalText implements encoding.TextUnmarshaler. It is strict (see Symbol).
func (s *Symbol) UnmarshalText(text []byte) error {
	sym, err := NewSymbol(string(text))
	if err != nil {
		return err
	}
	*s = sym
	return nil
}

func isASCIIAlphanumeric(r rune) bool {
	return isASCIILetter(r) || ('0' <= r && r <= '9')
}
