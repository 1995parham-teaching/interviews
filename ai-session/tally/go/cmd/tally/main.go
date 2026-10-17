package main

import (
	"fmt"
	"io"
	"os"

	"tally"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "usage: tally [file]")
		os.Exit(2)
	}

	var in io.Reader = os.Stdin
	if len(os.Args) == 2 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		defer f.Close()
		in = f
	}

	failures, err := tally.Run(in, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if failures > 0 {
		os.Exit(1)
	}
}
