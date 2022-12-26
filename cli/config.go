package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

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
	Projects            []Project
}

func newConfig() Config {
	return Config{
		AllowedLicenses:     []string{},
		PackageLicenseMap:   map[string]string{},
		Tags:                []string{},
		Projects:            []Project{},
		ScanVulnerabilities: true,
		ScanLicenses:        false,
	}
}

func (input ActionInput) toConfig() Config {
	project := Project{
		Name:            input.ProjectName,
		PackageManager:  input.PackageManager,
		TargetDirectory: input.TargetDirectory,
	}
	inputConfig := Config{
		ScanVulnerabilities: input.ScanLicenses,
		ScanLicenses:        input.ScanLicenses,
		Projects:            []Project{project},
	}

	return inputConfig
}

// loads the given filePath into a Config. If error returns Config{}.
func loadConfigFile(filePath string) (Config, error) {

	// read config file
	content, err := ioutil.ReadFile(filePath)

	// Deserialize
	configFile := Config{}
	if err == nil {
		err = yaml.Unmarshal(content, &configFile)
	}

	return configFile, err
}

// todo: TDD merge
func merge2Configs(baseConfig Config, overrideWith Config) Config {
	newConfig := Config{}
	// AllowedLicenses
	if overrideWith.AllowedLicenses != nil {
		newConfig.AllowedLicenses = overrideWith.AllowedLicenses
	} else if baseConfig.AllowedLicenses != nil {
		newConfig.AllowedLicenses = baseConfig.AllowedLicenses
	}

	// PackageLicenseMap
	if overrideWith.PackageLicenseMap != nil {
		newConfig.PackageLicenseMap = overrideWith.PackageLicenseMap
	} else if baseConfig.PackageLicenseMap != nil {
		newConfig.PackageLicenseMap = baseConfig.PackageLicenseMap
	}

	// Tags
	if overrideWith.Tags != nil {
		newConfig.Tags = overrideWith.Tags
	} else if baseConfig.Tags != nil {
		newConfig.Tags = baseConfig.Tags
	}

	// Projects
	if overrideWith.Projects != nil {
		newConfig.Projects = overrideWith.Projects
	} else if baseConfig.Projects != nil {
		newConfig.Projects = baseConfig.Projects
	}

	// ScanVulnerabilities
	newConfig.ScanVulnerabilities = overrideWith.ScanVulnerabilities

	// ScanLicenses
	newConfig.ScanLicenses = overrideWith.ScanLicenses

	return newConfig
}

func mergeConfigs(configs []Config) Config {
	nrConfigs := len(configs)
	switch nrConfigs {
	case 0:
		return newConfig()
	case 1:
		return configs[0]
	case 2:
		newConfig := merge2Configs(configs[0], configs[1])
		return newConfig
	default:
		newConfig := merge2Configs(configs[0], configs[1])
		configs := append([]Config{newConfig}, configs[2:]...)
		return mergeConfigs(configs)
	}
}

func BuildConfig(actionInput ActionInput) Config {
	configs := []Config{newConfig()}
	if actionInput.IgnoreConfigFile {
		fmt.Println("Ignore config file.")
		configs = append(configs, actionInput.toConfig())
	} else {
		// input to config
		configs = append(configs, actionInput.toConfig())
		// load config file
		fileConfig, err := loadConfigFile(actionInput.ConfigFile)
		if err != nil {
			log.Fatalf("Error loading config file %s:  \n", err)
			log.Fatal(err)
		}
		// fileConfig overrides input config
		configs = append(configs, fileConfig)
	}

	return mergeConfigs(configs)
}
