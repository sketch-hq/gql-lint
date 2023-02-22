package format

import (
	"fmt"

	"github.com/sketch-hq/gql-lint/output"
)

var UnusedFormatter Formatter = Formatter{
	StdoutFormat:   unusedStdOut,
	JsonFormat:     jsonOut,
	MarkdownFormat: unusedMarkdownOut,
}

func unusedStdOut(out output.Data) (string, error) {
	print := ""
	hasUnused := false
	out.Walk(func(schema string, f output.Field, i int) {
		if i == 0 {
			if hasUnused {
				print += "\n" // new line for each schema, except the first one.
			}
			print += fmt.Sprintln("Schema:", schema)
			hasUnused = true
		}
		print += fmt.Sprintf("  %s (line %d) is unused and can be removed \n", f.Field, f.Line)
	})
	return print, nil
}

func unusedMarkdownOut(out output.Data) (string, error) {
	print := ""
	hasUnused := false

	out.Walk(func(schema string, f output.Field, fieldIdx int) {
		if fieldIdx == 0 {
			if hasUnused {
				print += "\n" // new line for each schema, except the first one.
			}
			hasUnused = true
			print += fmt.Sprintln("Schema:", schema)
		}
		print += fmt.Sprintf("- %s (line `%d`)\n", f.Field, f.Line)
	})

	if !hasUnused {
		print += fmt.Sprintln("Nothing can be removed right now")
	}
	return print, nil
}
