package deps

import "github.com/guackamolly/zero-monitor/internal/data/models"

// Returns a list of all packages installed on the system, as reported by the package manager.
func List() ([]models.Package, error) {
	return list()
}
