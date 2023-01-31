package parser

import (
	"fmt"

	"github.com/vektah/gqlparser/v2/ast"
	gqlparser "github.com/vektah/gqlparser/v2/parser"
	gqlvalidator "github.com/vektah/gqlparser/v2/validator"
)

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

func ParseSchema(sources ...*ast.Source) (*ast.Schema, error) {
	schema, schemaerr := gqlvalidator.LoadSchema(sources...)
	if schemaerr != nil {
		return nil, schemaerr
	}

	return schema, nil
}

func ParseQueries(schema *ast.Schema, sources ...*ast.Source) (QueryFieldList, error) {
	fields := QueryFieldList{}
	for _, source := range sources {
		doc, queryerr := gqlparser.ParseQuery(source)
		if queryerr != nil {
			return nil, queryerr
		}
		// Ignoring errors here as there could be validation errors that we're not interested in.
		// We only call `.Validate` so the parser can populate the `Definition` fields.
		_ = gqlvalidator.Validate(schema, doc)
		fields = buildQueryTokens(doc, fields)
	}

	return fields, nil
}

func ParseDeprecatedFields(schema *ast.Schema) []SchemaField {
	var fields []SchemaField

	for name, definition := range schema.Types {
		if name == "__Directive" {
			continue
		}

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

			// This means the field doesn't belong to this schema
			if f.Definition == nil {
				continue
			}

			dep, depReason := isDeprecated(f.Definition.Directives)
			if !dep {
				fields = extractFields(f.SelectionSet, path, f.Definition.Type.Name(), fields)
				continue
			}

			parent := parentType
			if parent == "" {
				parent = parentPath
			}
			if parent == "" {
				parent = f.ObjectDefinition.Name
			}

			field := QueryField{
				Path:              path,
				SchemaPath:        fmt.Sprintf("%s.%s", parent, f.Name),
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
