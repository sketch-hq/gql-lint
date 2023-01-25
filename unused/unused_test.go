package unused

import (
	"testing"

	"github.com/matryer/is"
	"github.com/sketch-hq/gql-lint/schema"
)

func TestGetUnusedFields(t *testing.T) {
	is := is.New(t)

	schema, err := schema.Load("testdata/schemas/with_deprecations.gql")
	is.NoErr(err)

	unusedFields, err := GetUnusedFields(schema, []string{
		"testdata/queries/one/one.gql",
		"testdata/queries/deprecation/deprecation.gql",
	})
	is.NoErr(err)

	is.Equal(len(unusedFields), 1)
	is.Equal(unusedFields[0].Name, "Book.oldTitle")
}
