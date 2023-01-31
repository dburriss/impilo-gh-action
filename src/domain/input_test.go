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
	useConfigFile := false
	skipScanVulnerabilities := false
	scanLicenses := false
	opts := NewActionInputOpts()
	assert.EqualValues(t, configFile, opts.ConfigFile, fmt.Sprintf("Expected %s, instead got %s", configFile, opts.ConfigFile))
	assert.EqualValues(t, useConfigFile, opts.UseConfigFile, fmt.Sprintf("Expected %t, instead got %t", useConfigFile, opts.UseConfigFile))
	assert.EqualValues(t, packageManager, opts.PackageManager, fmt.Sprintf("Expected %s, instead got %s", "npm", opts.PackageManager))
	assert.EqualValues(t, projectName, opts.ProjectName, fmt.Sprintf("Expected %s, instead got %s", projectName, opts.ProjectName))
	assert.EqualValues(t, targetDir, opts.TargetDirectory, fmt.Sprintf("Expected %s, instead got %s", targetDir, opts.TargetDirectory))
	assert.EqualValues(t, skipScanVulnerabilities, opts.SkipScanVulnerabilities, fmt.Sprintf("Expected %t, instead got %t", skipScanVulnerabilities, opts.SkipScanVulnerabilities))
	assert.EqualValues(t, scanLicenses, opts.ScanLicenses, fmt.Sprintf("Expected %t, instead got %t", scanLicenses, opts.ScanLicenses))
}
