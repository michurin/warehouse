package text_test

import (
	"testing"

	"github.com/michurin/warehouse/go/chat2/text"
)

func TestSanitizeText(t *testing.T) {
	for n, tt := range []struct {
		input    string
		crop     int
		empty    string
		expected string
	}{
		// spaces
		{"just text", 100, "", "just text"},
		{" just\n\r text\b\b\bnext ", 100, "", "just text next"},
		// cropping
		{"crop text", 3, "", "cro"},
		{"crop text", 4, "", "crop"},
		{"crop text", 5, "", "crop"},
		{"crop text", 6, "", "crop t"},
		// fallback
		{"", 10, "fallback", "fallback"},
		{"text", 0, "fallback", "fallback"},
	} {
		msg := text.SanitizeText(tt.input, tt.crop, tt.empty)
		if msg != tt.expected {
			t.Errorf("Test #%d failted: %+v: %q", n, tt, msg)
		}
	}
}
