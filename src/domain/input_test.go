package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// I check that the generated defaults for NewActionInputOpts are as expected
func TestInputHasActionDefaults(t *testing.T) {
	configFile := "impilo.yml"
	packageManager := "npm"
	projectName := "app"
	targetDir := ""
	ignoreConfigFile := true
	scanVulnerabilities := true
	scanLicenses := false
	opts := NewActionInputOpts()
	assert.EqualValues(t, configFile, opts.ConfigFile, fmt.Sprintf("Expected %s, instead got %s", configFile, opts.ConfigFile))
	assert.EqualValues(t, ignoreConfigFile, opts.IgnoreConfigFile, fmt.Sprintf("Expected %t, instead got %t", ignoreConfigFile, opts.IgnoreConfigFile))
	assert.EqualValues(t, packageManager, opts.PackageManager, fmt.Sprintf("Expected %s, instead got %s", "npm", opts.PackageManager))
	assert.EqualValues(t, projectName, opts.ProjectName, fmt.Sprintf("Expected %s, instead got %s", projectName, opts.ProjectName))
	assert.EqualValues(t, targetDir, opts.TargetDirectory, fmt.Sprintf("Expected %s, instead got %s", targetDir, opts.TargetDirectory))
	assert.EqualValues(t, scanVulnerabilities, opts.ScanVulnerabilities, fmt.Sprintf("Expected %t, instead got %t", scanVulnerabilities, opts.ScanVulnerabilities))
	assert.EqualValues(t, scanLicenses, opts.ScanLicenses, fmt.Sprintf("Expected %t, instead got %t", scanLicenses, opts.ScanLicenses))
}
