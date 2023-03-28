package main

import "testing"

func Test(t *testing.T) {
	cases := []struct {
		input  string
		output string
	}{
		{
			input:  "babad",
			output: "bab",
		},
		{
			input:  "cbbd",
			output: "bb",
		},
	}

	for _, c := range cases {
		o := longestpalindrome(c.input)
		if o != c.output {
			t.Fail()
		}
	}
}
