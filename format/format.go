package format

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sketch-hq/gql-lint/output"
	"github.com/sketch-hq/gql-lint/utils"
)

const (
	JsonFormat     = "json"
	StdoutFormat   = "stdout"
	MarkdownFormat = "markdown"
	XcodeFormat    = "xcode"
	AnnotateFormat = "annotate"
)

type Formatter map[string]func(output.Data) (string, error)

func (f *Formatter) Format(format string, out output.Data) (string, error) {
	validFormats := f.validFormats()
	if !utils.Contains(validFormats, format) {
		return "", fmt.Errorf("%s is not a valid output format. Choose between %s", format, strings.Join(validFormats, ", "))
	}

	return (*f)[format](out)
}

func (f *Formatter) validFormats() []string {
	return utils.Keys(*f)
}

func jsonOut(out output.Data) (string, error) {
	bytes, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("failed to encode json: %s", err)
	}

	return string(bytes), nil
}

func xcodeOut(out output.Data) (string, error) {
	print := ""
	out.Walk(func(_ string, f output.Field, _ int) {
		reason := strings.ReplaceAll(f.DeprecationReason, "\n", " ")
		print += fmt.Sprintf("%s:%d: warning: ", f.File, f.Line)
		print += fmt.Sprintf("%s is deprecated ", f.Field)
		print += fmt.Sprintf("- Reason: %s", reason)
		print += fmt.Sprintln()
	})
	return print, nil
}
