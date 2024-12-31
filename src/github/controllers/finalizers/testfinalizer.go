package finalizers

import (
	"fmt"
	"log"
	"net/http"

	"github/controllers/tokenhandlers"
	"github/utils"

	"github.com/gin-gonic/gin"
)

func TestFinalize(c *gin.Context) {
	// Load environment variables
	envs, err := utils.Loadenv(".env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading environment variables"})
		return
	}

	// Access the necessary environment variables
	installationID, exists := envs["INSTALLATION_ID"]
	if !exists {
		log.Printf("INSTALLATION_ID not found in .env file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "INSTALLATION_ID not found in environment variables"})
		return
	}
	apiEndpoint := fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", installationID)

	// Initialize the TokenManager
	jwtToken := tokenhandlers.GetJWT()
	tokenhandlers.NewTokenManager(apiEndpoint, jwtToken)

	// Get the token from the TokenManager
	token, err := tokenhandlers.GetInstance().GetToken()
	if err != nil {
		log.Printf("Error getting token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting token"})
		return
	}

	// Call Finalize with the token and other parameters
	err = Finalize(token, "puneeth072003", "testing-CS", []TestsResponseFormat{
		{
			TestName:     "test_file_operations",
			TestFilePath: "tests/test_file_operations.py",
			ParentPath:   "file_operations.py",
			Code:         `import pytest\nimport os\nfrom file_operations import read_file, write_file\n\ndef test_read_file_valid():\n    test_filename = 'test_file.txt'\n    test_content = 'This is a test content.'\n    with open(test_filename, 'w') as f:\n        f.write(test_content)\n    read_content = read_file(test_filename)\n    assert read_content == test_content\n    os.remove(test_filename)\n\ndef test_read_file_not_found():\n    with pytest.raises(FileNotFoundError):\n        read_file('non_existent_file.txt')\n\ndef test_write_file_valid():\n    test_filename = 'test_write_file.txt'\n    test_content = 'Content to be written.'\n    write_file(test_filename, test_content)\n    with open(test_filename, 'r') as f:\n        written_content = f.read()\n    assert written_content == test_content\n    os.remove(test_filename)\n\n# Coughed up by CODESOURCERER`,
		},
		{
			TestName:     "test_main",
			TestFilePath: "tests/test_main.py",
			ParentPath:   "main.py",
			Code:         `import pytest\nimport os\nfrom main import main\nfrom unittest.mock import patch\n\n@patch('main.process_content', return_value='processed_content')\ndef test_main_integration(mock_process_content):\n    input_file = 'input.txt'\n    output_file = 'output.txt'\n    input_content = 'original content'\n    with open(input_file, 'w') as f:\n        f.write(input_content)\n    main()\n    with open(output_file, 'r') as f:\n        output_content = f.read()\n    assert output_content == 'processed_content'\n    os.remove(input_file)\n    os.remove(output_file)\n\n# Coughed up by CODESOURCERER`,
		},
	})
	if err != nil {
		log.Printf("Error finalizing: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finalizing"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully created draft PR from sandbox branch"})
}
