/*
Copyright © 2022 fleetingclarity <72276886+fleetingclarity@users.noreply.github.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path"
)

const (
	defaultRepoDir          = "docs/decisions"
	defaultConfigExt        = "yaml"
	defaultConfigName       = ".adr"
	uiInitFileExistsNonZero = `A configuration file already exists in the current directory. 
If you would like to re-initialize your adr repository please delete it first then rerun the init 
command. Additionally the strict flag was set so the exit code will be set to 1.`
	uiInitFileExists = `A configuration file already exists in the current directory. If you would 
like to re-initialize your adr repository please delete it first then rerun the init command.`
	uiInitInitializing = `Initializing adr repository...`
	uiInitSuccess      = `Success! You repository has been initialized and is ready to start tracking 
architecture decision records.`
)

var (
	dir    string
	strict bool
)

// NewInitCmd represents the init command
func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init [options]",
		Aliases: []string{"initialize", "start", "manage"},
		Short:   "Initialize a repository for management by the adr tool",
		Long: `Init (adr init) will initialize a repository to be managed by
the adr tool. 

It will first check for $HOME/.adr.yaml, and if found it will use applicable
settings from that configuration file during initialization.

It will then create a repository local version of .adr.yaml using the options
specified. When there are collisions the priority order of options is:
1. Command line options
2. Env vars
3. $HOME/.adr.yaml`,
		Run: runInit,
	}
	return cmd
}

func runInit(cmd *cobra.Command, args []string) {
	cmd.Println(uiInitInitializing)
	if config.UsingLocalConfig {
		// no clobbering, just exit with optional failure status
		if strict {
			cmd.Println(uiInitFileExistsNonZero)
			os.Exit(1)
		}
		cmd.Println(uiInitFileExists)
		os.Exit(0)
	}
	if config.Repository.Path == "" {
		config.Repository = &Repository{
			Path: dir,
		}
	}
	f, err := os.Create(path.Join(config.WorkingDirectory, config.CfgFileName+"."+config.CfgFileExt))
	err = WriteLocalConfig(&config, f)
	cobra.CheckErr(err)
	err = f.Close()
	cobra.CheckErr(err)
	err = config.EnsureRepositoryExists()
	if err != nil {
		cmd.Println(err)
	}
	cmd.Println(uiInitSuccess)
}

func (c *Config) EnsureRepositoryExists() error {
	if _, err := os.Stat(c.Repository.Path); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(path.Join(c.WorkingDirectory, c.Repository.Path), os.ModePerm)
		if err != nil {
			return errors.New(fmt.Sprintf("warning: unable to create the repository directory %s. You will likely need to create it manually", path.Join(c.WorkingDirectory, c.Repository.Path)))
		}
	}
	return nil
}

func WriteLocalConfig(c *Config, w io.Writer) error {
	o, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	_, err = w.Write(o)
	return err
}

func init() {
	initCmd := NewInitCmd()
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	initCmd.Flags().StringVarP(&dir, "repository", "r", defaultRepoDir, "Change the path that adrs will be stored in")
	initCmd.Flags().BoolVarP(&strict, "strict", "s", false, "Used for scripting to cause failure if we're already initialized")
}
