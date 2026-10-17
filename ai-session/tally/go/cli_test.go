package tally

import (
	"os"
	"strings"
	"testing"
)

func TestExampleTrip(t *testing.T) {
	f, err := os.Open("examples/trip.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var out strings.Builder
	failures, err := Run(f, &out)
	if err != nil || failures != 0 {
		t.Fatalf("Run = %d failures, %v", failures, err)
	}
	want := "#1 hotel: ali paid 90.00\n" +
		"#2 dinner: sara paid 24.00\n" +
		"#3 train tickets: reza paid 18.00\n" +
		"sara paid ali 10.00\n" +
		"ali   +36.00\n" +
		"reza  -24.00\n" +
		"sara  -12.00\n" +
		"reza -> ali 24.00\n" +
		"sara -> ali 12.00\n"
	if out.String() != want {
		t.Errorf("output:\n%s\nwant:\n%s", out.String(), want)
	}
}

func TestErrorsDoNotStopTheRun(t *testing.T) {
	in := strings.NewReader("member ali\npay ali mina 5\nmember sara\nbalances\n")
	var out strings.Builder
	failures, err := Run(in, &out)
	if err != nil || failures != 1 {
		t.Fatalf("Run = %d failures, %v; want 1 failure", failures, err)
	}
	want := "error: line 2: unknown member: mina\n" +
		"ali   0.00\n" +
		"sara  0.00\n"
	if out.String() != want {
		t.Errorf("output:\n%s\nwant:\n%s", out.String(), want)
	}
}
