package format

import (
	"fmt"
	"strings"

	"github.com/sketch-hq/gql-lint/output"
)

var DeprecationFormatter Formatter = Formatter{
	StdoutFormat:   deprecationStdOut,
	JsonFormat:     jsonOut,
	XcodeFormat:    xcodeOut,
	AnnotateFormat: deprecationAnnotateOut,
}

func deprecationStdOut(out output.Data) (string, error) {
	print := ""
	out.Walk(func(schema string, f output.Field, i int) {
		if i == 0 {
			print += fmt.Sprintln("Schema:", schema)
		}
		print += fmt.Sprintf("  %s is deprecated\n", f.Field)
		print += fmt.Sprintf("    File:   %s:%d\n", f.File, f.Line)
		reason := strings.ReplaceAll(f.DeprecationReason, "\n", " ")
		print += fmt.Sprintf("    Reason: %s\n", reason)
		print += fmt.Sprintln()
	})
	return print, nil
}

func deprecationAnnotateOut(out output.Data) (string, error) {
	print := ""
	var replacer = strings.NewReplacer(
		"\n", "\\n",
		"\"", "\\\"",
	)

	print += "["
	out.Walk(func(_ string, f output.Field, idx int) {
		if idx > 0 || (idx == 0 && print != "[") {
			print += ",\n" // add a comma if we're not the first element or when we start a new schema (not the first one)
		}
		reason := replacer.Replace(f.DeprecationReason)
		print += fmt.Sprintf("{ \"file\": \"%s\", \"line\": %d, ", f.File, f.Line)
		print += fmt.Sprintf("\"title\": \"%s is deprecated\", ", f.Field)
		print += fmt.Sprintf("\"message\": \"%s\", ", reason)
		print += "\"annotation_level\": \"warning\" }"

	})
	print += "]"
	return print, nil
}
