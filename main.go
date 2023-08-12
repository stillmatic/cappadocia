package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: watch <glob_pattern> <command> [args...]")
		return
	}

	globPattern := os.Args[1]
	cmdArgs := os.Args[2:]

	// Set up fsnotify watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Error creating watcher: %v\n", err)
		return
	}
	defer watcher.Close()

	// Initially find and watch all matching files
	files, err := filepath.Glob(globPattern)
	if err != nil {
		fmt.Printf("Error matching files: %v\n", err)
		return
	}

	for _, file := range files {
		err = watcher.Add(file)
		if err != nil {
			fmt.Printf("Error watching file %s: %v\n", file, err)
			continue
		}
	}
	fmt.Printf("Watching %d files matching %s\n", len(files), globPattern)

	// Handle events
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				fmt.Printf("File %s changed, running command: %v\n", event.Name, cmdArgs)
				cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					fmt.Printf("Error running command: %v\n", err)
				}
			}

			// If new files are added matching the pattern, start watching them
			newFiles, err := filepath.Glob(globPattern)
			for _, file := range newFiles {
				if err = watcher.Add(file); err != nil {
					fmt.Printf("Error watching new file %s: %v\n", file, err)
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("Watcher error: %v\n", err)
		}
	}
}
