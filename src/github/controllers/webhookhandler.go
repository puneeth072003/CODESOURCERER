package controllers

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WebhookHandler(c *gin.Context) {
	event := c.GetHeader("X-GitHub-Event")

	if event == "push" {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Printf("Unable to read request body: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to read request body",
			})
			return
		}
		log.Printf("Received push event: %s", string(body))
		c.JSON(http.StatusOK, gin.H{
			"message": "Received push event",
		})
	} else {
		c.Status(http.StatusNoContent)
	}
}
