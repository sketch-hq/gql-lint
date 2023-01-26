package sources_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/matryer/is"
	"github.com/sketch-hq/gql-lint/sources"
)

func TestLoadSchema(t *testing.T) {
	t.Run("load schema from http with missing deprecated directive", func(t *testing.T) {
		is := is.New(t)
		server := httpServer(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			serveJsonFixture(is, w, "schema-without-deprecated.json")
		})
		defer server.Close()

		_, err := sources.LoadSchema(server.URL)
		is.NoErr(err)
	})

	t.Run("load schema from http", func(t *testing.T) {
		is := is.New(t)
		server := httpServer(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			serveJsonFixture(is, w, "schema.json")
		})
		defer server.Close()

		_, err := sources.LoadSchema(server.URL)
		is.NoErr(err)
	})

	t.Run("error if url returns not 200 OK", func(t *testing.T) {
		is := is.New(t)
		server := httpServer(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		})
		defer server.Close()

		_, err := sources.LoadSchema(server.URL)
		is.True(err != nil)
	})

	t.Run("load schema from file", func(t *testing.T) {
		is := is.New(t)
		_, err := sources.LoadSchema("testdata/schema.gql")
		is.NoErr(err)
	})
}

func TestLoadQueries(t *testing.T) {
	run := func(t *testing.T, source string, expectError bool) {
		is := is.New(t)
		is.Helper()
		schema, err := sources.LoadSchema("testdata/schema.gql")
		is.NoErr(err)

		fields, err := sources.LoadQueries(schema, []string{source})
		if expectError {
			is.True(err != nil)
			is.True(fields == nil)
		} else {
			is.NoErr(err)
			is.True(len(fields) > 0)
		}
	}

	t.Run("from a file", func(t *testing.T) {
		run(t, "testdata/query.gql", false)
	})

	t.Run("errors if file not found", func(t *testing.T) {
		run(t, "testdata/file_not_found.gql", true)
	})

	t.Run("from a directory", func(t *testing.T) {
		run(t, "testdata/queries", false)
	})
}

func httpServer(handler http.HandlerFunc) *httptest.Server {
	server := httptest.NewServer(handler)
	return server
}

func serveJsonFixture(is *is.I, w io.Writer, file string) {
	is.Helper()
	f, err := os.Open("testdata/" + file)
	is.NoErr(err)

	_, err = io.Copy(w, f)
	is.NoErr(err)
}
