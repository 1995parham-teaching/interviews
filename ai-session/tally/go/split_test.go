package tally

import (
	"maps"
	"testing"
)

func TestParseSplit(t *testing.T) {
	cases := []struct {
		spec   string
		amount int64
		want   map[string]int64
	}{
		{"equal:ali,sara,reza", 9000, map[string]int64{"ali": 3000, "sara": 3000, "reza": 3000}},
		{"shares:ali=1,sara=2", 3000, map[string]int64{"ali": 1000, "sara": 2000}},
		{"exact:ali=20.00,sara=30.00", 5000, map[string]int64{"ali": 2000, "sara": 3000}},
	}
	for _, c := range cases {
		got, err := ParseSplit(c.spec, c.amount)
		if err != nil || !maps.Equal(got, c.want) {
			t.Errorf("ParseSplit(%q, %d) = %v, %v; want %v", c.spec, c.amount, got, err, c.want)
		}
	}
}

func TestExactMoreThanAmount(t *testing.T) {
	if _, err := ParseSplit("exact:ali=20.00,sara=40.00", 5000); err == nil {
		t.Error("parts above the amount were accepted")
	}
}

func TestParseSplitInvalid(t *testing.T) {
	for _, spec := range []string{"equal", "equal:", "halves:ali,sara", "shares:ali", "shares:ali=0"} {
		if _, err := ParseSplit(spec, 1000); err == nil {
			t.Errorf("ParseSplit(%q) succeeded, want an error", spec)
		}
	}
}
