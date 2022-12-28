package main

import (
	"dburriss/impilo_gh/domain"
	"dburriss/impilo_gh/npm"
)

func BuildCommands(config domain.Config) []domain.Command {
	commands := []domain.Command{}

	if config.ScanVulnerabilities {
		commands = append(commands, npm.MakeScanVulnerabilitiesCommand(config)...)
	}

	if config.ScanLicenses {

	}

	return commands
}
