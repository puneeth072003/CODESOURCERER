package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

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

			// Fetch PR comment to parse dependencies and context
			prComment, err := FetchPullRequestComment(repoOwner, repoName, pullRequestNumber)
			if err != nil {
				log.Printf("Unable to fetch pull request comment: %v", err)
			} else if prComment != "" {
				log.Printf("PR Comment: %s", prComment)
			} else {
				log.Println("No pull request comment found")
			}

			dependencies, context := ParseCommentForDependencies(prComment)
			log.Printf("Dependencies from PR Comment: %v", dependencies)
			log.Printf("Context: %s", context)

			mergeData := map[string]interface{}{
				"merge_id":     mergeID,
				"commit_sha":   commitSHA,
				"pull_request": pullRequestNumber,
				"files":        []map[string]interface{}{},
			}

			// Fetch changed files from the PR
			changedFiles, err := fetchPullRequestFiles(repoOwner, repoName, pullRequestNumber)
			if err != nil {
				log.Printf("Unable to fetch changed files: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Unable to fetch changed files",
				})
				return
			}

			// Process each changed file
			for _, file := range changedFiles {
				filePath := file["filename"].(string)
				fileContent, err := FetchFileContentFromGitHub(repoOwner, repoName, commitSHA, filePath)
				if err != nil {
					log.Printf("Unable to fetch file content for %s: %v", filePath, err)
					continue
				}

				// Filter dependencies relevant to the current file
				fileDependencies := filterDependenciesForFile(filePath, dependencies)

				mergeData["files"] = append(mergeData["files"].([]map[string]interface{}), map[string]interface{}{
					"path":         filePath,
					"content":      fileContent,
					"dependencies": fileDependencies,
				})

				log.Printf("Processed file: %s with dependencies: %v", filePath, fileDependencies)
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Pull request merged into 'testing' branch and files processed",
				"data":    mergeData,
			})
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
	if deps, found := dependencies[filePath]; found {
		return deps
	}
	return []string{}
}
