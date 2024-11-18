//go:build linux
// +build linux

package deps

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

var (
	supportedPkgManagers = []string{"dpkg", "yum", "dnf", "pacman"}
	systemPkgManager     string
)

func list() ([]models.Package, error) {
	pm, err := packageManager()
	if err != nil {
		return nil, err
	}

	switch pm {
	case "dpkg":
		return listDpkg()
	default:
		return nil, fmt.Errorf("function not implemented yet for pm: %s", pm)
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

func listDpkg() ([]models.Package, error) {
	cmd := exec.Command("dpkg", "-l")

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

		if !strings.HasPrefix(l, "ii ") {
			continue
		}

		parts := strings.Fields(l)
		if len(parts) < 3 {
			continue
		}

		name := parts[1]
		version := parts[2]
		description := strings.Join(parts[4:], " ")
		pkgs = append(pkgs, models.NewPackage(name, description, version))
	}

	return pkgs, nil
}
