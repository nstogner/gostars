package main

import (
	"go/build"
	"sort"

	"github.com/pkg/errors"
)

type importMap map[string]struct{}

func (i importMap) populate(path string) error {
	ctx := build.Default
	pkg, err := ctx.Import(path, ".", 0)
	if err != nil {
		return errors.Wrap(err, "importing dir")
	}
	if pkg.Goroot {
		return nil
	}

	for _, p := range pkg.Imports {
		i[p] = struct{}{}
		if err := i.populate(p); err != nil {
			return errors.Wrap(err, "finding")
		}
	}
	return nil
}

func filterAndOrder(m map[string]struct{}, filter func(string) bool) []string {
	// Filter on github packages.
	var paths []string
	for p := range m {
		if filter(p) {
			paths = append(paths, p)
		}
	}
	sort.Strings(paths)
	return paths
}
