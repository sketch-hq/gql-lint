package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/sketch-hq/gql-lint/output"
	"github.com/sketch-hq/gql-lint/parser"
)

var (
	queriesDir   string
	schemaFile   string
	outputFormat string
	depFlags     *flag.FlagSet
	diffFlags    *flag.FlagSet
)

func main() {
	os.Exit(run())
}

func run() int {
	depFlags = flag.NewFlagSet("deprecation", flag.ExitOnError)
	diffFlags = flag.NewFlagSet("diff", flag.ExitOnError)

	depFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s deprecation [<args] <queries directory>\n", path.Base(os.Args[0]))
		depFlags.PrintDefaults()
	}
	depFlags.StringVar(&schemaFile, "schema", "", "Server's schema file")
	depFlags.StringVar(&outputFormat, "output", "stdout", "Output format. Choose between json and stdout. Defaults is stdout.")

	diffFlags.StringVar(&outputFormat, "output", "stdout", "Output format. Choose between json and stdout. Defaults is stdout.")
	diffFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s diff [<args>] <json file> <json file>\n", path.Base(os.Args[0]))
		diffFlags.PrintDefaults()
	}

	if len(os.Args) < 2 {
		help("")
		return 0
	}

	switch os.Args[1] {
	case "deprecation":
		depFlags.Parse(os.Args[2:])
		queriesDir = depFlags.Arg(0)
		if queriesDir == "" {
			fmt.Fprint(os.Stderr, "You must specify a directory for queries and mutations\n\n")
			help(os.Args[1])
			return 1
		}
		if schemaFile == "" {
			fmt.Fprint(os.Stderr, "You must specify a schema file using -schema\n\n")
			help(os.Args[1])
			return 1
		}
		return runDeprecation()

	case "diff":
		diffFlags.Parse(os.Args[2:])
		if len(diffFlags.Args()) < 2 {
			fmt.Fprint(os.Stderr, "Expected two json files\n\n")
			help(os.Args[1])
			return 1
		}
		return runDiff(diffFlags.Arg(0), diffFlags.Arg(1))

	default:
		help(os.Args[1])
		return 0
	}
}

func help(subcommand string) {
	switch subcommand {
	case "deprecation":
		depFlags.Usage()
	case "diff":
		diffFlags.Usage()
	default:
		fmt.Printf("Usage: %s <command> [<args>]\n\n", path.Base(os.Args[0]))
		depFlags.Usage()
		fmt.Println()
		diffFlags.Usage()
	}

}

func runDeprecation() int {
	schema, err := parser.ParseSchemaFile(schemaFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse schema file %s: %s", schemaFile, err)
		return 1
	}

	queryFields, err := parser.ParseQueryDir(queriesDir, schema)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse files in %s: %s", queriesDir, err)
		return 1
	}

	switch outputFormat {
	case "stdout":
		deprecationStdOut(queryFields)

	case "json":
		deprecationJsonOut(queryFields)

	default:
		fmt.Fprintf(os.Stderr, "%s is not a valid output format. Choose between json and stdout", outputFormat)
		return 1
	}
	return 0
}

func runDiff(fileA string, fileB string) int {
	result, err := output.CompareFiles(fileA, fileB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to diff: %s\n", err)
		return 1
	}

	switch outputFormat {
	case "stdout":
		diffStdOut(fileA, fileB, result)

	case "json":
		diffJsonOut(result)

	default:
		fmt.Fprintf(os.Stderr, "%s is not a valid output format. Choose between json and stdout\n", outputFormat)
		return 1
	}
	return 0
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
