//go:build linux
// +build linux

package deps

import (
	"bytes"
	"io"
	"os/exec"
	"strings"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

func list() ([]models.Package, error) {
	return listDpkg()
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
