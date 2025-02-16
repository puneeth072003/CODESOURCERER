package contexts

import (
	pb "github.com/codesourcerer-bot/proto/generated"
	"github.com/google/generative-ai-go/genai"
)

func GetRegeratorContext(cache *pb.CachedContents) []*genai.Content {
	return []*genai.Content{
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("{\n  \"merge_id\": \"merge_1234\",\n  \"commit_sha\": \"abc123def456\",\n  \"pull_request\": 42,\n  \"context\": \"This PR implements factorial and combination functions and prints the combination result.\",\n  \"framework\": \"pytest\",\n  \"contexts\": [\n    {\n      \"path\": \"q1.py\",\n      \"content\": \"def factorial(n):\\n    if n == 0:\\n        return 1\\n    else:\\n        return n * factorial(n - 1)\"\n    },\n    {\n      \"path\": \"q2.py\",\n      \"content\": \"from q1 import factorial\\n\\ndef combinations(n, r):\\n    # Using float division to avoid integer division issues\\n    return factorial(n) / (factorial(r) * factorial(n - r))\",\n      \"dependencies\": [\n        {\n          \"name\": \"q1.py\",\n          \"content\": \"def factorial(n):\\n    if n == 0:\\n        return 1\\n    else:\\n        return n * factorial(n - 1)\"\n        }\n      ]\n    },\n    {\n      \"path\": \"q3.py\",\n      \"content\": \"from q2 import combinations\\n\\nn = 5\\nr = 2\\nresult = combinations(n, r)\\nprint(f\\\"Combinations of {n} items taken {r} at a time: {result}\\\")\",\n      \"dependencies\": [\n        {\n          \"name\": \"q2.py\",\n          \"content\": \"from q1 import factorial\\n\\ndef combinations(n, r):\\n    return factorial(n) / (factorial(r) * factorial(n - r))\"\n        }\n      ]\n    }\n  ],\n  \"tests\": [\n    {\n      \"testname\": \"test_q1\",\n      \"testfilepath\": \"tests/test_q1.py\",\n      \"parentpath\": \"q1.py\",\n      \"code\": \"import pytest\\nfrom q1 import factorial\\n\\ndef test_factorial_positive():\\n    assert factorial(5) == 120\\n\\ndef test_factorial_zero():\\n    assert factorial(0) == 1\\n\\ndef test_factorial_one():\\n    assert factorial(1) == 1\"\n    },\n    {\n      \"testname\": \"test_q2\",\n      \"testfilepath\": \"tests/test_q2.py\",\n      \"parentpath\": \"q2.py\",\n      \"code\": \"import pytest\\nfrom q2 import combinations\\n\\ndef test_combinations_valid_input():\\n    # Expected: 5C2 = 10.0\\n    assert combinations(5, 2) == 10.0\\n\\ndef test_combinations_edge_cases():\\n    assert combinations(0, 0) == 1.0\\n    assert combinations(5, 0) == 1.0\\n    assert combinations(5, 5) == 1.0\"\n    },\n    {\n      \"testname\": \"test_q3\",\n      \"testfilepath\": \"tests/test_q3.py\",\n      \"parentpath\": \"q3.py\",\n      \"code\": \"import pytest\\nfrom io import StringIO\\nimport sys\\n\\n\\ndef test_q3_output_correctness(capsys):\\n    from q3 import n, r, result\\n    old_stdout = sys.stdout\\n    sys.stdout = captured_output = StringIO()\\n    print(f\\\"Combinations of {n} items taken {r} at a time: {result}\\\")\\n    sys.stdout = old_stdout\\n    output = captured_output.getvalue().strip()\\n    expected_output = f\\\"Combinations of {n} items taken {r} at a time: {result}\\\"\\n    assert output == expected_output\"\n    }\n  ],\n  \"error\": \"Error Summary: The tests for q2 were failing due to using integer division instead of float division, and the test for q3 failed because stdout capture did not match the expected output format. Please adjust the tests to address these issues.\"\n}\n"),
			},
		},
		{
			Role: "model",
			Parts: []genai.Part{
				genai.Text("```json\n{\n  \"tests\": [\n    {\n      \"testname\": \"test_q1\",\n      \"path\": \"q1.py\",\n      \"tests\": [\n        {\n          \"testname\": \"test_factorial_positive\",\n          \"path\": \"q1.py\",\n          \"code\": \"import pytest\\nfrom q1 import factorial\\n\\ndef test_factorial_positive():\\n    assert factorial(5) == 120\"\n        },\n        {\n          \"testname\": \"test_factorial_zero\",\n          \"path\": \"q1.py\",\n          \"code\": \"import pytest\\nfrom q1 import factorial\\n\\ndef test_factorial_zero():\\n    assert factorial(0) == 1\"\n        },\n        {\n          \"testname\": \"test_factorial_one\",\n          \"path\": \"q1.py\",\n          \"code\": \"import pytest\\nfrom q1 import factorial\\n\\ndef test_factorial_one():\\n    assert factorial(1) == 1\"\n        },\n        {\n          \"testname\": \"test_factorial_negative\",\n          \"path\": \"q1.py\",\n          \"code\": \"import pytest\\nfrom q1 import factorial\\n\\ndef test_factorial_negative():\\n    with pytest.raises(RecursionError):\\n        factorial(-1)\"\n        }\n      ]\n    },\n    {\n      \"testname\": \"test_q2\",\n      \"path\": \"q2.py\",\n      \"tests\": [\n        {\n          \"testname\": \"test_combinations_valid_input\",\n          \"path\": \"q2.py\",\n          \"code\": \"import pytest\\nfrom q2 import combinations\\n\\ndef test_combinations_valid_input():\\n    assert combinations(5, 2) == 10.0\"\n        },\n        {\n          \"testname\": \"test_combinations_edge_cases\",\n          \"path\": \"q2.py\",\n          \"code\": \"import pytest\\nfrom q2 import combinations\\n\\ndef test_combinations_edge_cases():\\n    assert combinations(0, 0) == 1.0\\n    assert combinations(5, 0) == 1.0\\n    assert combinations(5, 5) == 1.0\"\n        },\n        {\n          \"testname\": \"test_combinations_invalid_input\",\n          \"path\": \"q2.py\",\n          \"code\": \"import pytest\\nfrom q2 import combinations\\n\\ndef test_combinations_invalid_input():\\n    with pytest.raises(ValueError):\\n        combinations(5, 6)\"\n        }\n      ]\n    },\n    {\n      \"testname\": \"test_q3\",\n      \"path\": \"q3.py\",\n      \"tests\": [\n        {\n          \"testname\": \"test_q3_output_correctness\",\n          \"path\": \"q3.py\",\n          \"code\": \"import pytest\\nfrom io import StringIO\\nimport sys\\n\\ndef test_q3_output_correctness(capsys):\\n    import q3\\n    captured = capsys.readouterr()\\n    expected_output = \\\"Combinations of 5 items taken 2 at a time: 10.0\\\\n\\\"\\n    assert captured.out == expected_output\"\n        }\n      ]\n    }\n  ]\n}\n```\n"),
			},
		},
	}
}
