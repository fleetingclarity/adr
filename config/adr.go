package config

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unicode"
)

type ADR struct {
	FormatName    string
	TitleTemplate string
	BodyTemplate  string
}

// New will create a new ADR using the in-memory configuration. It will determine the
// next number for the new ADR and write a new file to repoDir
func (a *ADR) New(repoDir string, values map[string]string) error {
	// 1. determine next number and pad with 0s
	n, err := next(repoDir)
	if err != nil {
		return err
	}
	ns := fmt.Sprintf("%03d", n)
	values["Title"] = Sanitize(values["Title"])
	values["Number"] = ns
	values["Date"] = fmt.Sprintf("%v-%v-%v", time.Now().Year(), time.Now().Month(), time.Now().Day())
	// 2. create go template
	t := template.New(fmt.Sprintf("%s-adr", a.FormatName))
	// 3. use title template to create new file
	tt, err := t.Parse(a.TitleTemplate)
	f, err := titledFile(tt, repoDir, values)
	defer f.Close()
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

// Sanitize ensures that the given string matches Title expectations (e.g. lowercase, no spaces, etc)
func Sanitize(title string) string {
	if title == "" {
		return "no-title-given"
	}
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return '-'
		}
		if unicode.IsUpper(r) {
			return unicode.ToLower(r)
		}
		if r == '-' || unicode.IsDigit(r) || unicode.IsLower(r) {
			return r
		}
		return -1
	}, title)
}

