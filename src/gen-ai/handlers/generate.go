package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	pb "github.com/codesourcerer-bot/proto/generated"

	"github.com/codesourcerer-bot/gen-ai/models"

	"github.com/google/generative-ai-go/genai"
)

// TODO: Handle the Configuration Neatly
func GetTestsFromAI(ctx context.Context, payload *pb.GithubContextRequest, model *genai.GenerativeModel) (*pb.GeneratedTestsResponse, error) {
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

	var result pb.GeneratedTestsResponse

	if err := json.Unmarshal([]byte(textPart), &result); err != nil {
		return nil, fmt.Errorf("Unable to unmarshal: %v", err)
	}

	return &result, nil

}
