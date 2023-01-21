package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const DEFAULT_BASE_BRANCH_NAME = "main"

func main() {
	var help bool
	flag.BoolVar(&help, "h", false, "help message")
	flag.BoolVar(&help, "help", false, "help message")
	flag.Parse()

	if help {
		fmt.Println(`Print GitHub Pull Request ID from commit hash.

USAGE
    git-get-pr-id <commit hash>

EXAMPLES
    $ gh-get-pr-id 123abc`)
		os.Exit(0)
	}

	if len(os.Args) != 2 {
		log.Fatal("error require commit has. e.g) `git-get-pr-id <commit hash>`")
	}
	commit := os.Args[1]

	var base string
	if os.Getenv("BASE_BRANCH") != "" {
		base = os.Getenv("BASE_BRANCH")
	} else {
		base = DEFAULT_BASE_BRANCH_NAME
	}

	ap, err := exec.Command("git", "rev-list", "--ancestry-path", fmt.Sprintf("%s..%s", commit, base)).Output()
	if err != nil {
		log.Fatalf("error exec git rev-list --ancestry-path %s..%s\nerr: %s", commit, base, err)
	}

	fp, err := exec.Command("git", "rev-list", "--first-parent", fmt.Sprintf("%s..%s", commit, base)).Output()
	if err != nil {
		log.Fatalf("error exec git rev-list --first-parent %s..%s\nerr: %s", commit, base, err)
	}

	sameCommits := []string{}
	for _, apc := range strings.Split(string(ap), "\n") {
		if apc == "" {
			continue
		}
		for _, fpc := range strings.Split(string(fp), "\n") {
			if fpc == "" {
				continue
			}

			if apc == fpc {
				sameCommits = append(sameCommits, fpc)
			}
		}
	}

	mergeCommit := sameCommits[len(sameCommits)-1]

	mlog, err := exec.Command("git", "log", "-1", "--format=%B", mergeCommit).Output()
	if err != nil {
		log.Fatalf("error exec: git log -1 --format=%%B %s\nerr: %s", mergeCommit, err)
	}

	msg := strings.Split(strings.TrimSpace(string(mlog)), " ")
	prID := strings.Trim(msg[3], "#")

	fmt.Println(prID)
}
