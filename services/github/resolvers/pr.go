package resolvers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type prEventBody map[string]interface{}

func NewPrBody(b io.ReadCloser) (*prEventBody, error) {
	body, err := io.ReadAll(b)
	if err != nil {
		log.Printf("Unable to unmarshal pull request event: %v", err)
		return nil, fmt.Errorf("Unable to unmarshal pull request event")
	}

	var prEvent prEventBody
	if err := json.Unmarshal(body, &prEvent); err != nil {
		log.Printf("Unable to unmarshal pull request event: %v", err)
		return nil, fmt.Errorf("Unable to process pull request event")
	}

	return &prEvent, nil

}

func (pr *prEventBody) GetPRStatus() (string, bool, string) {
	action := (*pr)["action"].(string)
	merged := (*pr)["pull_request"].(map[string]interface{})["merged"].(bool)
	baseBranch := (*pr)["pull_request"].(map[string]interface{})["base"].(map[string]interface{})["ref"].(string)

	return action, merged, baseBranch
}

// TODO: Get the default branch as well
func (pr *prEventBody) GetRepoInfo() (string, string) {
	repoOwner := (*pr)["repository"].(map[string]interface{})["owner"].(map[string]interface{})["login"].(string)
	repoName := (*pr)["repository"].(map[string]interface{})["name"].(string)

	return repoName, repoOwner
}

func (pr *prEventBody) GetPRInfo() (int, string) {
	pullRequestNumber := int((*pr)["number"].(float64))
	commitSHA := (*pr)["pull_request"].(map[string]interface{})["merge_commit_sha"].(string)

	return pullRequestNumber, commitSHA
}
