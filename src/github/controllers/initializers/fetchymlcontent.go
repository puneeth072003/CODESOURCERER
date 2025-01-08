package initializers

import (
	"log"

	"gopkg.in/yaml.v3"
)

// YMLConfig represents the overall structure of the YAML
type YMLConfig struct {
	Configuration ymlConfiguration `yaml:"configuration"`
	Environment   ymlEnvironment   `yaml:"environment"`
	Caching       ymlCaching       `yaml:"caching"`
}

// Configuration holds the YAML configuration fields
type ymlConfiguration struct {
	TestDirectory    string `yaml:"test-directory"`
	Comments         string `yaml:"comments"`
	TestingBranch    string `yaml:"testing-branch"`
	TestingFramework string `yaml:"testing-framework"`
	WaterMark        string `yaml:"water-mark"`
}

// Environment holds environment-specific configurations
type ymlEnvironment struct {
	PythonVersion string `yaml:"python-version"`
}

// Caching holds caching-related configurations
type ymlCaching struct {
	Enabled      bool   `yaml:"enabled"`
	RedisCaching string `yaml:"redis-caching"`
}

const configFilePath = "codesourcerer-config.yml"

var defaultConfig = YMLConfig{
	Configuration: ymlConfiguration{TestDirectory: "/tests", Comments: "on", TestingBranch: "testing", TestingFramework: "pytest", WaterMark: "on"},
	Environment:   ymlEnvironment{PythonVersion: "3.12"},
	Caching:       ymlCaching{Enabled: false, RedisCaching: "off"},
}

// FetchConfig fetches Application Config contents and returns it as a structure
func FetchConfig(owner, repo, commitSHA string) YMLConfig {
	// Fetch the file content from GitHub
	content, err := FetchFileContentFromGitHub(owner, repo, commitSHA, configFilePath)
	if err != nil {
		log.Printf("unable to find config file. Using the Default Configuration. Error: %v", err)
		return defaultConfig
	}

	// Parse YAML content
	var config YMLConfig
	err = yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		log.Printf("Failed to parse yml file. Using the Default Configuration. Error: %v", err)
		return defaultConfig
	}

	return config
}
