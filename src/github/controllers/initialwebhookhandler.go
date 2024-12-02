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

			// Fetch PR description and dependencies
			prDescription, err := FetchPullRequestDescription(repoOwner, repoName, pullRequestNumber)
			if err != nil {
				log.Printf("Unable to fetch pull request description: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Unable to fetch pull request description",
				})
				return
			}

			dependencies, context := ParsePRDescription(prDescription)
			log.Printf("Dependencies from PR Description: %v", dependencies)
			log.Printf("Context: %s", context)

			// Prepare the JSON response
			responseData := struct {
				MergeID     string `json:"merge_id"`
				CommitSHA   string `json:"commit_sha"`
				PullRequest int    `json:"pull_request"`
				Context     string `json:"context"`
				Files       []struct {
					Path         string `json:"path"`
					Content      string `json:"content"`
					Dependencies []struct {
						Name    string `json:"name"`
						Content string `json:"content"`
					} `json:"dependencies"`
				} `json:"files"`
			}{
				MergeID:     mergeID,
				CommitSHA:   commitSHA,
				PullRequest: pullRequestNumber,
				Context:     context,
			}

			// Fetch files changed in the PR
			changedFiles, err := FetchPullRequestFiles(repoOwner, repoName, pullRequestNumber)
			if err != nil {
				log.Printf("Unable to fetch changed files: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Unable to fetch changed files",
				})
				return
			}

			// Synchronous processing of files and dependencies
			for _, file := range changedFiles {
				filePath := file["filename"].(string)
				fileContent, err := FetchFileContentFromGitHub(repoOwner, repoName, commitSHA, filePath)
				if err != nil {
					log.Printf("Unable to fetch file content for %s: %v", filePath, err)
					fileContent = "Error fetching content"
				} else {
					log.Printf("Successfully fetched content for file: %s", filePath)
				}

				fileDependencies := FilterDependenciesForFile(filePath, dependencies)
				var formattedDeps []struct {
					Name    string `json:"name"`
					Content string `json:"content"`
				}

				for _, dep := range fileDependencies {
					depContent, err := FetchFileContentFromGitHub(repoOwner, repoName, commitSHA, dep)
					if err != nil {
						log.Printf("Unable to fetch content for dependency %s: %v", dep, err)
						depContent = "Error fetching content"
					} else {
						log.Printf("Successfully fetched content for dependency: %s", dep)
					}
					formattedDeps = append(formattedDeps, struct {
						Name    string `json:"name"`
						Content string `json:"content"`
					}{
						Name:    dep,
						Content: depContent,
					})
				}

				responseData.Files = append(responseData.Files, struct {
					Path         string `json:"path"`
					Content      string `json:"content"`
					Dependencies []struct {
						Name    string `json:"name"`
						Content string `json:"content"`
					} `json:"dependencies"`
				}{
					Path:         filePath,
					Content:      fileContent,
					Dependencies: formattedDeps,
				})
			}

			jsonData, err := json.MarshalIndent(responseData, "", "  ")
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			fmt.Println(string(jsonData))

			// Return the JSON response after all processing is complete
			c.JSON(http.StatusOK, responseData)
		} else {
			c.Status(http.StatusNoContent)
		}
	} else {
		c.Status(http.StatusNoContent)
	}
}
