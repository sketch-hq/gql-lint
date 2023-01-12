package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sketch-hq/gql-lint/parser"
)

var queriesDir string
var schemaFile string

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <query/mutation directory>\n", os.Args[0])
		flag.PrintDefaults()
	}

	depFlags := flag.NewFlagSet("deprecation", flag.ExitOnError)
	depFlags.StringVar(&schemaFile, "schema", "", "server's schema file")

	if len(os.Args) < 2 {
		fmt.Println("expected 'deprecation' subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "deprecation":
		depFlags.Parse(os.Args[2:])
		queriesDir = depFlags.Arg(0)
		if queriesDir == "" {
			fmt.Fprint(os.Stderr, "You must specify a directory for queries and mutations\n")
			os.Exit(1)
		}
		if schemaFile == "" {
			fmt.Fprint(os.Stderr, "You must specify a schema file\n")
			os.Exit(1)
		}
		runDeprecation()
	default:
		fmt.Println("expected 'deprecation' subcommand")
		os.Exit(1)
	}
}

func runDeprecation() {
	schema, err := parser.ParseSchemaFile(schemaFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse schema file %s: %s", schemaFile, err)
		os.Exit(1)
	}

	queryFields, err := parser.ParseQueryDir(queriesDir, schema)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse files in %s: %s", queriesDir, err)
		os.Exit(1)
	}

	for _, q := range queryFields {
		fmt.Printf("%s is deprecated\n", q.Path)
		fmt.Printf("  File:   %s:%d\n", q.File, q.Line)
		fmt.Printf("  Reason: %s\n", q.DeprecationReason)
		fmt.Println()
	}
}
