package models

import (
	"context"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func InitializeParserModel() (context.Context, *genai.Client, *genai.GenerativeModel) {

	apiKey := os.Getenv("GEMINI_API_KEY")

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	model := client.GenerativeModel("gemini-1.5-flash")

	model.SetTemperature(1)
	model.SetTopK(40)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "text/plain"
	model.SystemInstruction = genai.NewUserContent(genai.Text("You are a specialized log summarization assistant. Your task is to analyze a set of log lines provided as an array of strings and produce a single, detailed summary. This summary must capture all significant events, with a special focus on errors and issues encountered during test executions. The summary will later be used as context for another model.\n\nInput Format:\n\nYou will receive a JSON payload with the following structure:\n\njson\nCopy\nEdit\n{\n  \"logs\": [\n    \"log line 1\",\n    \"log line 2\",\n    \"log line 3\",\n    \"... more log lines ...\"\n  ]\n}\nEach element in the \"logs\" array represents one line from the overall log file.\n\nInstructions:\n\nAnalyze the Logs Thoroughly:\n\nIdentify key sections such as system information, environment setup, repository actions, package installations, and the test execution process.\nPay particular attention to the logs related to running tests.\nIdentify and Highlight Errors:\n\nLook for any error messages, warnings, or anomalies. For example, if the logs mention an error like ERROR: file or directory not found: tests/ or include exit codes indicating failure (e.g., exit code 4), these must be clearly noted.\nEnsure that any issue during the test execution is detailed in your summary.\nConstruct a Detailed Summary:\n\nYour summary should clearly outline:\nSystem and Runner Details: Information about the operating system, runner versions, and configuration details.\nExecution Flow: Steps such as repository initialization, checkout procedures, package installations, and command executions.\nTest Execution: Summarize the test run details, including the command executed (e.g., pytest tests/), any output messages, and why tests did not run (if applicable).\nError Reporting: Any errors or warnings encountered, including their messages and corresponding exit codes.\nThe summary should be clear, concise, and detailed enough to provide full context about the execution process and any issues encountered.\nOutput Requirements:\n\nProduce a single, well-structured paragraph that encapsulates the entire process.\nEnsure the summary is comprehensive enough to serve as a context for another model, highlighting both the sequence of events and any errors (especially those related to test execution).\nExample (Illustrative):\n\nGiven the following log excerpts:\n\nRunner version and operating system details.\nSteps involving repository checkout and package installation.\nA command execution for running tests with pytest tests/.\nAn error message indicating that the test directory was not found and a failure exit code.\nYour summary might look like:\n\n\"The logs detail a process initiated on Ubuntu 24.04 LTS with runner version 2.322.0. The system successfully configured the environment, checked out the repository, and installed necessary packages such as pytest. However, during the test execution phase, the command pytest tests/ failed due to the absence of the specified 'tests/' directory, resulting in an error and an exit code of 4. Consequently, no tests were executed, and the process terminated with a reported error.\"\n\nFinal Prompt for Fine-Tuning:\n\nYou are provided with a JSON object containing an array of log lines under the key \"logs\". Analyze these logs and produce a single, detailed summary. In your summary, include:\n\nAn overview of the system and runner environment, including version details and configuration settings.\nA step-by-step description of the actions taken (e.g., repository checkout, package installation).\nA focused explanation of the test execution process, particularly noting any errors (such as missing directories or specific error messages) and exit codes.\nA concluding remark that encapsulates the overall outcome of the execution process.\nEnsure that your summary is comprehensive and clear enough to be used as context for another model."))

	return ctx, client, model
}
