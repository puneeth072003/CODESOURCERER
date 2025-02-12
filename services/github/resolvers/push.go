package resolvers

import (
	"fmt"
	"log"

	pb "github.com/codesourcerer-bot/proto/generated"

	"github.com/codesourcerer-bot/github/lib/gh"
)

func PushNewBranchWithTests(installationToken, owner, repo, baseBranch, newBranch, cacheResult string, tests *pb.GeneratedTestsResponse) error {

	// Get GitHub client
	client, ctx := gh.GetClient(installationToken)

	// Create a new branch
	err := gh.CreateBranch(client, ctx, owner, repo, baseBranch, newBranch)
	if err != nil {
		log.Fatalf("Error creating branch: %v", err)
		return err
	}

	// Add the test files with content
	for _, testFile := range tests.Tests {
		err = gh.CreateFiles(client, ctx, owner, repo, newBranch, testFile.GetTestfilepath(), testFile.GetCode())
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

	prTitle := "chore: tests generated for the code added" // hardcoded for now
	var prBody string

	switch cacheResult {
	case "DISABLED":
		prBody = "### This is a draft PR created from the sandbox branch." // hardcoded for now

	case "DONE":
		prBody = fmt.Sprintf(`
		### This is a draft PR created from the sandbox branch.
		
		> This PR has been cached!
		`)

	case "ERROR":
		prBody = fmt.Sprintf(`
		### This is a draft PR created from the sandbox branch.
		
		> This PR could not be cached!
		`)
	}

	err = gh.CreatePR(client, ctx, owner, repo, prTitle, newBranch, defaultBranch, prBody)
	if err != nil {
		log.Fatalf("Error creating draft PR: %v", err)
		return err
	}

	fmt.Println("Successfully created draft PR from sandbox branch")
	return nil
}
