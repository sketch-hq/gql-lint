package unused

import (
	"testing"

	"github.com/matryer/is"
)

func TestGetUnusedFields(t *testing.T) {
	is := is.New(t)

	unusedFields, err := GetUnusedFields(
		"testdata/schemas/with_deprecations.gql",
		[]string{
			"testdata/queries/one.gql",
			"testdata/queries/deprecation.gql",
		},
	)
	is.NoErr(err)

	is.Equal(len(unusedFields), 1)
	is.Equal(unusedFields[0].Name, "Book.oldTitle")
}
