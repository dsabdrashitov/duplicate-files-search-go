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

func (r *CSVReader) nextToken() (empty bool, t string, eol bool, e error) {
	result := make([]byte, 0)
	quoted := false
	quote := false
	broken := false
	qcr := false
	for {
		c, err := r.br.Read()
		if err != nil {
			switch {
			case err != io.EOF:
				return false, "", false, err
			case broken:
				return false, "", false, ErrorFormat
			case qcr:
				return false, "", false, ErrorFormat
			case quote:
				return false, string(result), true, nil
			case quoted:
				return false, "", false, ErrorFormat
			case len(result) > 0:
				return false, string(result), true, nil
			default:
				return false, "", false, io.EOF
			}
		}
		if broken {
			// in broken state skip all till new line
			if c == NL {
				return false, "", false, ErrorFormat
			} else {
				continue
			}
		}
		if c == '"' && !quoted && len(result) == 0 {
			quoted = true
			continue
		}
		if quoted {
			if quote {
				switch {
				case c == NL:
					return false, string(result), true, nil
				case qcr:
					broken = true
					continue
				case c == CR:
					qcr = true
					continue
				case c == ',':
					return false, string(result), false, nil
				case c == '"':
					result = append(result, c)
					quote = false
					continue
				default:
					broken = true
					continue
				}
			} else {
				if c == '"' {
					quote = true
					continue
				} else {
					result = append(result, c)
					continue
				}
			}
		} else {
			switch {
			case c == ',':
				return false, string(result), false, nil
			case c == NL:
				if len(result) > 0 && result[len(result)-1] == CR {
					result = result[:len(result)-1]
				}
				if len(result) == 0 {
					return true, "", true, nil
				} else {
					return false, string(result), true, nil
				}
			case c == '"':
				broken = true
				continue
			default:
				result = append(result, c)
				continue
			}
		}
	}
}

func (r *CSVReader) NextRow() (int64, []string, error) {
	result := make([]string, 0)
	var offset int64
	for {
		if len(result) == 0 {
			offset = r.br.Offset()
		}
		empty, t, eol, err := r.nextToken()
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
