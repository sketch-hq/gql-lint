package ops

const (
	schemaFileFlagName   = "schema"
	schemaFileDefault    = ""
	outputFormatFlagName = "output"
	jsonFormat           = "json"
	stdoutFormat         = "stdout"
	xcodeFormat          = "xcode"
)

var flags = struct {
	outputFormat string
	schemaFile   string
}{}

func setFlagsToDefault() {
	flags.outputFormat = stdoutFormat
	flags.schemaFile = schemaFileDefault
}
