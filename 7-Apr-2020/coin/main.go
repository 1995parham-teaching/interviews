package main

import (
	"fmt"
	"math"
	"time"
)

// Coins returns all available types of coins
func Coins() []int {
	return []int{1, 5, 7, 10}
}

const one = 1

// MinimumCoinsProblem solution with memorization
type MinimumCoinsProblem struct {
	Table map[int]int
	Calls int
}

func main() {
	amount := 23

	problem := MinimumCoinsProblem{
		Table: make(map[int]int),
	}

	start := time.Now()

	fmt.Printf("%d needs %d coins\n", amount, problem.Minimum(amount, 0))

	fmt.Printf("Execution Time: %v\n", time.Since(start))

	fmt.Printf("Number of Calls: %d\n", problem.Calls)
}

// Minimum number of coins for n amount of money
func (p *MinimumCoinsProblem) Minimum(n int, number int) int {
	p.Calls++

	if n == 0 {
		return number
	}

	if val, ok := p.Table[n]; ok {
		return number + val
	}

	min := math.MaxInt32

	for _, v := range Coins() {
		if n < v {
			break
		}

		r := p.Minimum(n-v, number+one)

		if r < min {
			min = r
		}
	}

	p.Table[n] = min - number

	return min
}
