package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dsabdrashitov/duplicate-files-search-go/internal/bp"
	"github.com/dsabdrashitov/duplicate-files-search-go/internal/constants"
	"github.com/dsabdrashitov/duplicate-files-search-go/internal/futil"
)

func findNearestDb() (string, error) {
	fmt.Printf("No db path provided. Search for nearest above.\n")
	for p, done := bp.Must(os.Getwd()), ""; p != done; p, done = filepath.Dir(p), p {
		filename := filepath.Join(p, constants.DbFile)
		if bp.Must(futil.Exists(filename)) {
			fmt.Printf("Found db at '%v'", filename)
			return filename, nil
		}
	}
	return "", fmt.Errorf("no db file found")
}
