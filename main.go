package main

import (
	// "os"
	"bytes"
	"fmt"
	"log"
	"os/user"
	"time"
	"os/exec"
	"github.com/mattn/go-zglob"
	"path/filepath"
)


func findGitRepos(homedir string) []string {
	repos, err := zglob.Glob(homedir + "/**/.git")
	for i, repo := range repos {
		repos[i] = filepath.Dir(repo)
	}

    if err != nil {
        log.Fatal(err)
	}

	log.Println("Done")
	return repos
}

func checkGitRepoUpdated(dir string) (output string, err error) {
	cmd := exec.Command("git", "ls-files", "--modified", "--debug")
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return stderr.String(), err
	}
	output = string(out)
	return output, nil
}

type repoUpdateCheck struct {
	repo string
	output string
	err error
}

func main() {
	usr, err := user.Current()
    if err != nil {
        log.Fatal( err )
	}
	
	homedir := usr.HomeDir

	// refresh repos every 5 minutes
	// start routine which checks them every 10s
	// send done msg when we finish refreshing the next time

	findGitReposTick := time.NewTicker(15 * time.Minute)
	checkGitRepoUpdatedTick := time.NewTicker(1 * time.Minute)

	gitReposToTrack := make(chan []string)

	go func() {
		var repos []string
		repoLastCheckup := make(map[string]string)
		repoUpdateCheckChan := make(chan repoUpdateCheck) 
		ignore := make(map[string]bool)

		for {
			select {
			case <-checkGitRepoUpdatedTick.C:
				for _, repo := range repos {
					lastOutput, _ := repoLastCheckup[repo]

					if ignore[repo] == true {
						continue
					}

					go func(repo, lastOutput string) {
						output, err := checkGitRepoUpdated(repo)
						check := repoUpdateCheck{
							repo, 
							output,
							err,
						}
						repoUpdateCheckChan <- check
					}(repo, lastOutput)
				}

				break
			
			case check := <-repoUpdateCheckChan:
				if check.err != nil {
					fmt.Println("ERROR", check.err, check.output)
					ignore[check.repo] = true
					continue
					// delete(check.repo, repos)
				}

				lastOutput, ok := repoLastCheckup[check.repo]
				if !ok {
					// First time checking.      
					repoLastCheckup[check.repo] = check.output
					continue
				}

				changed := (lastOutput != check.output)
				if changed {
					fmt.Println("UPDATED", check.repo)
				}  
				repoLastCheckup[check.repo] = check.output
				break;

			case repos = <-gitReposToTrack:
				break
			}

		}
	}()

	refresh := func() {
		fmt.Println("EVENT Refreshing")
		repos := findGitRepos(homedir)
		fmt.Println("EVENT Refreshed")
		gitReposToTrack <- repos
	}

	// First run.
	refresh()

	for {
		select {
		case <-findGitReposTick.C:
			refresh()
		}
	}
}