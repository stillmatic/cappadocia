package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/alecthomas/repr"
	"github.com/fsnotify/fsnotify"
)

var CLI struct {
	Watch struct {
		GlobPattern string   `arg:"" required:"" name:"glob" help:"Glob pattern to watch files." type:"string"`
		Command     string   `arg:"" required:"" name:"command" help:"Command to run upon file changes." type:"string"`
		Args        []string `arg:"" required:"" name:"args" help:"Additional arguments for the command."`
	} `cmd:"" help:"Watch files and run a command upon changes."`
}

func main() {
	ctx := kong.Parse(&CLI)
	// fmt.Printf("Command: %s\n", ctx.Command())
	switch ctx.Command() {
	case "watch <glob> <command> <args>":
		globPattern := CLI.Watch.GlobPattern
		cmdArgs := append([]string{CLI.Watch.Command}, CLI.Watch.Args...)
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			fmt.Printf("Error creating watcher: %v\n", err)
			return
		}
		defer watcher.Close()

		files, err := filepath.Glob(globPattern)
		if err != nil {
			fmt.Printf("Error matching files: %v\n", err)
			return
		}

		if len(files) == 0 {
			fmt.Printf("No files found matching %s\n", globPattern)
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

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
					repr.Print(cmd)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					err := cmd.Run()
					if err != nil {
						fmt.Printf("Error running command: %v\n", err)
					}
				}

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

	default:
		panic(ctx.Command())
	}
}
