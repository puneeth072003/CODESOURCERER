package resolvers

import (
	"log"

	pb "github.com/codesourcerer-bot/proto/generated"

	"github.com/codesourcerer-bot/github/lib"
)

func PushNewBranchWithTests(owner, repo, baseBranch, newBranch, cacheResult string, tests *pb.GeneratedTestsResponse) error {

	// Get GitHub client
	client, ctx, err := lib.GetClient()
	if err != nil {
		log.Printf("Error creating branch: %v", err)
		return err
	}

	// Create a new branch
	err = lib.CreateBranch(client, ctx, owner, repo, baseBranch, newBranch)
	if err != nil {
		log.Printf("Error creating branch: %v", err)
		return err
	}

	// Add the test files with content
	for _, testFile := range tests.Tests {
		err = lib.CreateFiles(client, ctx, owner, repo, newBranch, testFile.GetTestfilepath(), testFile.GetCode())
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

	prTitle := "chore: tests generated for the code added"            // hardcoded for now
	prBody := "This is a draft PR created from the sandbox branch.\n" // hardcoded for now

	switch cacheResult {

	case "DONE":
		prBody += "This PR has been cached!"

	case "ERROR":
		prBody = "This PR could not be cached!"
	}

	err = lib.CreatePR(client, ctx, owner, repo, prTitle, newBranch, defaultBranch, prBody)
	if err != nil {
		log.Fatalf("Error creating draft PR: %v", err)
		return err
	}

	log.Println("Successfully created draft PR from sandbox branch")
	return nil
}

func CommitRetriedTests(owner, repo, branch string, tests *pb.GeneratedTestsResponse) error {
	client, ctx, err := lib.GetClient()
	if err != nil {
		return err
	}

	for _, testFile := range tests.Tests {
		err := lib.CreateFiles(client, ctx, owner, repo, branch, testFile.GetTestfilepath(), testFile.GetCode())
		if err != nil {
			log.Fatalf("Error creating file %s: %v", testFile.Testfilepath, err)
			return err
		}
	}

	return nil
}
