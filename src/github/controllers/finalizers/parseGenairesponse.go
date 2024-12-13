package finalizers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Test struct {
	TestName     string `json:"testname"`
	TestFilePath string `json:"testfilepath"`
	ParentPath   string `json:"parentpath"`
	Code         string `json:"code"`
}

type TestsContainer struct {
	Tests []Test `json:"tests"`
}

func ParseServer2Response(rawResponse string) ([]Test, error) {
	jsonString := `{
		"tests": [
		  {
			"testname": "test_q2",
			"testfilepath": "tests/test_q2.py",
			"parentpath": "q2.py",
			"code": "import pytest\nfrom q2 import combinations\n\ndef test_combinations_valid_input():\n    assert combinations(5, 2) == 10.0\n\ndef test_combinations_edge_cases():\n    assert combinations(0, 0) == 1.0\n    assert combinations(5, 0) == 1.0\n    assert combinations(5, 5) == 1.0\n\ndef test_combinations_invalid_input():\n    with pytest.raises(ValueError):\n       combinations(5, 6)\n\n# Coughed up by CODESOURCERER"
		  },
		  {
			"testname": "test_q3",
			"testfilepath": "tests/test_q3.py",
			"parentpath": "q3.py",
			"code": "import pytest\nfrom q3 import combinations\n\ndef test_q3_output_correctness(capsys):\n    import q3\n    captured = capsys.readouterr()\n    assert \"Combinations of 5 items taken 2 at a time: 10.0\" in captured.out\n\n# Coughed up by CODESOURCERER"
		  }
		]
	  }`
	// Create a variable to hold the parsed data
	var testsContainer TestsContainer

	// Parse the JSON string
	err := json.Unmarshal([]byte(jsonString), &testsContainer)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return testsContainer.Tests, nil
}

func TestParseServer2Response(c *gin.Context) {
	// Parse the response
	response, err := ParseServer2Response("")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the parsed response as JSON
	c.JSON(http.StatusOK, gin.H{"response": response})
}
