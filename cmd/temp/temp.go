package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dsabdrashitov/duplicate-files-search-go/internal/bp"
)

func main() {
	for p, done := bp.Must(os.Getwd()), ""; p != done; p, done = filepath.Dir(p), p {
		fmt.Println(p)
	}
}
