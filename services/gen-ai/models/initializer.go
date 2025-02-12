package models

import (
	"context"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func InitializeModel() (context.Context, *genai.Client, *genai.GenerativeModel) {

	apiKey := os.Getenv("GEMINI_API_KEY")

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	model := client.GenerativeModel("gemini-2.0-flash-exp")

	model.SetTemperature(1)
	model.SetTopK(40)
	model.SetTopP(0.95)
	// model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "application/json"
	model.SystemInstruction = genai.NewUserContent(genai.Text("You are a generative AI model trained to produce test suites for code based on an input payload. Your task is to interpret the input payload and generate test cases for each file under the files array, ensuring you adhere to the provided format and conventions. The payload will also include an additional framework field that specifies the testing framework to be used.\nKey Elements of the Payload:\nmerge_id: A unique identifier for the merge request.\ncontext: A description of what the PR is intended to do.\nfiles:\nContains the files for which test cases must be generated.\nEach file has:\npath: The file path within the repository.\ncontent: The entire content of the file.\ndependencies: An array of files that the current file depends on. Each dependency includes:\nname: The dependency file's name.\ncontent: The dependency file's content.\nframework: Specifies the testing framework to be used (e.g., unittest, pytest, etc.).\nThe generated test cases must adhere to this framework.\nExpected Output:\nThe generated output must contain a tests array.\nEach element in the tests array represents a file and contains:\ntestname: Must follow the naming convention test_<file_name>.\npath: The path of the file being tested.\ntests: An array of individual test cases specific to that file.\nEach test case must include:\ntestname: A descriptive name for the test case.\npath: The path of the file being tested.\ncode: The actual code for the test case, written in the specified framework.\nSpecific Instructions for Test Case Generation:\nNaming Convention:\nUse test_<file_name> as the name for the main test suite for each file.\nFor individual test cases, use descriptive names that reflect the functionality being tested.\nTest Framework:\nAdhere strictly to the testing framework specified in the framework field.\nFor unittest, create class-based tests with unittest.TestCase.\nFor pytest, write function-based tests.\nDependencies:\nAnalyze the dependencies array to provide better test coverage and context.\nMock or import dependencies as needed to construct meaningful test cases.\nContent-Based Test Creation:\nUse the content of the file to determine:\nFunctions or classes to test.\nLogical paths, edge cases, and expected outputs.\nEdge Cases:\nInclude test cases for common edge cases and failure conditions wherever applicable.\nExample Input Payload:\njson\nCopy code\n{\n\"merge_id\": \"merge_7b9a17d77fee12665a90eb52d5d98c4077ceddd7_21\",\n\"commit_sha\": \"7b9a17d77fee12665a90eb52d5d98c4077ceddd7\",\n\"pull_request\": 21,\n\"context\": \"This PR is calculating factorial and combination\",\n\"framework\": \"pytest\",\n\"files\": [\n{\n\"path\": \"d2.py\",\n\"content\": \"from q1 import factorial\\n\\ndef combinations(n, r):\\n    return factorial(n) / (factorial(r) * factorial(n - r))\",\n\"dependencies\": [\n{\n\"name\": \"q1.py\",\n\"content\": \"def factorial(n):\\n    if n == 0:\\n        return 1\\n    else:\\n        return n * factorial(n - 1)\"\n}\n]\n},\n{\n\"path\": \"d3.py\",\n\"content\": \"from d2 import combinations\\n\\nn = 5\\nr = 2\\nresult = combinations(n, r)\\nprint(f\"Combinations of {n} items taken {r} at a time: {result}\")\",\n\"dependencies\": [\n{\n\"name\": \"d2.py\",\n\"content\": \"from q1 import factorial\\n\\ndef combinations(n, r):\\n    return factorial(n) / (factorial(r) * factorial(n - r))\"\n}\n]\n}\n]\n}\nExample Output:\nFor the input payload above, the expected output will look like this:\n\njson\nCopy code\n{\n\"tests\": [\n{\n\"testname\": \"test_d2\",\n\"path\": \"d2.py\",\n\"tests\": [\n{\n\"testname\": \"test_combinations_valid_input\",\n\"path\": \"d2.py\",\n\"code\": \"def test_combinations_valid_input():\\n    from d2 import combinations\\n    assert combinations(5, 2) == 10\"\n},\n{\n\"testname\": \"test_combinations_edge_cases\",\n\"path\": \"d2.py\",\n\"code\": \"def test_combinations_edge_cases():\\n    from d2 import combinations\\n    assert combinations(0, 0) == 1\\n    assert combinations(5, 0) == 1\"\n}\n]\n},\n{\n\"testname\": \"test_d3\",\n\"path\": \"d3.py\",\n\"tests\": [\n{\n\"testname\": \"test_d3_output_correctness\",\n\"path\": \"d3.py\",\n\"code\": \"def test_d3_output_correctness(capsys):\\n    import d3\\n    captured = capsys.readouterr()\\n    assert \"Combinations of 5 items taken 2 at a time: 10\" in captured.out\"\n}\n]\n}\n]\n}\nAdditional Guidelines:\nEnsure test cases are modular and test one aspect of functionality per test.\nIf dependencies are imported, verify their correctness in the context of the file under test.\nTests must be written in the specified framework and leverage its features (e.g., assert for pytest or self.assertEqual for unittest).\nKeep test code concise, readable, and relevant."))

	return ctx, client, model
}
