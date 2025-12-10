package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	config := Config{
		Modules: []Module{
			{
				Name: "test",
				Glob: "test/**/*",
			},
			{
				FilesRe: "modules/(?P<scope>\\w+)/.*",
			},
		},
	}
	tests := []struct {
		name     string
		message  string
		files    []string
		expected string
	}{
		{
			name:     "simple commit without scope",
			message:  "fix: fixed the bug",
			files:    []string{},
			expected: "fix: fixed the bug",
		},
		{
			name:     "simple commit with scope",
			message:  "fix: fixed the bug",
			files:    []string{"test/some/file"},
			expected: "fix(test): fixed the bug",
		},

		{
			name:     "simple commit with 2 scope",
			message:  "fix: fixed the bug",
			files:    []string{"test/some/file", "modules/A/thing"},
			expected: "fix(A,test): fixed the bug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := doTheThing(config, tt.message, tt.files)
			if err != nil {
				t.Fail()
			}
			if msg != tt.expected {
				t.Errorf("doTheThing() outputs %s, want %s: %+v", msg, tt.expected, tt)
			}
		})
	}
}
