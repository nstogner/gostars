package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/pkg/errors"
)

func main() {
	// Parse inputs.
	var flags struct {
		help      bool
		threshold int
		json      bool
	}
	flag.BoolVar(&flags.help, "help", false, "show help")
	flag.IntVar(&flags.threshold, "threshold", -1, "report any projects with fewer stars than threshold, exits with status 1 if any exist")
	flag.BoolVar(&flags.json, "json", false, "output as json")
	flag.Parse()

	if flags.help {
		printUsage()
		os.Exit(0)
	}

	pkg := "."
	if args := flag.Args(); len(args) == 1 {
		pkg = args[0]
	}

	// Find all package imports.
	imap := make(importMap)
	if err := imap.populate(pkg); err != nil {
		log.Fatal(errors.Wrap(err, "finding imports"))
	}

	// Filter on github packages.
	var paths []string
	for p := range imap {
		if isGithubPath(p) {
			paths = append(paths, p)
		}
	}
	sort.Strings(paths)

	// Find and filter on stars.
	res := make([]result, 0)
	for _, p := range paths {
		stars, err := getGithubStars(p)
		if err != nil {
			if err == errNotFound {
				// Allow for non-published/private repos.
				continue
			}
			log.Fatal(errors.Wrap(err, "getting github stars"))
		}

		if flags.threshold == -1 || stars > flags.threshold {
			res = append(res, result{
				Path:  p,
				Stars: stars,
			})
		}
		time.Sleep(time.Second)
	}

	// Write output.
	if flags.json {
		if err := json.NewEncoder(os.Stdout).Encode(res); err != nil {
			log.Fatal(errors.Wrap(err, "outputting json"))
		}
	} else {
		for _, r := range res {
			if _, err := fmt.Printf("%v\t%v\n", r.Stars, r.Path); err != nil {
				// Utoh, the world is probably over as well.
				log.Fatal(errors.Wrap(err, "printing to stdout"))
			}
		}
	}

	// Set exit status.
	if flags.threshold > -1 && len(res) > 0 {
		os.Exit(1)
	}
}

type result struct {
	Path  string `json:"path"`
	Stars int    `json:"stars"`
}

func printUsage() {
	fmt.Println("gostars [options] [<package>]\n")
	flag.CommandLine.SetOutput(os.Stdout)
	flag.Usage()
}
