package main

import (
	"log"
	"os/user"
	"time"
	"github.com/mattn/go-zglob"
	"github.com/fsnotify/fsnotify"
)

func findGitRepos(homedir string) []string {
	repos, err := zglob.Glob(homedir + "/**/.git/index")

	// Filter non git dirs
	// ie: Fatal: Not a git repository
    if err != nil {
        log.Fatal(err)
	}

	log.Println("Done")
	return repos
}

func main() {
	log.Println("indexing git repos")

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
		sendPathsToWatcher := func() {
			count := 0

			for _, f := range findGitRepos(homedir) {
				count += 1
				repoChan <- f
			}

			log.Println("")
		}

		// First run.
		sendPathsToWatcher()

		// Loop
		for {
		select {
			case <- ticker.C:
				sendPathsToWatcher()
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