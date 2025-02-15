package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/codesourcerer-bot/github/resolvers"
	"github.com/gin-gonic/gin"
)

func WorkflowHandler(ctx *gin.Context) error {

	workflowBody, err := resolvers.NewWorkflowBody(ctx.Request.Body)
	if err != nil {
		return err
	}

	name, status, result := workflowBody.GetWorkflowDetails()

	if name != "Run Tests in Directory" || status != "completed" || result != "failure" {
		ctx.Status(http.StatusNoContent)
		return nil
	}

	owner, repoName := workflowBody.GetRepoDetails()
	jobUrl := workflowBody.GetWorkflowJobUrl()

	logs, err := resolvers.FetchLogs(jobUrl, owner, repoName)
	if err != nil {
		return err
	}

	parsedLogs := ""
	for _, log := range strings.Split(logs, "\n") {
		if len(log) > 29 {
			parsedLogs = parsedLogs + log[29:] + "\n"
		}
	}

	fmt.Printf("%s", parsedLogs)

	return nil
}
