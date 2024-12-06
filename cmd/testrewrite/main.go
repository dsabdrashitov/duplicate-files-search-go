package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dsabdrashitov/duplicate-files-search-go/pkg/csvdb"
)

const (
	TESTSIZE = 100000
)

func main() {
	var time0 time.Time

	filename := "data/db2.csv"
	fmt.Printf("Open %v\n", filename)

	db, err := csvdb.New(filename, csvdb.ColumnCountValidator{Count: 2})
	if err != nil {
		panic(err)
	}
	defer db.Close()

	vals := db.Get("a")
	print("a", vals)
	var last int64
	if vals != nil {
		last, err = strconv.ParseInt(vals[1], 10, 64)
		if err != nil {
			panic(err)
		}
	} else {
		last = 0
	}

	time0 = time.Now()
	if err := db.Rewrite(); err != nil {
		panic(err)
	}
	fmt.Printf("Rewrite0 took %v\n", time.Since(time0))

	time0 = time.Now()
	for range TESTSIZE {
		last++
		if err := db.Set("a", []string{"a", fmt.Sprint(last)}); err != nil {
			panic(err)
		}
		if err := db.Service(); err != nil {
			panic(err)
		}
	}
	fmt.Printf("Set1 took %v (%v per element)\n", time.Since(time0), time.Since(time0)/TESTSIZE)

	vals = db.Get("a")
	print("a", vals)

	time0 = time.Now()
	if err := db.Rewrite(); err != nil {
		panic(err)
	}
	fmt.Printf("Rewrite1 took %v\n", time.Since(time0))

	time0 = time.Now()
	for range TESTSIZE {
		last++
		if err := db.Set("a", []string{"a", fmt.Sprint(last)}); err != nil {
			panic(err)
		}
	}
	fmt.Printf("Set2 took %v (%v per element)\n", time.Since(time0), time.Since(time0)/TESTSIZE)

	vals = db.Get("a")
	print("a", vals)

	time0 = time.Now()
	if err := db.Rewrite(); err != nil {
		panic(err)
	}
	fmt.Printf("Rewrite2 took %v\n", time.Since(time0))

	fmt.Println("done")
}

func print(k string, v []string) {
	fmt.Printf("%q = %v\n", k, v)
}
