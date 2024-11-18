package deps

import "github.com/guackamolly/zero-monitor/internal/data/models"

// Returns a list of all packages installed on the system, as reported by the package manager.
// Error is not nil if the system package manager is not recognized, or an error has occurred when invoking
// the package manager.
func List() ([]models.Package, error) {
	return list()
}

// Returns the package manager which the system supports.
// Error is not nil if the system package manager is not recognized.
//
// Example: On Debian, it returns "dpkg".
func PackageManager() (string, error) {
	return packageManager()
}
