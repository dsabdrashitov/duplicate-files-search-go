package cvsdb

import (
	"io"
	"os"
)

type File struct {
	*os.File
}

func Open(filename string) (*File, error) {
	file, err := os.OpenFile(filename, os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	var this *File = &File{file}
	return this, nil
}

var readBuffer []byte = make([]byte, 1024)

func (this *File) Next() (string, error) {
	result := make([]byte, 1)
	for {
		n, err := this.Read(readBuffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return "", err
			}
		}
		result = append(result, readBuffer[:n]...)
	}
	return string(result), nil
}