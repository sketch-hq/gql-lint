package ops

const (
	schemaFileFlagName   = "schema"
	outputFormatFlagName = "output"
	jsonFormat           = "json"
	stdoutFormat         = "stdout"
	markdownFormat       = "markdown"
	xcodeFormat          = "xcode"
	annotateFormat       = "annotate"
	ignoreFlagName       = "ignore"
	verboseFlagName      = "verbose"
)

var flags = struct {
	outputFormat string
	schemaFiles  []string
	ignore       []string
	verbose      bool
}{}

// this is important for tests as these flags wont be reset between each test
// run unless we do it here.
func setFlagsToDefault() {
	flags.outputFormat = stdoutFormat
	flags.schemaFiles = []string{}
	flags.ignore = []string{}
	flags.verbose = false
}
