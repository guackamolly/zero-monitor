package internal

// Holds the program version value. This value is linked at build-time with the -X flag.
// See [tools/build].
var version string

// If the program is running on release mode. This value is linked at build-time with the -X flag.
// See [tools/build].
var release string = "false"

// Indicates the version that the program is running.
func Version() string {
	return version
}

// Indicates if the program is running on release mode.
func Release() bool {
	return release == "true"
}
