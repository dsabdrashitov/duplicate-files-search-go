package csvdb

import (
	"fmt"
	"io"
	"os"
)

type BufferedReader struct {
	file *os.File
	buf  []byte
	pos  int64
	done int64
	size int64
}

func NewBufferedReader(file *os.File) (*BufferedReader, error) {
	pos, err := offset(file)
	if err != nil {
		return nil, err
	}
	result := &BufferedReader{
		file: file,
		buf:  make([]byte, fileBufferSize),
		pos:  pos,
		done: 0,
		size: 0,
	}
	return result, nil
}

func (br *BufferedReader) Read() (byte, error) {
	if br.done < br.size {
		result := br.buf[br.done]
		br.done++
		return result, nil
	}
	if br.size == int64(len(br.buf)) {
		pos, err := offset(br.file)
		if err != nil {
			return 0, err
		}
		if pos != br.pos+int64(len(br.buf)) {
			return 0, IOError{fmt.Sprintf("Wrong file offset: %d. Expected: %d.", pos, br.pos+int64(len(br.buf)))}
		}
		br.pos = pos
		br.done = 0
		br.size = 0
	}
	r, err := br.file.Read(br.buf[br.size:])
	br.size = br.size + int64(r)
	switch {
	case br.done < br.size:
		result := br.buf[br.done]
		br.done++
		return result, nil
	case err != nil:
		return 0, err
	default:
		return 0, IOError{"File refused to read without error"}
	}
}

func (br *BufferedReader) ReadLine() (string, error) {
	buf := make([]byte, 0)
	for {
		b, err := br.Read()
		if err != nil {
			return string(buf), err
		}
		buf = append(buf, b)
		if b == 0x0A {
			buf = buf[:len(buf)-1]
			if len(buf) > 0 && buf[len(buf)-1] == 0x0D {
				buf = buf[:len(buf)-1]
			}
			return string(buf), nil
		}
	}
}

func (br *BufferedReader) Offset() int64 {
	return br.pos + br.done
}

func offset(file *os.File) (int64, error) {
	return file.Seek(0, io.SeekCurrent)
}
