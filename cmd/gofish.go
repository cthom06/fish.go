package main

import (
	"os"
	"io/ioutil"
	"github.com/cthom06/GoFish/fish"
	"bufio"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		os.Stderr.WriteString(os.Args[0] + ": no input file\n")
		os.Exit(1)
	}
	if f, e := os.Open(os.Args[1]); e == nil {
		if d, e := ioutil.ReadAll(f); e != nil {
			os.Stderr.WriteString(os.Args[0] + ": " + e.String() + "\n")
			os.Exit(2)
		} else {
			r := fish.NewRuntime(d)
			sargs := os.Args[2:]
			for i := 0; i < len(sargs); i++ {
				switch sargs[i] {
				case "-s", "--string":
					i++
					if i == len(sargs) {
						os.Stderr.WriteString(os.Args[0] + ": invalid arguments\n")
						os.Exit(4)
					}
					for _, v := range sargs[i] {
						r.Push(int(v))
					}
					for i < len(sargs) - 1 && sargs[i+1] != "-s" && sargs[i+1] != "--string" && sargs[i+1] != "-v" && sargs[i+1] != "--value" {
						i++
						for _, v := range sargs[i] {
							r.Push(int(v))
						}
					}
				case "-v", "--value":
					i++
					if i == len(sargs) {
						os.Stderr.WriteString(os.Args[0] + ": invalid arguments\n")
						os.Exit(4)
					}
					if v, e := strconv.Atoi(sargs[i]); e == nil {
						r.Push(v)
					} else {
						os.Stderr.WriteString(os.Args[0] + ": invalid arguments\n")
						os.Exit(4)
					}
					for i < len(sargs) - 1 && sargs[i+1] != "-s" && sargs[i+1] != "--string" && sargs[i+1] != "-v" && sargs[i+1] != "--value" {
						i++
						if v, e := strconv.Atoi(sargs[i]); e == nil {
							r.Push(v)
						} else {
							os.Stderr.WriteString(os.Args[0] + ": invalid arguments\n")
							os.Exit(4)
						}
					}
				default:
					os.Stderr.WriteString(os.Args[0] + ": invalid arguments\n")
				}
			}
			b := bufio.NewWriter(os.Stdout)
			if e := r.Run(os.Stdin, b); e != fish.NoError {
				b.Flush()
				os.Stderr.WriteString(e.String() + "\n")
				os.Exit(3)
			}
			b.Flush()
		}
	} else {
		os.Stderr.WriteString(os.Args[0] + ": " + e.String() + "\n")
		os.Exit(2)
	}
}
