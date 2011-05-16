package main

import (
	"os"
	"io/ioutil"
	"github.com/cthom06/GoFish/fish"
	"bufio"
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
			b := bufio.NewWriter(os.Stdout)
			if e := r.Run(os.Stdin, b); e != fish.NoError {
				b.Flush()
				os.Stderr.WriteString(e.String() + "\n")
				os.Exit(3)
			}
			b.Flush()
		}
	} else {
		os.Stderr.WriteString(e.String() + "\n")
		os.Exit(2)
	}
}
