package validators

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type workflowEventBody map[string]interface{}

func NewWorkflowBody(b io.ReadCloser) (*workflowEventBody, error) {
	body, err := io.ReadAll(b)
	if err != nil {
		log.Printf("Unable to unmarshal pull request event: %v", err)
		return nil, fmt.Errorf("Unable to unmarshal pull request event")
	}

	var prEvent workflowEventBody
	if err := json.Unmarshal(body, &prEvent); err != nil {
		log.Printf("Unable to unmarshal pull request event: %v", err)
		return nil, fmt.Errorf("Unable to process pull request event")
	}

	return &prEvent, nil

}

func (wf *workflowEventBody) GetWorkflowDetails() (string, string, string) {
	status := (*wf)["action"].(string)
	name := (*wf)["workflow_run"].(map[string]interface{})["name"].(string)
	result := (*wf)["workflow_run"].(map[string]interface{})["conclusion"]
	if result != nil {
		return name, status, result.(string)
	}

	return name, status, ""
}

func (wf *workflowEventBody) GetWorkflowJobUrl() string {
	url := (*wf)["workflow_run"].(map[string]interface{})["jobs_url"].(string)

	return url
}

func (wf *workflowEventBody) GetRepoDetails() (string, string, string) {
	owner := (*wf)["workflow_run"].(map[string]interface{})["repository"].(map[string]interface{})["owner"].(map[string]interface{})["login"].(string)
	repoName := (*wf)["workflow_run"].(map[string]interface{})["repository"].(map[string]interface{})["name"].(string)
	branchName := (*wf)["workflow_run"].(map[string]interface{})["head_branch"].(string)

	return owner, repoName, branchName
}
