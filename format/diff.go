package format

import (
	"fmt"
	"strings"

	"github.com/sketch-hq/gql-lint/output"
)

var DiffFormatter Formatter = Formatter{
	StdoutFormat: diffStdOut,
	JsonFormat:   jsonOut,
	XcodeFormat:  xcodeOut,
}

func diffStdOut(out output.Data) (string, error) {
	print := ""
	out.Walk(func(schema string, f output.Field, i int) {
		if i == 0 {
			print += fmt.Sprintln("Schema: ", schema)
		}
		reason := strings.ReplaceAll(f.DeprecationReason, "\n", " ")
		print += fmt.Sprintf("%s (%s)\n", f.Field, reason)
		print += fmt.Sprintf("  %s:%d\n", f.File, f.Line)
	})
	return print, nil
}
