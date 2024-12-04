package controllers

import (
	"encoding/json"
	"fmt"
	"github/controllers/generators"
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

			log.Printf("Response from Server 2: %s", server2Response)
			c.JSON(http.StatusOK, gin.H{
				"message": "Payload processed and forwarded successfully",
				"server2": server2Response,
			})

			// Generate test files
			tests := []struct {
				TestName    string
				TestPath    string
				ParentPath  string
				CodeContent string
			}{
				{
					TestName:   "test_t2",
					TestPath:   "tests/test_t2.py",
					ParentPath: "t2.py",
					CodeContent: `import pytest
				from t2 import combinations
				from q1 import factorial # Import dependency for testing
				
				def test_combinations_valid_input():
					assert combinations(5, 2) == 10.0
					assert combinations(10,3) == 120.0
					#Coughed up by CODESOURCERER
				
				def test_combinations_edge_cases():
					assert combinations(0, 0) == 1.0
					assert combinations(5, 0) == 1.0
					assert combinations(5, 5) == 1.0
					#Coughed up by CODESOURCERER
				
				def test_combinations_invalid_input():
					with pytest.raises(ValueError):
						combinations(5, 6) # r > n
					with pytest.raises(ValueError):
						combinations(-1, 2) # negative n
					with pytest.raises(ValueError):
						combinations(5, -2) # negative r
					#Coughed up by CODESOURCERER`,
				},
				{
					TestName:   "test_t3",
					TestPath:   "tests/test_t3.py",
					ParentPath: "t3.py",
					CodeContent: `import pytest
				from t3 import combinations #This import might fail if t2.py or q1.py have errors
				from io import StringIO
				import sys
				
				def test_t3_output_correctness():
					# Capture stdout
					old_stdout = sys.stdout
					redirected_output = StringIO()
					sys.stdout = redirected_output
					
					#Call the function being tested
					combinations(5,2)
				
					#Restore stdout and check output
					sys.stdout = old_stdout
					assert "Combinations of 5 items taken 2 at a time: 10.0" in redirected_output.getvalue()
					#Coughed up by CODESOURCERER`,
				},
			}

			err = generators.GenerateTestFiles(tests)
			if err != nil {
				log.Printf("Error generating test files: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error generating test files",
				})
				return
			}

			err = generators.CommitNewFilesToBranch("testing")
			if err != nil {
				log.Printf("Error committing new files to branch: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error committing new files to branch",
				})
				return
			}
			prTitle := "Add tests for file operations and main"
			prBody := "This PR adds tests for file operations and main module."
			newBranch := "testing"
			err = generators.CreateDraftPullRequest(repoOwner, repoName, newBranch, commitSHA, prTitle, prBody)
			if err != nil {
				log.Printf("Error creating draft PR: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error creating draft PR",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Test files generated and draft PR created successfully",
			})

		} else {
			c.Status(http.StatusNoContent)
		}
	} else {
		c.Status(http.StatusNoContent)
	}
}
