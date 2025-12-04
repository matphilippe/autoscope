package main

import (
	"reflect"
	"testing"
)

func TestModuleMatchesFile(t *testing.T) {
	tests := []struct {
		name       string
		module     Module
		filepath   string
		wantMatch  bool
		wantScopes []string
	}{
		{
			name: "simple glob match",
			module: Module{
				Name:  "my-scope",
				Files: "modules/my-scope/**/*",
			},
			filepath:   "modules/my-scope/main.tf",
			wantMatch:  true,
			wantScopes: []string{"my-scope"},
		},
		{
			name: "simple glob no match",
			module: Module{
				Name:  "my-scope",
				Files: "modules/my-scope/**/*",
			},
			filepath:   "modules/other-scope/main.tf",
			wantMatch:  false,
			wantScopes: nil,
		},
		{
			name: "regex module match",
			module: Module{
				FilesRe: `modules/(?P<scope>\w+)/.*`,
			},
			filepath:   "modules/heey/main.tf",
			wantMatch:  true,
			wantScopes: []string{"heey"},
		},
		{
			name: "regex module no match",
			module: Module{
				FilesRe: `modules/(?P<scope>\w+)/.*`,
			},
			filepath:   "lib/main.tf",
			wantMatch:  false,
			wantScopes: nil,
		},
		{
			name: "regex module multiple captures",
			module: Module{
				FilesRe: `(?P<scope>api|db)/(?P<scope>\w+)/.*`,
			},
			filepath:   "api/users/handler.go",
			wantMatch:  true,
			wantScopes: []string{"api", "users"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatch, gotScopes := tt.module.matchesFile(tt.filepath)
			if gotMatch != tt.wantMatch {
				t.Errorf("Module.matchesFile() match = %v, want %v", gotMatch, tt.wantMatch)
			}
			if !reflect.DeepEqual(gotScopes, tt.wantScopes) {
				t.Errorf("Module.matchesFile() scopes = %v, want %v", gotScopes, tt.wantScopes)
			}
		})
	}
}

func TestGetScopesForFiles(t *testing.T) {
	config := &Config{
		Modules: []Module{
			{
				Name:  "api",
				Files: "src/api/**/*",
			},
			{
				Name:  "db",
				Files: "src/db/**/*",
			},
			{
				FilesRe: `modules/(?P<scope>\w+)/.*`,
			},
		},
	}

	tests := []struct {
		name       string
		files      []string
		wantScopes []string
	}{
		{
			name:       "single file matches one module",
			files:      []string{"src/api/handler.go"},
			wantScopes: []string{"api"},
		},
		{
			name:       "multiple files same module",
			files:      []string{"src/api/handler.go", "src/api/model.go"},
			wantScopes: []string{"api"},
		},
		{
			name:       "multiple files different modules",
			files:      []string{"src/api/handler.go", "src/db/query.go"},
			wantScopes: []string{"api", "db"},
		},
		{
			name:       "file matches regex module",
			files:      []string{"modules/auth/main.tf"},
			wantScopes: []string{"auth"},
		},
		{
			name:       "mixed matches",
			files:      []string{"src/api/handler.go", "modules/cache/main.tf"},
			wantScopes: []string{"api", "cache"},
		},
		{
			name:       "no matches",
			files:      []string{"README.md"},
			wantScopes: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotScopes := config.getScopesForFiles(tt.files)
			if !reflect.DeepEqual(gotScopes, tt.wantScopes) {
				t.Errorf("Config.getScopesForFiles() = %v, want %v", gotScopes, tt.wantScopes)
			}
		})
	}
}
