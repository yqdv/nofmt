package main

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
)

func formatSrc(src []byte) []byte {
	fset := token.NewFileSet()
	srcAst, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	var out bytes.Buffer
	printer.Fprint(&out, fset, srcAst)
	return out.Bytes()
}

func main() {
	var err error
	var isStdin bool
	var isStdout bool
	var fpaths []string

	if len(os.Args) == 1 {
		isStdin = true
		isStdout = true
	} else if len(os.Args) == 2 {
		isStdin = false
		isStdout = true
		fpaths = []string{os.Args[1]}
	} else if len(os.Args) >= 3 && os.Args[1] == "-w" {
		isStdin = false
		isStdout = false
		fpaths = os.Args[2:]
	} else {
		fmt.Printf("Usage: %s -w myfile.go\n", os.Args[0])
		os.Exit(1)
	}

	for _, fpath := range fpaths {
		var src []byte
		if isStdin {
			src, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				panic(err)
			}
		} else {
			src, err = ioutil.ReadFile(fpath)
			if err != nil {
				panic(err)
			}
		}

		src = formatSrc(src)

		if isStdout {
			os.Stdout.Write(src)
		} else {
			err = ioutil.WriteFile(fpath, src, 0o644)
			if err != nil {
				panic(err)
			}
		}
	}
}
