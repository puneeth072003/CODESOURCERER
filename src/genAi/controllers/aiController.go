package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
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
	aiOutput := processAI(ctx, payload, model)
	// response := ResponsePayload{Output: aiOutput}
	response := aiOutput
	c.JSON(http.StatusOK, response)
}

func processAI(ctx context.Context, payload RequestPayload, model *genai.GenerativeModel) string {
	session := model.StartChat()
	session.History = []*genai.Content{
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("Key Elements of the Payload:\nmerge_id: A unique identifier for the merge request.\ncontext: A description of what the PR is intended to do.\nfiles:\nContains the files for which test cases must be generated.\nEach file has:\npath: The file path within the repository.\ncontent: The entire content of the file.\ndependencies: An array of files that the current file depends on. Each dependency includes:\nname: The dependency file's name.\ncontent: The dependency file's content.\nframework: Specifies the testing framework to be used (e.g., unittest, pytest, etc.).\nThe generated test cases must adhere to this framework.\nExpected Output:\nThe generated output must contain a tests array.\nEach element in the tests array represents a file and contains:\ntestname: Must follow the naming convention test_<file_name>.\npath: The path of the file being tested.\ntests: An array of individual test cases specific to that file.\nEach test case must include:\ntestname: A descriptive name for the test case.\npath: The path of the file being tested.\ncode: The actual code for the test case, written in the specified framework.\nSpecific Instructions for Test Case Generation:\nNaming Convention:\nUse test_<file_name> as the name for the main test suite for each file.\nFor individual test cases, use descriptive names that reflect the functionality being tested.\nTest Framework:\nAdhere strictly to the testing framework specified in the framework field.\nFor unittest, create class-based tests with unittest.TestCase.\nFor pytest, write function-based tests.\nDependencies:\nAnalyze the dependencies array to provide better test coverage and context.\nMock or import dependencies as needed to construct meaningful test cases.\nContent-Based Test Creation:\nUse the content of the file to determine:\nFunctions or classes to test.\nLogical paths, edge cases, and expected outputs.\nEdge Cases:\nInclude test cases for common edge cases and failure conditions wherever applicable.\nExample Input Payload:\njson\nCopy code\n{\n  \"merge_id\": \"merge_7b9a17d77fee12665a90eb52d5d98c4077ceddd7_21\",\n  \"commit_sha\": \"7b9a17d77fee12665a90eb52d5d98c4077ceddd7\",\n  \"pull_request\": 21,\n  \"context\": \"This PR is calculating factorial and combination\",\n  \"framework\": \"pytest\",\n  \"files\": [\n    {\n      \"path\": \"d2.py\",\n      \"content\": \"from q1 import factorial\\n\\ndef combinations(n, r):\\n    return factorial(n) / (factorial(r) * factorial(n - r))\",\n      \"dependencies\": [\n        {\n          \"name\": \"q1.py\",\n          \"content\": \"def factorial(n):\\n    if n == 0:\\n        return 1\\n    else:\\n        return n * factorial(n - 1)\"\n        }\n      ]\n    },\n    {\n      \"path\": \"d3.py\",\n      \"content\": \"from d2 import combinations\\n\\nn = 5\\nr = 2\\nresult = combinations(n, r)\\nprint(f\\\"Combinations of {n} items taken {r} at a time: {result}\\\")\",\n      \"dependencies\": [\n        {\n          \"name\": \"d2.py\",\n          \"content\": \"from q1 import factorial\\n\\ndef combinations(n, r):\\n    return factorial(n) / (factorial(r) * factorial(n - r))\"\n        }\n      ]\n    }\n  ]\n}\nExample Output:\nFor the input payload above, the expected output will look like this:\n\njson\nCopy code\n{\n  \"tests\": [\n    {\n      \"testname\": \"test_d2\",\n      \"path\": \"d2.py\",\n      \"tests\": [\n        {\n          \"testname\": \"test_combinations_valid_input\",\n          \"path\": \"d2.py\",\n          \"code\": \"def test_combinations_valid_input():\\n    from d2 import combinations\\n    assert combinations(5, 2) == 10\"\n        },\n        {\n          \"testname\": \"test_combinations_edge_cases\",\n          \"path\": \"d2.py\",\n          \"code\": \"def test_combinations_edge_cases():\\n    from d2 import combinations\\n    assert combinations(0, 0) == 1\\n    assert combinations(5, 0) == 1\"\n        }\n      ]\n    },\n    {\n      \"testname\": \"test_d3\",\n      \"path\": \"d3.py\",\n      \"tests\": [\n        {\n          \"testname\": \"test_d3_output_correctness\",\n          \"path\": \"d3.py\",\n          \"code\": \"def test_d3_output_correctness(capsys):\\n    import d3\\n    captured = capsys.readouterr()\\n    assert \\\"Combinations of 5 items taken 2 at a time: 10\\\" in captured.out\"\n        }\n      ]\n    }\n  ]\n}\nAdditional Guidelines:\nEnsure test cases are modular and test one aspect of functionality per test.\nIf dependencies are imported, verify their correctness in the context of the file under test.\nTests must be written in the specified framework and leverage its features (e.g., assert for pytest or self.assertEqual for unittest).\nKeep test code concise, readable, and relevant."),
			},
		},
		{
			Role: "model",
			Parts: []genai.Part{
				genai.Text("```json\n{\n  \"tests\": [\n    {\n      \"testname\": \"test_d2\",\n      \"path\": \"d2.py\",\n      \"tests\": [\n        {\n          \"testname\": \"test_combinations_valid_input\",\n          \"path\": \"d2.py\",\n          \"code\": \"def test_combinations_valid_input():\\n    from d2 import combinations\\n    assert combinations(5, 2) == 10.0\"\n        },\n        {\n          \"testname\": \"test_combinations_edge_cases\",\n          \"path\": \"d2.py\",\n          \"code\": \"def test_combinations_edge_cases():\\n    from d2 import combinations\\n    assert combinations(0, 0) == 1.0\\n    assert combinations(5, 0) == 1.0\\n    assert combinations(5, 5) == 1.0\"\n        },\n        {\n          \"testname\": \"test_combinations_invalid_input\",\n          \"path\": \"d2.py\",\n          \"code\": \"def test_combinations_invalid_input():\\n    from d2 import combinations\\n    from pytest import raises\\n    with raises(ValueError):\\n        combinations(5, 6)  #r > n\\n    with raises(ValueError):\\n        combinations(-1,2) #negative n\\n    with raises(ValueError):\\n        combinations(5,-2) #negative r\"\n        }\n      ]\n    },\n    {\n      \"testname\": \"test_d3\",\n      \"path\": \"d3.py\",\n      \"tests\": [\n        {\n          \"testname\": \"test_d3_output_correctness\",\n          \"path\": \"d3.py\",\n          \"code\": \"def test_d3_output_correctness(capsys):\\n    import d3\\n    captured = capsys.readouterr()\\n    assert \\\"Combinations of 5 items taken 2 at a time: 10.0\\\" in captured.out\"\n        }\n      ]\n    },\n    {\n      \"testname\": \"test_q1\",\n      \"path\": \"q1.py\",\n      \"tests\": [\n        {\n          \"testname\": \"test_factorial_positive\",\n          \"path\": \"q1.py\",\n          \"code\": \"def test_factorial_positive():\\n    from q1 import factorial\\n    assert factorial(5) == 120\"\n        },\n        {\n          \"testname\": \"test_factorial_zero\",\n          \"path\": \"q1.py\",\n          \"code\": \"def test_factorial_zero():\\n    from q1 import factorial\\n    assert factorial(0) == 1\"\n        },\n        {\n          \"testname\": \"test_factorial_negative\",\n          \"path\": \"q1.py\",\n          \"code\": \"def test_factorial_negative():\\n    from q1 import factorial\\n    from pytest import raises\\n    with raises(ValueError):\\n        factorial(-1)\"\n        }\n      ]\n    }\n  ]\n}\n```\n"),
			},
		},
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("{\n  \"merge_id\": \"merge_ab1234ef5678gh90ij12klmnop345qrstu678vwxy_45\",\n  \"commit_sha\": \"ab1234ef5678gh90ij12klmnop345qrstu678vwxy\",\n  \"pull_request\": 45,\n  \"context\": \"This PR provides utility functions for mathematical operations and their usage.\",\n  \"framework\": \"unittest\",\n  \"files\": [\n    {\n      \"path\": \"math_utils.py\",\n      \"content\": \"def add(a, b):\\n    return a + b\\n\\ndef multiply(a, b):\\n    return a * b\",\n      \"dependencies\": []\n    },\n    {\n      \"path\": \"main.py\",\n      \"content\": \"from math_utils import add, multiply\\n\\ndef calculate():\\n    x = 10\\n    y = 5\\n    return {\\n        'sum': add(x, y),\\n        'product': multiply(x, y)\\n    }\",\n      \"dependencies\": [\n        {\n          \"name\": \"math_utils.py\",\n          \"content\": \"def add(a, b):\\n    return a + b\\n\\ndef multiply(a, b):\\n    return a * b\"\n        }\n      ]\n    }\n  ]\n}\n"),
			},
		},
		{
			Role: "model",
			Parts: []genai.Part{
				genai.Text("```json\n{\n  \"tests\": [\n    {\n      \"testname\": \"test_math_utils\",\n      \"path\": \"math_utils.py\",\n      \"tests\": [\n        {\n          \"testname\": \"test_add\",\n          \"path\": \"math_utils.py\",\n          \"code\": \"import unittest\\nfrom math_utils import add\\n\\nclass TestMathUtils(unittest.TestCase):\\n    def test_add(self):\\n        self.assertEqual(add(2, 3), 5)\\n        self.assertEqual(add(-1, 1), 0)\\n        self.assertEqual(add(0, 0), 0)\"\n        },\n        {\n          \"testname\": \"test_multiply\",\n          \"path\": \"math_utils.py\",\n          \"code\": \"import unittest\\nfrom math_utils import multiply\\n\\nclass TestMathUtils(unittest.TestCase):\\n    def test_multiply(self):\\n        self.assertEqual(multiply(2, 3), 6)\\n        self.assertEqual(multiply(-1, 1), -1)\\n        self.assertEqual(multiply(0, 5), 0)\\n        self.assertEqual(multiply(5,0),0)\"\n        }\n      ]\n    },\n    {\n      \"testname\": \"test_main\",\n      \"path\": \"main.py\",\n      \"tests\": [\n        {\n          \"testname\": \"test_calculate\",\n          \"path\": \"main.py\",\n          \"code\": \"import unittest\\nfrom main import calculate\\n\\nclass TestMain(unittest.TestCase):\\n    def test_calculate(self):\\n        result = calculate()\\n        self.assertEqual(result['sum'], 15)\\n        self.assertEqual(result['product'], 50)\"\n        }\n      ]\n    }\n  ]\n}\n```\n"),
			},
		},
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("modify the path inside tests element as parentpath and aslo add another element called as testfilepath that teslls  the path in which final test file should be placed, which is usually in the same directory as the parent file but under folder named tests if not mentioned explictily in the incomming payload "),
			},
		},
		{
			Role: "model",
			Parts: []genai.Part{
				genai.Text("```json\n{\n  \"tests\": [\n    {\n      \"testname\": \"test_math_utils\",\n      \"testfilepath\": \"tests/test_math_utils.py\",\n      \"tests\": [\n        {\n          \"testname\": \"test_add\",\n          \"parentpath\": \"math_utils.py\",\n          \"code\": \"import unittest\\nfrom math_utils import add\\n\\nclass TestMathUtils(unittest.TestCase):\\n    def test_add(self):\\n        self.assertEqual(add(2, 3), 5)\\n        self.assertEqual(add(-1, 1), 0)\\n        self.assertEqual(add(0, 0), 0)\"\n        },\n        {\n          \"testname\": \"test_multiply\",\n          \"parentpath\": \"math_utils.py\",\n          \"code\": \"import unittest\\nfrom math_utils import multiply\\n\\nclass TestMathUtils(unittest.TestCase):\\n    def test_multiply(self):\\n        self.assertEqual(multiply(2, 3), 6)\\n        self.assertEqual(multiply(-1, 1), -1)\\n        self.assertEqual(multiply(0, 5), 0)\\n        self.assertEqual(multiply(5,0),0)\"\n        }\n      ]\n    },\n    {\n      \"testname\": \"test_main\",\n      \"testfilepath\": \"tests/test_main.py\",\n      \"tests\": [\n        {\n          \"testname\": \"test_calculate\",\n          \"parentpath\": \"main.py\",\n          \"code\": \"import unittest\\nfrom main import calculate\\n\\nclass TestMain(unittest.TestCase):\\n    def test_calculate(self):\\n        result = calculate()\\n        self.assertEqual(result['sum'], 15)\\n        self.assertEqual(result['product'], 50)\"\n        }\n      ]\n    }\n  ]\n}\n```\n"),
			},
		},
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("please generate only single test object  for each file menioned in files and also each object inted of containing element called as tests make it as code"),
			},
		},
		{
			Role: "model",
			Parts: []genai.Part{
				genai.Text("```json\n{\n  \"tests\": [\n    {\n      \"testname\": \"test_math_utils\",\n      \"testfilepath\": \"tests/test_math_utils.py\",\n      \"parentpath\": \"math_utils.py\",\n      \"code\": \"import unittest\\nfrom math_utils import add, multiply\\n\\nclass TestMathUtils(unittest.TestCase):\\n    def test_add(self):\\n        self.assertEqual(add(2, 3), 5)\\n        self.assertEqual(add(-1, 1), 0)\\n        self.assertEqual(add(0, 0), 0)\\n    def test_multiply(self):\\n        self.assertEqual(multiply(2, 3), 6)\\n        self.assertEqual(multiply(-1, 1), -1)\\n        self.assertEqual(multiply(0, 5), 0)\\n        self.assertEqual(multiply(5,0),0)\"\n    },\n    {\n      \"testname\": \"test_main\",\n      \"testfilepath\": \"tests/test_main.py\",\n      \"parentpath\": \"main.py\",\n      \"code\": \"import unittest\\nfrom main import calculate\\n\\nclass TestMain(unittest.TestCase):\\n    def test_calculate(self):\\n        result = calculate()\\n        self.assertEqual(result['sum'], 15)\\n        self.assertEqual(result['product'], 50)\"\n    }\n  ]\n}\n```\n"),
			},
		},
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("good everything is seems to be ok but one thing please use import statements relative to parent file path that you mentioned "),
			},
		},
		{
			Role: "model",
			Parts: []genai.Part{
				genai.Text("```json\n{\n  \"tests\": [\n    {\n      \"testname\": \"test_math_utils\",\n      \"testfilepath\": \"tests/test_math_utils.py\",\n      \"parentpath\": \"math_utils.py\",\n      \"code\": \"import unittest\\nfrom .math_utils import add, multiply\\n\\nclass TestMathUtils(unittest.TestCase):\\n    def test_add(self):\\n        self.assertEqual(add(2, 3), 5)\\n        self.assertEqual(add(-1, 1), 0)\\n        self.assertEqual(add(0, 0), 0)\\n    def test_multiply(self):\\n        self.assertEqual(multiply(2, 3), 6)\\n        self.assertEqual(multiply(-1, 1), -1)\\n        self.assertEqual(multiply(0, 5), 0)\\n        self.assertEqual(multiply(5,0),0)\"\n    },\n    {\n      \"testname\": \"test_main\",\n      \"testfilepath\": \"tests/test_main.py\",\n      \"parentpath\": \"main.py\",\n      \"code\": \"import unittest\\nfrom .main import calculate\\n\\nclass TestMain(unittest.TestCase):\\n    def test_calculate(self):\\n        result = calculate()\\n        self.assertEqual(result['sum'], 15)\\n        self.assertEqual(result['product'], 50)\"\n    }\n  ]\n}\n```\n"),
			},
		},
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("{\n  \"merge_id\": \"merge_cd456klm789gh123opq890rst456uvwxyz123abc_101\",\n  \"context\": \"This PR introduces utility functions for string transformations and their usage in a main script.\",\n  \"framework\": \"pytest\",\n  \"test_directory\": \"src/puneeth/tests/\",\n  \"files\": [\n    {\n      \"path\": \"string_utils.py\",\n      \"content\": \"def to_uppercase(string):\\n    return string.upper()\\n\\ndef to_lowercase(string):\\n    return string.lower()\",\n      \"dependencies\": []\n    },\n    {\n      \"path\": \"main.py\",\n      \"content\": \"from string_utils import to_uppercase, to_lowercase\\n\\ndef transform_string(input_str, transform_type):\\n    if transform_type == 'upper':\\n        return to_uppercase(input_str)\\n    elif transform_type == 'lower':\\n        return to_lowercase(input_str)\\n    else:\\n        raise ValueError(\\\"Invalid transform type\\\")\",\n      \"dependencies\": [\n        {\n          \"name\": \"string_utils.py\",\n          \"content\": \"def to_uppercase(string):\\n    return string.upper()\\n\\ndef to_lowercase(string):\\n    return string.lower()\"\n        }\n      ]\n    }\n  ]\n}"),
			},
		},
		{
			Role: "model",
			Parts: []genai.Part{
				genai.Text("```json\n{\n  \"tests\": [\n    {\n      \"testname\": \"test_string_utils\",\n      \"testfilepath\": \"src/puneeth/tests/test_string_utils.py\",\n      \"parentpath\": \"string_utils.py\",\n      \"code\": \"import pytest\\nfrom .string_utils import to_uppercase, to_lowercase\\n\\ndef test_to_uppercase():\\n    assert to_uppercase('hello') == 'HELLO'\\n    assert to_uppercase('') == ''\\n    assert to_uppercase('HeLlO') == 'HELLO'\\n\\ndef test_to_lowercase():\\n    assert to_lowercase('WORLD') == 'world'\\n    assert to_lowercase('') == ''\\n    assert to_lowercase('WoRlD') == 'world'\"\n    },\n    {\n      \"testname\": \"test_main\",\n      \"testfilepath\": \"src/puneeth/tests/test_main.py\",\n      \"parentpath\": \"main.py\",\n      \"code\": \"import pytest\\nfrom .main import transform_string\\n\\ndef test_transform_string_upper():\\n    assert transform_string('test', 'upper') == 'TEST'\\n\\ndef test_transform_string_lower():\\n    assert transform_string('TEST', 'lower') == 'test'\\n\\ndef test_transform_string_invalid():\\n    with pytest.raises(ValueError):\\n        transform_string('test', 'invalid')\"\n    }\n  ]\n}\n```\n"),
			},
		},
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("Generate test cases for the given files while adhering to the following rules:\n\nUse the testing framework specified in the framework field from the input payload.\nThe test file's location must adhere to the path specified in test_directory. If not provided, default to tests/<module_name>/test_<file_name>.py.\nEnsure that every function in the file (parentpath) has corresponding tests, covering both normal and edge cases.\nMock or account for all dependencies mentioned in the input payload.\nInclude detailed assertions for all possible scenarios, including typical inputs, edge cases, and exceptions.\nFormat the generated test code according to PEP-8 standards.\nInclude meaningful comments explaining the purpose of each test.\nValidate and report any missing test_directory, ensuring fallback to a default structure.\nEnsure compatibility with Python 3.8+ unless otherwise specified.\n\nInstruction for Import Paths:\nWhen generating test cases, ensure that the import paths are consistent with the project structure. Use absolute imports for clarity and reliability. Include __init__.py in all relevant directories to make them Python packages. For example:\n\nIf the source file is at src/puneeth/string_utils.py and the test file is at tests/string_operations/test_string_utils.py, the import should be:\nfrom src.puneeth.string_utils import to_uppercase\n"),
			},
		},
		{
			Role: "model",
			Parts: []genai.Part{
				genai.Text("```json\n{\n  \"tests\": [\n    {\n      \"testname\": \"test_string_utils\",\n      \"testfilepath\": \"src/puneeth/tests/test_string_utils.py\",\n      \"parentpath\": \"string_utils.py\",\n      \"code\": \"# coding=utf-8\\n\\\"\\\"\\\"Unit tests for string_utils.py\\\"\\\"\\\"\\nimport unittest\\nfrom src.puneeth.string_utils import to_uppercase, to_lowercase\\n\\nclass TestStringUtils(unittest.TestCase):\\n    def test_to_uppercase_normal_case(self):\\n        # Test with a normal string\\n        self.assertEqual(to_uppercase('hello'), 'HELLO')\\n\\n    def test_to_uppercase_empty_string(self):\\n        # Test with an empty string\\n        self.assertEqual(to_uppercase(''), '')\\n\\n    def test_to_uppercase_mixed_case(self):\\n        # Test with a mixed case string\\n        self.assertEqual(to_uppercase('HeLlO'), 'HELLO')\\n\\n    def test_to_uppercase_special_characters(self):\\n        # Test with special characters\\n        self.assertEqual(to_uppercase('Hello, world!'), 'HELLO, WORLD!')\\n\\n    def test_to_lowercase_normal_case(self):\\n        # Test with a normal string\\n        self.assertEqual(to_lowercase('WORLD'), 'world')\\n\\n    def test_to_lowercase_empty_string(self):\\n        # Test with an empty string\\n        self.assertEqual(to_lowercase(''), '')\\n\\n    def test_to_lowercase_mixed_case(self):\\n        # Test with a mixed case string\\n        self.assertEqual(to_lowercase('WoRlD'), 'world')\\n\\n    def test_to_lowercase_special_characters(self):\\n        # Test with special characters\\n        self.assertEqual(to_lowercase('Hello, world!'), 'hello, world!')\\n\\n\"\n    },\n    {\n      \"testname\": \"test_main\",\n      \"testfilepath\": \"src/puneeth/tests/test_main.py\",\n      \"parentpath\": \"main.py\",\n      \"code\": \"# coding=utf-8\\n\\\"\\\"\\\"Unit tests for main.py\\\"\\\"\\\"\\nimport unittest\\nfrom src.puneeth.main import transform_string\\nfrom src.puneeth.string_utils import to_uppercase, to_lowercase #Needed to test error handling\\n\\nclass TestMain(unittest.TestCase):\\n    def test_transform_string_upper(self):\\n        # Test with upper transformation\\n        self.assertEqual(transform_string('test', 'upper'), 'TEST')\\n\\n    def test_transform_string_lower(self):\\n        # Test with lower transformation\\n        self.assertEqual(transform_string('TEST', 'lower'), 'test')\\n\\n    def test_transform_string_invalid(self):\\n        # Test with invalid transformation type\\n        with self.assertRaises(ValueError) as context:\\n            transform_string('test', 'invalid')\\n        self.assertEqual(str(context.exception), \\\"Invalid transform type\\\")\\n\\n    def test_transform_string_empty(self):\\n        #Test with empty string\\n        self.assertEqual(transform_string('', 'upper'), '')\\n        self.assertEqual(transform_string('', 'lower'), '')\\n\\n    def test_transform_string_special_chars(self):\\n        #Test with special characters\\n        self.assertEqual(transform_string('Te$t!', 'upper'), 'TE$T!')\\n        self.assertEqual(transform_string('TE$T!', 'lower'), 'te$t!')\\n\"\n    }\n  ]\n}\n```\n"),
			},
		},
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("comments should be generated based on the \"comments\" element in incomming payload it will be \"on\" or \"off\", if it is on then generate comment lines in the code and if off then just don't generate comments in the test code.\n\nAnd mandatorliy add this comment at the end of each test generated \"#Coughed up by CODESOURCERER\""),
			},
		},
		{
			Role: "model",
			Parts: []genai.Part{
				genai.Text("```json\n{\n  \"tests\": [\n    {\n      \"testname\": \"test_string_utils\",\n      \"testfilepath\": \"src/puneeth/tests/test_string_utils.py\",\n      \"parentpath\": \"string_utils.py\",\n      \"code\": \"# coding=utf-8\\n\\\"\\\"\\\"Unit tests for string_utils.py\\\"\\\"\\\"\\nimport unittest\\nfrom src.puneeth.string_utils import to_uppercase, to_lowercase\\n\\nclass TestStringUtils(unittest.TestCase):\\n    def test_to_uppercase_normal_case(self):\\n        # Test with a normal string\\n        self.assertEqual(to_uppercase('hello'), 'HELLO')\\n        #Coughed up by CODESOURCERER\\n\\n    def test_to_uppercase_empty_string(self):\\n        # Test with an empty string\\n        self.assertEqual(to_uppercase(''), '')\\n        #Coughed up by CODESOURCERER\\n\\n    def test_to_uppercase_mixed_case(self):\\n        # Test with a mixed case string\\n        self.assertEqual(to_uppercase('HeLlO'), 'HELLO')\\n        #Coughed up by CODESOURCERER\\n\\n    def test_to_uppercase_special_characters(self):\\n        # Test with special characters\\n        self.assertEqual(to_uppercase('Hello, world!'), 'HELLO, WORLD!')\\n        #Coughed up by CODESOURCERER\\n\\n    def test_to_lowercase_normal_case(self):\\n        # Test with a normal string\\n        self.assertEqual(to_lowercase('WORLD'), 'world')\\n        #Coughed up by CODESOURCERER\\n\\n    def test_to_lowercase_empty_string(self):\\n        # Test with an empty string\\n        self.assertEqual(to_lowercase(''), '')\\n        #Coughed up by CODESOURCERER\\n\\n    def test_to_lowercase_mixed_case(self):\\n        # Test with a mixed case string\\n        self.assertEqual(to_lowercase('WoRlD'), 'world')\\n        #Coughed up by CODESOURCERER\\n\\n    def test_to_lowercase_special_characters(self):\\n        # Test with special characters\\n        self.assertEqual(to_lowercase('Hello, world!'), 'hello, world!')\\n        #Coughed up by CODESOURCERER\\n\\n\"\n    },\n    {\n      \"testname\": \"test_main\",\n      \"testfilepath\": \"src/puneeth/tests/test_main.py\",\n      \"parentpath\": \"main.py\",\n      \"code\": \"# coding=utf-8\\n\\\"\\\"\\\"Unit tests for main.py\\\"\\\"\\\"\\nimport unittest\\nfrom src.puneeth.main import transform_string\\nfrom src.puneeth.string_utils import to_uppercase, to_lowercase #Needed to test error handling\\n\\nclass TestMain(unittest.TestCase):\\n    def test_transform_string_upper(self):\\n        # Test with upper transformation\\n        self.assertEqual(transform_string('test', 'upper'), 'TEST')\\n        #Coughed up by CODESOURCERER\\n\\n    def test_transform_string_lower(self):\\n        # Test with lower transformation\\n        self.assertEqual(transform_string('TEST', 'lower'), 'test')\\n        #Coughed up by CODESOURCERER\\n\\n    def test_transform_string_invalid(self):\\n        # Test with invalid transformation type\\n        with self.assertRaises(ValueError) as context:\\n            transform_string('test', 'invalid')\\n        self.assertEqual(str(context.exception), \\\"Invalid transform type\\\")\\n        #Coughed up by CODESOURCERER\\n\\n    def test_transform_string_empty(self):\\n        #Test with empty string\\n        self.assertEqual(transform_string('', 'upper'), '')\\n        self.assertEqual(transform_string('', 'lower'), '')\\n        #Coughed up by CODESOURCERER\\n\\n    def test_transform_string_special_chars(self):\\n        #Test with special characters\\n        self.assertEqual(transform_string('Te$t!', 'upper'), 'TE$T!')\\n        self.assertEqual(transform_string('TE$T!', 'lower'), 'te$t!')\\n        #Coughed up by CODESOURCERER\\n\"\n    }\n  ]\n}\n```\n"),
			},
		},
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("{\n  \"merge_id\": \"merge_ab123xy456op789gh456uvw123klm890rst_102\",\n  \"context\": \"This PR introduces functions for basic mathematical operations and their usage in a calculator script.\",\n  \"framework\": \"pytest\",\n  \"test_directory\": \"src/puneeth/tests/\",\n  \"comments\": \"off\",\n  \"files\": [\n    {\n      \"path\": \"math_utils.py\",\n      \"content\": \"def add(a, b):\\n    return a + b\\n\\ndef subtract(a, b):\\n    return a - b\",\n      \"dependencies\": []\n    },\n    {\n      \"path\": \"calculator.py\",\n      \"content\": \"from math_utils import add, subtract\\n\\ndef calculate(operation, a, b):\\n    if operation == 'add':\\n        return add(a, b)\\n    elif operation == 'subtract':\\n        return subtract(a, b)\\n    else:\\n        raise ValueError(\\\"Invalid operation\\\")\",\n      \"dependencies\": [\n        {\n          \"name\": \"math_utils.py\",\n          \"content\": \"def add(a, b):\\n    return a + b\\n\\ndef subtract(a, b):\\n    return a - b\"\n        }\n      ]\n    }\n  ]\n}"),
			},
		},
		{
			Role: "model",
			Parts: []genai.Part{
				genai.Text("```json\n{\n  \"tests\": [\n    {\n      \"testname\": \"test_math_utils\",\n      \"testfilepath\": \"src/puneeth/tests/test_math_utils.py\",\n      \"parentpath\": \"math_utils.py\",\n      \"code\": \"import pytest\\nfrom src.puneeth.math_utils import add, subtract\\n\\ndef test_add():\\n    assert add(2, 3) == 5\\n    assert add(-1, 1) == 0\\n    assert add(0, 0) == 0\\n    assert add(100,200) == 300\\n    assert add(-5,-5) == -10\\n    #Coughed up by CODESOURCERER\\n\\ndef test_subtract():\\n    assert subtract(5, 2) == 3\\n    assert subtract(1, 1) == 0\\n    assert subtract(0, 0) == 0\\n    assert subtract(200,100) == 100\\n    assert subtract(5,-5) == 10\\n    #Coughed up by CODESOURCERER\"\n    },\n    {\n      \"testname\": \"test_calculator\",\n      \"testfilepath\": \"src/puneeth/tests/test_calculator.py\",\n      \"parentpath\": \"calculator.py\",\n      \"code\": \"import pytest\\nfrom src.puneeth.calculator import calculate\\n\\ndef test_calculate_add():\\n    assert calculate('add', 2, 3) == 5\\n    assert calculate('add', -1, 1) == 0\\n    assert calculate('add', 0, 0) == 0\\n    #Coughed up by CODESOURCERER\\n\\ndef test_calculate_subtract():\\n    assert calculate('subtract', 5, 2) == 3\\n    assert calculate('subtract', 1, 1) == 0\\n    assert calculate('subtract', 0, 0) == 0\\n    #Coughed up by CODESOURCERER\\n\\ndef test_calculate_invalid_operation():\\n    with pytest.raises(ValueError) as excinfo:\\n        calculate('invalid', 2, 3)\\n    assert str(excinfo.value) == \\\"Invalid operation\\\"\\n    #Coughed up by CODESOURCERER\"\n    }\n  ]\n}\n```\n"),
			},
		},
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("Handling Imports:\n\nFor Files in Root Directory:\nIf the file being tested (parentpath) is located in the root directory of the repository, imports in the generated test code should directly reference the file as from <filename> import <functions>. Example: from math_utils import add.\nFor Files in Subdirectories:\nIf the file being tested is inside a subdirectory, imports should be relative to the repository structure using appropriate paths. Example: from src.puneeth.math_utils import add.\nEnsure that the import path aligns with the structure provided in the test_directory field, and the parentpath of the file."),
			},
		},
		{
			Role: "model",
			Parts: []genai.Part{
				genai.Text("```json\n{\n  \"tests\": [\n    {\n      \"testname\": \"test_math_utils\",\n      \"testfilepath\": \"src/puneeth/tests/test_math_utils.py\",\n      \"parentpath\": \"src/puneeth/math_utils.py\",\n      \"code\": \"import pytest\\nfrom src.puneeth.math_utils import add, subtract\\n\\ndef test_add():\\n    assert add(2, 3) == 5\\n    assert add(-1, 1) == 0\\n    assert add(0, 0) == 0\\n    assert add(100,200) == 300\\n    assert add(-5,-5) == -10\\n    #Coughed up by CODESOURCERER\\n\\ndef test_subtract():\\n    assert subtract(5, 2) == 3\\n    assert subtract(1, 1) == 0\\n    assert subtract(0, 0) == 0\\n    assert subtract(200,100) == 100\\n    assert subtract(5,-5) == 10\\n    #Coughed up by CODESOURCERER\"\n    },\n    {\n      \"testname\": \"test_calculator\",\n      \"testfilepath\": \"src/puneeth/tests/test_calculator.py\",\n      \"parentpath\": \"src/puneeth/calculator.py\",\n      \"code\": \"import pytest\\nfrom src.puneeth.calculator import calculate\\n\\ndef test_calculate_add():\\n    assert calculate('add', 2, 3) == 5\\n    assert calculate('add', -1, 1) == 0\\n    assert calculate('add', 0, 0) == 0\\n    #Coughed up by CODESOURCERER\\n\\ndef test_calculate_subtract():\\n    assert calculate('subtract', 5, 2) == 3\\n    assert calculate('subtract', 1, 1) == 0\\n    assert calculate('subtract', 0, 0) == 0\\n    #Coughed up by CODESOURCERER\\n\\ndef test_calculate_invalid_operation():\\n    with pytest.raises(ValueError) as excinfo:\\n        calculate('invalid', 2, 3)\\n    assert str(excinfo.value) == \\\"Invalid operation\\\"\\n    #Coughed up by CODESOURCERER\"\n    }\n  ]\n}\n```\n"),
			},
		},
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("{\n  \"merge_id\": \"merge_ef789xyz123abc456opq890rst123uvw456klm789_104\",\n  \"context\": \"This PR introduces functionality to manipulate lists and integrates it into a script for statistical computations.\",\n  \"framework\": \"pytest\",\n  \"test_directory\": \"tests/\",\n  \"comments\": \"off\",\n  \"files\": [\n    {\n      \"path\": \"list_utils.py\",\n      \"content\": \"def get_max(numbers):\\n    if not numbers:\\n        raise ValueError(\\\"Empty list\\\")\\n    return max(numbers)\\n\\ndef get_min(numbers):\\n    if not numbers:\\n        raise ValueError(\\\"Empty list\\\")\\n    return min(numbers)\",\n      \"dependencies\": []\n    },\n    {\n      \"path\": \"statistics/calculate_stats.py\",\n      \"content\": \"from list_utils import get_max, get_min\\n\\ndef calculate_range(numbers):\\n    return get_max(numbers) - get_min(numbers)\\n\\ndef calculate_average(numbers):\\n    if not numbers:\\n        raise ValueError(\\\"Empty list\\\")\\n    return sum(numbers) / len(numbers)\",\n      \"dependencies\": [\n        {\n          \"name\": \"list_utils.py\",\n          \"content\": \"def get_max(numbers):\\n    if not numbers:\\n        raise ValueError(\\\"Empty list\\\")\\n    return max(numbers)\\n\\ndef get_min(numbers):\\n    if not numbers:\\n        raise ValueError(\\\"Empty list\\\")\\n    return min(numbers)\"\n        }\n      ]\n    }\n  ]\n}\n"),
			},
		},
		{
			Role: "model",
			Parts: []genai.Part{
				genai.Text("```json\n{\n  \"tests\": [\n    {\n      \"testname\": \"test_list_utils\",\n      \"testfilepath\": \"tests/test_list_utils.py\",\n      \"parentpath\": \"list_utils.py\",\n      \"code\": \"import pytest\\nfrom list_utils import get_max, get_min\\n\\ndef test_get_max():\\n    assert get_max([1, 5, 2, 8, 3]) == 8\\n    assert get_max([-1, -5, -2, -8, -3]) == -1\\n    with pytest.raises(ValueError):\\n        get_max([])\\n    #Coughed up by CODESOURCERER\\n\\ndef test_get_min():\\n    assert get_min([1, 5, 2, 8, 3]) == 1\\n    assert get_min([-1, -5, -2, -8, -3]) == -8\\n    with pytest.raises(ValueError):\\n        get_min([])\\n    #Coughed up by CODESOURCERER\"\n    },\n    {\n      \"testname\": \"test_calculate_stats\",\n      \"testfilepath\": \"tests/test_calculate_stats.py\",\n      \"parentpath\": \"statistics/calculate_stats.py\",\n      \"code\": \"import pytest\\nfrom statistics.calculate_stats import calculate_range, calculate_average\\n\\ndef test_calculate_range():\\n    assert calculate_range([1, 5, 2, 8, 3]) == 7\\n    assert calculate_range([-1, -5, -2, -8, -3]) == 7\\n    with pytest.raises(ValueError):\\n        calculate_range([])\\n    #Coughed up by CODESOURCERER\\n\\ndef test_calculate_average():\\n    assert calculate_average([1, 5, 2, 8, 3]) == 3.8\\n    assert calculate_average([-1, -5, -2, -8, -3]) == -3.8\\n    with pytest.raises(ValueError):\\n        calculate_average([])\\n    #Coughed up by CODESOURCERER\"\n    }\n  ]\n}\n```\n"),
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error serializing payload: %v\n", err)
		return "Error processing AI request"
	}

	payloadString := string(payloadBytes)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(15*time.Second))
	defer cancel()
	response, err := session.SendMessage(ctx, genai.Text(payloadString))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "Error: Model response timed out."
		}
		log.Printf("Error generating response: %v", err)
		return fmt.Sprintf("Error generating response: %v", err)
	}

	if len(response.Candidates) == 0 {
		fmt.Println("Model did not generate any response.")
		return "failure"
	} else {
		responseBytes, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			fmt.Println("Error marshalling response:", err)
		} else {
			fmt.Println(string(responseBytes))
		}
		var stringResponse string
		if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
			finalResponse := response.Candidates[0].Content.Parts[0]
			fmt.Println(reflect.TypeOf(finalResponse))
			jsonData, _ := json.Marshal(finalResponse)
			stringResponse = string(jsonData)
			fmt.Println(stringResponse)
			fmt.Println(reflect.TypeOf(stringResponse))
			if len(stringResponse) > 12 {
				cleaned := stringResponse[10 : len(stringResponse)-6]

				fmt.Println("Cleaned String")
				var result map[string]interface{}
				err := json.Unmarshal([]byte(cleaned), &result)
				if err != nil {
					fmt.Println("Error parsing JSON")
				} else {
					fmt.Println("Parsed JSON")
				}
			} else {
				fmt.Println("Input string is too short to slice!")
			}
		}
		return stringResponse
	}
}
