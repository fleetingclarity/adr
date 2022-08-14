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

func writeAndClose(p string, c string) error {
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(c))
	return err
}

func Test_Init(t *testing.T) {
	// https://github.com/carolynvs/stingoftheviper/blob/main/main_test.go
	workDir, err := os.MkdirTemp("", "adr-wrk")
	require.NoError(t, err, "error creating a temp directory for testing")
	startDir, err := os.Getwd()
	require.NoError(t, err, "error getting the current working directory")
	defer os.Chdir(startDir)
	defer os.RemoveAll(workDir)
	err = os.Chdir(workDir)
	require.NoError(t, err, "error changing directory to the temp directory")
	testHomeDir, err := os.MkdirTemp("", "adr-home")
	require.NoError(t, err, "error creating a temp home directory for testing")
	defer os.RemoveAll(testHomeDir)
	realHome, err := os.UserHomeDir()
	require.NoError(t, err, "error getting the user home directory ($HOME)")
	os.Setenv("HOME", testHomeDir)
	defer func() {
		os.Setenv("HOME", realHome)
	}()
	configFile := defaultConfigName + "." + defaultConfigExt

	// test pick up content from base config in home directory and use it to initialize a local repo
	t.Run("initialize with settings from base config", func(t *testing.T) {
		// create a base config file
		expectedContent := "repository:\n    path: .\n"
		err = os.Mkdir(path.Join(testHomeDir, defaultConfigName), 0755)
		require.NoError(t, err, "error creating the config dir under the temp home")
		err := writeAndClose(path.Join(testHomeDir, defaultConfigName, configFile), expectedContent)
		require.NoError(t, err, "error creating a config file for testing in the temp home dir")
		defer os.Remove(path.Join(testHomeDir, configFile))
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
		defer os.Remove(path.Join(workDir, configFile))
		require.NoError(t, err, "error reading the created config file in the test working directory")
		assert.Equal(t, expectedContent, string(createdFileContents), "should match")
	})
	// test using defaults when no global config in home directory
	t.Run("initialize with defaults", func(t *testing.T) {
		expectedContent := "repository:\n    path: " + defaultRepoDir + "\n"
		cmdOutput := &bytes.Buffer{}
		cmd := NewInitCmd()
		cmd.SetOut(cmdOutput)
		cmd.SetErr(cmdOutput)
		err = cmd.Execute()
		require.NoError(t, err, "error during command execution")
		createdFileContents, err := os.ReadFile(path.Join(workDir, configFile))
		require.NoError(t, err, "error reading the contents of the created file")
		assert.Equal(t, expectedContent, string(createdFileContents), "should have a file with default contents")
	})

	// todo: test no file modification when local config exists
	// todo: test no file modification and exit status is non-zero when local config exists
	// ?todo: test flags and env vars?
}
