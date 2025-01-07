package controllers

import (
	"encoding/json"
	"fmt"
	"github/controllers/initializers"
	"github/handlers"
	"sync"

	"io"
	"log"
	"net/http"

	pb "protobuf/generated"

	"github.com/gin-gonic/gin"
)

func getFileContents(fileContents []map[string]interface{}, repoOwner, repoName, commitSHA string) <-chan *pb.SourceFilePayload {
	outChan := make(chan *pb.SourceFilePayload)

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

			outChan <- &pb.SourceFilePayload{
				Path:    filePath,
				Content: fileContent,
			}
		}
		close(outChan)
	}()

	return outChan
}

func getDependencyContents(fileChan <-chan *pb.SourceFilePayload, dependencies map[string][]string, repoOwner, repoName, commitSHA string) <-chan *pb.SourceFilePayload {
	outChan := make(chan *pb.SourceFilePayload)

	go func() {
		for f := range fileChan {
			fileDependencies := FilterDependenciesForFile(f.Path, dependencies)
			var wg sync.WaitGroup
			depChan := make(chan *pb.SourceFileDependencyPayload, len(fileDependencies))

			for _, dep := range fileDependencies {
				wg.Add(1)

				go func(channel chan<- *pb.SourceFileDependencyPayload, dep string) {
					defer wg.Done()

					depContent, err := initializers.FetchFileContentFromGitHub(repoOwner, repoName, commitSHA, dep)
					if err != nil {
						log.Printf("Unable to fetch content for dependency %s: %v", dep, err)
						depContent = "Error fetching content"
					} else {
						log.Printf("Successfully fetched content for dependency: %s", dep)
					}

					channel <- &pb.SourceFileDependencyPayload{
						Name:    dep,
						Content: depContent,
					}
				}(depChan, dep)
			}

			wg.Wait()
			close(depChan)

			var deps []*pb.SourceFileDependencyPayload

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
	responseymldata := initializers.FetchConfig(repoOwner, repoName, commitSHA)

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
	payload := pb.GithubContextRequest{
		MergeId:       mergeID,
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
		payload.Files = append(payload.Files, f)
	}

	log.Println("##### Constructed payload:", payload.String()) // basically string form of unsigned int data

	generatedTests, err := handlers.GetGeneratedTestsFromGenAI(&payload)
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
	log.Printf("Response from Server 2: %v", generatedTests.String())
	c.JSON(http.StatusOK, gin.H{
		"message": "Payload processed and forwarded successfully",
		"server2": &generatedTests,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Test files generated and draft PR created successfully",
	})

}
