package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	
	configContent := `modules:
  - name: api
    files: src/api/**/*
  - name: db
    files: src/db/**/*
  - filesRe: modules/(?P<scope>\w+)/.*
`
	
	err := os.WriteFile(filepath.Join(tmpDir, ".svscope.yaml"), []byte(configContent), 0644)
	if err != nil {
		t.Fatal(err)
	}
	
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	
	err = os.MkdirAll(filepath.Join(tmpDir, "src/api"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	
	err = os.WriteFile(filepath.Join(tmpDir, "src/api/handler.go"), []byte("package api"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	
	cmd = exec.Command("git", "add", "src/api/handler.go")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	
	commitMsgFile := filepath.Join(tmpDir, "COMMIT_EDITMSG")
	err = os.WriteFile(commitMsgFile, []byte("fix: fixed the bug\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)
	
	config, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	
	changedFiles, err := getChangedFiles()
	if err != nil {
		t.Fatal(err)
	}
	
	if len(changedFiles) != 1 || changedFiles[0] != "src/api/handler.go" {
		t.Errorf("Expected [src/api/handler.go], got %v", changedFiles)
	}
	
	scopes := config.getScopesForFiles(changedFiles)
	if len(scopes) != 1 || scopes[0] != "api" {
		t.Errorf("Expected [api], got %v", scopes)
	}
	
	commitMsg, _ := os.ReadFile(commitMsgFile)
	newMsg := addScopesToCommit(string(commitMsg), scopes)
	
	expected := "fix(api): fixed the bug\n"
	if newMsg != expected {
		t.Errorf("Expected %q, got %q", expected, newMsg)
	}
}
