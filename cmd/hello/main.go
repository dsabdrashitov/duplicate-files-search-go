package main

import (
	"fmt"
	"os"

	"github.com/dsabdrashitov/duplicate-files-search-go/pkg/csvdb"
)

func main() {
	fmt.Println("Hello, world!")
	filename := "data/data.txt"
	fmt.Printf("Open %q\n", filename)
	f, err := os.OpenFile(filename, os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.Seek(64, 0)
	if err != nil {
		panic(err)
	}
	_, err = f.Write([]byte(csvdb.ErroneousRow))
	if err != nil {
		panic(err)
	}
	fmt.Println("done")
}
