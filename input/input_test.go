package input

import (
	"testing"

	"github.com/matryer/is"
)

func TestQueryFiles(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		name    string
		args    []string
		ignore  []string
		include []string
		want    []string
	}{
		{
			name: "expands directories",
			args: []string{"testdata/*.gql"},
			want: []string{"testdata/one.gql", "testdata/two.gql"},
		},
		{
			name: "doesn't duplicate files",
			args: []string{"testdata/one.gql", "testdata/*.gql"},
			want: []string{"testdata/one.gql", "testdata/two.gql"},
		},
		{
			name: "accepts single files",
			args: []string{"testdata/one.gql"},
			want: []string{"testdata/one.gql"},
		},
		{
			name: "expands nested directories",
			args: []string{"testdata/**/*.gql"},
			want: []string{"testdata/one.gql", "testdata/two.gql", "testdata/nested/one.gql", "testdata/nested/two.gql"},
		},
		{
			name:   "expands nested directories",
			args:   []string{"testdata/**/*.gql"},
			ignore: []string{"testdata/**/one.gql"},
			want:   []string{"testdata/two.gql", "testdata/nested/two.gql"},
		},
		{
			name:   "expands nested directories",
			args:   []string{"testdata/**/*.gql"},
			ignore: []string{"./testdata/**/one.gql"},
			want:   []string{"testdata/two.gql", "testdata/nested/two.gql"},
		},
		{
			name:    "removes everything not matching an include filter",
			args:    []string{"testdata/**/*.gql"},
			include: []string{"testdata/one.gql"},
			want:    []string{"testdata/one.gql"},
		},
		{
			name:    "removes everything not matching a nested directory include filter",
			args:    []string{"testdata/**/*.gql"},
			include: []string{"./**/one.gql"},
			want:    []string{"testdata/one.gql", "testdata/nested/one.gql"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExpandGlobs(tt.args, tt.include, tt.ignore)
			is.NoErr(err)

			is.Equal(got, tt.want)
		})
	}
}
