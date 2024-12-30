package initializers

import (
	"gopkg.in/yaml.v3"
)

// Config holds the YAML configuration fields
type Config struct {
	TestDirectory    string `yaml:"test-directory"`
	Comments         string `yaml:"comments"`
	TestingBranch    string `yaml:"testing-branch"`
	TestingFramework string `yaml:"testing-framework"`
	WaterMark        string `yaml:"water-mark"`
	RedisCaching     string `yaml:"redis-caching"`
}

// FetchAndReturnYAMLContents fetches YAML contents and returns it as a structure
func FetchAndReturnYAMLContents(owner, repo, commitSHA, filePath string) (Config, error) {
	// Fetch the file content from GitHub
	content, err := FetchFileContentFromGitHub(owner, repo, commitSHA, filePath)
	if err != nil {
		return Config{}, err
	}

	// Parse YAML content
	var config Config
	err = yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
