package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	bp "github.com/dsabdrashitov/duplicate-files-search-go/internal/boilerplate"
	"github.com/dsabdrashitov/duplicate-files-search-go/internal/futil"
	"github.com/dsabdrashitov/duplicate-files-search-go/pkg/csvdb"
)

const (
	DbFile = ".dfsdb"
	DFS    = "dfs"
)

func subcommandInit(args []string) error {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	fs.Parse(args)
	tail := fs.Args()
	var p string
	switch len(tail) {
	case 0:
		p = bp.Must(os.Getwd())
		fmt.Printf("No directory provided. Use current: %v.\n", p)
	case 1:
		p = tail[0]
	default:
		return fmt.Errorf("too many arguments (%v); DB directory expected", len(tail))
	}
	p = bp.Must(filepath.Abs(p))
	if !bp.Must(futil.Exists(p)) {
		return fmt.Errorf("path '%v' not exists", p)
	}
	if !bp.Must(futil.IsDirectory(p)) {
		return fmt.Errorf("path '%v' is not directory", p)
	}
	f := filepath.Join(p, DbFile)
	if bp.Must(futil.Exists(f)) {
		return fmt.Errorf("database at '%v' exists", f)
	}
	db := bp.Must(csvdb.New(f, csvdb.ColumnCountValidator{Count: 3}))
	db.Service()
	fmt.Printf("Database at '%v' initiated.\n", f)
	return nil
}

func subcommandScan(args []string) error {
	fmt.Println(args)
	fmt.Println("scan!!!")
	return nil
}

func subcommandHelp(args []string) error {
	panic("unimplemented")
}

func processSubcommand(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("subcommand expected; try '%v help'", DFS)
	}
	switch args[0] {
	case "init":
		return subcommandInit(args[1:])
	case "scan":
		return subcommandScan(os.Args[1:])
	case "help":
		return subcommandHelp(os.Args[1:])
	default:
		return fmt.Errorf("unknown subcommand '%v'", args[0])
	}
}

func main() {
	if err := processSubcommand(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
