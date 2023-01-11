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
	flag.StringVar(&schemaFile, "schema", "", "Schema file")
	flag.Parse()
	queriesDir = flag.Arg(0)
	if queriesDir == "" {
		fmt.Fprint(os.Stderr, "You must specify a directory for queries and mutations\n")
		os.Exit(1)
	}

	if schemaFile == "" {
		fmt.Fprint(os.Stderr, "You must specify a schema file\n")
		os.Exit(1)
	}

	run()
}

func run() {
	schemaFields, err := parser.ParseSchemaFile(schemaFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse schema file %s: %s", schemaFile, err)
		os.Exit(1)
	}
	// a, err := json.MarshalIndent(schemaFields, "", "  ")
	// if err != nil {
	// 	fmt.Println("error:", err)
	// }
	// fmt.Print(string(a))

	queryFields, err := parser.ParseQueryDir(queriesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse files in %s: %s", queriesDir, err)
		os.Exit(1)
	}

	// b, err := json.MarshalIndent(queryFields, "", "  ")
	// if err != nil {
	// 	fmt.Println("error:", err)
	// }
	// fmt.Print(string(b))

	for qk, qv := range queryFields {
		if sv, ok := schemaFields[qk]; ok {
			fmt.Printf("%s is deprecated\n", qv.Path)
			fmt.Printf("  File:   %s:%d\n", qv.Position.Src.Name, qv.Position.Line)
			fmt.Printf("  Reason: test%s\n", sv.DeprecationReason)
			fmt.Println()
		}
	}

}

// only fields and enum values

// Mac: separate directories
// FE: @capacitorapi directive

// type Type struct {
// 	Name   string
// 	Kind   ast.Operation
// 	Fields []string
// }

// type Token struct {
// 	Name         string
// 	Kind         string
// 	Position     *ast.Position
// 	IsDeprecated bool
// }

// type ResultLookup map[string]Token

// func main() {
// 	schema, err := parseSchema("sketchql.gql")
// 	if err != nil {
// 		fmt.Printf("Failed to parse schema: %s\n", err)
// 		return
// 	}

// 	clientTokens := make(ResultLookup)
// 	// tokensForDir("../cloud-frontend/packages/gql-types/graphql/", schema, clientTokens)
// 	// tokensForDir("test", schema, clientTokens)
// 	// tokensForDir("/tmp/schema", schema, clientTokens)
// 	tokensForDir("/Users/oh/sketch/Sketch/Modules/SketchCloudKit/Source/Resources/SketchQL Queries/SketchQL", schema, clientTokens)
// 	// fmt.Printf("Client result: %+v\n", clientTokens)
// 	b, err := json.MarshalIndent(clientTokens, "", "  ")
// 	if err != nil {
// 		fmt.Println("error:", err)
// 	}
// 	fmt.Print(string(b))
// }

// func tokensForDir(startDir string, schema *ast.Schema, result ResultLookup) {
// 	files := findClientFiles(startDir)
// 	queryContent := strings.Builder{}
// 	for _, file := range files {
// 		c, err := os.ReadFile(file)
// 		if err != nil {
// 			fmt.Printf("Could not read file %s\n", file)
// 			continue
// 		}
// 		_, _ = queryContent.Write(c)
// 	}

// 	doc, err := parseQueryString(queryContent.String())
// 	if err != nil {
// 		fmt.Printf("Could not parse queries %s\n", err)
// 		return
// 	}
// 	// ignoring errors as they often contains irrelevant issues for us, such as
// 	// duplicate named operations, etc
// 	errs := validator.Validate(schema, doc)
// 	if errs != nil {
// 		fmt.Printf("Could not validate queries: %s\n", errs.Error())
// 		return
// 	}

// 	buildQueryTokens(doc, result)
// }

// func parseSchema(file string) (*ast.Schema, error) {
// 	queryContents, err := os.ReadFile(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	source := &ast.Source{Name: file, Input: string(queryContents)}
// 	doc, schemaerr := gqlparser.LoadSchema(source)
// 	if schemaerr != nil {
// 		return nil, schemaerr
// 	}
// 	return doc, nil
// }

// func parseQueryString(query string) (*ast.QueryDocument, error) {
// 	source := &ast.Source{Name: "query", Input: query}
// 	doc, queryerr := parser.ParseQuery(source)
// 	if queryerr != nil {
// 		return nil, queryerr
// 	}

// 	return doc, nil
// }

// func parseQuery(file string, schema *ast.Schema) (*ast.QueryDocument, error) {
// 	queryContents, err := os.ReadFile(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	source := &ast.Source{Name: file, Input: string(queryContents)}
// 	doc, queryerr := parser.ParseQuery(source)
// 	if queryerr != nil {
// 		return nil, queryerr
// 	}
// 	return doc, nil
// }

// func findClientFiles(startDir string) []string {
// 	files := []string{}
// 	filepath.WalkDir(startDir, func(path string, d fs.DirEntry, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if d.IsDir() {
// 			return nil
// 		}
// 		files = append(files, path)
// 		return nil
// 	})
// 	return files
// }

// func buildSchemaTokens(schema *ast.Schema, result ResultLookup) {
// 	for _, t := range schema.Types {
// 		for _, field := range t.Fields {
// 			token := Token{
// 				Kind:     string(t.Kind),
// 				Name:     field.Name,
// 				Position: t.Position,
// 			}
// 			path := t.Name + "." + field.Name
// 			result[path] = token
// 		}
// 	}
// }

// func buildQueryTokens(query *ast.QueryDocument, clientTokens ResultLookup) {
// 	for _, o := range query.Operations {
// 		extractFields(o, o.SelectionSet, "", clientTokens)
// 	}
// }

// func extractFields(op *ast.OperationDefinition, set ast.SelectionSet, parentPath string, result ResultLookup) {
// 	for _, s := range set {
// 		switch f := s.(type) {
// 		case *ast.Field:
// 			token := Token{
// 				Kind:     string(op.Operation),
// 				Name:     f.Name,
// 				Position: f.Position,
// 			}
// 			path := ""
// 			if parentPath != "" {
// 				path = parentPath + "." + f.Name
// 			} else {
// 				path = op.Name + "." + f.Name
// 			}
// 			result[path] = token
// 			extractFields(op, f.SelectionSet, path, result)
// 		case *ast.FragmentSpread:
// 			token := Token{
// 				Kind:     string(op.Operation),
// 				Name:     f.Name,
// 				Position: f.Position,
// 			}
// 			path := ""
// 			if parentPath != "" {
// 				path = parentPath + "." + f.Name
// 			} else {
// 				path = op.Name + "." + f.Name
// 			}
// 			result[path] = token
// 			extractFields(op, f.Definition.SelectionSet, path, result)
// 		}
// 	}
// }
