package main

import (
	"fmt"
	"os"
)

//go:generate go run ./gen/gen.go

type Project struct {
	Name                    string
	Tags                    []string
	PackageManager          string
	TargetDirectory         string
	AllowedLicenses         []string
	SkipScanVulnerabilities bool
	SkipScanLicenses        bool
}

type Config struct {
	AllowedLicenses     []string
	PackageLicenseMap   map[string]string
	Tags                []string
	ScanVulnerabilities bool
	ScanLicenses        bool
}

func newConfig() Config {
	return Config{
		AllowedLicenses:     []string{},
		PackageLicenseMap:   map[string]string{},
		Tags:                []string{},
		ScanVulnerabilities: true,
		ScanLicenses:        false,
	}
}

func main() {
	// input
	args := os.Args[1:]
	actionInput := NewActionInput(args)
	fmt.Println("Input", actionInput)

	// load config file if exists
	config := newConfig()
	if !actionInput.IgnoreConfigFile {
		fmt.Println("Config", config)
	} else {
		fmt.Println("Ignore config file.")
	}

	// override with input

	// parse to commands

	// execute commands
}
