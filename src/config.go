package main

import (
	"dburriss/impilo_gh/domain"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// loads the given filePath into a Config. If error returns Config{}.
func loadConfigFile(filePath string) (domain.Config, error) {

	// read config file
	content, err := ioutil.ReadFile(filePath)

	// Deserialize
	configFile := domain.Config{}
	if err == nil {
		err = yaml.Unmarshal(content, &configFile)
	}

	return configFile, err
}

func BuildConfig(actionInput domain.ActionInput) domain.Config {
	configs := []domain.Config{domain.NewConfig()}
	if actionInput.IgnoreConfigFile {
		fmt.Println("Ignore config file.")
		configs = append(configs, actionInput.ToConfig())
	} else {
		// input to config
		configs = append(configs, actionInput.ToConfig())
		// load config file
		fileConfig, err := loadConfigFile(actionInput.ConfigFile)
		if err != nil {
			log.Fatalf("Error loading config file %s:  \n", err)
			log.Fatal(err)
		}
		// fileConfig overrides input config
		configs = append(configs, fileConfig)
	}

	return domain.MergeConfigs(configs)
}
