package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codesourcerer-bot/github/connections"
	"github.com/codesourcerer-bot/github/lib"
	"github.com/codesourcerer-bot/github/resolvers"
	"github.com/codesourcerer-bot/github/validators"
	pb "github.com/codesourcerer-bot/proto/generated"
	"github.com/gin-gonic/gin"
)

func WorkflowHandler(ctx *gin.Context) error {

	workflowBody, err := validators.NewWorkflowBody(ctx.Request.Body)
	if err != nil {
		return err
	}

	name, status, result := workflowBody.GetWorkflowDetails()

	if name != "Run Tests in Directory" || status != "completed" || result == "" {
		ctx.Status(http.StatusNoContent)
		return nil
	}

	owner, repoName, branchName := workflowBody.GetRepoDetails()
	jobUrl := workflowBody.GetWorkflowJobUrl()

	cacheKey := fmt.Sprintf("%s/%s/tree/%s", owner, repoName, branchName)

	if result == "success" {
		if ok, err := connections.DeleteContextAndTestsToDatabase(cacheKey); err != nil || !ok {
			log.Printf("unable to delete cache: %v", err)
			ctx.JSON(http.StatusAccepted, gin.H{"error": "unable to delete cach"})
			return nil
		}
		ctx.JSON(http.StatusAccepted, gin.H{"message": "cache has been cleared due to workflow success"})
		return nil
	}

	if result != "failure" {
		ctx.Status(http.StatusNoContent)
		return nil
	}

	if isRetryExhausted, err := connections.GetRetryExhaustionStatus(cacheKey); err != nil {
		return fmt.Errorf("unable to fetch retry count: %v", isRetryExhausted)
	} else if isRetryExhausted {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "retry has been exhausted"})
		return nil
	}

	logs, err := lib.FetchLogs(jobUrl, owner, repoName)
	if err != nil {
		return err
	}

	cache, err := connections.GetContextAndTestsFromDatabase(cacheKey)
	if err != nil {
		return err
	}

	payload := &pb.RetryMechanismPayload{
		Cache: cache,
		Logs:  logs,
	}

	generatedTests, err := connections.GetRetriedTestsFromGenAI(payload)
	if err != nil {
		log.Printf("Error from GenAI Service: %v", err)
		return fmt.Errorf("error forwarding payload to GenAI Service")
	}

	if ok, err := connections.SetContextAndTestsToDatabase(cacheKey, cache.GetContexts(), generatedTests.GetTests()); err != nil || !ok {
		log.Printf("unable to update cache: %v", err)
		return fmt.Errorf("unable to update cache")
	}

	if err != nil {
		log.Printf("Error getting token: %v", err)
		return fmt.Errorf("error getting token")
	}

	if err = resolvers.CommitRetriedTests(owner, repoName, branchName, generatedTests); err != nil {
		log.Printf("unable to commit test files: %v", err)
		return fmt.Errorf("unable to commit test files")
	}

	return nil
}
