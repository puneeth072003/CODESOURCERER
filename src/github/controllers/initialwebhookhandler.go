package controllers

import (
	"encoding/json"
	"fmt"
	"github/controllers/finalizers"
	"github/controllers/initializers"
	"sync"

	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	MergeID       string `json:"merge_id"`
	Context       string `json:"context"`
	Framework     string `json:"framework"`
	TestDirectory string `json:"test_directory"`
	Comments      string `json:"comments"`
	Files         []File `json:"files"`
}

type File struct {
	Path         string       `json:"path"`
	Content      string       `json:"content"`
	Dependencies []Dependency `json:"dependencies"`
}

type Dependency struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func getFileContents(fileContents []map[string]interface{}, repoOwner, repoName, commitSHA string) <-chan File {
	outChan := make(chan File)

	go func() {
		for _, f := range fileContents {
			filePath := f["filename"].(string)

			fileContent, err := initializers.FetchFileContentFromGitHub(repoOwner, repoName, commitSHA, filePath)
			if err != nil {
				log.Printf("Unable to fetch file content for %s: %v", filePath, err)
				fileContent = "Error fetching content"
			} else {
				log.Printf("Successfully fetched content for file: %s", filePath)
			}

			outChan <- File{
				Path:    filePath,
				Content: fileContent,
			}
		}
		close(outChan)
	}()

	return outChan
}

func getDependencyContents(fileChan <-chan File, dependencies map[string][]string, repoOwner, repoName, commitSHA string) <-chan File {
	outChan := make(chan File)

	go func() {
		for f := range fileChan {
			fileDependencies := FilterDependenciesForFile(f.Path, dependencies)
			var wg sync.WaitGroup
			depChan := make(chan Dependency, len(fileDependencies))

			for _, dep := range fileDependencies {
				wg.Add(1)

				go func(channel chan<- Dependency, dep string) {
					defer wg.Done()

					depContent, err := initializers.FetchFileContentFromGitHub(repoOwner, repoName, commitSHA, dep)
					if err != nil {
						log.Printf("Unable to fetch content for dependency %s: %v", dep, err)
						depContent = "Error fetching content"
					} else {
						log.Printf("Successfully fetched content for dependency: %s", dep)
					}

					channel <- Dependency{
						Name:    dep,
						Content: depContent,
					}
				}(depChan, dep)
			}

			wg.Wait()
			close(depChan)

			var deps []Dependency

			for d := range depChan {
				deps = append(deps, d)
			}

			f.Dependencies = deps
			outChan <- f
		}
		close(outChan)
	}()

	return outChan
}

func WebhookHandler(c *gin.Context) {
	event := c.GetHeader("X-GitHub-Event")

	if event != "pull_request" {
		c.Status(http.StatusNoContent)
		return
	}

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

	if action != "closed" || !merged || baseBranch != "testing" {
		c.Status(http.StatusNoContent)
		return
	}

	log.Printf("Pull request merged into 'testing' branch")

	repoOwner := prEvent["repository"].(map[string]interface{})["owner"].(map[string]interface{})["login"].(string) //require it later
	repoName := prEvent["repository"].(map[string]interface{})["name"].(string)                                     //require it later
	pullRequestNumber := int(prEvent["number"].(float64))
	commitSHA := prEvent["pull_request"].(map[string]interface{})["merge_commit_sha"].(string) //require it later

	mergeID := fmt.Sprintf("merge_%s_%d", commitSHA, pullRequestNumber)

	// yml content fetch
	log.Printf("fetching content from yaml file of repository")
	responseymldata, err := initializers.FetchAndReturnYAMLContents(repoOwner, repoName, commitSHA, "codesourcerer-config.yml")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// log the responseymldata
	log.Printf("YAML Data Retrieved: %+v", responseymldata)

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
	responseData := Response{
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

	// Pipeline Pattern of Gorountine
	fileChan := getFileContents(changedFiles, repoOwner, repoName, commitSHA)
	fileChan = getDependencyContents(fileChan, dependencies, repoOwner, repoName, commitSHA)

	for f := range fileChan {
		responseData.Files = append(responseData.Files, f)
	}

	jsonData, err := json.MarshalIndent(responseData, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	log.Println("##### Constructed payload:", string(jsonData)) // basically string form of unsigned int data

	server2URL := "http://localhost:3001/process"
	server2Response, err := SendPayload(server2URL, string(jsonData))
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

	installationToken := "your_installation_token"    // Replace with actual token
	owner := "your_repo_owner"                        // Replace with actual owner
	repo := "your_repo_name"                          // Replace with actual repo name
	filePath := "path/to/your/file.txt"               // Replace with actual file path
	fileContent := "This is the content of the file." // Replace with actual file content

	// Finalize the request
	err = finalizers.Finalize(installationToken, owner, repo, filePath, fileContent)
	if err != nil {
		log.Printf("Error finalizing the request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error finalizing",
		})
		return
	}

}
