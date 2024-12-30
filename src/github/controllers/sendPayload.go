package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SendPayload(url string, payload string) (map[string]interface{}, error) {

	// Create the HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		log.Printf("Unable to create HTTP request: %v", err)
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client with a timeout
	client := &http.Client{Timeout: 30 * time.Second} // Adjusted timeout to handle long response times
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Non-200 response from Server 2: %d - %s", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	// Read and parse the response body
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		log.Printf("Error decoding JSON response: %v", err)
		return nil, err
	}

	return responseBody, nil
}

func TestSendPayload(c *gin.Context) {
	samplePayload := `{
						"merge_id": "merge_uvw456rst789xyz123abc890klm567def234_107",
						"context": "This PR adds utility functions for date formatting and integrates these into a scheduling module.",
						"framework": "pytest",
						"test_directory": "tests/",
						"comments": "off",
						"files": [
							{
							"path": "date_utils.py",
							"content": "from datetime import datetime\n\ndef format_date(date):\n    return date.strftime('%Y-%m-%d')\n\ndef parse_date(date_string):\n    return datetime.strptime(date_string, '%Y-%m-%d')",
							"dependencies": []
							},
							{
							"path": "scheduling/schedule_manager.py",
							"content": "from date_utils import format_date, parse_date\n\ndef get_formatted_date_for_today():\n    return format_date(datetime.now())",
							"dependencies": [
								{
								"name": "date_utils.py",
								"content": "from datetime import datetime\n\ndef format_date(date):\n    return date.strftime('%Y-%m-%d')\n\ndef parse_date(date_string):\n    return datetime.strptime(date_string, '%Y-%m-%d')"
								}
							]
							}
						]
					}`

	url := "http://localhost:3001/process"

	response, err := SendPayload(url, samplePayload)
	if err != nil {
		log.Printf("Error in SendPayload: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Response from server: %v", response)
	c.JSON(http.StatusOK, gin.H{"response": response})
}
