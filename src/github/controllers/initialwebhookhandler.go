package controllers

import (
	"encoding/json"
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

		// Check if the pull request is closed and merged
		action := prEvent["action"].(string)
		merged := prEvent["pull_request"].(map[string]interface{})["merged"].(bool)                                     // boolean representing the status of PR
		baseBranch := prEvent["pull_request"].(map[string]interface{})["base"].(map[string]interface{})["ref"].(string) // name of the base branch of the PR

		if action == "closed" && merged && baseBranch == "testing" {
			log.Printf("Pull request merged into 'testing' branch")
			c.JSON(http.StatusOK, gin.H{
				"message": "Pull request merged into 'testing' branch",
			})
		} else {
			c.Status(http.StatusNoContent)
		}
	} else {
		c.Status(http.StatusNoContent)
	}
}
