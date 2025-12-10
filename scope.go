package main

import (
	"fmt"
	"regexp"
	"slices"

	"github.com/bmatcuk/doublestar/v4"
)

type Module struct {
	Name       string `yaml:"name"`
	Glob       string `yaml:"glob"`
	FilesRe    string `yaml:"filesRe"`
	compiledRe *regexp.Regexp
}

type Config struct {
	Modules []Module `yaml:"modules"`
}

func (mod *Module) IsValid() error {
	if mod.Name == "" && mod.Glob != "" {
		return fmt.Errorf("if it has a glob, it must have a name")
	}
	if mod.Name != "" && mod.Glob == "" {
		return fmt.Errorf("if it has a name, it must have a glob")
	}
	if mod.Name != "" && mod.FilesRe != "" {
		return fmt.Errorf("either define name and glob, OR filesRe")
	}
	return nil
}

func (c *Config) IsValid() error {
	for _, mod := range c.Modules {
		if err := mod.IsValid(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Module) matchesFile(filepath string) (string, error) {
	if m.Name != "" {
		return m.matchesGlob(filepath)
	}
	if m.FilesRe != "" {
		return m.matchesRegex(filepath)
	}
	return "", nil
}

func (m *Module) matchesGlob(path string) (string, error) {
	matched, err := doublestar.Match(m.Glob, path)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", nil
	}
	return m.Name, nil
}

func (m *Module) matchesRegex(path string) (string, error) {
	if m.compiledRe == nil {
		re, err := regexp.Compile(m.FilesRe)
		if err != nil {
			return "", err
		}
		m.compiledRe = re
	}

	matches := m.compiledRe.FindStringSubmatch(path)
	if matches == nil {
		return "", nil
	}

	scopes := []string{}
	subexpNames := m.compiledRe.SubexpNames()
	for i, name := range subexpNames {
		if i > 0 && name == "scope" && i < len(matches) {
			scopes = append(scopes, matches[i])
		}
	}
	if len(scopes) >= 2 {
		return "", fmt.Errorf("match error: got 2 matches on %s for file %s: %+v", m.FilesRe, path, matches)
	}
	return scopes[0], nil
}

func (c *Config) getScopesForFiles(files []string) ([]string, error) {
	scopeSet := make(map[string]bool)
	for _, file := range files {
		for i := range c.Modules {
			scope, err := c.Modules[i].matchesFile(file)
			if err != nil {
				return nil, err
			}
			if scope != "" {
				scopeSet[scope] = true
			}
		}
	}

	scopes := make([]string, 0, len(scopeSet))
	for k := range scopeSet {
		scopes = append(scopes, k)
	}
	slices.Sort(scopes)
	return scopes, nil
}
