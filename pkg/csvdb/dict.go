package csvdb

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type dbEntry struct {
	pos  int64
	vals []string
}

type DB struct {
	fields int
	file   *os.File
	fsize  int64
	m      map[string]dbEntry
	msize  int64
}

func New(filename string, fields int) (*DB, error) {
	f, err := os.OpenFile(filename, os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	result := &DB{
		fields: fields,
		file:   f,
		m:      make(map[string]dbEntry),
		msize:  0,
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
		if len(line) != 1+1 && len(line) != 1+fields+1 {
			// ignore lines with broken format
			continue
		}
		if len(line[len(line)-1]) != 0 {
			// it's error, last column should be empty
			continue
		}
		if len(line) == 2 {
			// kill
			result.unset(line[0])
		} else {
			// add
			result.set(line[0], dbEntry{pos, line[1 : 1+fields]})
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
		if len(vals) != db.fields {
			return errors.New(fmt.Sprintf("%v fields received, %v expected", len(vals), db.fields))
		}
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

func (db *DB) writeDelete(key string) error {
	data := encode(key, nil)
	n, err := db.file.Write(data)
	db.fsize += int64(n)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) writeSet(key string, vals []string) (dbEntry, error) {
	data := encode(key, vals)
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
		db.msize -= dataSize(key, existing.vals)
	}
	db.msize += dataSize(key, entry.vals)
	db.m[key] = entry
}

func (db *DB) unset(key string) {
	entry, ok := db.m[key]
	if ok {
		db.msize -= dataSize(key, entry.vals)
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

func dataSize(key string, vals []string) int64 {
	return int64(len(encode(key, vals)))
}

func encode(key string, vals []string) []byte {
	result := make([]byte, 0)
	result = append(result, []byte(escape(key))...)
	result = append(result, []byte(",")...)
	if vals != nil {
		for _, v := range vals {
			result = append(result, []byte(escape(v))...)
			result = append(result, []byte(",")...)
		}
	}
	result = append(result, NL)
	return result
}
