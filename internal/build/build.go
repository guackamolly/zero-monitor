package internal

import (
	"flag"
	"os"
)

// Holds the program version value. This value is linked at build-time with the -X flag.
// See [tools/build].
var version string

func init() {
	flag.BoolFunc("version", "Prints build version", func(s string) error {
		println(version)

		os.Exit(0)
		return nil
	})

	flag.Parse()
}
