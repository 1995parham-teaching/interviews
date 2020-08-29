package main

import "fmt"

const N = 10000

func main() {
	sum := 0

	m := New()
	d := make(chan string)

	for i := 0; i < N; i++ {
		go func(mutex Mutex, done chan string) {
			mutex.Lock()
			defer mutex.Unlock()

			sum++
			if sum == N {
				done <- "done"
			}
		}(m, d)
	}

	<-d

	fmt.Println(sum)
}

type Mutex struct {
	l chan string
}

func New() Mutex {
	return Mutex{l: make(chan string, 1)}
}

func (m Mutex) Lock() {
	m.l <- "lock"
}

func (m Mutex) Unlock() {
	<-m.l
}
