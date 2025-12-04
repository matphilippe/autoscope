package main

import (
	"testing"
)

func TestParseCommitMessage(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType string
		wantScope string
		wantDesc string
	}{
		{
			name:      "simple commit without scope",
			input:     "fix: fixed the bug",
			wantType:  "fix",
			wantScope: "",
			wantDesc:  "fixed the bug",
		},
		{
			name:      "commit with scope",
			input:     "fix(api): fixed the bug",
			wantType:  "fix",
			wantScope: "api",
			wantDesc:  "fixed the bug",
		},
		{
			name:      "commit with multiple scopes",
			input:     "feat(api,db): add new feature",
			wantType:  "feat",
			wantScope: "api,db",
			wantDesc:  "add new feature",
		},
		{
			name:      "commit without type",
			input:     "just a commit message",
			wantType:  "",
			wantScope: "",
			wantDesc:  "just a commit message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotScope, gotDesc := parseCommitMessage(tt.input)
			if gotType != tt.wantType {
				t.Errorf("parseCommitMessage() type = %v, want %v", gotType, tt.wantType)
			}
			if gotScope != tt.wantScope {
				t.Errorf("parseCommitMessage() scope = %v, want %v", gotScope, tt.wantScope)
			}
			if gotDesc != tt.wantDesc {
				t.Errorf("parseCommitMessage() desc = %v, want %v", gotDesc, tt.wantDesc)
			}
		})
	}
}

func TestAddScopeToCommit(t *testing.T) {
	tests := []struct {
		name      string
		commit    string
		newScopes []string
		want      string
	}{
		{
			name:      "add scope to commit without scope",
			commit:    "fix: fixed the bug",
			newScopes: []string{"api"},
			want:      "fix(api): fixed the bug",
		},
		{
			name:      "add scope to commit with existing scope",
			commit:    "fix(api): fixed the bug",
			newScopes: []string{"db"},
			want:      "fix(api,db): fixed the bug",
		},
		{
			name:      "add multiple scopes",
			commit:    "feat: new feature",
			newScopes: []string{"api", "db"},
			want:      "feat(api,db): new feature",
		},
		{
			name:      "don't add duplicate scope",
			commit:    "fix(api): fixed the bug",
			newScopes: []string{"api"},
			want:      "fix(api): fixed the bug",
		},
		{
			name:      "don't add duplicate scope in multiple",
			commit:    "fix(api,db): fixed the bug",
			newScopes: []string{"db", "cache"},
			want:      "fix(api,db,cache): fixed the bug",
		},
		{
			name:      "handle non-conventional commit",
			commit:    "just a message",
			newScopes: []string{"api"},
			want:      "just a message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addScopesToCommit(tt.commit, tt.newScopes)
			if got != tt.want {
				t.Errorf("addScopesToCommit() = %v, want %v", got, tt.want)
			}
		})
	}
}
