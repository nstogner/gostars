package main

import (
	"go/build"

	"github.com/pkg/errors"
)

type importMap map[string]struct{}

func (i importMap) find(path string) error {
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
		if err := i.find(p); err != nil {
			return errors.Wrap(err, "finding")
		}
	}
	return nil
}
