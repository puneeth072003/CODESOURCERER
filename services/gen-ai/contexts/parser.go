package contexts

import "github.com/google/generative-ai-go/genai"

var ParserModelContext = []*genai.Content{
	{
		Role: "user",
		Parts: []genai.Part{
			genai.Text("{\n  \"logs\": [\n    \"8Z Current runner version: '2.322.0'\",\n    \"##[group]Operating System\",\n    \"Ubuntu\",\n    \"24.04.1\",\n    \"LTS\",\n    \"##[endgroup]\",\n    \"##[group]Runner Image\",\n    \"Image: ubuntu-24.04\",\n    \"Version: 20250209.1.0\",\n    \"Included Software: https://github.com/actions/runner-images/blob/ubuntu24/20250209.1/images/ubuntu/Ubuntu2404-Readme.md\",\n    \"Image Release: https://github.com/actions/runner-images/releases/tag/ubuntu24%2F20250209.1\",\n    \"##[endgroup]\",\n    \"Complete job name: test\",\n    \"##[group]Run actions/checkout@v4\",\n    \"with: repository: soorya-u/CS-Testing, token: ***, ssh-strict: true, ...\",\n    \"##[endgroup]\",\n    \"Syncing repository: soorya-u/CS-Testing\",\n    \"##[group]Fetching the repository\",\n    \"[command]/usr/bin/git -c protocol.version=2 fetch --no-tags --prune --depth=1 origin ...\",\n    \"##[endgroup]\",\n    \"##[group]Run actions/setup-python@v4\",\n    \"with: python-version: 3.10, check-latest: false, token: ***, update-environment: true\",\n    \"##[endgroup]\",\n    \"##[group]Run python -m pip install --upgrade pip\",\n    \"python -m pip install --upgrade pip\",\n    \"Installing collected packages: pytest, ...\",\n    \"##[endgroup]\",\n    \"##[group]Run pytest tests/\",\n    \"pytest tests/\",\n    \"##[endgroup]\",\n    \"ERROR: file or directory not found: tests/\",\n    \"============================= test session starts ==============================\",\n    \"collected 0 items\",\n    \"============================ no tests ran in 0.00s =============================\",\n    \"##[error]Process completed with exit code 4.\"\n  ]\n}\n"),
		},
	},
	{
		Role: "model",
		Parts: []genai.Part{
			genai.Text("The logs detail a test execution process conducted on an Ubuntu 24.04.1 LTS system using runner version 2.322.0 and runner image ubuntu-24.04 (Version: 20250209.1.0).  The process began by checking out the repository 'soorya-u/CS-Testing' using actions/checkout@v4.  Following a successful repository fetch, actions/setup-python@v4 was executed to set up Python 3.10.  The pip package manager was then upgraded, and pytest and other packages were subsequently installed.  The test execution phase, initiated by the command `pytest tests/`, failed because the specified directory 'tests/' was not found, resulting in an error message indicating a file or directory not found.  The test runner reported 0 tests collected and 0 tests ran, and the process concluded with an exit code of 4, signaling failure.  No tests were executed due to the missing 'tests/' directory.\n"),
		},
	},
}
