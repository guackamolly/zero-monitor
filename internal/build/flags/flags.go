package autoload

import (
	"flag"
	"fmt"
	"os"

	build "github.com/guackamolly/zero-monitor/internal/build"
)

var verbose = flag.Bool("verbose", false, "Runs the program in verbose mode (all logs)")

func init() {
	flag.BoolFunc("version", "Prints build version", func(s string) error {
		fmt.Println(build.Version())

		os.Exit(0)
		return nil
	})

	flag.Parse()
}

// Indicates if the program is running on verbose mode.
func Verbose() bool {
	return *verbose
}
