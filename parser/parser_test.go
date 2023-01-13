package parser

import (
	"testing"
)

func TestParseSchemaFile(t *testing.T) {
	schema, err := ParseSchemaFile("./fixtures/with_deprecations.gql")

	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}

	deprecatedField := schema.Types["Book"].Fields[1]

	if deprecatedField.Name != "title" {
		t.Fatalf("expected field '%s' to be 'title", deprecatedField.Name)
	}

	deprecated := false
	for _, directive := range deprecatedField.Directives {
		if directive.Name == "deprecated" {
			if directive.Arguments[0].Name != "reason" {
				t.Fatalf("Expected reason")
			}

			if directive.Arguments[0].Value.String() != `"untitled books are better"` {
				t.Fatalf("unexpected description: %s", directive.Arguments[0].Value.String())
			}
			deprecated = true
		}
	}

	if !deprecated {
		t.Fatalf("%s expected to be deprecated", deprecatedField.Name)
	}

	field := schema.Types["Book"].Fields[2]

	if field.Name != "author" {
		t.Fatalf("expected field '%s' to be 'title", field.Name)
	}

	deprecated = false
	for _, directive := range field.Directives {
		if directive.Name == "deprecated" {
			deprecated = true
		}
	}

	if deprecated {
		t.Fatalf("%s isn't expected to be deprecated", field.Name)
	}
}

func TestParseSchemaFile_NotFound(t *testing.T) {
	_, err := ParseSchemaFile("./fixtures/not_found.gql")

	if err == nil {
		t.Fatalf("Expected error, got none")
	}
}

func TestParseQueryDir(t *testing.T) {
	schema, err := ParseSchemaFile("./fixtures/with_deprecations.gql")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	fields, err := ParseQueryDir("./fixtures/queries", schema)

	if len(fields) != 1 {
		t.Fatalf("expected only one deprecation, got %d", len(fields))
	}

	field := fields[0]

	if field.Path == "author.book.title" {
		t.Fatalf("expected 'author.book.title', got %s", field.Path)
	}

	if !field.IsDeprecated {
		t.Fatalf("expected '%s' to be deprecated", field.Path)
	}

	if field.File != "fixtures/queries/deprecation.gql" || field.Line != 7 {
		t.Fatalf("unexpected location for deprecation warning: %s:%d", field.File, field.Line)
	}
}
