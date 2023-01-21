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
	var commit string
	flag.StringVar(&commit, "commit", "", "commit hash")
	flag.StringVar(&commit, "c", "", "commit hash")
	flag.Parse()

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
