package main

import (
	"flag"
	"os"
	"slices"

	"github.com/kardianos/service"
)

const (
	ExitCodeNoServiceSystemAvailable   = 60
	ExitCodeRequiresElevatedPrivileges = 61
	ExitCodeNotStarted                 = 62
)

var (
	flagName        = flag.String("name", "", "name of the service to launch at startup")
	flagDescription = flag.String("description", "", "small description of what the service does")
	flagExecPath    = flag.String("exec", "", "absolute path of the program to execute at startup")
	flagUser        = flag.String("user", "", "id of the user to launch the service")
)

func init() {
	flag.BoolFunc("help", "prints program usage", func(s string) error {
		println("Utility program that helps creating a system service for launching a program at startup.")
		return nil
	})
	flag.Parse()

	if slices.Contains([]string{*flagName, *flagDescription, *flagExecPath, *flagUser}, "") {
		flag.Usage()
		os.Exit(2)
	}
}

func main() {
	s := launchService()
	if isInstalled(s) {
		println("service already installed")
		return
	}

	err := s.Install()
	if err != nil {
		println("requires elevated privileges")
		os.Exit(ExitCodeRequiresElevatedPrivileges)
	}

	err = s.Start()
	if err != nil {
		println("service not started: " + err.Error())
		os.Exit(ExitCodeNotStarted)
	}

	println("service installed and managed by: " + s.Platform())
}

func isInstalled(s service.Service) bool {
	_, err := s.Status()
	return err != service.ErrNotInstalled
}

func launchService() service.Service {
	cfg := &service.Config{
		Name:        *flagName,
		DisplayName: *flagName,
		Description: *flagDescription,
		Executable:  *flagExecPath,
		UserName:    *flagUser,
	}

	s, err := service.New(nil, cfg)
	if err == service.ErrNoServiceSystemDetected {
		println("no service manager available")
		os.Exit(ExitCodeNoServiceSystemAvailable)
	}

	return s
}
