package controllers

import (
	"encoding/json"
	"fmt"
	"github/controllers/initializers"

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

			repoOwner := prEvent["repository"].(map[string]interface{})["owner"].(map[string]interface{})["login"].(string) //require it later
			repoName := prEvent["repository"].(map[string]interface{})["name"].(string)                                     //require it later
			pullRequestNumber := int(prEvent["number"].(float64))
			commitSHA := prEvent["pull_request"].(map[string]interface{})["merge_commit_sha"].(string) //require it later

			mergeID := fmt.Sprintf("merge_%s_%d", commitSHA, pullRequestNumber)

			// Fetch PR description and dependencies
			prDescription, err := initializers.FetchPullRequestDescription(repoOwner, repoName, pullRequestNumber)
			if err != nil {
				log.Printf("Unable to fetch pull request description: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Unable to fetch pull request description",
				})
				return
			}

			dependencies, context := initializers.ParsePRDescription(prDescription)
			log.Printf("Dependencies from PR Description: %v", dependencies)
			log.Printf("Context: %s", context)

			// Initialize the responseData structure
			responseData := struct {
				MergeID       string `json:"merge_id"`
				Context       string `json:"context"`
				Framework     string `json:"framework"`
				TestDirectory string `json:"test_directory"`
				Comments      string `json:"comments"`
				Files         []struct {
					Path         string `json:"path"`
					Content      string `json:"content"`
					Dependencies []struct {
						Name    string `json:"name"`
						Content string `json:"content"`
					} `json:"dependencies"`
				} `json:"files"`
			}{
				MergeID:       mergeID,
				Context:       context,
				Framework:     "pytest", // Hardcoded framework
				TestDirectory: "tests/", // Hardcoded test directory
				Comments:      "off",    // Hardcoded comments setting
			}

			// Fetch files changed in the PR
			changedFiles, err := initializers.FetchPullRequestFiles(repoOwner, repoName, pullRequestNumber)
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
				fileContent, err := initializers.FetchFileContentFromGitHub(repoOwner, repoName, commitSHA, filePath)
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
					depContent, err := initializers.FetchFileContentFromGitHub(repoOwner, repoName, commitSHA, dep)
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

			// jsonData, err := json.MarshalIndent(responseData, "", "  ")
			// if err != nil {
			// 	fmt.Println("Error:", err)
			// 	return
			// }

			// fmt.Println(string(jsonData))

			log.Printf("Generated payload: %+v", responseData)

			server2URL := "http://localhost:3001/process"
			server2Response, err := SendPayload(server2URL, responseData)
			if err != nil {
				log.Printf("Error sending payload to Server 2: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error forwarding payload to Server 2",
					"error":   err.Error(),
				})
				return
			}

			log.Print("Waiting for the response from Server 2...")
			// Now we wait for responseData
			log.Printf("Response from Server 2: %v", server2Response)
			c.JSON(http.StatusOK, gin.H{
				"message": "Payload processed and forwarded successfully",
				"server2": server2Response,
			})

		} else {
			c.Status(http.StatusNoContent)
		}
	} else {
		c.Status(http.StatusNoContent)
	}
}
