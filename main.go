package main

import (
	"fmt"
	"io"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Hello, world!")
	f, err := os.OpenFile("data/data.txt", os.O_CREATE, 0777)
	check(err)
	defer func() { check(f.Close()) }()

	buf := make([]byte, 1)
	b := make([]byte, 0)
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		check(err)
		b = append(b, buf[:n]...)
	}

	fmt.Println("Read:\n" + string(b))

	n, err := f.Write([]byte("ok"))
	check(err)
	fmt.Println(n)
}
