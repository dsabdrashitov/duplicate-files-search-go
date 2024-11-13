package main

import (
	"fmt"

	"github.com/dsabdrashitov/duplicate-files-search-go/pkg/cvsdb"
)

func main() {
	filename := "data/data.txt"
	fmt.Printf("Open %q\n", filename)
	f, err := cvsdb.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s, err := f.Read()
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
}
