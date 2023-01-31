package unused

import (
	"testing"

	"github.com/matryer/is"
	"github.com/sketch-hq/gql-lint/sources"
)

func TestGetUnusedFields(t *testing.T) {
	is := is.New(t)

	schema, err := sources.LoadSchema("testdata/schemas/with_deprecations.gql")
	is.NoErr(err)

	t.Run("reports unused fields", func(t *testing.T) {
		is := is.New(t)
		unusedFields, err := GetUnusedFields(schema, []string{
			"testdata/queries/one/one.gql",
			"testdata/queries/deprecation/deprecation.gql",
		})
		is.NoErr(err)

		is.Equal(len(unusedFields), 1)
		is.Equal(unusedFields[0].Name, "Book.oldTitle")
	})

	t.Run("sorts fields aphabetically", func(t *testing.T) {
		is := is.New(t)
		unusedFields, err := GetUnusedFields(schema, []string{
			"testdata/queries/one/one.gql",
		})
		is.NoErr(err)

		is.Equal(len(unusedFields), 2)

		is.Equal(unusedFields[0].Name, "Book.oldTitle")
		is.Equal(unusedFields[1].Name, "Book.title")
	})
}
