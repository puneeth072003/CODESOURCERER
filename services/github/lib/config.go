package lib

import (
	"log"

	pb "github.com/codesourcerer-bot/proto/generated"

	"gopkg.in/yaml.v3"
)

// YMLConfig represents the overall structure of the YAML
type YMLConfig struct {
	Configuration ymlConfiguration  `yaml:"configuration"`
	Environment   ymlEnvironment    `yaml:"environment"`
	Caching       ymlCaching        `yaml:"caching"`
	Extras        map[string]string `yml:"extras"`
}

// Configuration holds the YAML configuration fields
type ymlConfiguration struct {
	TestDirectory    string `yaml:"test-directory"`
	Comments         bool   `yaml:"comments"`
	TestingBranch    string `yaml:"testing-branch"`
	TestingFramework string `yaml:"testing-framework"`
	WaterMark        bool   `yaml:"water-mark"`
}

// Environment holds environment-specific configurations
type ymlEnvironment struct {
	PythonVersion float32 `yaml:"python-version"`
}

// Caching holds caching-related configurations
type ymlCaching struct {
	Enabled      bool `yaml:"enabled"`
	RedisCaching bool `yaml:"redis-caching"`
}

const configFilePath = "codesourcerer-config.yml"

var defaultConfig = YMLConfig{
	Configuration: ymlConfiguration{TestDirectory: "/tests", Comments: true, TestingBranch: "testing", TestingFramework: "pytest", WaterMark: true},
	Environment:   ymlEnvironment{PythonVersion: 3.12},
	Caching:       ymlCaching{Enabled: false, RedisCaching: false},
	Extras:        nil,
}

// FetchConfig fetches Application Config contents and returns it as a structure
func FetchYmlConfig(owner, repo, commitSHA string) YMLConfig {
	// Fetch the file content from GitHub
	content, err := FetchFileFromGitHub(owner, repo, commitSHA, configFilePath)
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

func GetGenerationOptions(ymlConfig YMLConfig) *pb.Configuration {

	basicConfig := pb.BasicConfig{
		TestDirectory:    ymlConfig.Configuration.TestDirectory,
		Comments:         ymlConfig.Configuration.Comments,
		TestingFramework: ymlConfig.Configuration.TestingFramework,
		WaterMark:        ymlConfig.Configuration.WaterMark,
	}

	return &pb.Configuration{
		Configuration: &basicConfig,
		Extras:        ymlConfig.Extras,
	}
}
