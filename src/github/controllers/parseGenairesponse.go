package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Test struct {
	TestName     string `json:"testname"`
	TestFilePath string `json:"testfilepath"`
	ParentPath   string `json:"parentpath"`
	Code         string `json:"code"`
}

func ParseServer2Response(rawResponse string) (map[string]interface{}, error) {
	// Step 1: Parse the raw response to extract the output field
	type Server2Response struct {
		Output string `json:"output"`
	}
	var server2Response Server2Response
	err := json.Unmarshal([]byte(rawResponse), &server2Response)
	if err != nil {
		log.Printf("Error unmarshaling Server 2 response: %v", err)
		return nil, fmt.Errorf("invalid JSON structure: %v", err)
	}

	// Step 2: Check if the output field exists
	if server2Response.Output == "" {
		return nil, fmt.Errorf("missing 'output' field in response")
	}

	// Step 3: Unescape the JSON string inside the `output` field
	unescapedOutput := strings.ReplaceAll(server2Response.Output, `\\`, `\`)

	// Step 4: Parse the unescaped string into a valid JSON object
	var parsedOutput map[string]interface{}
	err = json.Unmarshal([]byte(unescapedOutput), &parsedOutput)
	if err != nil {
		log.Printf("Error parsing output JSON: %v", err)
		return nil, fmt.Errorf("invalid 'output' JSON: %v", err)
	}

	return parsedOutput, nil
}

func TestParseServer2Response(c *gin.Context) {
	// Define a slice of tests
	response := []Test{
		{
			TestName:     "test_file_operations",
			TestFilePath: "tests/test_file_operations.py",
			ParentPath:   "file_operations.py",
			Code: `import pytest
from file_operations import read_file, write_file
import tempfile
import os

def test_read_file():
    with tempfile.NamedTemporaryFile(mode="w", delete=False) as f:
        f.write("test content")
        filepath = f.name
    assert read_file(filepath) == "test content"
    os.remove(filepath)
    # Coughed up by CODESOURCERER

def test_read_file_nonexistent():
    with pytest.raises(FileNotFoundError):
        read_file('nonexistent_file.txt')
    # Coughed up by CODESOURCERER

def test_write_file():
    with tempfile.NamedTemporaryFile(mode="w", delete=False) as f:
        filepath = f.name
        write_file(filepath, "test content")
    with open(filepath, 'r') as f:
        assert f.read() == "test content"
    os.remove(filepath)
    # Coughed up by CODESOURCERER

def test_write_file_empty():
    with tempfile.NamedTemporaryFile(mode="w", delete=False) as f:
        filepath = f.name
        write_file(filepath, "")
    with open(filepath, 'r') as f:
        assert f.read() == ""
    os.remove(filepath)
    # Coughed up by CODESOURCERER`,
		},
		{
			TestName:     "test_main",
			TestFilePath: "tests/test_main.py",
			ParentPath:   "main.py",
			Code: `import pytest
from main import main
from unittest.mock import patch
from file_operations import write_file
import tempfile
import os

@patch('main.process_content')
@patch('main.write_file')
@patch('main.read_file')
def test_main(mock_read_file, mock_write_file, mock_process_content):
    mock_read_file.return_value = "test content"
    mock_process_content.return_value = "modified content"
    with tempfile.NamedTemporaryFile(mode="w", delete=False) as input_file:
        input_file.write("test content")
        input_filepath = input_file.name
    main()
    mock_read_file.assert_called_once_with(input_filepath)
    mock_process_content.assert_called_once_with("test content")
    mock_write_file.assert_called_once_with('output.txt', "modified content")
    os.remove(input_filepath)
    os.remove('output.txt')
    # Coughed up by CODESOURCERER`,
		},
	}

	// Return the response as JSON
	c.JSON(http.StatusOK, gin.H{"tests": response})
}
