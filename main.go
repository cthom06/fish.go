package main

import (
	"os"
	"io/ioutil"
	"./fish"
)

func main() {
	if len(os.Args) < 2 {
		os.Stderr.WriteString("No file\n")
		os.Exit(1)
	}
	if f, e := os.Open(os.Args[1]); e == nil {
		if d, e := ioutil.ReadAll(f); e != nil {
			os.Stderr.WriteString(e.String() + "\n")
			os.Exit(2)
		} else {
			r := fish.NewRuntime(d)
			if e := r.Run(os.Stdin, os.Stdout, nil); e != fish.NoError {
				os.Stdout.WriteString(e.String() + "\n")
				os.Exit(3)
			} else os.Exit(0)
		}
	} else {
		os.Stderr.WriteString(e.String() + "\n")
		os.Exit(2)
	}
}
