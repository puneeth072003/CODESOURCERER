package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

type ResponsePayload struct {
	Output string `json:"output"`
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

	// Send the raw JSON output
	c.Data(http.StatusOK, "application/json", aiOutput)
}

func processAI(ctx context.Context, payload RequestPayload, model *genai.GenerativeModel) ([]byte, error) {
	session := model.StartChat()
	session.History = []*genai.Content{
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("Key Elements of the Payload:\nInput Fields:\nmerge_id: A unique identifier for the merge request.\ncontext: Describes the pull request's purpose and what it introduces or changes.\nframework: Specifies the testing framework to use (e.g., unittest, pytest).\ntest_directory (optional): Specifies the directory where test files should be placed. Defaults to tests/.\ncomments: Determines whether comments are included in the generated test code. Possible values:\n\"on\": Include descriptive comments in the test code.\n\"off\": Exclude comments entirely.\nfiles:\npath: Path of the file in the repository.\ncontent: Complete content of the file.\ndependencies: An array of files that the current file depends on, containing:\nname: Dependency file name.\ncontent: Dependency file's content.\nOutput Format:\nExpected Structure\njson\nCopy code\n{\n  \"tests\": [\n    {\n      \"testname\": \"<main test suite name>\",\n      \"testfilepath\": \"<generated test file path>\",\n      \"parentpath\": \"<original file path>\",\n      \"code\": \"<entire test code>\"\n    }\n  ]\n}\ntestname: Follows the naming convention test_<file_name>.\ntestfilepath: Full path to the generated test file. Default is tests/ directory, but it should respect the provided test_directory field if specified.\nparentpath: Original file path in the repository.\ncode: The complete test code written in the specified framework.\nTest Case Generation Instructions:\nFramework:\n\nUse the framework specified in the framework field (unittest or pytest).\nEnsure compatibility with Python 3.8+ unless explicitly stated otherwise.\nFile Locations and Imports:\n\ntestfilepath: Place generated tests in the directory specified by test_directory. If not provided, default to placing tests in tests/ relative to the original file.\nImports:\nFor files in the root directory, use direct imports like from <filename> import <functions/classes>.\nFor files in subdirectories, use absolute imports based on the repository structure.\nTest Coverage:\n\nGenerate tests for all functions/classes in the original file, covering:\nTypical inputs.\nEdge cases.\nException handling (where applicable).\nEnsure meaningful assertions and robust coverage.\nMock dependencies as needed to simulate their behavior.\nCode Style:\n\nFormat all test code according to PEP-8 standards.\nKeep code modular and concise.\nComments:\n\nControlled by the comments field:\n\"on\": Add descriptive comments explaining the purpose of each test and key code sections.\n\"off\": Exclude comments entirely.\nAlways include this comment at the end of the test code:\n# Coughed up by CODESOURCERER.\nNaming Conventions:\n\nMain test suite: test_<file_name> (e.g., test_date_utils for date_utils.py).\nIndividual test cases: Use descriptive names indicating functionality (e.g., test_format_date_valid_input).\nDefault Behavior:\n\nIf test_directory is missing, default to placing test files under tests/<module_name>/.\nEnsure __init__.py files are present in all relevant directories for Python package compatibility.\nExample Input:\njson\nCopy code\n{\n  \"merge_id\": \"merge_uvw456rst789xyz123abc890klm567def234_107\",\n  \"context\": \"This PR adds utility functions for date formatting and integrates these into a scheduling module.\",\n  \"framework\": \"pytest\",\n  \"test_directory\": \"tests/\",\n  \"comments\": \"off\",\n  \"files\": [\n    {\n      \"path\": \"date_utils.py\",\n      \"content\": \"from datetime import datetime\\n\\ndef format_date(date):\\n    return date.strftime('%Y-%m-%d')\\n\\ndef parse_date(date_string):\\n    return datetime.strptime(date_string, '%Y-%m-%d')\",\n      \"dependencies\": []\n    },\n    {\n      \"path\": \"scheduling/schedule_manager.py\",\n      \"content\": \"from date_utils import format_date, parse_date\\n\\ndef get_formatted_date_for_today():\\n    return format_date(datetime.now())\",\n      \"dependencies\": [\n        {\n          \"name\": \"date_utils.py\",\n          \"content\": \"from datetime import datetime\\n\\ndef format_date(date):\\n    return date.strftime('%Y-%m-%d')\\n\\ndef parse_date(date_string):\\n    return datetime.strptime(date_string, '%Y-%m-%d')\"\n        }\n      ]\n    }\n  ]\n}\nExample Output:\njson\nCopy code\n{\n  \"tests\": [\n    {\n      \"testname\": \"test_date_utils\",\n      \"testfilepath\": \"tests/test_date_utils.py\",\n      \"parentpath\": \"date_utils.py\",\n      \"code\": \"import pytest\\nfrom date_utils import format_date, parse_date\\n\\ndef test_format_date():\\n    date = datetime(2024, 12, 12)\\n    assert format_date(date) == '2024-12-12'\\n\\ndef test_parse_date():\\n    date_string = '2024-12-12'\\n    assert parse_date(date_string) == datetime(2024, 12, 12)\\n\\n# Coughed up by CODESOURCERER\"\n    },\n    {\n      \"testname\": \"test_schedule_manager\",\n      \"testfilepath\": \"tests/scheduling/test_schedule_manager.py\",\n      \"parentpath\": \"scheduling/schedule_manager.py\",\n      \"code\": \"import pytest\\nfrom scheduling.schedule_manager import get_formatted_date_for_today\\n\\ndef test_get_formatted_date_for_today(mocker):\\n    mock_date = mocker.patch('scheduling.schedule_manager.datetime')\\n    mock_date.now.return_value = datetime(2024, 12, 12)\\n    assert get_formatted_date_for_today() == '2024-12-12'\\n\\n# Coughed up by CODESOURCERER\"\n    }\n  ]\n}\n"),
			},
		},
		{
			Role: "model",
			Parts: []genai.Part{
				genai.Text("```json\n{\n  \"tests\": [\n    {\n      \"testname\": \"test_date_utils\",\n      \"testfilepath\": \"tests/test_date_utils.py\",\n      \"parentpath\": \"date_utils.py\",\n      \"code\": \"import pytest\\nfrom datetime import datetime\\nfrom date_utils import format_date, parse_date\\n\\ndef test_format_date_valid_input():\\n    date = datetime(2024, 1, 15)\\n    assert format_date(date) == '2024-01-15'\\n\\ndef test_format_date_edge_case():\\n    date = datetime(1999, 12, 31)\\n    assert format_date(date) == '1999-12-31'\\n\\ndef test_parse_date_valid_input():\\n    date_string = '2024-01-15'\\n    assert parse_date(date_string) == datetime(2024, 1, 15)\\n\\ndef test_parse_date_edge_case():\\n    date_string = '1999-12-31'\\n    assert parse_date(date_string) == datetime(1999, 12, 31)\\n\\n\\n\\ndef test_parse_date_invalid_format():\\n    with pytest.raises(ValueError):\\n        parse_date('2024/01/15')\\n\\n# Coughed up by CODESOURCERER\"\n    },\n    {\n      \"testname\": \"test_schedule_manager\",\n      \"testfilepath\": \"tests/scheduling/test_schedule_manager.py\",\n      \"parentpath\": \"scheduling/schedule_manager.py\",\n      \"code\": \"import pytest\\nfrom datetime import datetime\\nfrom scheduling.schedule_manager import get_formatted_date_for_today\\n\\ndef test_get_formatted_date_for_today(mocker):\\n    mocked_datetime = mocker.patch('scheduling.schedule_manager.datetime')\\n    mocked_datetime.now.return_value = datetime(2024, 1, 15)\\n    assert get_formatted_date_for_today() == '2024-01-15'\\n\\n    mocked_datetime.now.return_value = datetime(2024, 12, 31)\\n    assert get_formatted_date_for_today() == '2024-12-31'\\n\\n# Coughed up by CODESOURCERER\"\n    }\n  ]\n}\n```\n"),
			},
		},
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("{\n  \"merge_id\": \"merge_abcd1234efgh5678ijkl9101mnopqrstuvwx_45\",\n  \"context\": \"This PR introduces math utility functions for basic operations and integrates them into a calculator module.\",\n  \"framework\": \"pytest\",\n  \"test_directory\": \"tests/\",\n  \"comments\": \"on\",\n  \"files\": [\n    {\n      \"path\": \"math_utils.py\",\n      \"content\": \"def add(a, b):\\n    return a + b\\n\\ndef subtract(a, b):\\n    return a - b\\n\\ndef divide(a, b):\\n    if b == 0:\\n        raise ValueError(\\\"Cannot divide by zero\\\")\\n    return a / b\",\n      \"dependencies\": []\n    },\n    {\n      \"path\": \"calculator/calc_engine.py\",\n      \"content\": \"from math_utils import add, subtract, divide\\n\\ndef calculate(expression):\\n    # A simple parser for 'a op b' expressions\\n    parts = expression.split()\\n    a = int(parts[0])\\n    op = parts[1]\\n    b = int(parts[2])\\n\\n    if op == '+':\\n        return add(a, b)\\n    elif op == '-':\\n        return subtract(a, b)\\n    elif op == '/':\\n        return divide(a, b)\\n    else:\\n        raise ValueError(\\\"Unsupported operation\\\")\",\n      \"dependencies\": [\n        {\n          \"name\": \"math_utils.py\",\n          \"content\": \"def add(a, b):\\n    return a + b\\n\\ndef subtract(a, b):\\n    return a - b\\n\\ndef divide(a, b):\\n    if b == 0:\\n        raise ValueError(\\\"Cannot divide by zero\\\")\\n    return a / b\"\n        }\n      ]\n    }\n  ]\n}"),
			},
		},
		{
			Role: "model",
			Parts: []genai.Part{
				genai.Text("```json\n{\n  \"tests\": [\n    {\n      \"testname\": \"test_math_utils\",\n      \"testfilepath\": \"tests/test_math_utils.py\",\n      \"parentpath\": \"math_utils.py\",\n      \"code\": \"# tests/test_math_utils.py\\nimport pytest\\nfrom math_utils import add, subtract, divide\\n\\n\\n# Test case for add function with valid inputs\\ndef test_add_valid_inputs():\\n    # Test adding two positive numbers\\n    assert add(5, 3) == 8\\n    # Test adding a positive and a negative number\\n    assert add(5, -3) == 2\\n    # Test adding two negative numbers\\n    assert add(-5, -3) == -8\\n    # Test adding zero\\n    assert add(5, 0) == 5\\n\\n\\n# Test case for subtract function with valid inputs\\ndef test_subtract_valid_inputs():\\n    # Test subtracting a smaller number from a larger number\\n    assert subtract(5, 3) == 2\\n    # Test subtracting a larger number from a smaller number\\n    assert subtract(3, 5) == -2\\n    # Test subtracting a negative number\\n    assert subtract(5, -3) == 8\\n    # Test subtracting zero\\n    assert subtract(5, 0) == 5\\n\\n\\n# Test case for divide function with valid inputs\\ndef test_divide_valid_inputs():\\n    # Test dividing two positive numbers\\n    assert divide(6, 3) == 2\\n    # Test dividing a number by 1\\n    assert divide(6, 1) == 6\\n    # Test dividing a negative number by a positive number\\n    assert divide(-6, 3) == -2\\n    # Test dividing a positive number by a negative number\\n    assert divide(6, -3) == -2\\n    # Test dividing negative numbers\\n    assert divide(-6, -3) == 2\\n\\n\\n# Test case for divide function with zero division\\ndef test_divide_by_zero():\\n    # Test dividing by zero, expecting a ValueError\\n    with pytest.raises(ValueError, match=\\\"Cannot divide by zero\\\"):\\n        divide(5, 0)\\n\\n# Coughed up by CODESOURCERER\"\n    },\n    {\n      \"testname\": \"test_calc_engine\",\n      \"testfilepath\": \"tests/calculator/test_calc_engine.py\",\n      \"parentpath\": \"calculator/calc_engine.py\",\n      \"code\": \"# tests/calculator/test_calc_engine.py\\nimport pytest\\nfrom calculator.calc_engine import calculate\\n\\n\\n# Test case for calculate function with valid addition expression\\ndef test_calculate_add_valid_expression():\\n    # Test valid addition expression\\n    assert calculate(\\\"5 + 3\\\") == 8\\n\\n\\n# Test case for calculate function with valid subtraction expression\\ndef test_calculate_subtract_valid_expression():\\n    # Test valid subtraction expression\\n    assert calculate(\\\"5 - 3\\\") == 2\\n\\n\\n# Test case for calculate function with valid division expression\\ndef test_calculate_divide_valid_expression():\\n    # Test valid division expression\\n    assert calculate(\\\"6 / 3\\\") == 2\\n\\n\\n# Test case for calculate function with unsupported operation\\ndef test_calculate_unsupported_operation():\\n    # Test an unsupported operation, expecting a ValueError\\n    with pytest.raises(ValueError, match=\\\"Unsupported operation\\\"):\\n        calculate(\\\"5 * 3\\\")\\n\\n\\n# Test case for calculate function with division by zero\\ndef test_calculate_divide_by_zero():\\n  with pytest.raises(ValueError, match = \\\"Cannot divide by zero\\\"):\\n    calculate(\\\"5 / 0\\\")\\n\\n\\n# Test case for calculate with non-int inputs\\ndef test_calculate_non_int_inputs():\\n    with pytest.raises(ValueError):\\n        calculate(\\\"a + b\\\")\\n\\n# Coughed up by CODESOURCERER\"\n    }\n  ]\n}\n```\n"),
			},
		},
	}
	// 	payloadBytes, err := json.Marshal(payload)
	// 	if err != nil {
	// 		fmt.Printf("Error serializing payload: %v\n", err)
	// 		return "Error processing AI request"
	// 	}

	// 	payloadString := string(payloadBytes)
	// 	ctx, cancel := context.WithTimeout(ctx, time.Duration(15*time.Second))
	// 	defer cancel()
	// 	response, err := session.SendMessage(ctx, genai.Text(payloadString))
	// 	if err != nil {
	// 		if errors.Is(err, context.DeadlineExceeded) {
	// 			return "Error: Model response timed out."
	// 		}
	// 		log.Printf("Error generating response: %v", err)
	// 		return fmt.Sprintf("Error generating response: %v", err)
	// 	}

	// 	if len(response.Candidates) == 0 {
	// 		fmt.Println("Model did not generate any response.")
	// 		return "failure"
	// 	} else {
	// 		_, err := json.MarshalIndent(response, "", "  ")
	// 		if err != nil {
	// 			fmt.Println("Error marshalling response:", err)
	// 		} else {
	// 			fmt.Println("Response processed")
	// 		}
	// 		var stringResponse string
	// 		var cleaned string
	// 		if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
	// 			finalResponse := response.Candidates[0].Content.Parts[0]
	// 			fmt.Println(reflect.TypeOf(finalResponse))
	// 			jsonData, _ := json.Marshal(finalResponse)
	// 			stringResponse = string(jsonData)
	// 		}
	// 		if len(stringResponse) > 12 {
	// 			cleaned = stringResponse[10 : len(stringResponse)-6]
	// 			var result map[string]interface{}
	// 			err := json.Unmarshal([]byte(cleaned), &result)
	// 			if err != nil {
	// 				fmt.Println("Error parsing JSON")
	// 			} else {
	// 				fmt.Println("Parsed JSON")
	// 			}
	// 		} else {
	// 			fmt.Println("Input string is too short to slice!")
	// 		}
	// 		return cleaned
	// 	}
	// }
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

	// Marshal the content of the response to JSON format
	jsonData, err := json.Marshal(response.Candidates[0].Content.Parts[0])
	if err != nil {
		return nil, fmt.Errorf("error marshaling response to JSON: %v", err)
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, jsonData, "", "  ") // Indentation with 2 spaces
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
	}

	fmt.Println(prettyJSON.String())
	return jsonData, nil
}
