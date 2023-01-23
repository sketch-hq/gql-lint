package input

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestReadArgs(t *testing.T) {
	is := is.New(t)

	args := []string{"arg1", "arg2"}
	got, err := ReadArgs(args)
	is.NoErr(err)

	is.Equal(got, args)
}

func TestReadArgsPiped(t *testing.T) {
	simulatePiping("mocked", func() {
		is := is.New(t)

		got, err := ReadArgs([]string{"arg1", "arg2"})
		is.NoErr(err)

		is.Equal(got, []string{"mocked"})
	})
}

func simulatePiping(input string, test func()) {
	tmpfile, err := ioutil.TempFile("", input)

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(input)); err != nil {
		log.Fatal(err)
	}
	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = tmpfile

	test()

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
}

func TestQueryFiles(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		name string
		args []string
		want []string
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryFiles(tt.args)
			is.NoErr(err)

			is.Equal(got, tt.want)
		})
	}
}
