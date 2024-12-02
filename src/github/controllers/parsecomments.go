package controllers

import "strings"

// Function to parse the PR comment for dependencies
func ParseCommentForDependencies(comment string) (map[string][]string, string) {
	dependencies := make(map[string][]string)
	var context string

	// Split comment into lines
	lines := strings.Split(comment, "\n")

	// Parse dependencies
	var currentFile string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "@file:") {
			// Extract the file path
			currentFile = strings.TrimSpace(strings.TrimPrefix(line, "@file:"))
		} else if strings.HasPrefix(line, "@dependencies:") {
			// Extract dependencies for the current file
			dependencyList := strings.TrimSpace(strings.TrimPrefix(line, "@dependencies:"))
			if currentFile != "" {
				dependencies[currentFile] = append(dependencies[currentFile], strings.Split(dependencyList, ",")...)
			}
		} else if strings.HasPrefix(line, "@context:") {
			// Extract the context
			context = strings.TrimSpace(strings.TrimPrefix(line, "@context:"))
		}
	}

	return dependencies, context
}
