package main

import (
	"reflect"
	"testing"
)

func TestModuleMatchesFile(t *testing.T) {
	tests := []struct {
		name      string
		module    Module
		filepath  string
		shouldErr bool
		expected  string
	}{
		{
			name: "simple glob match",
			module: Module{
				Name: "my-scope",
				Glob: "modules/my-scope/**/*",
			},
			filepath: "modules/my-scope/main.tf",
			expected: "my-scope",
		},
		{
			name: "simple glob no match",
			module: Module{
				Name: "my-scope",
				Glob: "modules/my-scope/**/*",
			},
			filepath: "modules/other-scope/main.tf",
			expected: "",
		},
		{
			name: "regex module match",
			module: Module{
				FilesRe: `modules/(?P<scope>\w+)/.*`,
			},
			filepath: "modules/heey/main.tf",
			expected: "heey",
		},
		{
			name: "regex module no match",
			module: Module{
				FilesRe: `modules/(?P<scope>\w+)/.*`,
			},
			filepath: "lib/main.tf",
			expected: "",
		},
		{
			name: "regex module multiple captures should be no no",
			module: Module{
				FilesRe: `(?P<scope>api|db)/(?P<scope>\w+)/.*`,
			},
			filepath:  "api/users/handler.go",
			shouldErr: true,
			expected:  "",
		},
		{
			name: "wrong glob should be nono",
			module: Module{
				Name: "my-scope",
				Glob: "[]a]", // Got it from a doublestar test
			},
			filepath:  "api/users/handler.go",
			shouldErr: true,
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := tt.module.matchesFile(tt.filepath)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("expecting this test to err, but got (%s, %v) (%+v) \n", match, err, tt)
				}
			} else {
				if match != tt.expected {
					t.Errorf("Module.matchesFile() match = %v, want %v (%+v) \n", match, tt.expected, tt)
				}
			}
		})
	}
}

func TestGetScopesForFiles(t *testing.T) {
	config := &Config{
		Modules: []Module{
			{
				Name: "api",
				Glob: "src/api/**/*",
			},
			{
				Name: "db",
				Glob: "src/db/**/*",
			},
			{
				FilesRe: `modules/(?P<scope>\w+)/.*`,
			},
		},
	}

	tests := []struct {
		name     string
		files    []string
		expected []string
	}{
		{
			name:     "single file matches one module",
			files:    []string{"src/api/handler.go"},
			expected: []string{"api"},
		},
		{
			name:     "multiple files same module",
			files:    []string{"src/api/handler.go", "src/api/model.go"},
			expected: []string{"api"},
		},
		{
			name:     "multiple files different modules",
			files:    []string{"src/api/handler.go", "src/db/query.go"},
			expected: []string{"api", "db"},
		},
		{
			name:     "file matches regex module",
			files:    []string{"modules/auth/main.tf"},
			expected: []string{"auth"},
		},
		{
			name:     "mixed matches",
			files:    []string{"src/api/handler.go", "modules/cache/main.tf"},
			expected: []string{"api", "cache"},
		},
		{
			name:     "mixed matches unsorted",
			files:    []string{"modules/cache/main.tf", "src/api/handler.go"},
			expected: []string{"api", "cache"},
		},
		{
			name:     "no matches",
			files:    []string{"README.md"},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotScopes, err := config.getScopesForFiles(tt.files)
			if err != nil {
				t.Errorf("(%+v) ran into error: %v", tt, err)
			}
			if !reflect.DeepEqual(gotScopes, tt.expected) {
				t.Errorf("Config.getScopesForFiles() = %v, want %v", gotScopes, tt.expected)
			}
		})
	}
}
