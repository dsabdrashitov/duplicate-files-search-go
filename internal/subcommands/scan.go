package subcommands

import (
	"fmt"

	"github.com/dsabdrashitov/duplicate-files-search-go/internal/bp"
	"github.com/dsabdrashitov/duplicate-files-search-go/internal/constants"
	"github.com/dsabdrashitov/duplicate-files-search-go/internal/futil"
	"github.com/dsabdrashitov/duplicate-files-search-go/pkg/csvdb"
)

func ScanPaths(dbPath string, paths []string) error {
	if !bp.Must(futil.Exists(dbPath)) {
		return fmt.Errorf("database '%v' not exists", dbPath)
	}
	if !bp.Must(futil.IsFile(dbPath)) {
		return fmt.Errorf("'%v' is not file", dbPath)
	}
	db := bp.Must(csvdb.New(dbPath, csvdb.ColumnCountValidator{Count: constants.ColumnsCount}))
	defer db.Close()
	for _, path := range paths {
		fmt.Printf("Scan '%v' into '%v'\n", path, dbPath)
	}
	return nil
}

func Scan(dbPath string, path string) error {
	return ScanPaths(dbPath, []string{path})
}
