package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dsabdrashitov/duplicate-files-search-go/internal/bp"
	"github.com/dsabdrashitov/duplicate-files-search-go/internal/constants"
	"github.com/dsabdrashitov/duplicate-files-search-go/internal/subcommands"
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
	return subcommands.Init(p)
}

func subcommandScan(args []string) error {
	fs := flag.NewFlagSet("scan", flag.ExitOnError)
	var db string
	fs.StringVar(&db, "db", "", "path to db file")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if db == "" {
		var err error
		if db, err = findNearestDb(); err != nil {
			return err
		}
	}
	tail := fs.Args()
	if len(tail) == 0 {
		return fmt.Errorf("no paths to scan provided")
	}
	return subcommands.ScanPaths(db, tail)
}

func subcommandHelp(args []string) error {
	panic("unimplemented")
}

func processSubcommand(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("subcommand expected; try '%v help'", constants.DFS)
	}
	switch args[0] {
	case "init":
		return subcommandInit(args[1:])
	case "scan":
		return subcommandScan(args[1:])
	case "help":
		return subcommandHelp(args[1:])
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
