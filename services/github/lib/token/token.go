package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type TokenManager struct {
	token       string
	expiration  time.Time
	mu          sync.Mutex // To ensure only one refresh at a time
	refreshing  bool       // To indicate ongoing refresh
	apiEndpoint string     // GitHub API endpoint
	jwtToken    string     // Your app's JWT for authentication
}

// Initialize the TokenManager (Singleton)
var once sync.Once
var instance *TokenManager

func NewTokenManager(apiEndpoint, jwtToken string) *TokenManager {
	once.Do(func() {
		instance = &TokenManager{
			apiEndpoint: apiEndpoint,
			jwtToken:    jwtToken,
		}
		instance.StartProactiveRefresh(1 * time.Minute)
	})
	return instance
}

func (tm *TokenManager) GetToken() (string, error) {
	tm.mu.Lock()

	// Check if the token is still valid
	if time.Now().Before(tm.expiration) {
		defer tm.mu.Unlock()
		return tm.token, nil
	}

	// If already refreshing, wait for it to complete
	if tm.refreshing {
		tm.mu.Unlock()
		time.Sleep(100 * time.Millisecond) // Wait briefly before retrying
		return tm.GetToken()
	}

	// Otherwise, refresh the token
	tm.refreshing = true
	tm.mu.Unlock()

	err := tm.refreshToken()
	if err != nil {
		tm.mu.Lock()
		tm.refreshing = false
		tm.mu.Unlock()
		return "", err
	}

	tm.mu.Lock()
	tm.refreshing = false
	defer tm.mu.Unlock()
	return tm.token, nil
}

func (tm *TokenManager) refreshToken() error {
	// Call GitHub API to generate a new token
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("POST", tm.apiEndpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tm.jwtToken))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return errors.New("failed to refresh token: " + resp.Status)
	}

	var response struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}
	tm.token = response.Token
	tm.expiration = time.Now().Add(10 * time.Minute) // Example: 10 minutes validity

	tm.mu.Lock()
	tm.refreshing = false
	tm.mu.Unlock()

	return nil
}

func (tm *TokenManager) StartProactiveRefresh(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			if time.Now().After(tm.expiration.Add(-2 * time.Minute)) {
				_ = tm.refreshToken()
			}
		}
	}()
}

// Exported function to get the singleton instance
func GetInstance() *TokenManager {
	if instance == nil {
		panic("TokenManager is not initialized. Call Initialize first.")
	}
	return instance
}
