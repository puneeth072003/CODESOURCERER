package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func SendPayload(url string, payload interface{}) (string, error) {
	// Marshal the payload into JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Unable to marshal payload: %v", err)
		return "", err
	}

	// Create the HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Unable to create HTTP request: %v", err)
		return "", err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client with a timeout
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Non-200 response from Server 2: %d - %s", resp.StatusCode, string(bodyBytes))
		return "", fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	// Read and return the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return "", err
	}

	return string(responseBody), nil
}
