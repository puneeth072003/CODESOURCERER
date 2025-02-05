package controllers

import (
	"fmt"
	"log"
	"net/http"

	pb "github.com/codesourcerer-bot/proto/generated"

	"github.com/codesourcerer-bot/github/handlers"
	"github.com/codesourcerer-bot/github/lib/gh"
	"github.com/codesourcerer-bot/github/lib/token"
	"github.com/codesourcerer-bot/github/resolvers"
	"github.com/codesourcerer-bot/github/utils"
	"github.com/gin-gonic/gin"
)

func PullRequestHandler(c *gin.Context) error {

	prBody, err := resolvers.NewPrBody(c.Request.Body)
	if err != nil {
		return err
	}

	repoName, repoOwner := prBody.GetRepoInfo()

	action, merged, baseBranch := prBody.GetPRStatus()

	// TODO: Replace the condition later (replace hardcoded baseBranch Value)
	if action != "closed" || !merged || baseBranch != "testing" {
		c.Status(http.StatusNoContent)
		return nil
	}

	pullRequestNumber, commitSHA := prBody.GetPRInfo()

	ymlConfig := gh.FetchYmlConfig(repoOwner, repoName, commitSHA)

	prDescription, err := gh.FetchPullRequestDescription(repoOwner, repoName, pullRequestNumber)
	if err != nil {
		log.Printf("Unable to fetch pull request description: %v", err)
		return fmt.Errorf("unable to fetch pull request description")
	}

	dependencies, context := utils.ParsePRDescription(prDescription)

	changedFiles, err := gh.FetchPullRequestFiles(repoOwner, repoName, pullRequestNumber)
	if err != nil {
		log.Printf("Unable to fetch changed files: %v", err)
		return fmt.Errorf("unable to fetch changed files")
	}

	fileChan := resolvers.GetFileContents(changedFiles, repoOwner, repoName, commitSHA)
	fileChan = resolvers.GetDependencyContents(fileChan, dependencies, repoOwner, repoName, commitSHA)

	mergeID := fmt.Sprintf("merge_%s_%d", commitSHA, pullRequestNumber)
	genConfig := gh.GetGenerationOptions(ymlConfig)
	payload := pb.GithubContextRequest{
		MergeId: mergeID,
		Context: context,
		Config:  genConfig,
	}

	for f := range fileChan {
		payload.Files = append(payload.Files, f)
	}

	generatedTests, err := handlers.GetGeneratedTestsFromGenAI(&payload)
	if err != nil {
		log.Printf("Error sending payload to GenAI Service: %v", err)
		return fmt.Errorf("error forwarding payload to GenAI Service")
	}

	token, err := token.GetInstance().GetToken()
	if err != nil {
		log.Printf("Error getting token: %v", err)
		return fmt.Errorf("error getting token")
	}

	err = resolvers.PushNewBranchWithTests(token, repoOwner, repoName, generatedTests)
	if err != nil {
		log.Printf("Error finalizing: %v", err)
		return fmt.Errorf("error finalizing")
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Pull request has been raised"})

	return nil

}
