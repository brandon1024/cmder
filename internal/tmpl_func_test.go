package internal

import (
	"slices"
	"testing"
)

func TestWrapLines(t *testing.T) {
	tests := []struct {
		name   string
		length int
		in     string
		want   []string
	}{
		// Basic wrapping
		{"empty", 10,
			"",
			nil},
		{"fits exactly", 10,
			"abc defghi",
			[]string{"abc defghi"}},
		{"one rune over", 10,
			"abc defghij",
			[]string{"abc", "defghij"}},
		{"breaks at last space that fits", 10,
			"The quick brown fox jumps over the lazy dog",
			[]string{"The quick", "brown fox", "jumps over", "the lazy", "dog"}},
		{"surrounding whitespace trimmed", 10,
			"   padded   ",
			[]string{"padded"}},
		{"double space at break", 2,
			"a  b",
			[]string{"a", "b"}},

		// Long words
		{"long word kept whole", 10,
			"abcdefghijklmnop qr",
			[]string{"abcdefghijklmnop", "qr"}},
		{"long word alone", 10,
			"https://example.com/a/very/long/path",
			[]string{"https://example.com/a/very/long/path"}},
		{"url between words", 10,
			"see https://example.com/a/very/long/path ok",
			[]string{"see", "https://example.com/a/very/long/path", "ok"}},

		// Line feeds
		{"newline is hard break", 20,
			"first line\nsecond line",
			[]string{"first line", "second line"}},
		{"blank line preserved", 10,
			"a\n\nb",
			[]string{"a", "", "b"}},
		{"trailing newline", 10,
			"abc\n",
			[]string{"abc"}},
		{"crlf", 10,
			"abc\r\ndef\r\n",
			[]string{"abc", "def"}},
		{"wrap per paragraph", 10,
			"The quick brown fox\njumps over the lazy dog",
			[]string{"The quick", "brown fox", "jumps over", "the lazy", "dog"}},

		// UTF-8 (lengths are in runes, not bytes)
		{"multibyte fits by runes", 14,
			"Grüße aus Köln",
			[]string{"Grüße aus Köln"}},
		{"multibyte wraps by runes", 9,
			"Grüße aus Köln",
			[]string{"Grüße aus", "Köln"}},
		{"multibyte long word", 5,
			"日本語のテキスト 次の行",
			[]string{"日本語のテキスト", "次の行"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wrapLines(tt.in, tt.length)

			if !slices.Equal(got, tt.want) {
				t.Errorf("wrapLines(%q, %d)\n got: %q\nwant: %q", tt.in, tt.length, got, tt.want)
			}
		})
	}
}
