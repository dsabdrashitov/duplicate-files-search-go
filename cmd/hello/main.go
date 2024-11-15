package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")
	line := "\"hello,\"\"abacaba\", x \"\n"
	var x string
	n, err := fmt.Sscanf(line, "%q", &x)
	fmt.Println(n, err)
	fmt.Println(x)
}
