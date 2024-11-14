package main

import (
	"fmt"
	"io"
	"os"

	"github.com/dsabdrashitov/duplicate-files-search-go/pkg/csvdb"
)

func trim(s string) string {
	if len(s) > 0 && s[len(s)-1] == 0x0A {
		s = s[:len(s)-1]
	}
	if len(s) > 0 && s[len(s)-1] == 0x0D {
		s = s[:len(s)-1]
	}
	return s
}

func main() {
	filename := "data/data.txt"
	fmt.Printf("Open %q\n", filename)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	reader, err := csvdb.NewBufferedReader(f)
	if err != nil {
		panic(err)
	}
	for {
		offset0 := reader.Offset()
		line, err := reader.ReadLine()
		offset1 := reader.Offset()
		fmt.Printf("Read line: '%s' from %d to %d\n", trim(line), offset0, offset1)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("done")
}
