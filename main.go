package main

import (
	"log"
	"os/user"
	"time"
	"github.com/mattn/go-zglob"
	"github.com/fsnotify/fsnotify"
)

func findGitRepos(homedir string) []string {
	repos, err := zglob.Glob(homedir + "/**/.git")

    if err != nil {
        log.Fatal(err)
	}
	log.Println("Done")
	return repos
}

func main() {
	usr, err := user.Current()
    if err != nil {
        log.Fatal( err )
	}
	
	homedir := usr.HomeDir

	reposToTrack := make(map[string]bool)
	repoChan := make(chan string)

	ticker := time.NewTicker(5 * time.Minute)
	quit := make(chan struct{})
	go func() {
		// First run.
		for _, repo := range findGitRepos(homedir) {
			repoChan <- repo
		}

		// Loop
		for {
		select {
			case <- ticker.C:
				for _, repo := range findGitRepos(homedir) {
					repoChan <- repo
				}
			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()
	
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	for {
		select {
		case repo := <-repoChan:
			if _, ok := reposToTrack[repo]; !ok {
				err = watcher.Add(repo)
				if err != nil {
					log.Fatal(err)
				} else {
					reposToTrack[repo] = true
				}
			}
		case <-done:
			return
		}
	}
}