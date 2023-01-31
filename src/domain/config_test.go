package domain

import (
	"reflect"
	"testing"
)

func canBeNil(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return true
	default:
		return false
	}
}

func TestNewConfigShouldInitAllFields(t *testing.T) {
	config := NewConfig()
	v := reflect.ValueOf(config)
	n := v.NumField()

	for i := 0; i < n; i++ {
		field := v.Field(i)
		if canBeNil(field) && field.IsNil() {
			tp := reflect.TypeOf(config)
			name := tp.Field(i).Name
			t.Errorf("Expected %s to not be nil", name)
		}
	}
}

func TestMergeSetsAllowedLicenses(t *testing.T) {
	expected := []string{"MIT"}
	baseConfig := NewConfig()
	overrideWith := NewConfig()
	overrideWith.AllowedLicenses = expected

	config := MergeConfigs([]Config{baseConfig, overrideWith})

	if !reflect.DeepEqual(config.AllowedLicenses, expected) {
		t.Errorf("Expected %+q, instead got %+q", expected, config.AllowedLicenses)
	}
}

func TestMergePackageLicenseMap(t *testing.T) {
	expected := map[string]string{"go": "BSD-3"}
	baseConfig := NewConfig()
	overrideWith := NewConfig()
	overrideWith.PackageLicenseMap = expected

	config := MergeConfigs([]Config{baseConfig, overrideWith})

	if !reflect.DeepEqual(config.PackageLicenseMap, expected) {
		t.Errorf("Expected %+q, instead got %+q", expected, config.PackageLicenseMap)
	}
}

func TestMergeTags(t *testing.T) {
	expected := []string{"app:impilo"}
	baseConfig := NewConfig()
	overrideWith := NewConfig()
	overrideWith.Tags = expected

	config := MergeConfigs([]Config{baseConfig, overrideWith})

	if !reflect.DeepEqual(config.Tags, expected) {
		t.Errorf("Expected %+q, instead got %+q", expected, config.Tags)
	}
}

func TestMergeProjects(t *testing.T) {
	project := Project{
		Name: "test-app",
	}
	expected := []Project{project}
	baseConfig := NewConfig()
	overrideWith := NewConfig()
	overrideWith.Projects = expected

	config := MergeConfigs([]Config{baseConfig, overrideWith})

	if !reflect.DeepEqual((config.Projects), expected) {
		t.Errorf("Expected %+v, instead got %+v", expected, config.Projects)
	}
}

func TestMergeScanVulnerabilities(t *testing.T) {
	expected := true
	baseConfig := NewConfig()
	baseConfig.ScanVulnerabilities = false
	overrideWith := NewConfig()
	overrideWith.ScanVulnerabilities = expected

	config := MergeConfigs([]Config{baseConfig, overrideWith})

	if config.ScanVulnerabilities != expected {
		t.Errorf("Expected %+v, instead got %+v", expected, config.ScanVulnerabilities)
	}
}

func TestMergeScanLicenses(t *testing.T) {
	expected := true
	baseConfig := NewConfig()
	baseConfig.ScanLicenses = false
	overrideWith := NewConfig()
	overrideWith.ScanLicenses = expected

	config := MergeConfigs([]Config{baseConfig, overrideWith})

	if config.ScanLicenses != expected {
		t.Errorf("Expected %+v, instead got %+v", expected, config.ScanLicenses)
	}
}

func TestToConfigProjectName(t *testing.T) {
	expected := "test-app"
	input := NewActionInput([]string{})
	input.ProjectName = expected

	config := input.ToConfig()
	actual := config.Projects[0].Name

	if actual != expected {
		t.Errorf("Expected %+v, instead got %+v", expected, actual)
	}
}

func TestToConfigPackageManager(t *testing.T) {
	expected := "nuget"
	input := NewActionInput([]string{})
	input.PackageManager = expected

	config := input.ToConfig()
	actual := config.Projects[0].PackageManager

	if actual != expected {
		t.Errorf("Expected %+v, instead got %+v", expected, actual)
	}
}

func TestToConfigTargetDirectory(t *testing.T) {
	expected := "nuget"
	input := NewActionInput([]string{})
	input.TargetDirectory = expected

	config := input.ToConfig()
	actual := config.Projects[0].TargetDirectory

	if actual != expected {
		t.Errorf("Expected %+v, instead got %+v", expected, actual)
	}
}

func TestToConfigScanVulnerabilities(t *testing.T) {
	expected := true
	input := NewActionInput([]string{})
	input.SkipScanVulnerabilities = true

	config := input.ToConfig()
	actual := config.ScanVulnerabilities

	if actual != expected {
		t.Errorf("Expected %+v, instead got %+v", expected, actual)
	}
}

func TestToConfigScanLicenses(t *testing.T) {
	expected := true
	input := NewActionInput([]string{})
	input.ScanLicenses = true

	config := input.ToConfig()
	actual := config.ScanLicenses

	if actual != expected {
		t.Errorf("Expected %+v, instead got %+v", expected, actual)
	}
}
