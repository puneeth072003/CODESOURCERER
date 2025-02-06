package resolvers

import (
	"log"
	"sync"

	"github.com/codesourcerer-bot/github/lib/gh"
	"github.com/codesourcerer-bot/github/utils"

	pb "github.com/codesourcerer-bot/proto/generated"
)

func GetFileContents(fileContents []map[string]interface{}, repoOwner, repoName, commitSHA string) <-chan *pb.SourceFilePayload {
	outChan := make(chan *pb.SourceFilePayload)

	go func() {
		for _, f := range fileContents {
			filePath := f["filename"].(string)

			fileContent, err := gh.FetchFileFromGitHub(repoOwner, repoName, commitSHA, filePath)
			if err != nil {
				log.Printf("Unable to fetch file content for %s: %v", filePath, err)
				fileContent = "Error fetching content"
			} else {
				log.Printf("Successfully fetched content for file: %s", filePath)
			}

			outChan <- &pb.SourceFilePayload{
				Path:    filePath,
				Content: fileContent,
			}
		}
		close(outChan)
	}()

	return outChan
}

func GetDependencyContents(fileChan <-chan *pb.SourceFilePayload, dependencies map[string][]string, repoOwner, repoName, commitSHA string) <-chan *pb.SourceFilePayload {
	outChan := make(chan *pb.SourceFilePayload)

	go func() {
		for f := range fileChan {
			fileDependencies := utils.FilterDependenciesForFile(f.Path, dependencies)
			var wg sync.WaitGroup
			depChan := make(chan *pb.SourceFileDependencyPayload, len(fileDependencies))

			for _, dep := range fileDependencies {
				wg.Add(1)

				go func(channel chan<- *pb.SourceFileDependencyPayload, dep string) {
					defer wg.Done()

					depContent, err := gh.FetchFileFromGitHub(repoOwner, repoName, commitSHA, dep)
					if err != nil {
						log.Printf("Unable to fetch content for dependency %s: %v", dep, err)
						depContent = "Error fetching content"
					} else {
						log.Printf("Successfully fetched content for dependency: %s", dep)
					}

					channel <- &pb.SourceFileDependencyPayload{
						Name:    dep,
						Content: depContent,
					}
				}(depChan, dep)
			}

			wg.Wait()
			close(depChan)

			var deps []*pb.SourceFileDependencyPayload

			for d := range depChan {
				deps = append(deps, d)
			}

			f.Dependencies = deps
			outChan <- f
		}
		close(outChan)
	}()

	return outChan
}
