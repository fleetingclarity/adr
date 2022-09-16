package config

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"text/template"
	"time"
)

type ADR struct {
	FormatName    string
	Sections      []string
	TitleTemplate string
	BodyTemplate  string
}

func (a *ADR) Name() string {
	return a.FormatName
}

// New will create a new ADR using the in-memory configuration. It will determine the
// next number for the new ADR and write a new file to repoDir
func (a *ADR) New(repoDir string, values map[string]any) error {
	// 1. determine next number and pad with 0s
	n, err := a.next(repoDir)
	if err != nil {
		return err
	}
	ns := fmt.Sprintf("%03d", n)
	values["Number"] = ns
	values["Date"] = fmt.Sprintf("%v-%v-%v", time.Now().Year(), time.Now().Month(), time.Now().Day())
	// 2. create go template
	t := template.New(fmt.Sprintf("%s-adr", a.FormatName))
	// 3. use title template to create new file
	tt, err := t.Parse(a.TitleTemplate)
	f, err := a.titledFile(tt, repoDir, values)
	if err != nil {
		return err
	}
	// 4. execute body template and write to created file
	bt, err := t.Parse(a.BodyTemplate)
	err = bt.Execute(f, values)
	if err != nil {
		return err
	}
	return nil
}

func (a *ADR) titledFile(t *template.Template, repoDir string, v map[string]any) (io.Writer, error) {
	pathBuffer := bytes.NewBufferString("")
	err := t.Execute(pathBuffer, v)
	if err != nil {
		return nil, err
	}
	f, err := os.Create(path.Join(repoDir, pathBuffer.String()))
	if err != nil {
		return nil, err
	}
	return f, err
}

func (a *ADR) next(dir string) (int, error) {
	var matches []string
	err := filepath.WalkDir(dir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match("[0-9]*.md", filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, filepath.Base(path))
		}
		return nil
	})
	if err != nil {
		return -1, err
	}
	if len(matches) > 0 {
		sort.Strings(matches)
		c := matches[len(matches)-1]
		r := regexp.MustCompile(`^\d+`)
		r.Longest()
		m := r.FindString(c)
		v, err := strconv.Atoi(m)
		if err != nil {
			return -1, err
		}
		return v + 1, nil
	} else {
		return 1, nil
	}
}

const (
	defaultTitleTemplate = "{{ .Number }}-{{ .Title }}.md"
	defaultBodyTemplate  = `# {{ .Number }}-{{ .Title }}
Date: {{ .Date }}

## Status
Placeholder value

## Context
Placeholder value

## Decision
Placeholder value

## Consequences
Placeholder value

`
)
