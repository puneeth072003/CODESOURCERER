package partials

import (
	"log"
	"net/http"
	"os"

	"github.com/codesourcerer-bot/github/resolvers"
	"github.com/codesourcerer-bot/github/utils"
	"github.com/gin-gonic/gin"

	pb "github.com/codesourcerer-bot/proto/generated"
)

func TestFinalize(c *gin.Context) {

	// Access the necessary environment variables
	installationID := os.Getenv("INSTALLATION_ID")
	if installationID == "" {
		log.Printf("INSTALLATION_ID not found in .env file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "INSTALLATION_ID not found in environment variables"})
		return
	}

	// Create test files in the format expected by pb.GeneratedTestsResponse
	testFiles := []*pb.TestFilePayload{
		{
			Testname:     "test_file_operations",
			Testfilepath: "tests/test_file_operations.py",
			Parentpath:   "file_operations.py",
			Code:         `import pytest\nimport os\nfrom file_operations import read_file, write_file\n\ndef test_read_file_valid():\n    test_filename = 'test_file.txt'\n    test_content = 'This is a test content.'\n    with open(test_filename, 'w') as f:\n        f.write(test_content)\n    read_content = read_file(test_filename)\n    assert read_content == test_content\n    os.remove(test_filename)\n\ndef test_read_file_not_found():\n    with pytest.raises(FileNotFoundError):\n        read_file('non_existent_file.txt')\n\ndef test_write_file_valid():\n    test_filename = 'test_write_file.txt'\n    test_content = 'Content to be written.'\n    write_file(test_filename, test_content)\n    with open(test_filename, 'r') as f:\n        written_content = f.read()\n    assert written_content == test_content\n    os.remove(test_filename)\n\n# Coughed up by CODESOURCERER`,
		},
		{
			Testname:     "test_main",
			Testfilepath: "tests/test_main.py",
			Parentpath:   "main.py",
			Code:         `import pytest\nimport os\nfrom main import main\nfrom unittest.mock import patch\n\n@patch('main.process_content', return_value='processed_content')\ndef test_main_integration(mock_process_content):\n    input_file = 'input.txt'\n    output_file = 'output.txt'\n    input_content = 'original content'\n    with open(input_file, 'w') as f:\n        f.write(input_content)\n    main()\n    with open(output_file, 'r') as f:\n        output_content = f.read()\n    assert output_content == 'processed_content'\n    os.remove(input_file)\n    os.remove(output_file)\n\n# Coughed up by CODESOURCERER`,
		},
	}

	// Create the GeneratedTestsResponse
	generatedTestsResponse := &pb.GeneratedTestsResponse{
		Tests: testFiles,
	}

	newBranch := utils.GetRandomBranch()

	// Call Finalize with the token and other parameters
	err := resolvers.PushNewBranchWithTests("puneeth072003", "testing-CS", "testing", newBranch, "DISABLED", generatedTestsResponse)
	if err != nil {
		log.Printf("Error finalizing: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finalizing"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Test files generated and draft PR created successfully",
	})
}
