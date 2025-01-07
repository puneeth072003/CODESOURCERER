package finalizers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"

	pb "protobuf/generated"
)

type TestsResponseFormat struct {
	TestName     string `json:"testname"`
	TestFilePath string `json:"testfilepath"`
	ParentPath   string `json:"parentpath"`
	Code         string `json:"code"`
}

func Finalize(installationToken string, owner string, repo string, testspointer *pb.GeneratedTestsResponse) error {
	// Step 1: Convert the tests to a slice of TestsResponseFormat
	var testFiles []TestsResponseFormat
	for _, test := range testspointer.Tests {
		testFiles = append(testFiles, TestsResponseFormat{
			TestName:     test.Testname,
			TestFilePath: test.Testfilepath,
			ParentPath:   test.Parentpath,
			Code:         test.Code,
		})
	}

	// Step 2: Get GitHub client
	client, ctx := GetClient(installationToken)

	// Step 3: Generate a random branch name
	randomString, err := generateRandomString(5)
	if err != nil {
		log.Fatalf("Error generating random string: %v", err)
		return err
	}
	newBranchName := "CS-sandbox-" + randomString

	// Step 4: Create a new branch
	err = CreateBranch(client, ctx, owner, repo, newBranchName)
	if err != nil {
		log.Fatalf("Error creating branch: %v", err)
		return err
	}

	// Step 5: Add the test files with content
	for _, testFile := range testFiles {
		// Decode the test file content (if necessary)
		actualString, err := strconv.Unquote(`"` + testFile.Code + `"`)
		if err != nil {
			log.Fatalf("Error decoding string: %v", err)
			return err
		}

		// Create the file in the repository
		err = CreateFiles(client, ctx, owner, repo, newBranchName, testFile.TestFilePath, actualString)
		if err != nil {
			log.Fatalf("Error creating file %s: %v", testFile.TestFilePath, err)
			return err
		}
	}

	// Step 6: Create a pull request from the new branch
	prTitle := "chore: tests generated for the code added"          // hardcoded for now
	baseBranch := "main"                                            // hardcoded for now
	prBody := "This is a draft PR created from the sandbox branch." // hardcoded for now
	err = CreatePR(client, ctx, owner, repo, prTitle, newBranchName, baseBranch, prBody)
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
