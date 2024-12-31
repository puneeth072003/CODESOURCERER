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
	err = Finalize(token, "puneeth072003", "testing-CS", "sample.txt", "hi there")
	if err != nil {
		log.Printf("Error finalizing: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finalizing"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully created draft PR from sandbox branch"})
}
