package main

import (
	"fmt"
)

func main() {
	report()
}

func report() {
	var arr [5]int

	arr[1] = 1378
	fmt.Printf("arr[1]: %d\n", arr[1])

	slice := arr[1:3]
	slice[0] = 1373

	fmt.Printf("arr[1]: %d\n", arr[1])

	//arr[5] = 3
	//slice[2] = 3

	fmt.Printf("len: %d\n", len(slice))
	fmt.Printf("cap: %d\n", cap(slice))
	slice = append(slice, 3)
	fmt.Printf("len: %d\n", len(slice))
	fmt.Printf("cap: %d\n", cap(slice))

	slice = arr[:] // Works fine
	//arr = slice[0 : 5] Compile time error

	fmt.Printf("slice[0]: %d\n", slice[0])
}
