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
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path"
)

const (
	defaultRepoDir    = "docs/decisions"
	defaultConfigExt  = "yaml"
	defaultConfigName = ".adr"
)

var (
	dir    string
	strict bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
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

func runInit(cmd *cobra.Command, args []string) {
	fmt.Println("init called")
	if config.CfgFile == "" {
		config.CfgFile = path.Join(config.WorkingDirectory, defaultConfigName+"."+defaultConfigExt)
	}
	if localConfigExists(config.CfgFile) {
		// no clobbering, just exit with optional failure status
		fmt.Printf("file '%s' exists is true\n", config.CfgFile)
		if strict {
			fmt.Printf(`The configuration file %s already exists. If you want
to re-initialize your adr repository please delete it first then rerun the init command.
Additionally the strict flag was set so the exit code will be set to 1`, config.CfgFile)
			os.Exit(1)
		}
		fmt.Printf(`The configuration file %s already exists. If you want
to re-initialize your adr repository please delete it first then rerun the init command`, config.CfgFile)
		os.Exit(0)
	}
	fmt.Printf("just for giggles the file is %s\n", config.CfgFile)
	fmt.Println("let's continue")
	config.Repository = &Repository{
		RelativePath: dir,
	}
	f, err := os.Create(config.CfgFile)
	err = WriteLocalConfig(&config, f)
	cobra.CheckErr(err)
	fmt.Printf("Success! The configuration file at %s will be used to manage the repository located at %s\n", config.CfgFile, dir)
}

func localConfigExists(f string) bool {
	_, err := os.Stat(f)
	return err == nil
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
