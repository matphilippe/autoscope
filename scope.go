package main

import (
	"regexp"

	"github.com/bmatcuk/doublestar/v4"
)

type Module struct {
	Name    string `yaml:"name"`
	Files   string `yaml:"files"`
	FilesRe string `yaml:"filesRe"`
	
	compiledRe *regexp.Regexp
}

type Config struct {
	Modules []Module `yaml:"modules"`
}

func (m *Module) matchesFile(filepath string) (bool, []string) {
	if m.FilesRe != "" {
		return m.matchesRegex(filepath)
	}
	
	if m.Files != "" {
		return m.matchesGlob(filepath)
	}
	
	return false, nil
}

func (m *Module) matchesGlob(path string) (bool, []string) {
	matched, err := doublestar.Match(m.Files, path)
	if err != nil || !matched {
		return false, nil
	}
	
	return true, []string{m.Name}
}

func (m *Module) matchesRegex(path string) (bool, []string) {
	if m.compiledRe == nil {
		re, err := regexp.Compile(m.FilesRe)
		if err != nil {
			return false, nil
		}
		m.compiledRe = re
	}
	
	matches := m.compiledRe.FindStringSubmatch(path)
	if matches == nil {
		return false, nil
	}
	
	var scopes []string
	subexpNames := m.compiledRe.SubexpNames()
	for i, name := range subexpNames {
		if i > 0 && name == "scope" && i < len(matches) {
			scopes = append(scopes, matches[i])
		}
	}
	
	return true, scopes
}

func (c *Config) getScopesForFiles(files []string) []string {
	scopeSet := make(map[string]bool)
	var scopes []string
	
	for _, file := range files {
		for i := range c.Modules {
			matched, moduleScopes := c.Modules[i].matchesFile(file)
			if matched {
				for _, scope := range moduleScopes {
					if !scopeSet[scope] {
						scopeSet[scope] = true
						scopes = append(scopes, scope)
					}
				}
			}
		}
	}
	
	if scopes == nil {
		return []string{}
	}
	
	return scopes
}
