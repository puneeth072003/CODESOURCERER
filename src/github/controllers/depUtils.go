package controllers

import "log"

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
