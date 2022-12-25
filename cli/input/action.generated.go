// Code generated by go run ./gen.go; DO NOT EDIT.

// Package input contains functionality for handling the CLi input
package input

import (
	"strconv"
)

// ActionInput represents schema of action.yaml for generating input for arguments
type ActionInput struct {
	
	ConfigFile string
	IgnoreConfigFile bool
	ProjectDirectory string
	ScanVulnerabilities bool
	ScanLicenses bool
}

// NewActionInput creates a new ActionInput instance
func NewActionInput(args []string) ActionInput {
	argCount := len(args)
	
	var configfile = "impilo.yml"
	if argCount > 0 {
		configfile = args[0] 
		
	} 
	var ignoreconfigfile = false
	if argCount > 1 {
		
		tmp,bErr := strconv.ParseBool(args[1])
		if bErr != nil {
			ignoreconfigfile = tmp 
		}
	} 
	var projectdirectory = ""
	if argCount > 2 {
		projectdirectory = args[2] 
		
	} 
	var scanvulnerabilities = true
	if argCount > 3 {
		
		tmp,bErr := strconv.ParseBool(args[3])
		if bErr != nil {
			scanvulnerabilities = tmp 
		}
	} 
	var scanlicenses = true
	if argCount > 4 {
		
		tmp,bErr := strconv.ParseBool(args[4])
		if bErr != nil {
			scanlicenses = tmp 
		}
	} 

	return ActionInput{
		
		ConfigFile: configfile,
		IgnoreConfigFile: ignoreconfigfile,
		ProjectDirectory: projectdirectory,
		ScanVulnerabilities: scanvulnerabilities,
		ScanLicenses: scanlicenses,
	}
}
