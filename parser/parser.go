package parser

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
	gqlparser "github.com/vektah/gqlparser/v2/parser"
	gqlvalidator "github.com/vektah/gqlparser/v2/validator"
)

const deprecated = `directive @deprecated(
	reason: String = "No longer supported"
) on FIELD_DEFINITION | ARGUMENT_DEFINITION | ENUM_VALUE | INPUT_FIELD_DEFINITION`

type SchemaField struct {
	Name string
}

type QueryField struct {
	Path              string
	SchemaPath        string
	IsDeprecated      bool
	DeprecationReason string
	File              string
	Line              int
}

type QueryFieldList []QueryField

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}

func ParseQuerySource(source string, schema *ast.Schema) (QueryFieldList, error) {
	isDir, err := isDirectory(source)
	if err != nil {
		return nil, err
	}

	if isDir {
		return parseQueryDir(source, schema)
	} else {
		return parseQueryList(source, schema)
	}
}

func parseQueryDir(dir string, schema *ast.Schema) (QueryFieldList, error) {
	files := findQueryFiles(dir)
	if len(files) == 0 {
		return QueryFieldList{}, fmt.Errorf("no query files found in %s", dir)
	}
	return queryTokensFromFiles(files, schema)
}

func parseQueryList(file string, schema *ast.Schema) (QueryFieldList, error) {
	queryPaths, err := getLinesFromFile(file)
	if err != nil {
		return nil, err
	}
	return queryTokensFromFiles(queryPaths, schema)
}

func ParseDeprecatedFields(schema *ast.Schema) []SchemaField {
	var fields []SchemaField

	for name, definition := range schema.Types {
		if ok, _ := isDeprecated(definition.Directives); ok {
			fields = append(fields, SchemaField{Name: name})
		}

		for _, field := range definition.Fields {
			if ok, _ := isDeprecated(field.Directives); ok {
				fields = append(fields, SchemaField{Name: fmt.Sprintf("%s.%s", name, field.Name)})
			}
		}
	}

	return fields
}

func ParseSchema(name string, contents string, hasBuiltin bool) (*ast.Schema, error) {
	sources := []*ast.Source{
		// When parsing downloaded schemas built types are included. If we
		// don't set bultin=true the validator will complain about types using
		// reserved "__" names
		{Name: name, Input: string(contents), BuiltIn: hasBuiltin},
	}
	// workaround as absinthe based graphql servers does NOT include the deprecated
	// directive in the schema.
	if !strings.Contains(contents, "directive @deprecated") {
		sources = append(sources, &ast.Source{Input: deprecated, BuiltIn: true})
	}

	contents = strings.Replace(contents, `Represents a plan "type"`, "Represents a plan type", -1)
	fmt.Println(contents)

	schema, schemaerr := gqlvalidator.LoadSchema(sources...)
	if schemaerr != nil {
		return nil, schemaerr
	}

	return schema, nil
}

func findQueryFiles(startDir string) []string {
	files := []string{}
	filepath.WalkDir(startDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files
}

func getLinesFromFile(file string) ([]string, error) {
	var lines []string
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func parseQueryFile(file string) (*ast.QueryDocument, error) {
	queryContents, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	source := &ast.Source{Name: file, Input: string(queryContents)}
	doc, queryerr := gqlparser.ParseQuery(source)
	if queryerr != nil {
		return nil, queryerr
	}
	return doc, nil
}

func queryTokensFromFiles(files []string, schema *ast.Schema) (QueryFieldList, error) {
	fields := QueryFieldList{}
	for _, file := range files {
		doc, err := parseQueryFile(file)
		if err != nil {
			return nil, err
		}
		// Ignoring errors here as there could be validation errors that we're not interested in.
		// We only call `.Validate` so the parser can populate the `Definition` fields.
		_ = gqlvalidator.Validate(schema, doc)
		fields = buildQueryTokens(doc, fields)
	}

	return fields, nil
}

func buildQueryTokens(query *ast.QueryDocument, fields QueryFieldList) QueryFieldList {
	for _, o := range query.Operations {
		fields = extractFields(o.SelectionSet, "", "", fields)
	}
	for _, f := range query.Fragments {
		fields = extractFields(f.SelectionSet, f.TypeCondition, "", fields)
	}

	return fields
}

func extractFields(set ast.SelectionSet, parentPath string, parentType string, fields QueryFieldList) QueryFieldList {
	for _, s := range set {
		switch f := s.(type) {
		case *ast.Field:
			path := ""
			if parentPath != "" {
				path = parentPath + "." + f.Name
			} else {
				path = f.Name
			}

			dep, depReason := isDeprecated(f.Definition.Directives)
			if !dep {
				fields = extractFields(f.SelectionSet, path, f.Definition.Type.Name(), fields)
				continue
			}

			field := QueryField{
				Path:              path,
				SchemaPath:        fmt.Sprintf("%s.%s", parentType, f.Name),
				File:              f.Position.Src.Name,
				Line:              f.Position.Line,
				IsDeprecated:      dep,
				DeprecationReason: depReason,
			}

			fields = append(fields, field)
			fields = extractFields(f.SelectionSet, path, f.Definition.Type.Name(), fields)
		}
	}
	return fields
}

func isDeprecated(dl ast.DirectiveList) (bool, string) {
	for _, d := range dl {
		if d.Name == "deprecated" {
			if reason, ok := d.ArgumentMap(make(map[string]interface{}))["reason"]; ok {
				return true, reason.(string)
			}
			return true, ""
		}
	}

	return false, ""
}
