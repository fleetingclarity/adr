package cmd

import (
	"bytes"
	conf "github.com/fleetingclarity/adr/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

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
func setup() (startDir, workDir, configFile string) {
	workDir, _ = os.MkdirTemp("", "adr-wrk")
	startDir, _ = os.Getwd()
	os.Chdir(workDir)
	configFile = conf.DefaultConfigName + "." + conf.DefaultConfigExt
	return startDir, workDir, configFile
}

// cleanup is a test helper to tear down the test env
func cleanup(startDir, workDir string) {
	defer os.Chdir(startDir)
	defer os.RemoveAll(workDir)
}

// until I get a better understanding of cobra/viper let's use this to guard our user interface
func Test_InitHappyPathOutput(t *testing.T) {
	startDir, workDir, _ := setup()
	defer cleanup(startDir, workDir)
	cmdOutput := &bytes.Buffer{}
	cmd := NewInitCmd()
	cmd.SetOut(cmdOutput)
	cmd.SetErr(cmdOutput)
	err := cmd.Execute()
	require.NoError(t, err, "error during cmd execution")
	output := cmdOutput.String()
	expectedOutput := uiInitInitializing + "\n" + uiInitSuccess + "\n"
	assert.Equal(t, expectedOutput, output, "expected to break if more output is added to init cmd")
}

// test using defaults when no global config in home directory
func Test_InitWithDefaults(t *testing.T) {
	startDir, workDir, configFile := setup()
	defer cleanup(startDir, workDir)
	expectedContent := "repository:\n    path: " + conf.DefaultRepositoryDir + "\n"
	cmdOutput := &bytes.Buffer{}
	cmd := NewInitCmd()
	cmd.SetOut(cmdOutput)
	cmd.SetErr(cmdOutput)
	err := cmd.Execute()
	require.NoError(t, err, "error during command execution")
	createdFileContents, err := os.ReadFile(path.Join(workDir, configFile))
	require.NoError(t, err, "error reading the contents of the created file")
	assert.Contains(t, string(createdFileContents), expectedContent, "should have a file with default contents")
	err = os.RemoveAll(path.Join(workDir, "docs"))
}

// test no file modification when local config exists
func Test_InitNoChangesWhenLocalConfigExists(t *testing.T) {
	startDir, workDir, configFile := setup()
	defer cleanup(startDir, workDir)
	expectedContent := "repository:\n    path: " + conf.DefaultRepositoryDir + "\n"
	err := writeAndClose(path.Join(workDir, configFile), expectedContent)
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
