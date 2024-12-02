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

// formatDependencies converts dependencies into a detailed slice with name and content.
// func formatDependencies(dependencies []string, owner, repo, commitSHA string) []map[string]string {
// 	var formattedDependencies []map[string]string

// 	for _, dependency := range dependencies {
// 		depContent, err := FetchFileContentFromGitHub(owner, repo, commitSHA, dependency)
// 		if err != nil {
// 			log.Printf("Unable to fetch content for dependency %s: %v", dependency, err)
// 			depContent = "Error fetching content"
// 		}

// 		formattedDependencies = append(formattedDependencies, map[string]string{
// 			"name":    dependency,
// 			"content": depContent,
// 		})
// 	}

// 	return formattedDependencies
// }
