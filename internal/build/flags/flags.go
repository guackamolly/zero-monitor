package autoload

import (
	"flag"
	"fmt"
	"os"

	build "github.com/guackamolly/zero-monitor/internal/build"
)

var flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

var verbose *bool
var inviteLink *string

func init() {
	flags.BoolFunc("version", "Prints build version", func(s string) error {
		fmt.Println(build.Version())

		os.Exit(0)
		return nil
	})

	verbose = flags.Bool("verbose", false, "Runs the program in verbose mode (all logs)")
}

func WithNodeFlags() {
	inviteLink = flags.String("invite-link", "", "Uses the invite-link to connect to a new network")
	flags.Parse(os.Args[1:])
}

func WithMasterFlags() {
	flags.Parse(os.Args[1:])
}

// Indicates if the program is running on verbose mode.
func Verbose() bool {
	return *verbose
}

// Invite link of the network the node should connect to.
func InviteLink() string {
	return *inviteLink
}
