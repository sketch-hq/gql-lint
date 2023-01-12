package parser

import (
	"testing"
)

func TestParseSchemaFile(t *testing.T) {
	schema, err := ParseSchemaFile("./fixtures/with_deprecations.gql")

	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}

	deprecatedField := schema.Types["Book"].Fields[0]

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

	field := schema.Types["Book"].Fields[1]

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
