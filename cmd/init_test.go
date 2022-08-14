package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

func Test_WriteLocalConfig(t *testing.T) {
	type test struct {
		name           string
		in             *Config
		expected       string
		shouldNotMatch bool
	}
	tests := []test{
		{name: "happy . for local", in: &Config{Repository: &Repository{Path: "."}}, expected: "repository:\n    path: .\n", shouldNotMatch: false},
		{name: "sad . for local", in: &Config{Repository: &Repository{Path: "."}}, expected: "repository:\n    path:.\n", shouldNotMatch: true},
		{name: "happy anything", in: &Config{Repository: &Repository{Path: "anything"}}, expected: "repository:\n    path: anything\n", shouldNotMatch: false},
		{name: "sad anything", in: &Config{Repository: &Repository{Path: "anything"}}, expected: "repository:\n    path:anything\n", shouldNotMatch: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			err := WriteLocalConfig(tt.in, b)
			assert.NoError(t, err, "No test in this table should generate errors")
			if tt.shouldNotMatch {
				assert.NotEqual(t, string(b.Bytes()), tt.expected)
			} else {
				assert.Equal(t, string(b.Bytes()), tt.expected)
			}
		})
	}
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

// setup is a test helper to create the required test env
func setup() (startDir, workDir, testHomeDir, configFile string) {
	workDir, _ = os.MkdirTemp("", "adr-wrk")
	startDir, _ = os.Getwd()
	os.Chdir(workDir)
	testHomeDir, _ = os.MkdirTemp("", "adr-home")
	os.Setenv("HOME", testHomeDir)
	os.MkdirAll(path.Join(testHomeDir, defaultConfigName), os.ModePerm)
	configFile = defaultConfigName + "." + defaultConfigExt
	return startDir, workDir, testHomeDir, configFile
}

// cleanup is a test helper to tear down the test env
func cleanup(startDir, workDir, testHomeDir string) {
	defer os.Chdir(startDir)
	defer os.RemoveAll(workDir)
	defer os.RemoveAll(testHomeDir)
}

// test pick up content from base config in home directory and use it to initialize a local repo
func Test_InitWithHomeConfig(t *testing.T) {
	startDir, workDir, testHomeDir, configFile := setup()
	defer cleanup(startDir, workDir, testHomeDir)
	// create a base config file
	expectedContent := "repository:\n    path: .\n"
	err := writeAndClose(path.Join(testHomeDir, defaultConfigName, configFile), expectedContent)
	require.NoError(t, err, "error creating a config file for testing in the temp home dir")
	cmdOutput := &bytes.Buffer{}
	cmd := NewInitCmd()
	cmd.SetOut(cmdOutput)
	cmd.SetErr(cmdOutput)
	err = cmd.Execute()
	require.NoError(t, err, "error during cmd execution")
	output := cmdOutput.String()
	expectedOutput := uiInitInitializing + "\n" + uiInitSuccess + "\n"
	// todo: this should be a separate test for ui
	assert.Equal(t, expectedOutput, output, "expected to break if more output is added to init cmd")
	createdFileContents, err := os.ReadFile(path.Join(workDir, configFile))
	require.NoError(t, err, "error reading the created config file in the test working directory")
	assert.Equal(t, expectedContent, string(createdFileContents), "should match")
	err = os.Remove(path.Join(testHomeDir, defaultConfigName, configFile))
	err = os.Remove(path.Join(workDir, configFile))
}

// test using defaults when no global config in home directory
func Test_InitWithDefaults(t *testing.T) {
	startDir, workDir, testHomeDir, configFile := setup()
	defer cleanup(startDir, workDir, testHomeDir)
	expectedContent := "repository:\n    path: " + defaultRepoDir + "\n"
	cmdOutput := &bytes.Buffer{}
	cmd := NewInitCmd()
	cmd.SetOut(cmdOutput)
	cmd.SetErr(cmdOutput)
	err := cmd.Execute()
	require.NoError(t, err, "error during command execution")
	createdFileContents, err := os.ReadFile(path.Join(workDir, configFile))
	require.NoError(t, err, "error reading the contents of the created file")
	assert.Equal(t, expectedContent, string(createdFileContents), "should have a file with default contents")
	err = os.RemoveAll(path.Join(workDir, "docs"))
}

// test no file modification when local config exists
func Test_InitNoChangesWhenLocalConfigExists(t *testing.T) {
	startDir, workDir, testHomeDir, configFile := setup()
	defer cleanup(startDir, workDir, testHomeDir)
	homeConfigContent := "repository:\n    path: somepath\n"
	err := writeAndClose(path.Join(testHomeDir, defaultConfigName, configFile), homeConfigContent)
	defer os.Remove(path.Join(testHomeDir, defaultConfigName, configFile))
	require.NoError(t, err, "error writing a home config")
	expectedContent := "repository:\n    path: " + defaultRepoDir + "\n"
	err = writeAndClose(path.Join(workDir, configFile), expectedContent)
	require.NoError(t, err, "error writing the working directory config file")
	cmdOutput := &bytes.Buffer{}
	cmd := NewInitCmd()
	cmd.SetOut(cmdOutput)
	cmd.SetErr(cmdOutput)
	err = cmd.Execute()
	require.NoError(t, err, "error during command execution")
	contentAfterRun, err := os.ReadFile(path.Join(workDir, configFile))
	require.NoError(t, err, "error reading the contents of the workdir config file")
	assert.Equal(t, expectedContent, string(contentAfterRun), "should not have been modified")
	os.Remove(path.Join(workDir, configFile))
}

// ?todo: test flags and env vars?
