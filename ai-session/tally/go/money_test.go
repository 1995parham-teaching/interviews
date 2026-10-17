package tally

import "testing"

func TestParseAmount(t *testing.T) {
	cases := map[string]int64{"12.50": 1250, "7": 700, "0.05": 5}
	for text, want := range cases {
		got, err := ParseAmount(text)
		if err != nil || got != want {
			t.Errorf("ParseAmount(%q) = %d, %v; want %d", text, got, err, want)
		}
	}
}

func TestParseAmountRejectsGarbage(t *testing.T) {
	for _, text := range []string{"", "abc", "-3", "1,5", "0"} {
		if _, err := ParseAmount(text); err == nil {
			t.Errorf("ParseAmount(%q) succeeded, want an error", text)
		}
	}
}

func TestFormatAmount(t *testing.T) {
	cases := map[int64]string{1250: "12.50", 5: "0.05", -300: "-3.00"}
	for cents, want := range cases {
		if got := FormatAmount(cents); got != want {
			t.Errorf("FormatAmount(%d) = %q, want %q", cents, got, want)
		}
	}
}

func TestFormatSigned(t *testing.T) {
	cases := map[int64]string{4000: "+40.00", 0: "0.00"}
	for cents, want := range cases {
		if got := FormatSigned(cents); got != want {
			t.Errorf("FormatSigned(%d) = %q, want %q", cents, got, want)
		}
	}
}
