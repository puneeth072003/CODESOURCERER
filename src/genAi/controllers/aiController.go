package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"genAi/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
)

type Dependency struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type File struct {
	Path         string       `json:"path"`
	Content      string       `json:"content"`
	Dependencies []Dependency `json:"dependencies"`
}
type RequestPayload struct {
	MergeID       string `json:"merge_id"`
	Context       string `json:"context"`
	Framework     string `json:"framework"`
	TestDirectory string `json:"test_directory"`
	Comments      string `json:"comments"`
	Files         []File `json:"files"`
}

type Test struct {
	TestName     string `json:"testname"`
	TestFilePath string `json:"testfilepath"`
	ParentPath   string `json:"parentpath"`
	Code         string `json:"code"`
}

type GenAIResponse struct {
	Tests []Test `json:"tests"`
}

func ProcessAIRequest(ctx context.Context, c *gin.Context, client *genai.Client, model *genai.GenerativeModel) {
	var payload RequestPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format", "details": err.Error()})
		return
	}

	aiOutput, err := processAI(ctx, payload, model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process AI request", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &aiOutput)
}

func processAI(ctx context.Context, payload RequestPayload, model *genai.GenerativeModel) (*GenAIResponse, error) {
	session := model.StartChat()
	session.History = models.SessionHistory

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing payload: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	response, err := session.SendMessage(ctx, genai.Text(string(payloadBytes)))
	if err != nil {
		return nil, fmt.Errorf("error generating response: %v", err)
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("model did not generate any response")
	}

	part := response.Candidates[0].Content.Parts[0]

	textPart, ok := part.(genai.Text)
	if !ok {
		return nil, fmt.Errorf("unable to convert to text type: %v", err)
	}

	var result GenAIResponse

	if err := json.Unmarshal([]byte(textPart), &result); err != nil {
		return nil, fmt.Errorf("Unable to unmarshal: %v", err)
	}

	return &result, nil

}
