package utils

import (
	"log"
	"strings"
)

// Function to filter dependencies for a specific file
func FilterDependenciesForFile(filePath string, dependencies map[string][]string) []string {
	// Check if specific dependencies are mentioned for the file
	if deps, exists := dependencies[filePath]; exists && len(deps) > 0 {
		return deps
	}

	// Default to no dependencies if not specified
	log.Printf("No specific dependencies found for file: %s. Using the file itself.", filePath)
	return []string{}
}

func ParsePRDescription(description string) (map[string][]string, string) {
	lines := strings.Split(description, "\n")
	dependencies := make(map[string][]string)
	var context string
	var currentFile string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "$file:"):
			// Extract the file path
			currentFile = strings.TrimSpace(strings.TrimPrefix(line, "$file:"))
		case strings.HasPrefix(line, "$dependencies:"):
			// Extract dependencies for the current file
			if currentFile != "" {
				dependencyList := strings.TrimSpace(strings.TrimPrefix(line, "$dependencies:"))
				dependencies[currentFile] = append(dependencies[currentFile], strings.Split(dependencyList, ",")...)
			}
		case strings.HasPrefix(line, "$context:"):
			// Extract the context
			context = strings.TrimSpace(strings.TrimPrefix(line, "$context:"))
		}
	}

	return dependencies, context
}
