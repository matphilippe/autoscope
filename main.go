package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <commit-msg-file>\n", os.Args[0])
		os.Exit(1)
	}

	commitMsgFile := os.Args[1]

	config, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	changedFiles, err := getChangedFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting changed files: %v\n", err)
		os.Exit(1)
	}

	scopes := config.getScopesForFiles(changedFiles)

	if len(scopes) == 0 {
		os.Exit(0)
	}

	commitMsg, err := os.ReadFile(commitMsgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading commit message: %v\n", err)
		os.Exit(1)
	}

	newCommitMsg := addScopesToCommit(string(commitMsg), scopes)

	err = os.WriteFile(commitMsgFile, []byte(newCommitMsg), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing commit message: %v\n", err)
		os.Exit(1)
	}
}

func loadConfig() (*Config, error) {
	data, err := os.ReadFile(".svscope.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			data, err = os.ReadFile(".svscope.yml")
			if err != nil {
				return nil, fmt.Errorf("config file not found (.svscope.yaml or .svscope.yml)")
			}
		} else {
			return nil, err
		}
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func getChangedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var files []string
	for _, line := range lines {
		if line != "" {
			files = append(files, line)
		}
	}

	return files, nil
}
