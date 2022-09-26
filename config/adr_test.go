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
			sut := NewDefaultConfig().ADR
			err = sut.New(path.Join(workDir, DefaultRepositoryDir), v)
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
			actual := Sanitize(tt.in)
			assert.Equal(t, tt.expected, actual, tt.msg)
		})
	}
}

func Test_FindFileByNumber(t *testing.T) {
	type test struct {
		name   string
		search int
		files  []string
		eName  string
		msg    string
	}
	tests := []test{
		{name: "Find simple only one match, first file", search: 1, files: []string{"001-first.md", "002-second.md"}, eName: "001-first.md", msg: "Should find the first file"},
		{name: "Find simple only one match, second file", search: 2, files: []string{"001-first.md", "002-second.md"}, eName: "002-second.md", msg: "Should find the first file"},
		{name: "Find when search digit exists more than once, smaller", search: 1, files: []string{"001-first.md", "011-second.md"}, eName: "001-first.md", msg: "Should find the first file"},
		{name: "Find when search digit exists more than once, larger", search: 11, files: []string{"001-first.md", "011-second.md"}, eName: "011-second.md", msg: "Should find the second file"},
		{name: "Potential multiple matches", search: 11, files: []string{"011-first.md", "111-second.md"}, eName: "011-first.md", msg: "Should find the first file"},
		{name: "Same pattern extra digits", search: 10, files: []string{"101-first.md", "010-second.md"}, eName: "010-second.md", msg: "Should find the second file"},
	}
	for _, tt := range tests {
		startDir, workDir, err := setup()
		assert.NoError(t, err)
		for i, f := range tt.files {
			_ = writeAndClose(path.Join(DefaultRepositoryDir, f), fmt.Sprintf("contents for file %d\n", i))
		}
		actual, err := Find(DefaultRepositoryDir, tt.search)
		assert.NoError(t, err)
		assert.Equal(t, path.Join(DefaultRepositoryDir, tt.eName), actual)
		cleanup(startDir, workDir)
	}
}

func Test_UpdateStatus(t *testing.T) {
	startDir, workDir, err := setup()
	assert.NoError(t, err)
	expected := "Approved"
	sut := NewDefaultConfig()
	repoDir := path.Join(workDir, DefaultRepositoryDir)
	err = sut.New(repoDir, map[string]string{"Title": "asdf"})
	assert.NoError(t, err)
	err = UpdateStatus(repoDir+"/001-asdf.md", expected)
	assert.NoError(t, err)
	actualBytes, err := os.ReadFile(path.Join(repoDir, "001-asdf.md"))
	assert.NoError(t, err)
	actualContents := string(actualBytes)
	assert.Contains(t, actualContents, expected, "We should have found the replacement text")
	//fmt.Println(actualContents)
	cleanup(startDir, workDir)
}

func Test_Supersede(t *testing.T) {
	startDir, workDir, err := setup()
	handleHarnessErr(t, err)
	c := NewDefaultConfig()
	repoDir := path.Join(workDir, DefaultRepositoryDir)
	err = c.New(repoDir, map[string]string{"Title": "first"})
	handleHarnessErr(t, err)
	err = c.New(repoDir, map[string]string{"Title": "second"})
	handleHarnessErr(t, err)
	// begin test
	lp := &LinkPair{
		SourceNum: 1,
		TargetNum: 2,
		SourceMsg: "some note",
		BackMsg:   "quick something",
		RepoDir:   repoDir,
	}
	err = Supersede(lp)
	assert.NoError(t, err)
	supersededBytes, err := os.ReadFile(path.Join(DefaultRepositoryDir, "001-first.md"))
	supersededContents := string(supersededBytes)
	assert.Contains(t, supersededContents, ": "+lp.SourceMsg, "simple sanity check failed, we should see our msg in the superseded file")
	targetBytes, err := os.ReadFile(path.Join(DefaultRepositoryDir, "002-second.md"))
	targetContents := string(targetBytes)
	assert.Contains(t, targetContents, ": "+lp.BackMsg, "simple sanity check failed, we should see our back message in the superseding file")
	cleanup(startDir, workDir)
}

func handleHarnessErr(t *testing.T, e error) {
	if e != nil {
		t.Fatal(e)
	}
}
