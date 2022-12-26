package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Report interface {
	Run()
}

type ScanVulnerabilitiesStdOutReport struct {
	rawData string
}

func (data ScanVulnerabilitiesStdOutReport) Run() {
	fmt.Println("Scan vulnerabilities report")
	fmt.Println(data.rawData)
}

type Command interface {
	Execute() []Report
}

type ScanVulnerabilitiesCommand struct {
	PackageManager  string
	TargetDirectory string
}

func (data ScanVulnerabilitiesCommand) Execute() []Report {
	fmt.Println("Package manager", data.PackageManager)
	fmt.Println("Target directory", data.TargetDirectory)
	cmdName := "npm"
	cmdArgs := []string{"audit", "--json", "--silent"}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Dir = data.TargetDirectory
	var outbuf, errbuf strings.Builder // or bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	println(cmd.String())

	err := cmd.Run()
	stderr := errbuf.String()
	if len(stderr) > 0 && err != nil {
		log.Fatalf("Command %s failed.\n", cmdName)
		log.Fatalln(err)
		panic("Scan Vulnerabilities failed.")
	}

	var reports []Report
	//json, err := gabs.ParseJSON([]byte(outbuf.String()))
	var stdoutReport Report = ScanVulnerabilitiesStdOutReport{
		rawData: outbuf.String(),
	}
	reports = append(reports, stdoutReport)
	return reports
}

func makeScanVulnerabilitiesCommand(config Config) []Command {
	commands := []Command{}
	for _, project := range config.Projects {
		var cmd Command = ScanVulnerabilitiesCommand{
			TargetDirectory: project.TargetDirectory,
			PackageManager:  project.PackageManager,
		}
		commands = append(commands, cmd)
	}

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
