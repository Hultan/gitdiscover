package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

func main() {
	if len(os.Args)>1 {
		for i := 1; i < len(os.Args); i++ {
			arg := os.Args[i]

			switch arg {
			case "search":
				fmt.Printf("Search path : %v\n", os.Args[i+1])
				i++
			default:
				unknownCommand(arg)
				os.Exit(1) // Unknown command
			}
			i++
		}
	}

	config := NewConfig()
	config.Load()

	for _, p := range config.Paths {
		gitPath := path.Join(p, ".git")
		if _, err := os.Stat(gitPath); os.IsNotExist(err) {
			fmt.Printf("%v : Not a git directory! err=%v\n", p, err)
		} else {
			gs := getGitStatus(p)
			fmt.Printf("%v : %s", p,gs)
		}
	}

	os.Exit(0)
}

func unknownCommand(arg string) {
	fmt.Printf("Unknown command : %v\n", arg)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("\tgitdiscover search [path]")
}

func getGitStatus(path string) string {
	cmd := exec.Command("/home/per/bin/gitprompt-go")
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		return "failed to check git status"
	}
	return string(out)
}