package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WebhookController(ctx *gin.Context) {

	event := ctx.GetHeader("X-GitHub-Event")
	var err error
	if event == "pull_request" {
		err = PullRequestHandler(ctx)
	}
	// TODO: Handle Action Workflow

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
