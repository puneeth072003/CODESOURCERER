package lib

import (
	"fmt"
	"net/http"
)

func configureJsonHeadersWithAuth(req *http.Request, token string) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Accept", "application/vnd.github+json")
}

func configureRawHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/vnd.github.v3.raw")

}

func configureJsonHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/vnd.github.v3+json")
}
