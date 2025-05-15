package storage

import (
	"testing"
	"unicode/utf8"
)

func TestGenerateXMLID(t *testing.T) {
	tests := []struct {
		name     string
		length   int64
		wantLen  int
		wantOnly string
	}{
		{
			name:     "Default length",
			length:   LengthXMLID,
			wantLen:  int(LengthXMLID),
			wantOnly: "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789",
		},
		{
			name:     "Zero length",
			length:   0,
			wantLen:  0,
			wantOnly: "",
		},
		{
			name:     "Custom length",
			length:   15,
			wantLen:  15,
			wantOnly: "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateXMLID(tt.length)

			if utf8.RuneCountInString(got) != tt.wantLen {
				t.Errorf("GenerateXMLID() length = %d, want %d", utf8.RuneCountInString(got), tt.wantLen)
			}

			if tt.length > 0 {
				got2 := GenerateXMLID(tt.length)
				if got == got2 {
					t.Error("GenerateXMLID() returns the same value on subsequent calls, want different values")
				}
			}
		})
	}
}
