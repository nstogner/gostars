package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
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
		flag.CommandLine.SetOutput(os.Stdout)
		printUsage()
		os.Exit(0)
	}

	pkg := "."
	if len(flag.Args()) == 2 {
		pkg = flag.Args()[1]
	}

	// Find all package imports.
	imap := make(importMap)
	if err := imap.find(pkg); err != nil {
		log.Fatal(err)
	}

	// Filter on github packages.
	var paths []string
	for p := range imap {
		if isGithubPath(p) {
			paths = append(paths, p)
		}
	}

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

		if flags.threshold > -1 {
			if stars < flags.threshold {
				res = append(res, result{
					Path:  p,
					Stars: stars,
				})
			}
		} else {
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
			log.Fatal(errors.Wrap(err, "json"))
		}
	} else {
		for _, r := range res {
			fmt.Printf("%v\t%v\n", r.Stars, r.Path)
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
	flag.Usage()
}
