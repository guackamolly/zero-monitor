//go:build windows
// +build windows

package deps

//TODO: UNTESTED! NEED TESTER FRIENDS WITH A WINDOWS MACHINE :)

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

var (
	supportedPkgManagers = []string{"winget"} // Only consider winget for Windows
	systemPkgManager     string
)

func list() ([]models.Package, error) {
	pm, err := packageManager()
	if err != nil {
		return nil, err
	}

	switch pm {
	case "winget":
		return listWinget()
	default:
		return nil, fmt.Errorf("unsupported package manager: %s", pm)
	}
}

func packageManager() (string, error) {
	if len(systemPkgManager) != 0 {
		return systemPkgManager, nil
	}

	for _, pm := range supportedPkgManagers {
		if _, err := exec.LookPath(pm); err == nil {
			return pm, nil
		}
	}

	return "", fmt.Errorf("winget not found")
}

func listWinget() ([]models.Package, error) {
	// Use winget list --format json to get a list of installed packages in JSON format
	cmd := exec.Command("winget", "list", "--format", "json")

	bs, err := cmd.Output()
	if err != nil {
		// Check for specific error (e.g., winget not found) and handle accordingly
		return nil, fmt.Errorf("failed to execute winget: %w", err)
	}

	var packages []wingetPackage
	err = json.Unmarshal(bs, &packages)
	if err != nil {
		return nil, fmt.Errorf("failed to parse winget output: %w", err)
	}

	convertedPkgs := []models.Package{}
	for _, pkg := range packages {
		convertedPkgs = append(convertedPkgs, models.NewPackage(pkg.ID, "", pkg.Version))
	}

	return convertedPkgs, nil
}

// Winget package structure (may need adjustments based on actual winget output)
type wingetPackage struct {
	ID      string `json:"id"`
	Version string `json:"version"`
	// ... other fields as needed
}
