package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/codesourcerer-bot/gen-ai/contexts"
	"github.com/google/generative-ai-go/genai"
)

func getParsedLogsFromAI(c context.Context, payload []string, model *genai.GenerativeModel) (genai.Part, error) {

	session := model.StartChat()
	session.History = contexts.ParserModelContext

	c, cancel := context.WithTimeout(c, 15*time.Second)
	defer cancel()

	logs := make([]genai.Part, len(payload))
	for _, p := range payload {
		logs = append(logs, genai.Text(string(p)))
	}

	response, err := session.SendMessage(c, logs...)
	if err != nil {
		return nil, fmt.Errorf("error generating response: %v", err)
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("model did not generate any response")
	}

	part := response.Candidates[0].Content.Parts[0]

	return part, nil

}
