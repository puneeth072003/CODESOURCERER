package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ParseServer2Response(rawResponse string) (map[string]interface{}, error) {
	// Step 1: Parse the raw response to extract the output field
	type Server2Response struct {
		Output string `json:"output"`
	}
	var server2Response Server2Response
	err := json.Unmarshal([]byte(rawResponse), &server2Response)
	if err != nil {
		log.Printf("Error unmarshaling Server 2 response: %v", err)
		return nil, err
	}

	// Step 2: Unescape the JSON string inside the `output` field
	unescapedOutput := strings.ReplaceAll(server2Response.Output, `\\`, `\`)

	// Step 3: Parse the unescaped string into a valid JSON object
	var parsedOutput map[string]interface{}
	err = json.Unmarshal([]byte(unescapedOutput), &parsedOutput)
	if err != nil {
		log.Printf("Error parsing output JSON: %v", err)
		return nil, err
	}

	return parsedOutput, nil
}

func TestParseServer2Response(c *gin.Context) {
	// Example raw response for testing
	rawResponse := `{"output":"{\\n  \\\"tests\\\": [\\n    {\\n      \\\"testname\\\": \\\"test_g2\\\",\\n      \\\"testfilepath\\\": \\\"tests/test_g2.py\\\",\\n      \\\"parentpath\\\": \\\"g2.py\\\",\\n      \\\"code\\\": \\\"import pytest\\\\nfrom g2 import combinations\\\\nfrom q1 import factorial #importing factorial to test edge cases\\\\n\\\\ndef test_combinations_valid_input():\\\\n    assert combinations(5, 2) == 10.0\\\\n    #Coughed up by CODESOURCERER\\\\n\\\\ndef test_combinations_edge_cases():\\\\n    assert combinations(0, 0) == 1.0\\\\n    assert combinations(5, 0) == 1.0\\\\n    assert combinations(5, 5) == 1.0\\\\n    #Coughed up by CODESOURCERER\\\\n\\\\ndef test_combinations_invalid_input():\\\\n    with pytest.raises(ValueError):\\\\n        combinations(5, 6) #r \\u003e n\\\\n    with pytest.raises(ValueError):\\\\n        combinations(-1,2) #negative n\\\\n    with pytest.raises(ValueError):\\\\n        combinations(5,-2) #negative r\\\\n    #Coughed up by CODESOURCERER\\\\n\\\\ndef test_factorial_recursive_call(): #testing factorial function which is dependency of combinations\\\\n    assert factorial(5) == 120\\\\n    assert factorial(0) == 1\\\\n    with pytest.raises(RecursionError):\\\\n        factorial(-1)\\\\n    #Coughed up by CODESOURCERER\\\"\\n    },\\n    {\\n      \\\"testname\\\": \\\"test_g3\\\",\\n      \\\"testfilepath\\\": \\\"tests/test_g3.py\\\",\\n      \\\"parentpath\\\": \\\"g3.py\\\",\\n      \\\"code\\\": \\\"import pytest\\\\nfrom g2 import combinations\\\\n\\\\ndef test_g3_output_correctness(capsys):\\\\n    import g3\\\\n    captured = capsys.readouterr()\\\\n    assert \\\\\\\"Combinations of 5 items taken 2 at a time: 10.0\\\\\\\" in captured.out\\\\n    #Coughed up by CODESOURCERER\\\"\\n    }\\n  ]\\n}\\n"}`

	// Call the parsing function
	parsedJSON, err := ParseServer2Response(rawResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to parse response: %v", err),
		})
		return
	}

	// Return the parsed JSON as the response
	c.JSON(http.StatusOK, parsedJSON)
}
