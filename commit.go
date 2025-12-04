package main

import (
	"regexp"
	"strings"
)

var conventionalCommitRe = regexp.MustCompile(`^([a-z]+)(?:\(([^)]+)\))?:\s*(.*)$`)

func parseCommitMessage(msg string) (commitType, scope, description string) {
	msg = strings.TrimSpace(msg)
	matches := conventionalCommitRe.FindStringSubmatch(msg)
	
	if matches == nil {
		return "", "", msg
	}
	
	return matches[1], matches[2], matches[3]
}

func addScopesToCommit(commit string, newScopes []string) string {
	lines := strings.SplitN(commit, "\n", 2)
	firstLine := lines[0]
	rest := ""
	if len(lines) > 1 {
		rest = "\n" + lines[1]
	}
	
	commitType, existingScope, description := parseCommitMessage(firstLine)
	
	if commitType == "" {
		return commit
	}
	
	existingScopes := make(map[string]bool)
	if existingScope != "" {
		for _, s := range strings.Split(existingScope, ",") {
			existingScopes[strings.TrimSpace(s)] = true
		}
	}
	
	var allScopes []string
	if existingScope != "" {
		allScopes = append(allScopes, strings.Split(existingScope, ",")...)
	}
	
	for _, newScope := range newScopes {
		if !existingScopes[newScope] {
			allScopes = append(allScopes, newScope)
			existingScopes[newScope] = true
		}
	}
	
	var newFirstLine string
	if len(allScopes) == 0 {
		newFirstLine = commitType + ": " + description
	} else {
		newFirstLine = commitType + "(" + strings.Join(allScopes, ",") + "): " + description
	}
	
	return newFirstLine + rest
}
