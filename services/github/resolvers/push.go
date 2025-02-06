package resolvers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	pb "github.com/codesourcerer-bot/proto/generated"

	"github.com/codesourcerer-bot/github/lib/gh"
)

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func PushNewBranchWithTests(installationToken string, owner string, repo string, baseBranch string, tests *pb.GeneratedTestsResponse) error {

	// Step 1: Get GitHub client
	client, ctx := gh.GetClient(installationToken)

	// Step 2: Generate a random branch name
	randomString, err := generateRandomString(5)
	if err != nil {
		log.Fatalf("Error generating random string: %v", err)
		return err
	}
	newBranchName := "tests/CS-sandbox-" + randomString

	// Step 3: Create a new branch
	err = gh.CreateBranch(client, ctx, owner, repo, baseBranch, newBranchName)
	if err != nil {
		log.Fatalf("Error creating branch: %v", err)
		return err
	}

	// Step 4: Add the test files with content
	for _, testFile := range tests.Tests {
		// Directly use the test file content without unquoting
		err = gh.CreateFiles(client, ctx, owner, repo, newBranchName, testFile.Testfilepath, testFile.Code)
		if err != nil {
			log.Fatalf("Error creating file %s: %v", testFile.Testfilepath, err)
			return err
		}
	}

	repoInfo, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return err
	}

	defaultBranch := repoInfo.GetDefaultBranch()

	// Step 5: Create a pull request from the new branch
	prTitle := "chore: tests generated for the code added"          // hardcoded for now
	prBody := "This is a draft PR created from the sandbox branch." // hardcoded for now
	err = gh.CreatePR(client, ctx, owner, repo, prTitle, newBranchName, defaultBranch, prBody)
	if err != nil {
		log.Fatalf("Error creating draft PR: %v", err)
		return err
	}

	fmt.Println("Successfully created draft PR from sandbox branch")
	return nil
}
