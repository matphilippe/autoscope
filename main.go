package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func getConfigFile() string {
	path := os.Getenv("AUTOSCOPE_CONFIG")
	if path == "" {
		path = ".autoscope.yaml"
	}
	return path
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <commit-msg> file1 file2 file3 ...\n", os.Args[0])
		os.Exit(1)
	}
	msg := os.Args[1]
	files := os.Args[2:]
	config, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v, exiting\n", err)
		os.Exit(1)
	}
	out, err := doTheThing(*config, msg, files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error matching modules: %v", err)
		os.Exit(1)
	}
	fmt.Printf("%s", out)
}

func doTheThing(config Config, msg string, files []string) (string, error) {
	scopes, err := config.getScopesForFiles(files)
	if err != nil {
		return "", err
	}
	if len(scopes) == 0 {
		return msg, nil
	}
	return addScopesToCommit(msg, scopes), nil
}

func loadConfig() (*Config, error) {
	var config Config
	file := getConfigFile()
	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file `%s` not found", file)
		} else {
			return nil, err
		}
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	err = config.IsValid()
	if err != nil {
		return nil, err
	}
	return &config, nil
}
