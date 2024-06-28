/**
"""
Implement a music shuffling data structure with a Shuffle() method

Inputs:
1. List of N integers representing songs
2. Integer K < N, representing the cooldown for a song after each shuffle

The Shuffle method should randomly pick a song from the songs list.
After a song is picked, it should not be picked again for the next K shuffles.


Example:
[1, 2, 3, 4, 5]
2

> s = Shuffler([1,2,3,4,5], 2)
> s.shuffle() -> 1
> s.shuffle() -> 3
> s.shuffle() -> 2
> s.shuffle() -> 1
"""


**/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	s := NewShuffler([]int{1, 2, 3, 4, 5}, 2)
	fmt.Println(s.Shuffle())
	fmt.Println(s.Shuffle())
	fmt.Println(s.Shuffle())
	fmt.Println(s.Shuffle())
	fmt.Println(s.Shuffle())
	fmt.Println(s.Shuffle())
}

type Shuffler struct {
	songs  []int
	picked []int
	k      int
	rnd    *rand.Rand
}

func NewShuffler(songs []int, k int) *Shuffler {
	return &Shuffler{
		songs:  songs,
		picked: make([]int, len(songs)),
		k:      k,
		rnd:    rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (s *Shuffler) Shuffle() int {
	c := s.rnd.Intn(len(s.songs))

	for s.picked[c] > 0 {
		c = (c + 1) % len(s.songs)
	}

	for i := 0; i < len(s.songs); i++ {
		if s.picked[i] > 0 {
			s.picked[i]--
		}
	}

	s.picked[c] = s.k

	return s.songs[c]
}
