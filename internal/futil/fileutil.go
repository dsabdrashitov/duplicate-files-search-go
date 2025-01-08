package futil

import (
	"os"

	bp "github.com/dsabdrashitov/duplicate-files-search-go/internal/boilerplate"
)

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func IsDirectory(path string) bool {
	fileInfo := bp.Must(os.Stat(path))
	switch mode := fileInfo.Mode(); {
	case mode.IsDir():
		return true
	default:
		return false
	}
}

func IsFile(path string) bool {
	fileInfo := bp.Must(os.Stat(path))
	switch mode := fileInfo.Mode(); {
	case mode.IsRegular():
		return true
	default:
		return false
	}
}
