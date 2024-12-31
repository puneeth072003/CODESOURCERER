package finalizers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
)

func Finalize(installationToken string, owner string, repo string, filePath string, fileContent string) error {

	// Get GitHub client
	client, ctx := GetClient(installationToken)

	// Generate a random branch name
	randomString, err := generateRandomString(5)
	if err != nil {
		log.Fatalf("Error generating random string: %v", err)
		return err
	}
	newBranchName := "CS-sandbox-" + randomString

	// Create a new branch
	err = CreateBranch(client, ctx, owner, repo, newBranchName)
	if err != nil {
		log.Fatalf("Error creating branch: %v", err)
		return err
	}

	// Add a sample file with content
	// filePath := "sample.txt"
	// fileContent := "This is a sample file content."
	err = CreateFiles(client, ctx, owner, repo, newBranchName, filePath, fileContent)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
		return err
	}

	// Draft a pull request from the new branch
	prTitle := "chore: tests generated for the code added"          // hardcoded for now
	baseBranch := "testing"                                         // hardcoded for now
	prBody := "This is a draft PR created from the sandbox branch." // hardcoded for now
	err = CreateDraftPR(client, ctx, owner, repo, prTitle, newBranchName, baseBranch, prBody)
	if err != nil {
		log.Fatalf("Error creating draft PR: %v", err)
		return err
	}

	fmt.Println("Successfully created draft PR from sandbox branch")
	return nil
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
