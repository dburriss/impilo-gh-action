package npm

import (
	"dburriss/impilo_gh/domain"
	"testing"
)

func TestNpmVulnerabilitiesCommandReturnsReports(t *testing.T) {
	input := audit
	var cmd domain.Command = ScanNpmVulnerabilitiesCommand{PackageManager: "npm", TargetDirectory: ".", Load: func() string { return input }}

	reports := cmd.Execute()

	if len(reports) == 0 {
		t.Error("Expected more than zero reports")
	}
}

const audit = `
	{
		"auditReportVersion": 2,
		"vulnerabilities": {
		  "@nuxt/builder": {
			"name": "@nuxt/builder",
			"severity": "high",
			"isDirect": false,
			"via": [
			  "@nuxt/webpack"
			],
			"effects": [],
			"range": ">=2.14.0",
			"nodes": [
			  "node_modules/@nuxt/builder"
			],
			"fixAvailable": true
		  }
		},
		"metadata": {
			"vulnerabilities": {
			  "info": 0,
			  "low": 0,
			  "moderate": 0,
			  "high": 1,
			  "critical": 0,
			  "total": 1
			}
		}
	}
	`
