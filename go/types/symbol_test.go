package types

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestNewSymbol(t *testing.T) {
	tests := []struct {
		in      string
		want    string
		wantErr error
	}{
		{in: "BTC/USD", want: "BTC/USD"},

		// Case is canonicalised to uppercase.
		{in: "btc/usd", want: "BTC/USD"},
		{in: "Eth/Usd", want: "ETH/USD"},
		{in: "1inch/usd", want: "1INCH/USD"},

		// Digits are valid in legs.
		{in: "1INCH/USD", want: "1INCH/USD"},
		{in: "0G/USD", want: "0G/USD"},
		{in: "2Z/USD", want: "2Z/USD"},
		{in: "AI16Z/USDC", want: "AI16Z/USDC"},

		// Modern codes that start with X/Z are not legacy.
		{in: "XRP/USD", want: "XRP/USD"},
		{in: "ZEC/EUR", want: "ZEC/EUR"},

		{in: "", wantErr: ErrEmpty},
		{in: "BTC/", wantErr: ErrEmpty},
		{in: "/USD", wantErr: ErrEmpty},

		{in: "XXBTZUSD", wantErr: ErrLegacyPrefix},
		{in: "ZZ-anything", wantErr: ErrLegacyPrefix},
		{in: "XBT/USD", wantErr: ErrLegacyPrefix},
		{in: "BTC/ZUSD", wantErr: ErrLegacyPrefix},
		{in: "xbt/usd", wantErr: ErrLegacyPrefix},

		{in: "BTCUSD", wantErr: ErrMissingSeparator},
		{in: "!!!", wantErr: ErrMissingSeparator},

		{in: "BTC/USD:KDE", wantErr: ErrNonAlphanumeric}, // WS v2 venue suffix
		{in: "B!C/USD", wantErr: ErrNonAlphanumeric},
		{in: "BTC/US D", wantErr: ErrNonAlphanumeric},
		{in: "BTC/USD/EUR", wantErr: ErrNonAlphanumeric},
	}
	for _, tt := range tests {
		got, err := NewSymbol(tt.in)
		if tt.wantErr != nil {
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewSymbol(%q) error = %v, want %v", tt.in, err, tt.wantErr)
			}
			if !got.IsZero() {
				t.Errorf("NewSymbol(%q) = %q on error, want zero value", tt.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("NewSymbol(%q) unexpected error: %v", tt.in, err)
			continue
		}
		if got.String() != tt.want {
			t.Errorf("NewSymbol(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestSymbolJSON(t *testing.T) {
	btcusd, err := NewSymbol("BTC/USD")
	if err != nil {
		t.Fatal(err)
	}

	b, err := json.Marshal(map[Symbol]int{btcusd: 1})
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"BTC/USD":1}` {
		t.Errorf("marshal map = %s, want {\"BTC/USD\":1}", b)
	}

	var v struct {
		Pair Symbol `json:"pair"`
	}
	if err := json.Unmarshal([]byte(`{"pair":"btc/usd"}`), &v); err != nil {
		t.Fatal(err)
	}
	if v.Pair != btcusd {
		t.Errorf("decoded pair = %q, want BTC/USD", v.Pair)
	}

	b, err = json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"pair":"BTC/USD"}` {
		t.Errorf("marshal struct = %s, want {\"pair\":\"BTC/USD\"}", b)
	}

	if err := json.Unmarshal([]byte(`{"pair":"XBT/USD"}`), &v); err == nil {
		t.Errorf("decoding legacy XBT/USD succeeded, want error")
	}
}
