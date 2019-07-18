package chainup_test

import (
	"testing"

	"chainup.dev/chainup"
	"chainup.dev/lib/test"
)

func TestProviderTypeIsValid(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"digitalocean", true},
		{"", false},
		{"superprovider", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			net := chainup.NewProviderType(tt.name)

			got := net.IsValid()
			test.AssertBoolEqual(t, "ProviderType.IsValid()", got, tt.valid)
			test.AssertStringsEqual(t, "ProviderType.String()", net.String(), tt.name)
		})
	}
}
