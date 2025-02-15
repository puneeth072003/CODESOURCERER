package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WebhookController(ctx *gin.Context) {

	event := ctx.GetHeader("X-GitHub-Event")

	var err error

	switch event {
	case "pull_request":
		if err = PullRequestHandler(ctx); err == nil {
			return
		}

	case "workflow_run":
		if err = WorkflowHandler(ctx); err == nil {
			return
		}
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	ctx.Status(http.StatusNoContent)
}
