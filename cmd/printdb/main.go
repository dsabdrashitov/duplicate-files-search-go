package main

import (
	"fmt"
	"io"
	"os"

	"github.com/dsabdrashitov/duplicate-files-search-go/pkg/csvdb"
)

func main() {
	filename := "data/data.txt"
	fmt.Printf("Open %q\n", filename)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	reader, err := csvdb.NewCSVReader(f)
	if err != nil {
		panic(err)
	}
	for {
		offset, a, err := reader.NextRow()
		if err == io.EOF {
			break
		}
		if err != nil {
			switch err {
			case csvdb.ErrorFormat:
				fmt.Printf("Error at %d: %v\n", offset, err)
				continue
			default:
				panic(err)
			}
		}
		fmt.Printf("At %d read line with size %d:\n", offset, len(a))
		for i, t := range a {
			fmt.Printf("%d: '%s'\n", i, t)
		}
	}
	fmt.Println("done")
}
