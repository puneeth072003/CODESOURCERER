package resolvers

import (
	"fmt"
	"log"

	"github.com/codesourcerer-bot/github/connections"
	pb "github.com/codesourcerer-bot/proto/generated"
)

func CachePullRequest(shouldCache bool, repoOwner, repoName, newBranch string, context []*pb.SourceFilePayload, tests []*pb.TestFilePayload) string {

	var cacheResult string
	if shouldCache {
		cacheKey := fmt.Sprintf("%s/%s/tree/%s", repoOwner, repoName, newBranch)
		ok, err := connections.SetContextAndTestsToDatabase(cacheKey, context, tests)
		if err != nil || !ok {
			log.Printf("unable to cache contexts and tests: %v", err)
			cacheResult = "ERROR"
		} else {
			cacheResult = "DONE"
		}
	} else {
		cacheResult = "DISABLED"
	}

	return cacheResult
}
