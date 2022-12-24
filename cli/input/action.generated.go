// Code generated by go run ./gen.go; DO NOT EDIT.

// Package input contains functionality for handling the CLi input
package input

// ActionInput represents schema of action.yaml for generating input for arguments
type ActionInput struct {
	
	Ignoreconfigfile string
	Projectdirectory string
	Scanvulnerabilities string
	Scanlicenses string
	Configfile string
}

// NewActionInput creates a new ActionInput instance
func NewActionInput(args []string) ActionInput {
	argCount := len(args)
	
	var ignoreconfigfile = "false"
	if argCount > 0 {
		ignoreconfigfile = args[0]
	} 
	var projectdirectory = "./"
	if argCount > 1 {
		projectdirectory = args[1]
	} 
	var scanvulnerabilities = "true"
	if argCount > 2 {
		scanvulnerabilities = args[2]
	} 
	var scanlicenses = "true"
	if argCount > 3 {
		scanlicenses = args[3]
	} 
	var configfile = "impilo.yml"
	if argCount > 4 {
		configfile = args[4]
	} 

	return ActionInput{
		
		Ignoreconfigfile: ignoreconfigfile,
		Projectdirectory: projectdirectory,
		Scanvulnerabilities: scanvulnerabilities,
		Scanlicenses: scanlicenses,
		Configfile: configfile,
	}
}
