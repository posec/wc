package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var cFlag = flag.Bool("c", false, "count bytes")
var lFlag = flag.Bool("l", false, "count lines")
var wFlag = flag.Bool("w", false, "count words")

var exitStatus int

func main() {

	flag.Parse()
	if !*cFlag && !*wFlag && !*lFlag {
		*cFlag = true
		*wFlag = true
		*lFlag = true
	}

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"-"}
	}
	for _, f := range args {
		wc(f)
	}
	os.Exit(exitStatus)
}

func wc(filename string) {
	var err error
	var in *os.File

	if filename == "-" {
		// XCU7: Says what about "wc -"?
		in = os.Stdin
	} else {
		in, err = os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			exitStatus = 2
			return
		}
		defer in.Close()
	}

	var nBytes, nWords, nLines int
	inWord := false
	buffer := make([]byte, 16384)
	for {
		count, err := in.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			exitStatus = 2
			return
		}

		nBytes += count
		for i := 0; i < count; i += 1 {
			if buffer[i] == ' ' ||
				buffer[i] == '\t' ||
				buffer[i] == '\n' {
				inWord = false
			} else {
				if !inWord {
					nWords += 1
				}
				inWord = true
			}

			if buffer[i] == '\n' {
				nLines += 1
			}
		}
	}

	padSpace := ""
	if *lFlag {
		fmt.Printf("%s%d", padSpace, nLines)
		padSpace = " "
	}
	if *wFlag {
		fmt.Printf("%s%d", padSpace, nWords)
		padSpace = " "
	}
	if *cFlag {
		fmt.Printf("%s%d", padSpace, nBytes)
		padSpace = " "
	}
	fmt.Printf("\n")
}