func titledFile(t *template.Template, repoDir string, v map[string]string) (io.WriteCloser, error) {
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

func next(dir string) (int, error) {
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

// Find will return the repoDir/NNN-file-name.md for the specified ADR number if it exists in the repoDir
func Find(repoDir string, num int) (string, error) {
	if num > 999 {
		return "", errors.New("the adr tool does not support 4 digit records, please create a Github issue if you require over a thousand records")
	}
	var matches []string
	err := filepath.WalkDir(repoDir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(fmt.Sprintf("%03d*.md", num), filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	// always take the first match... hopefully there's only one
	return matches[0], nil
}

// UpdateStatus will search for the Status section and replace the existing status with the 'to' status
func UpdateStatus(path, to string) error {
	r, err := os.Open(path) // closed at end of function
	if err != nil {
		return err
	}
	w, err := os.Create(path + ".tmp")
	if err != nil {
		return err
	}
	defer w.Close()
	scanner := bufio.NewScanner(r)
	err = replaceSectionContent(scanner, w, "Status", to)
	if err != nil {
		return err
	}
	err = r.Close()
	if err != nil {
		return err
	}
	err = os.Rename(path+".tmp", path)
	if err != nil {
		return err
	}
	return nil
}

// replaceSectionContent will find the identified section if it exists and replace all content until the following section
func replaceSectionContent(scanner *bufio.Scanner, wc io.WriteCloser, sectionPattern string, replacement string) error {
	newSection := "## "
	inSection := false
	replacementWritten := false
	defer wc.Close()
	for scanner.Scan() {
		line := scanner.Text()
		if inSection && replacementWritten && strings.Contains(line, newSection) {
			inSection = false
		}
		if inSection && replacementWritten {
			line = ""
		}
		if inSection && !replacementWritten {
			caser := cases.Title(language.AmericanEnglish)
			line = caser.String(replacement)
			replacementWritten = true
		}
		if line == newSection+sectionPattern {
			// our next line should be what we're looking to replace
			inSection = true
		}
		b := []byte(line + "\n")
		_, err := wc.Write(b)
		if err != nil {
			return err
		}
	}
	return nil
}

// appendToSection will insert content just before the section following the identified section
func appendToSection(scanner *bufio.Scanner, wc io.WriteCloser, sectionPattern string, newContent string) error {
	newSection := "## "
	inSection := false
	appended := false
	defer wc.Close()
	for scanner.Scan() {
		line := scanner.Text()
		if inSection && !appended && strings.Contains(line, newSection) {
			line = newContent + "\n" + line
			inSection = false
			appended = true
		}
		if line == newSection+sectionPattern {
			// we're in the section, now just need to identify the start of the next section to enable insert
			inSection = true
		}
		b := []byte(line + "\n")
		_, err := wc.Write(b)
		if err != nil {
			return err
		}
	}
	return nil
}

type LinkPair struct {
	SourceNum int
	TargetNum int
	SourceMsg string
	BackMsg   string
	RepoDir   string
}

// Link will use the LinkPair to insert links into the 'Status' section
func Link(p *LinkPair) error {
	sp, err := Find(p.RepoDir, p.SourceNum)
	if err != nil {
		return err
	}
	tp, err := Find(p.RepoDir, p.TargetNum)
	if err != nil {
		return err
	}
	sbase := path.Base(sp)
	tbase := path.Base(tp)
	err = appendForLink(sp, fmt.Sprintf("[Links to %s: %s](./%s)", sbase, p.SourceMsg, tbase))
	if err != nil {
		return err
	}
	err = appendForLink(tp, fmt.Sprintf("[Links to %s: %s](./%s)", tbase, p.BackMsg, sbase))
	if err != nil {
		return err
	}
	return nil
}

// Supersede will change the Source status to 'Superseded' and link to the Target. It will also append the
// 'Supersedes' backlink to the Target status
func Supersede(p *LinkPair) error {
	sp, err := Find(p.RepoDir, p.SourceNum)
	if err != nil {
		return err
	}
	tp, err := Find(p.RepoDir, p.TargetNum)
	if err != nil {
		return err
	}
	tFile, err := os.Open(tp)
	if err != nil {
		return err
	}
	tWC, err := os.Create(tp + ".tmp")
	if err != nil {
		return err
	}
	sbase := path.Base(sp)
	tbase := path.Base(tp)
	var smsg, tmsg string
	if len(p.SourceMsg) > 0 {
		smsg = ": " + p.SourceMsg
	}
	if len(p.BackMsg) > 0 {
		tmsg = ": " + p.BackMsg
	}
	err = UpdateStatus(sp, "Superseded")
	if err != nil {
		return err
	}
	sFile, err := os.Open(sp)
	if err != nil {
		return err
	}
	sWC, err := os.Create(sp + ".tmp")
	if err != nil {
		return err
	}
	sScanner := bufio.NewScanner(sFile)
	err = appendToSection(sScanner, sWC, "Status", fmt.Sprintf("[Superseded by %s%s](./%s)", tbase, smsg, tbase))
	if err != nil {
		_ = sFile.Close()
		return err
	}
	_ = sFile.Close()
	err = os.Rename(sp+".tmp", sp)
	if err != nil {
		return err
	}
	tScanner := bufio.NewScanner(tFile)
	err = appendToSection(tScanner, tWC, "Status", fmt.Sprintf("[Supersedes %s%s](./%s)", sbase, tmsg, sbase))
	if err != nil {
		_ = tFile.Close()
		return err
	}
	_ = tFile.Close()
	err = os.Rename(tp+".tmp", tp)
	if err != nil {
		return err
	}
	return nil
}

func appendForLink(path, newContent string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(f)
	wc, err := os.Create(path + ".tmp")
	err = appendToSection(scanner, wc, "Status", newContent)
	if err != nil {
		return err
	}
	f.Close()
	err = os.Rename(path+".tmp", path)
	if err != nil {
		return err
	}
	return nil
}

const (
	defaultTitleTemplate = "{{ .Number }}-{{ .Title }}.md"
	defaultBodyTemplate  = `# {{ .Number }}-{{ .Title }}
Date: {{ .Date }}

## Status
Proposed

## Context
Describe the environment. What forces are exerting pressure on this decision? What are you trying to accomplish?

## Decision
Describe the decision but don't be too verbose. 1 or 2 pages of the details that matter. The audience is future team members.

## Consequences
Describe the effect of the decision. What are you trading off? What is good, bad, or even deferred to another day?

`
)
