package main

import (
	"context"
	"log"
	"os"

	"genAi/controllers"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	model.SetTemperature(0.9)
	model.SetTopP(0.5)
	model.SetTopK(20)
	model.SetMaxOutputTokens(100)
	model.SystemInstruction = genai.NewUserContent(genai.Text("You are a generative AI model trained to produce test suites for code based on an input payload. Your task is to interpret the input payload and generate test cases for each file under the files array, ensuring you adhere to the provided format and conventions. The payload will also include an additional framework field that specifies the testing framework to be used."))

	// Initialize the Gin router
	router := gin.Default()

	router.GET("/ping", controllers.Pong)
	router.POST("/process", func(c *gin.Context) {
		controllers.ProcessAIRequest(ctx, c, client, model)
	})

	// Start the server on port 3001
	log.Println("Server running on port 3001")
	router.Run(":3001")
}

func ExampleGenerativeModel_GenerateContent_config() {

}
