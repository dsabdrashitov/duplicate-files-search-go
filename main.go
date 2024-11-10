package main

import (
	"fmt"
	"io"
	"os"
	"strings"
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

	content := string(b)
	fmt.Println("Read:\n" + content)

	_, err = f.Seek(1, 0)
	check(err)
	n, err := f.Write([]byte("ok\nk"))
	check(err)
	fmt.Println(n)

	curPos := strings.Index(content, "$")
	if curPos != -1 {
		fmt.Printf("Cut %d bytes at the end of file.\n", len(content)-curPos)
		check(f.Truncate(int64(curPos)))
	} else {
		fmt.Println("No @ in file, so keep it intact.")
	}
}
