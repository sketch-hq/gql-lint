package introspection

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/wundergraph/graphql-go-tools/pkg/astprinter"
	wunderintro "github.com/wundergraph/graphql-go-tools/pkg/introspection"
)

//go:embed introspection.gql
var introspectionQuery string

type jsonBody struct {
	Data json.RawMessage `json:"data"`
}

func Load(url string) ([]byte, error) {
	jsonSchema, err := fetch(url)
	if err != nil {
		return []byte{}, err
	}
	schema, err := jsonToSDL(jsonSchema)
	if err != nil {
		return []byte{}, err
	}

	return schema, nil
}

// Schema will fetch a schema at the given url and unwrap the response and
// return the schema inside the `data` key in the response.
func fetch(url string) ([]byte, error) {
	q := map[string]string{
		"query":         introspectionQuery,
		"operationName": "IntrospectionQuery",
	}
	queryBytes, err := json.Marshal(q)
	if err != nil {
		return nil, fmt.Errorf("failed to encode the request body: %w", err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(queryBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to download schema: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response body: %w", err)
	}

	var body jsonBody
	err = json.Unmarshal(b, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return body.Data, nil
}

// jsonToSDL converts a json schema representation to a SDL string representation
func jsonToSDL(j []byte) ([]byte, error) {
	converter := wunderintro.JsonConverter{}
	doc, err := converter.GraphQLDocument(bytes.NewReader(j))
	if err != nil {
		return []byte{}, err
	}

	outWriter := &bytes.Buffer{}
	err = astprinter.PrintIndent(doc, nil, []byte("  "), outWriter)
	if err != nil {
		return []byte{}, err
	}

	return outWriter.Bytes(), nil
}
