package banner

import (
	"fmt"

	"github.com/labstack/gommon/color"
)

var banner = fmt.Sprintf(
	`
    @@
   @  @
@@@    @    @@@
        @  @
         @@

One-click lightweight server monitor tool based on ZeroMQ protocol.
%s
___________________________________________________________________
`,
	color.Blue("https://github.com/guackamolly/zero-monitor"),
)

func Print() {
	println(banner)
}
