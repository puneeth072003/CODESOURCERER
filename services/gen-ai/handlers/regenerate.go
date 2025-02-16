package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/codesourcerer-bot/gen-ai/contexts"
	pb "github.com/codesourcerer-bot/proto/generated"
	"github.com/google/generative-ai-go/genai"
)

func generateRetriedTestsFromAI(ctx context.Context, parsedLogs genai.Part, cache *pb.CachedContents, model *genai.GenerativeModel) (*pb.GeneratedTestsResponse, error) {
	session := model.StartChat()
	session.History = contexts.GetRegeratorContext(cache)

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	response, err := session.SendMessage(ctx, parsedLogs)
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
