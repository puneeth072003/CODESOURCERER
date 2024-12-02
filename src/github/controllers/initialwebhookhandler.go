package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

func WebhookHandler(c *gin.Context) {
	event := c.GetHeader("X-GitHub-Event")

	if event == "pull_request" {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Printf("Unable to read request body: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to read request body",
			})
			return
		}

		var prEvent map[string]interface{}
		if err := json.Unmarshal(body, &prEvent); err != nil {
			log.Printf("Unable to unmarshal pull request event: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to process pull request event",
			})
			return
		}

		action := prEvent["action"].(string)
		merged := prEvent["pull_request"].(map[string]interface{})["merged"].(bool)
		baseBranch := prEvent["pull_request"].(map[string]interface{})["base"].(map[string]interface{})["ref"].(string)

		if action == "closed" && merged && baseBranch == "testing" {
			log.Printf("Pull request merged into 'testing' branch")

			repoOwner := prEvent["repository"].(map[string]interface{})["owner"].(map[string]interface{})["login"].(string)
			repoName := prEvent["repository"].(map[string]interface{})["name"].(string)
			pullRequestNumber := int(prEvent["number"].(float64))
			commitSHA := prEvent["pull_request"].(map[string]interface{})["merge_commit_sha"].(string)

			mergeID := fmt.Sprintf("merge_%s_%d", commitSHA, pullRequestNumber)

			// Fetch PR description and dependencies
			prDescription, err := FetchPullRequestDescription(repoOwner, repoName, pullRequestNumber)
			if err != nil {
				log.Printf("Unable to fetch pull request description: %v", err)
			}

			dependencies, context := ParsePRDescription(prDescription)
			log.Printf("Dependencies from PR Description: %v", dependencies)
			log.Printf("Context: %s", context)

			mergeData := map[string]interface{}{
				"merge_id":     mergeID,
				"commit_sha":   commitSHA,
				"pull_request": pullRequestNumber,
				"context":      context,
				"files":        []map[string]interface{}{},
			}

			changedFiles, err := fetchPullRequestFiles(repoOwner, repoName, pullRequestNumber)
			if err != nil {
				log.Printf("Unable to fetch changed files: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Unable to fetch changed files",
				})
				return
			}

			// Use WaitGroup to wait for all file processing
			var wg sync.WaitGroup
			fileResults := make(chan map[string]interface{}, len(changedFiles))

			for _, file := range changedFiles {
				wg.Add(1)

				go func(file map[string]interface{}) {
					defer wg.Done()

					filePath := file["filename"].(string)
					fileContent, err := FetchFileContentFromGitHub(repoOwner, repoName, commitSHA, filePath)
					if err != nil {
						log.Printf("Unable to fetch file content for %s: %v", filePath, err)
						fileContent = "Error fetching content"
					}

					fileDependencies := filterDependenciesForFile(filePath, dependencies)
					formattedDeps := formatDependencies(fileDependencies, repoOwner, repoName, commitSHA)

					fileResults <- map[string]interface{}{
						"path":         filePath,
						"content":      fileContent,
						"dependencies": formattedDeps,
					}
				}(file)
			}

			wg.Wait()
			close(fileResults)

			// Collect all file results
			for result := range fileResults {
				mergeData["files"] = append(mergeData["files"].([]map[string]interface{}), result)
			}

			preprocessedData := gin.H{
				"message": "Pull request merged into 'testing' branch and files processed",
				"data":    mergeData,
			}

			fmt.Println(preprocessedData)
			c.JSON(http.StatusOK, preprocessedData)
		} else {
			c.Status(http.StatusNoContent)
		}
	} else {
		c.Status(http.StatusNoContent)
	}
}

// Fetch the list of changed files in the pull request
func fetchPullRequestFiles(owner, repo string, prNumber int) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/files", owner, repo, prNumber)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API responded with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var files []map[string]interface{}
	if err := json.Unmarshal(body, &files); err != nil {
		return nil, err
	}

	return files, nil
}

// Function to filter dependencies for a specific file
func filterDependenciesForFile(filePath string, dependencies map[string][]string) []string {
	// Check if specific dependencies are mentioned for the file
	if deps, exists := dependencies[filePath]; exists && len(deps) > 0 {
		return deps
	}

	// Default to no dependencies if not specified
	log.Printf("No specific dependencies found for file: %s. Using the file itself.", filePath)
	return []string{}
}

// formatDependencies converts dependencies into a detailed slice with name and content.
func formatDependencies(dependencies []string, owner, repo, commitSHA string) []map[string]string {
	var formattedDependencies []map[string]string

	for _, dependency := range dependencies {
		depContent, err := FetchFileContentFromGitHub(owner, repo, commitSHA, dependency)
		if err != nil {
			log.Printf("Unable to fetch content for dependency %s: %v", dependency, err)
			depContent = "Error fetching content"
		}

		formattedDependencies = append(formattedDependencies, map[string]string{
			"name":    dependency,
			"content": depContent,
		})
	}

	return formattedDependencies
}
