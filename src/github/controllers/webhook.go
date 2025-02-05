package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WebhookController(ctx *gin.Context) {

	event := ctx.GetHeader("X-GitHub-Event")
	var err error
	if event == "pull_request" {
		if err = PullRequestHandler(ctx); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	// TODO: Handle Action Workflow

	ctx.Status(http.StatusNoContent)
	return
}
