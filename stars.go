package main

import (
	"time"

	"github.com/pkg/errors"
)

type importStars struct {
	Path  string `json:"path"`
	Stars int    `json:"stars"`
}

func fetchAndFilterStars(paths []string, starFunc func(importPath string) (stars int, err error), threshold int) ([]importStars, error) {
	res := make([]importStars, 0)

	for _, p := range paths {
		stars, err := getGithubStars(p)
		if err != nil {
			if err == errNotFound {
				// Allow for non-published/private repos.
				continue
			}
			return nil, errors.Wrap(err, "getting stars")
		}

		if threshold == -1 || stars > threshold {
			res = append(res, importStars{
				Path:  p,
				Stars: stars,
			})
		}

		// Dont get ratelimited.
		time.Sleep(time.Second)
	}

	return res, nil
}
