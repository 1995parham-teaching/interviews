package main

import "testing"

func TestShuffle(t *testing.T) {
	cases := []struct {
		input []int
	}{
		{
			input: []int{1, 2, 3},
		},
		{
			input: []int{1, 1, 1},
		},
		{
			input: []int{1, 1, 2, 2},
		},
	}

	for _, c := range cases {
		o := Shuffle(c.input)
		m := make(map[int]int)

		for _, v := range c.input {
			m[v]++
		}

		for _, v := range o {
			m[v]--
		}

		for _, v := range m {
			if v != 0 {
				t.Fail()
			}
		}
	}
}

func TestPartition(t *testing.T) {
	cases := []struct {
		input    []int
		k        int
		expected [][]int
	}{
		{
			input:    []int{1, 2, 3},
			k:        1,
			expected: [][]int{{1, 2, 3}},
		},
		{
			input:    []int{1, 1, 1},
			k:        3,
			expected: [][]int{{1}, {1}, {1}},
		},
		{
			input:    []int{1, 1, 2, 2},
			k:        2,
			expected: [][]int{{1, 1}, {2, 2}},
		},
	}

	for _, c := range cases {
		o := Partition(c.input, c.k)

		if len(o) != len(c.expected) {
			t.Fail()
		}

		for i := range o {
			for j := range o[i] {
				if o[i][j] != c.expected[i][j] {
					t.Fail()
				}
			}
		}
	}
}
