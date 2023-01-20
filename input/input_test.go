package input

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/cobra"
)

func TestReadArgs(t *testing.T) {
	is := is.New(t)

	args := []string{"arg1", "arg2"}
	got := ReadArgs(&cobra.Command{}, args)
	is.Equal(got, args)
}

type MockCommand struct {
	Return string
}

func (c *MockCommand) InOrStdin() io.Reader {
	return strings.NewReader(c.Return)
}

func TestReadArgsPiped(t *testing.T) {
	simulatePiping(func() {
		is := is.New(t)

		cmd := &MockCommand{Return: "mocked"}

		got := ReadArgs(cmd, []string{"arg1", "arg2"})

		is.Equal(got, []string{"mocked"})
	})
}

func simulatePiping(test func()) {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte("mocked_value")); err != nil {
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
