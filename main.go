package main

import (
	"fmt"
	"log"
	"os/user"
	"time"
	"path/filepath"
	"os"
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

func getDirsTree(dirPath string) (dirs []string) {
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})

	return dirs
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
		sendPathsToWatcher := func() {
			count := 0

			for _, repo := range findGitRepos(homedir) {
				dir := filepath.Dir(repo)

				// dirs, _ := zglob.Glob(dir + "/**/")
				dirs := getDirsTree(dir)

				count += len(dirs)

				repoChan <- dir
				for _, subdir := range dirs {
					fmt.Println(count, subdir)
					repoChan <- subdir
				}
			}

			fmt.Println(count)
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