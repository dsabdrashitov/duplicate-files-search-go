package csvdb

import (
	"io"
	"os"
	"slices"
)

const (
	addCommand    = ""
	deleteCommand = "-"
)

type dbEntry struct {
	pos  int64
	vals []string
}

type RowValidator interface {
	Sum(k string, v []string) string
	Check(k string, v []string) bool
}

type ColumnCountValidator struct {
	Count int
}

func (validator ColumnCountValidator) Check(k string, v []string) bool {
	return len(v) == validator.Count
}

func (validator ColumnCountValidator) Sum(k string, v []string) string {
	return ""
}

type DB struct {
	file      *os.File
	validator RowValidator
	fsize     int64
	m         map[string]dbEntry
	msize     int64
}

func New(filename string, validator RowValidator) (*DB, error) {
	f, err := os.OpenFile(filename, os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	result := &DB{
		file:      f,
		validator: validator,
		m:         make(map[string]dbEntry),
		msize:     0,
	}
	reader, err := NewCSVReader(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	for {
		pos, line, err := reader.NextRow()
		if err == io.EOF {
			break
		}
		if err == ErrorFormat {
			continue
		}
		if err != nil {
			f.Close()
			return nil, err
		}
		if len(line) < 3 {
			// ignore lines with broken format
			continue
		}
		switch line[1] {
		case addCommand:
			k, v := line[0], line[2:len(line)-1]
			if validator.Check(k, v) && validator.Sum(k, v) == line[len(line)-1] {
				// add
				result.set(line[0], dbEntry{pos, line[2 : len(line)-1]})
			}
		case deleteCommand:
			if len(line) == 3 && validator.Sum(line[0], nil) == line[2] {
				// delete
				result.unset(line[0])
			}
		default:
			// ignore lines with broken format
		}
	}
	result.fsize = reader.br.Offset()
	return result, nil
}

func (db *DB) Close() error {
	return db.file.Close()
}

func (db *DB) Get(key string) []string {
	entry, ok := db.m[key]
	if ok {
		return entry.vals
	} else {
		return nil
	}
}

func (db *DB) Set(key string, vals []string) error {
	if vals == nil {
		// delete
		_, ok := db.m[key]
		if ok {
			err := db.writeDelete(key)
			if err != nil {
				return err
			}
			db.unset(key)
		}
	} else {
		// set
		existing, ok := db.m[key]
		if ok && unchanged(vals, existing.vals) {
			return nil
		}
		entry, err := db.writeSet(key, vals)
		if err != nil {
			return err
		}
		db.set(key, entry)
	}
	return nil
}

func (db *DB) Keys() []string {
	result := make([]string, len(db.m))
	i := 0
	for k := range db.m {
		result[i] = k
		i += 1
	}
	return result
}

func (db *DB) Rewrite() error {
	if db.fsize < db.msize {
		return ErrorFileTooSmall
	}
	eline := []byte(ErroneousRow)
	db.file.Write(eline)
	db.fsize = db.fsize + int64(len(eline))
	for k, e := range db.m {
		_, err := db.writeSet(k, e.vals)
		if err != nil {
			return err
		}
	}
	if _, err := db.file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	db.fsize = 0
	keys := db.Keys()
	slices.Sort(keys)
	for _, k := range keys {
		vals := db.m[k].vals
		db.unset(k)
		if err := db.Set(k, vals); err != nil {
			return err
		}
	}
	if err := db.file.Truncate(db.fsize); err != nil {
		return err
	}
	return nil
}

func (db *DB) Service() error {
	if db.fsize > serviceMultiplier*db.msize {
		return db.Rewrite()
	} else {
		return nil
	}
}

func (db *DB) Flush() error {
	return db.file.Sync()
}

func (db *DB) writeDelete(key string) error {
	data := encode(key, deleteCommand, nil, db.validator.Sum(key, nil))
	n, err := db.file.Write(data)
	db.fsize += int64(n)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) writeSet(key string, vals []string) (dbEntry, error) {
	data := encode(key, addCommand, vals, db.validator.Sum(key, vals))
	cur, err := offset(db.file)
	if err != nil {
		return dbEntry{}, err
	}
	if cur != db.fsize {
		return dbEntry{}, ErrorOffset{actual: cur, expected: db.fsize}
	}
	n, err := db.file.Write(data)
	db.fsize += int64(n)
	if err != nil {
		return dbEntry{}, err
	}
	return dbEntry{cur, vals}, nil
}

func (db *DB) set(key string, entry dbEntry) {
	existing, ok := db.m[key]
	if ok {
		db.msize -= dataSize(key, existing.vals, db.validator.Sum(key, existing.vals))
	}
	db.msize += dataSize(key, entry.vals, db.validator.Sum(key, entry.vals))
	db.m[key] = entry
}

func (db *DB) unset(key string) {
	entry, ok := db.m[key]
	if ok {
		db.msize -= dataSize(key, entry.vals, db.validator.Sum(key, entry.vals))
		delete(db.m, key)
	}
}

func unchanged(updated []string, existing []string) bool {
	if len(updated) != len(existing) {
		return false
	}
	for i, s := range updated {
		if s != existing[i] {
			return false
		}
	}
	return true
}

func dataSize(key string, vals []string, sum string) int64 {
	return int64(len(encode(key, addCommand, vals, sum)))
}

func encode(key string, command string, vals []string, sum string) []byte {
	result := make([]byte, 0)
	result = append(result, []byte(escape(key))...)
	result = append(result, []byte(","+command+",")...)
	for _, v := range vals {
		result = append(result, []byte(escape(v))...)
		result = append(result, []byte(",")...)
	}
	result = append(result, []byte(escape(sum))...)
	result = append(result, NL)
	return result
}
