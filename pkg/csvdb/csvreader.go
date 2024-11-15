package csvdb

import (
	"io"
	"os"
)

const (
	CR = 0x0D
	NL = 0x0A
)

type CSVReader struct {
	br *BufferedReader
}

func NewCSVReader(file *os.File) (*CSVReader, error) {
	br, err := NewBufferedReader(file)
	if err != nil {
		return nil, err
	}
	result := &CSVReader{
		br: br,
	}
	return result, nil
}

func (r *CSVReader) ReadRow() (int64, []string, error) {
	result := make([]string, 0)
	var offset int64
	for {
		if len(result) == 0 {
			offset = r.br.Offset()
		}
		empty, t, eol, err := r.readNext()
		if err != nil {
			return offset, nil, err
		}
		if eol {
			if len(result) == 0 && empty {
				continue
			} else {
				result = append(result, t)
				return offset, result, nil
			}
		}
		result = append(result, t)
	}
}

func (r *CSVReader) readNext() (empty bool, token string, eol bool, e error) {
	buf := make([]byte, 0)
	b, err := r.br.Read()
	if err != nil {
		return false, "", false, err
	}
	if b == '"' {
		skippedCR := false
		quote := false
		for {
			b, err = r.br.Read()
			if err == io.EOF {
				if quote {
					return false, "", false, ErrorNoNewline
				} else {
					return false, "", false, ErrorNoQuote
				}
			}
			if err != nil {
				return false, "", false, err
			}

			if b == '"' {
				if quote {
					if skippedCR {
						return false, "", false, ErrorNoNewline
					}
					buf = append(buf, b)
					quote = false
				} else {
					quote = true
				}
			} else {
				if quote {
					if b == NL {
						return false, string(buf), true, nil
					}
					if skippedCR {
						return false, "", false, ErrorNoNewline
					}
					if b == CR {
						skippedCR = true
					} else if b == ',' {
						return false, string(buf), false, nil
					} else {
						return false, "", false, ErrorNoComma
					}
				} else {
					buf = append(buf, b)
				}
			}
		}
	} else {
		for {
			if b == ',' {
				return len(buf) == 0, string(buf), false, nil
			}
			if b == 0x0A {
				if len(buf) > 0 && buf[len(buf)-1] == 0x0D {
					buf = buf[:len(buf)-1]
				}
				return len(buf) == 0, string(buf), true, nil
			}
			buf = append(buf, b)

			b, err = r.br.Read()
			if err == io.EOF {
				return false, "", false, ErrorNoNewline
			}
			if err != nil {
				return false, "", false, err
			}
		}
	}
}
