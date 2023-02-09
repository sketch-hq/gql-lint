package ops

const (
	schemaFileFlagName   = "schema"
	schemaFileDefault    = ""
	outputFormatFlagName = "output"
	jsonFormat           = "json"
	stdoutFormat         = "stdout"
	markdownFormat       = "markdown"
	xcodeFormat          = "xcode"
	ignoreFlagName       = "ignore"
	verboseFlagName      = "verbose"
)

var flags = struct {
	outputFormat string
	schemaFile   string
	schemaFiles  []string
	ignore       []string
	verbose      bool
}{}

// this is important for tests as these flags wont be reset between each test
// run unless we do it here.
func setFlagsToDefault() {
	flags.outputFormat = stdoutFormat
	flags.schemaFile = schemaFileDefault
	flags.schemaFiles = []string{}
	flags.ignore = []string{}
	flags.verbose = false
}
