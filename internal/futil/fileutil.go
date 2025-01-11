package futil

import (
	"os"
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

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	switch mode := fileInfo.Mode(); {
	case mode.IsDir():
		return true, nil
	default:
		return false, nil
	}
}

func IsFile(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	switch mode := fileInfo.Mode(); {
	case mode.IsRegular():
		return true, nil
	default:
		return false, nil
	}
}
