package format_test

import (
	"testing"

	"github.com/sketch-hq/gql-lint/format"
	"github.com/sketch-hq/gql-lint/output"
)

type testDefinition struct {
	name     string
	format   string
	expected string
}

var data = output.Data{
	"http://example.com/schema": []output.Field{
		{
			Field:             "Article.view",
			File:              "somefile.graphql",
			Line:              13,
			DeprecationReason: "Please migrate to Article.permissions",
		},
		{
			Field: "Article.author",
			File:  "somefile.graphql",
			Line:  25,
			DeprecationReason: `Please migrate to Article.owner.
See "https://example.com/migration" for more information.`,
		},
	},
	"http://other.example.com/schema": []output.Field{
		{
			Field:             "Book.title",
			File:              "other.graphql",
			Line:              10,
			DeprecationReason: "Replace with to Book.headline",
		},
	},
}

func runTests(t *testing.T, formatter format.Formatter, tests []testDefinition) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := formatter.Format(tt.format, data)
			if err != nil {
				t.Errorf("error = %v", err)
				return
			}
			if r != tt.expected {
				t.Errorf("got\n%v\nwant\n%v", r, tt.expected)
				return
			}
		})
	}
}

// This test the json and xcode formats for other formatters too
func TestDeprecationFormatter(t *testing.T) {

	tests := []testDefinition{
		{
			name:     "deprecation json output",
			format:   format.JsonFormat,
			expected: `{"http://example.com/schema":[{"field":"Article.view","file":"somefile.graphql","line":13,"reason":"Please migrate to Article.permissions"},{"field":"Article.author","file":"somefile.graphql","line":25,"reason":"Please migrate to Article.owner.\nSee \"https://example.com/migration\" for more information."}],"http://other.example.com/schema":[{"field":"Book.title","file":"other.graphql","line":10,"reason":"Replace with to Book.headline"}]}`,
		},
		{
			name:   "deprecation annotate output",
			format: format.AnnotateFormat,
			expected: `[{ "file": "somefile.graphql", "line": 13, "title": "Article.view is deprecated", "message": "Please migrate to Article.permissions", "annotation_level": "warning" },
{ "file": "somefile.graphql", "line": 25, "title": "Article.author is deprecated", "message": "Please migrate to Article.owner.\nSee \"https://example.com/migration\" for more information.", "annotation_level": "warning" },
{ "file": "other.graphql", "line": 10, "title": "Book.title is deprecated", "message": "Replace with to Book.headline", "annotation_level": "warning" }]`,
		},
		{
			name:   "deprecation xcode output",
			format: format.XcodeFormat,
			expected: `somefile.graphql:13: warning: Article.view is deprecated - Reason: Please migrate to Article.permissions
somefile.graphql:25: warning: Article.author is deprecated - Reason: Please migrate to Article.owner. See "https://example.com/migration" for more information.
other.graphql:10: warning: Book.title is deprecated - Reason: Replace with to Book.headline
`,
		},
		{
			name:   "deprecation stdout output",
			format: format.StdoutFormat,
			expected: `Schema: http://example.com/schema
  Article.view is deprecated
    File:   somefile.graphql:13
    Reason: Please migrate to Article.permissions

  Article.author is deprecated
    File:   somefile.graphql:25
    Reason: Please migrate to Article.owner. See "https://example.com/migration" for more information.

Schema: http://other.example.com/schema
  Book.title is deprecated
    File:   other.graphql:10
    Reason: Replace with to Book.headline

`,
		},
	}
	runTests(t, format.DeprecationFormatter, tests)
}

func TestDiffFormatter(t *testing.T) {

	tests := []testDefinition{
		{
			name:   "diff stdout output",
			format: format.StdoutFormat,
			expected: `Schema:  http://example.com/schema
Article.view (Please migrate to Article.permissions)
  somefile.graphql:13
Article.author (Please migrate to Article.owner. See "https://example.com/migration" for more information.)
  somefile.graphql:25
Schema:  http://other.example.com/schema
Book.title (Replace with to Book.headline)
  other.graphql:10
`,
		},
	}
	runTests(t, format.DiffFormatter, tests)
}

func TestUnusedFormatter(t *testing.T) {

	tests := []testDefinition{
		{
			name:   "unused stdout output",
			format: format.StdoutFormat,
			expected: `Schema: http://example.com/schema
  Article.view (line 13) is unused and can be removed 
  Article.author (line 25) is unused and can be removed 

Schema: http://other.example.com/schema
  Book.title (line 10) is unused and can be removed 
`,
		},
		{
			name:   "unused markdown output",
			format: format.MarkdownFormat,
			expected: `Schema: http://example.com/schema
- Article.view (line ` + "`13`" + `)
- Article.author (line ` + "`25`" + `)

Schema: http://other.example.com/schema
- Book.title (line ` + "`10`" + `)
`,
		},
	}
	runTests(t, format.UnusedFormatter, tests)
}
