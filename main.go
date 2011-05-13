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
		if e := r.Run(os.Stdin, os.Stdout, nil); e != fish.NoError {
			os.Stdout.WriteString(e.(string) + "\n")
			os.Exit(3)
		} else os.Exit(0)
	} else {
		os.Stderr.WriteString(e.String() + "\n")
		os.Exit(2)
	}
}
