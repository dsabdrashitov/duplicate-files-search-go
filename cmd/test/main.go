package main

import (
	"fmt"

	"github.com/dsabdrashitov/duplicate-files-search-go/pkg/csvdb"
)

func main() {
	filename := "data/db.csv"
	fmt.Printf("Open %v\n", filename)

	db, err := csvdb.New(filename, 2)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	vals := db.Get("a")
	print("a", vals)

	err = db.Set("a", []string{"v1", "v2"})
	if err != nil {
		panic(err)
	}

	vals = db.Get("a")
	print("a", vals)

	err = db.Set("b", nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("done")
}

func print(k string, v []string) {
	fmt.Printf("%q = %v\n", k, v)
}
