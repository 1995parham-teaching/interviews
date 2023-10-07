package main

import "fmt"

func main() {
	cases := []struct {
		input  string
		output int
	}{
		{
			input:  "(()",
			output: 2,
		},
		{
			input:  ")()())",
			output: 4,
		},
	}

	for _, c := range cases {
		if longestValidParentheses(c.input) != c.output {
			fmt.Printf("there output should be %d for %q\n", c.output, c.input)
		} else {
			fmt.Printf("%q passed successfully\n", c.input)
		}
	}
}

func longestValidParentheses(s string) int {
	longest_seq := 0

	for j := 0; j < len(s); j++ {
		current_seq := 0
		current_state := 0

		for i := j; i < len(s); i++ {
			current_seq += 1
			if s[i] == '(' {
				current_state += 1
			} else {
				current_state -= 1
			}
			if current_state == 0 {
				if longest_seq < current_seq {
					longest_seq = current_seq
				}
			}
			if current_state < 0 {
				current_state = 0
				current_seq = 0
			}
		}
	}

	return longest_seq
}
