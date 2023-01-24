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

const DownloadFailed = 10
const ParseSDLFailed = 20

//go:embed introspection.gql
var introspectionQuery string

type Error struct {
	Code int
	Err  error
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) ReturnCode() int {
	return e.Code
}

type jsonBody struct {
	Data json.RawMessage `json:"data"`
}

// Load will fetch a schema at the given url and unwrap the response and
// attempt to convert the schema from JSON to SDL.
func Load(url string) ([]byte, error) {
	jsonSchema, err := fetch(url)
	if err != nil {
		return []byte{}, &Error{Err: err, Code: DownloadFailed}
	}
	schema, err := jsonToSDL(jsonSchema)
	if err != nil {
		return []byte{}, &Error{Err: err, Code: ParseSDLFailed}
	}

	return schema, nil
}

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
