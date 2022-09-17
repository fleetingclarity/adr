package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"testing"
)

// setup is a test helper to create the required test env and initialized adr repo
func setup() (string, string, error) {
	workDir, _ := os.MkdirTemp("", "adr-wrk")
	startDir, _ := os.Getwd()
	err := os.Chdir(workDir)
	if err != nil {
		return "", "", err
	}
	c := NewDefaultConfig()
	err = c.CreateAndWrite()
	if err != nil {
		return "", "", err
	}
	err = c.EnsureRepositoryExists()
	if err != nil {
		return "", "", err
	}
	return startDir, workDir, nil
}

// writeAndClose is a test helper for creating temp files
func writeAndClose(p string, c string) error {
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(c))
	return err
}

// cleanup is a test helper to tear down the test env
func cleanup(startDir, workDir string) {
	defer os.Chdir(startDir)
	defer os.RemoveAll(workDir)
}

func Test_NumberIdentification(t *testing.T) {
	type test struct {
		name     string
		in       []string
		expected string
	}

	tests := []test{
		{name: "Simple case", in: []string{"001-a.md", "002-b.md"}, expected: "003-test-title.md"},
		{name: "No existing", in: []string{}, expected: "001-test-title.md"},
		{name: "Skip numbers", in: []string{"001-a.md", "003-c.md", "008-j.md"}, expected: "009-test-title.md"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startDir, workDir, err := setup()
			assert.NoError(t, err)
			for i, fs := range tt.in {
				p := path.Join(workDir, DefaultRepositoryDir, fs)
				_, err := os.Create(p)
				assert.NoError(t, err)
				err = writeAndClose(p, fmt.Sprintf("test file %d\n", i))
				assert.NoError(t, err)
			}
			v := make(map[string]string)
			v["Title"] = "test-title"
			cut := NewDefaultConfig().ADR
			err = cut.New(path.Join(workDir, DefaultRepositoryDir), v)
			assert.NoError(t, err)
			var files []string
			filepath.WalkDir(path.Join(workDir, DefaultRepositoryDir), func(s string, d fs.DirEntry, e error) error {
				if e != nil {
					return e
				}
				if !d.IsDir() {
					files = append(files, path.Base(s))
				}
				return nil
			})
			assert.Contains(t, files, tt.expected)
			cleanup(startDir, workDir)
		})
	}
}

func Test_TitleSanitizing(t *testing.T) {
	type test struct {
		name     string
		in       string
		expected string
		msg      string
	}
	tests := []test{
		{name: "Handle punctuation", in: "!some thing & with pu,nctuation", expected: "some-thing--with-punctuation", msg: "Expect two dashes between thing and with"},
		{name: "Handle uppercase", in: "UPPERCASE", expected: "uppercase", msg: "Simple upper to lower should never go wrong"},
		{name: "Handle dashes", in: "title-with-dashes", expected: "title-with-dashes", msg: "Not expected to do anything to dashes"},
		{name: "All cases", in: "*tItl,e-WitH mix!@", expected: "title-with-mix", msg: "Nobody would ever do this right?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cut := NewDefaultConfig()
			actual := cut.Sanitize(tt.in)
			assert.Equal(t, tt.expected, actual, tt.msg)
		})
	}
}
