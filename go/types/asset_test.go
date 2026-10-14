package types

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestNormaliseAsset(t *testing.T) {
	tests := []struct {
		raw, want string
	}{
		{"XXBT", "BTC"}, // not XBT: the altname is deprecated too
		{"XBT", "BTC"},
		{"ZUSD", "USD"},
		{"XETH", "ETH"},
		{"ZEUR", "EUR"},
		{"KFEE", "FEE"},
		{"USDC", "USDC"},
		{"BABY", "BABY"},
		{"UNKNWN", "UNKNWN"},
	}
	for _, tt := range tests {
		if got := normaliseAsset(tt.raw); got != tt.want {
			t.Errorf("normaliseAsset(%q) = %q, want %q", tt.raw, got, tt.want)
		}
	}
}

func TestNewAssetCode(t *testing.T) {
	tests := []struct {
		in      string
		want    string
		wantErr error
	}{
		{in: "BTC", want: "BTC"},
		{in: "XRP", want: "XRP"}, // modern codes that start with X/Z are valid
		{in: "XLM", want: "XLM"},
		{in: "ZEC", want: "ZEC"},
		{in: "btc", want: "BTC"},
		{in: "Usdc", want: "USDC"},
		{in: "", wantErr: ErrEmpty},
		{in: "XXBT", wantErr: ErrLegacyPrefix},
		{in: "ZUSD", wantErr: ErrLegacyPrefix},
		{in: "XBT", wantErr: ErrLegacyPrefix},
		{in: "xxbt", wantErr: ErrLegacyPrefix},
		{in: "zusd", wantErr: ErrLegacyPrefix},
		{in: "XBT.B", wantErr: ErrNonAlphabetic},
		{in: "BTC2", wantErr: ErrNonAlphabetic},
		{in: "BTÇ", wantErr: ErrNonAlphabetic}, // non-ASCII letter
	}
	for _, tt := range tests {
		got, err := NewAssetCode(tt.in)
		if tt.wantErr != nil {
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewAssetCode(%q) error = %v, want %v", tt.in, err, tt.wantErr)
			}
			if !got.IsZero() {
				t.Errorf("NewAssetCode(%q) = %q on error, want zero value", tt.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("NewAssetCode(%q) unexpected error: %v", tt.in, err)
			continue
		}
		if got.String() != tt.want {
			t.Errorf("NewAssetCode(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestAssetCodeFromWire(t *testing.T) {
	tests := []struct {
		raw, want string
	}{
		{"XXBT", "BTC"},
		{"XBT", "BTC"},
		{"XBT.B", "XBT.B"}, // unknown codes pass through verbatim
	}
	for _, tt := range tests {
		if got := assetCodeFromWire(tt.raw).String(); got != tt.want {
			t.Errorf("assetCodeFromWire(%q) = %q, want %q", tt.raw, got, tt.want)
		}
	}
}

func TestAssetCodeJSONMapKey(t *testing.T) {
	btc, err := NewAssetCode("BTC")
	if err != nil {
		t.Fatal(err)
	}

	b, err := json.Marshal(map[AssetCode]int{btc: 42})
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"BTC":42}` {
		t.Errorf("marshal = %s, want {\"BTC\":42}", b)
	}

	var m map[AssetCode]int
	if err := json.Unmarshal([]byte(`{"XXBT":1}`), &m); err != nil {
		t.Fatal(err)
	}
	if got := m[btc]; got != 1 {
		t.Errorf("decoded map = %v, want legacy XXBT key normalised to BTC", m)
	}
}

func TestAssetCodeJSONValue(t *testing.T) {
	var v struct {
		Asset AssetCode `json:"asset"`
	}
	if err := json.Unmarshal([]byte(`{"asset":"ZUSD"}`), &v); err != nil {
		t.Fatal(err)
	}
	if v.Asset.String() != "USD" {
		t.Errorf("decoded asset = %q, want USD", v.Asset)
	}

	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"asset":"USD"}` {
		t.Errorf("marshal = %s, want {\"asset\":\"USD\"}", b)
	}
}
