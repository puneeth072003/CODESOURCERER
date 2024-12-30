package initializers

import (
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

// FetchAndReturnYAMLContents fetches YAML contents and returns it as a structure
func FetchAndReturnYAMLContents(owner, repo, commitSHA, filePath string) (YMLConfig, error) {
	// Fetch the file content from GitHub
	content, err := FetchFileContentFromGitHub(owner, repo, commitSHA, filePath)
	if err != nil {
		return YMLConfig{}, err
	}

	// Parse YAML content
	var config YMLConfig
	err = yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		return YMLConfig{}, err
	}

	return config, nil
}
