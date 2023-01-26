package ops

const (
	schemaFileFlagName   = "schema"
	schemaFileDefault    = ""
	outputFormatFlagName = "output"
	jsonFormat           = "json"
	stdoutFormat         = "stdout"
	xcodeFormat          = "xcode"
	ignoreFlagName       = "ignore"
)

var flags = struct {
	outputFormat string
	schemaFile   string
	ignore       []string
}{}

func setFlagsToDefault() {
	flags.outputFormat = stdoutFormat
	flags.schemaFile = schemaFileDefault
}
