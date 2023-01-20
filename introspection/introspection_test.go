package introspection_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
	"github.com/sketch-hq/gql-lint/introspection"
)

func TestLoad(t *testing.T) {
	t.Run("returns error on connection issue", func(t *testing.T) {
		_, err := introspection.Load("http://127.0.0.1:invalid/")
		if err == nil {
			t.Errorf("expected error, got: %s", err)
		}
	})

	t.Run("returns error if json schema is invalid", func(t *testing.T) {
		server := httpServer(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			fmt.Fprintf(w, "invalid json")
		})
		defer server.Close()

		_, err := introspection.Load(server.URL)
		if err == nil {
			t.Errorf("expected error, got: %s", err)
		}
	})

	t.Run("returns schema as SDL", func(t *testing.T) {
		is := is.New(t)
		server := httpServer(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			fmt.Fprintf(w, `
			{
				"data": {
					"__schema": {
						"queryType": {
							"name": "Query"
						},
						"mutationType": {
							"name": "Mutation"
						}
					}
				}
			}
			`)
		})
		defer server.Close()

		schema, err := introspection.Load(server.URL)
		is.NoErr(err)

		expectedSchema := `schema {
    query: Query
    mutation: Mutation
}`

		is.Equal(string(schema), expectedSchema)
	})
}

func httpServer(handler http.HandlerFunc) *httptest.Server {
	server := httptest.NewServer(handler)

	return server
}
