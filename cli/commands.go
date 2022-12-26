package main

import (
	"fmt"
	"log"
	"os/exec"
)

type Report interface {
	Run()
}

type ScanVulnerabilitiesStdOutReport struct {
}

func (data ScanVulnerabilitiesStdOutReport) Report() {
	fmt.Println("Scan vulnerabilities report")
}

type Command interface {
	Execute() []Report
}

type ScanVulnerabilitiesCommand struct {
	packageManager  string
	targetDirectory string
}

func (data ScanVulnerabilitiesCommand) Execute() []Report {
	fmt.Println("Package manager", data.packageManager)
	fmt.Println("target directory", data.targetDirectory)
	cmdName := "npm"
	cmdArgs := []string{""}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Dir = data.targetDirectory
	err := cmd.Run()

	if err != nil {
		log.Fatalf("Command %s failed.\n", cmdName)
		log.Fatalln(err)
		panic("Scan Vulnerabilities failed.")
	}

	var result []Report

	return result
}

func makeScanVulnerabilitiesCommand(config Config) []Command {
	commands := []Command{}

	var cmd Command = ScanVulnerabilitiesCommand{}

	commands = append(commands, cmd)

	return commands
}

func BuildCommands(config Config) []Command {
	commands := []Command{}

	if config.ScanVulnerabilities {
		commands = append(commands, makeScanVulnerabilitiesCommand(config)...)
	}

	if config.ScanLicenses {

	}

	return commands
}
