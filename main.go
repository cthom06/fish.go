package main

import (
	"os"
	"./fish"
)

func main() {
	if len(os.Args) < 2 {
		os.Stderr.WriteString("No file\n")
		os.Exit(1)
	}
	if f, e := os.Open(os.Args[1]); e == nil {
		r := fish.NewRuntime(f)
		r.Run(os.Stdin, os.Stdout, false)
		os.Exit(3)
	} else {
		os.Stderr.WriteString(e.String() + "\n")
		os.Exit(2)
	}
}
