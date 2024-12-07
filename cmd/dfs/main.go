package main

import (
	"fmt"
	"os"
	"path/filepath"

	BP "github.com/dsabdrashitov/duplicate-files-search-go/internal/boilerplate"
	"github.com/dsabdrashitov/duplicate-files-search-go/internal/futil"
	"github.com/dsabdrashitov/duplicate-files-search-go/pkg/csvdb"
	"github.com/urfave/cli"
)

const (
	DbFile = ".dfsdb"
)

func initCommand(c *cli.Context) error {
	if len(c.Args().Tail()) != 0 {
		return fmt.Errorf("%v extra arguments", len(c.Args().Tail()))
	}
	var p string
	if c.Args().Present() {
		p = c.Args().First()
	} else {
		p = "."
	}
	p = BP.Must(filepath.Abs(p))
	if !BP.Must(futil.Exists(p)) {
		return fmt.Errorf("path '%v' not exists", p)
	}
	if !futil.IsDirectory(p) {
		return fmt.Errorf("path '%v' is not directory", p)
	}
	f := filepath.Join(p, DbFile)
	if BP.Must(futil.Exists(f)) {
		return fmt.Errorf("database at '%v' exists", f)
	}
	db := BP.Must(csvdb.New(f, csvdb.ColumnCountValidator{Count: 3}))
	db.Service()
	fmt.Printf("Database at '%v' initiated.\n", f)
	return nil
}

func scanCommand(c *cli.Context) error {
	fmt.Println(c.Args())
	fmt.Println("scan!!!")
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "dfs"
	app.Description = "Tool for collecting hashes of files."
	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "Create new empty database, if none.",
			Action: initCommand,
		},
		{
			Name:   "scan",
			Usage:  "Scan specified path.",
			Action: scanCommand,
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
