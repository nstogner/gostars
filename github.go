package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

var errNotFound = errors.New("not found")

func isGithubPath(path string) bool {
	return strings.HasPrefix(path, "github.com")
}

func getGithubStars(path string) (int, error) {
	split := strings.Split(path, "/")
	if n := len(split); n < 3 {
		return 0, fmt.Errorf("expected github import path to be at least 3 parts, got %v", n)
	}

	user, repo := split[1], split[2]
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", user, repo)
	resp, err := http.Get(url)
	if err != nil {
		return 0, errors.Wrap(err, "get request")
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		if resp.StatusCode == http.StatusNotFound {
			return 0, errNotFound
		}
		bdy, _ := ioutil.ReadAll(resp.Body)
		return 0, fmt.Errorf("expected 2XX status code, got %v: %s", resp.StatusCode, string(bdy))
	}

	var body struct {
		Stars int `json:"stargazers_count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return 0, errors.Wrap(err, "decoding github response to json")
	}

	return body.Stars, nil
}
