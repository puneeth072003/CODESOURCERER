package models

import (
	"context"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func InitializeRetryModel() (context.Context, *genai.Client, *genai.GenerativeModel) {

	apiKey := os.Getenv("GEMINI_API_KEY")

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	model := client.GenerativeModel("gemini-2.0-flash-exp")

	// TODO: Model Config Here

	return ctx, client, model
}
