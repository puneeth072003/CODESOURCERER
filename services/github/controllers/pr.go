package controllers

import (
	"fmt"
	"log"
	"net/http"

	pb "github.com/codesourcerer-bot/proto/generated"

	"github.com/codesourcerer-bot/github/connections"
	"github.com/codesourcerer-bot/github/lib"
	"github.com/codesourcerer-bot/github/resolvers"
	"github.com/codesourcerer-bot/github/utils"
	"github.com/codesourcerer-bot/github/validators"
	"github.com/gin-gonic/gin"
)

func PullRequestHandler(c *gin.Context) error {

	prBody, err := validators.NewPrBody(c.Request.Body)
	if err != nil {
		return err
	}

	repoName, repoOwner := prBody.GetRepoInfo()

	action, merged, baseBranch := prBody.GetPRStatus()

	if action != "closed" || !merged {
		c.Status(http.StatusNoContent)
		return nil
	}

	pullRequestNumber, commitSHA := prBody.GetPRInfo()

	ymlConfig := lib.FetchYmlConfig(repoOwner, repoName, commitSHA)

	if baseBranch != ymlConfig.Configuration.TestingBranch {
		c.Status(http.StatusNoContent)
		return nil
	}

	prDescription, err := lib.FetchPullRequestDescription(repoOwner, repoName, pullRequestNumber)
	if err != nil {
		log.Printf("Unable to fetch pull request description: %v", err)
		return fmt.Errorf("unable to fetch pull request description")
	}

	dependencies, context := utils.ParsePRDescription(prDescription)

	changedFiles, err := lib.FetchPullRequestFiles(repoOwner, repoName, pullRequestNumber)
	if err != nil {
		log.Printf("Unable to fetch changed files: %v", err)
		return fmt.Errorf("unable to fetch changed files")
	}

	fileChan := resolvers.GetFileContents(changedFiles, repoOwner, repoName, commitSHA)
	fileChan = resolvers.GetDependencyContents(fileChan, dependencies, repoOwner, repoName, commitSHA)

	mergeID := fmt.Sprintf("merge_%s_%d", commitSHA, pullRequestNumber)
	genConfig := lib.GetGenerationOptions(ymlConfig)
	payload := pb.GithubContextRequest{
		MergeId: mergeID,
		Context: context,
		Config:  genConfig,
	}

	for f := range fileChan {
		payload.Files = append(payload.Files, f)
	}

	generatedTests, err := connections.GetGeneratedTestsFromGenAI(&payload)
	if err != nil {
		log.Printf("Error sending payload to GenAI Service: %v", err)
		return fmt.Errorf("error forwarding payload to GenAI Service")
	}

	newBranch := utils.GetRandomBranch()

	cacheResult := resolvers.CachePullRequest(ymlConfig.Caching.Enabled, repoOwner, repoName, newBranch, payload.GetFiles(), generatedTests.GetTests())

	err = resolvers.PushNewBranchWithTests(repoOwner, repoName, ymlConfig.Configuration.TestingBranch, newBranch, cacheResult, generatedTests)
	if err != nil {
		log.Printf("Error finalizing: %v", err)
		return fmt.Errorf("error finalizing")
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Pull request has been raised"})

	return nil

}
