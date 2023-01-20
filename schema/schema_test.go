package schema_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/matryer/is"
	"github.com/sketch-hq/gql-lint/schema"
)

func TestLoad(t *testing.T) {
	t.Run("load schema from http with missing deprecated directive", func(t *testing.T) {
		is := is.New(t)
		server := httpServer(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			serveJsonFixture(is, w, "schema-without-deprecated.json")
		})
		defer server.Close()

		_, err := schema.Load(server.URL)
		is.NoErr(err)
	})

	t.Run("load schema from http", func(t *testing.T) {
		is := is.New(t)
		server := httpServer(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			serveJsonFixture(is, w, "schema.json")
		})
		defer server.Close()

		_, err := schema.Load(server.URL)
		is.NoErr(err)
	})

	t.Run("load schema from file", func(t *testing.T) {
		is := is.New(t)
		_, err := schema.Load("testdata/schema.gql")
		is.NoErr(err)
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
