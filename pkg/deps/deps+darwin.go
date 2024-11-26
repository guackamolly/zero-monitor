//go:build darwin
// +build darwin

package deps

//TODO: UNTESTED! NEED TESTER FRIENDS WITH A MAC MACHINE :)

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

var (
	supportedPkgManagers = []string{"brew"}
	systemPkgManager     string
)

func list() ([]models.Package, error) {
	pm, err := packageManager()
	if err != nil {
		return nil, err
	}

	switch pm {
	case "brew":
		return listBrew()
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

	return "", fmt.Errorf("system package manager is not recognized")
}

func listBrew() ([]models.Package, error) {
	// Use brew list to get installed packages
	cmd := exec.Command("brew", "list", "--formula") // List installed formulas

	bs, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	pkgs := []models.Package{}
	buf := bytes.NewBuffer(bs)
	for {
		l, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		// Extract name and version from brew list output format
		parts := strings.Fields(tr(l, "\t", " ")) // Handle potential tabs
		if len(parts) < 2 {
			continue
		}

		name := parts[0]
		version := parts[1]
		pkgs = append(pkgs, models.NewPackage(name, "", version)) // Description not available from brew list
	}

	return pkgs, nil
}

// Helper function to handle potential tabs in brew output (optional)
func tr(s string, old, new string) string {
	return strings.Replace(s, old, new, -1)
}
