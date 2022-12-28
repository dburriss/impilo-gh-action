package domain

// Input is generated
// TODO: move type to this package

// ====================
// Config
// ====================

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

func NewConfig() Config {
	return Config{
		AllowedLicenses:     []string{},
		PackageLicenseMap:   map[string]string{},
		Tags:                []string{},
		Projects:            []Project{},
		ScanVulnerabilities: true,
		ScanLicenses:        false,
	}
}

func (input ActionInput) ToConfig() Config {
	project := Project{
		Name:            input.ProjectName,
		PackageManager:  input.PackageManager,
		TargetDirectory: input.TargetDirectory,
	}
	inputConfig := Config{
		ScanVulnerabilities: input.ScanVulnerabilities,
		ScanLicenses:        input.ScanLicenses,
		Projects:            []Project{project},
	}

	return inputConfig
}

func Merge2Configs(baseConfig Config, overrideWith Config) Config {
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

func MergeConfigs(configs []Config) Config {
	nrConfigs := len(configs)
	switch nrConfigs {
	case 0:
		return NewConfig()
	case 1:
		return configs[0]
	case 2:
		newConfig := Merge2Configs(configs[0], configs[1])
		return newConfig
	default:
		newConfig := Merge2Configs(configs[0], configs[1])
		configs := append([]Config{newConfig}, configs[2:]...)
		return MergeConfigs(configs)
	}
}
