package subcommands

import (
	"fmt"
	"path/filepath"

	"github.com/dsabdrashitov/duplicate-files-search-go/internal/bp"
	"github.com/dsabdrashitov/duplicate-files-search-go/internal/constants"
	"github.com/dsabdrashitov/duplicate-files-search-go/internal/futil"
	"github.com/dsabdrashitov/duplicate-files-search-go/pkg/csvdb"
)

func Init(p string) error {
	p = bp.Must(filepath.Abs(p))
	if !bp.Must(futil.Exists(p)) {
		return fmt.Errorf("path '%v' not exists", p)
	}
	if !bp.Must(futil.IsDirectory(p)) {
		return fmt.Errorf("path '%v' is not directory", p)
	}
	f := filepath.Join(p, constants.DbFile)
	if bp.Must(futil.Exists(f)) {
		return fmt.Errorf("database at '%v' exists", f)
	}
	db := bp.Must(csvdb.New(f, csvdb.ColumnCountValidator{Count: constants.ColumnsCount}))
	defer db.Close()
	db.Service()
	fmt.Printf("Database at '%v' initiated.\n", f)
	return nil
}
