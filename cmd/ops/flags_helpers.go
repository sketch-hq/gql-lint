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
	includeFlagName      = "include"
	verboseFlagName      = "verbose"
)

var flags = struct {
	outputFormat string
	schemaFiles  []string
	ignore       []string
	include      []string
	verbose      bool
}{}

// this is important for tests as these flags wont be reset between each test
// run unless we do it here.
func setFlagsToDefault() {
	flags.outputFormat = stdoutFormat
	flags.schemaFiles = []string{}
	flags.ignore = []string{}
	flags.include = []string{}
	flags.verbose = false
}
