// This is a special package that allows embedding the public folder as embedded file system.
// This package should only be imported by master, otherwise it will increase the final binary size.
package public

import (
	"embed"
)

//go:embed *
var FS embed.FS
