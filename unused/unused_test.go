package unused_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/sketch-hq/gql-lint/unused"
)

func TestGetUnusedFields(t *testing.T) {
	is := is.New(t)

	schema := "testdata/schemas/with_deprecations.gql"

	t.Run("reports unused fields", func(t *testing.T) {
		is := is.New(t)
		unusedFields, err := unused.GetUnusedFields([]string{schema}, []string{
			"testdata/queries/one/one.gql",
			"testdata/queries/deprecation/deprecation.gql",
		}, false)
		is.NoErr(err)

		is.Equal(len(unusedFields), 1)
		is.Equal(unusedFields[schema][0].Field, "Book.oldTitle")
	})

	t.Run("sorts fields aphabetically", func(t *testing.T) {
		is := is.New(t)
		unusedFields, err := unused.GetUnusedFields([]string{schema}, []string{
			"testdata/queries/one/one.gql",
		}, false)
		is.NoErr(err)

		is.Equal(len(unusedFields[schema]), 2)

		is.Equal(unusedFields[schema][0].Field, "Book.oldTitle")
		is.Equal(unusedFields[schema][1].Field, "Book.title")
	})
}
