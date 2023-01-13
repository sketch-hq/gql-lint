package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/sketch-hq/gql-lint/output"
	"github.com/sketch-hq/gql-lint/parser"
)

var queriesDir string
var schemaFile string
var outputFormat string

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "expected 'deprecation' subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "deprecation":
		depFlags := flag.NewFlagSet("deprecation", flag.ExitOnError)
		depFlags.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage: %s <query/mutation directory>\n", os.Args[0])
			depFlags.PrintDefaults()
		}
		depFlags.StringVar(&schemaFile, "schema", "", "Server's schema file")
		depFlags.StringVar(&outputFormat, "output", "stdout", "Output format. Choose between json and stdout. Defaults is stdout.")
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

	case "diff":
		diffFlags := flag.NewFlagSet("diff", flag.ExitOnError)
		diffFlags.StringVar(&outputFormat, "output", "stdout", "Output format. Choose between json and stdout. Defaults is stdout.")
		diffFlags.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage: %s <json file> <json file>\n", os.Args[0])
			diffFlags.PrintDefaults()
		}
		diffFlags.Parse(os.Args[2:])
		if len(diffFlags.Args()) < 2 {
			fmt.Fprint(os.Stderr, "expected two json files to be given")
		}
		runDiff(diffFlags.Arg(0), diffFlags.Arg(1))

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

	switch outputFormat {
	case "stdout":
		deprecationStdOut(queryFields)

	case "json":
		deprecationJsonOut(queryFields)

	default:
		fmt.Fprintf(os.Stderr, "%s is not a valid output format. Choose between json and stdout", outputFormat)
		os.Exit(1)
	}
}

func runDiff(fileA string, fileB string) {
	result, err := output.CompareFiles(fileA, fileB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to diff: %s", err)
		os.Exit(1)
	}

	switch outputFormat {
	case "stdout":
		diffStdOut(fileA, fileB, result)

	case "json":
		diffJsonOut(result)

	default:
		fmt.Fprintf(os.Stderr, "%s is not a valid output format. Choose between json and stdout\n", outputFormat)
		os.Exit(1)
	}
}

func deprecationStdOut(queryFields parser.QueryFieldList) {
	for _, q := range queryFields {
		fmt.Printf("%s is deprecated\n", q.Path)
		fmt.Printf("  File:   %s:%d\n", q.File, q.Line)
		fmt.Printf("  Reason: %s\n", q.DeprecationReason)
		fmt.Println()
	}
}

func deprecationJsonOut(queryFields parser.QueryFieldList) {
	out := output.Data{}

	for _, q := range queryFields {
		f := output.Field{
			Field:             q.Path,
			File:              q.File,
			Line:              q.Line,
			DeprecationReason: q.DeprecationReason,
		}
		out = append(out, f)
	}
	bytes, err := json.Marshal(out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to encode json: %s\n", err)
		os.Exit(1)
	}

	fmt.Print(string(bytes))
}

func diffStdOut(_ string, fileB string, out output.Data) {
	if len(out) == 0 {
		return
	}

	for _, f := range out {
		fmt.Printf("%s (%s)\n", f.Field, f.DeprecationReason)
		fmt.Printf("  %s:%d\n", f.File, f.Line)
	}
}

func diffJsonOut(out output.Data) {
	bytes, err := json.Marshal(out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to encode json: %s\n", err)
		os.Exit(1)
	}

	fmt.Print(string(bytes))
}
