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
	"github.com/spf13/cobra"
	"os"
	"path"
)

const (
	defaultRepoDir    = "docs/decisions"
	defaultConfigExt  = "yaml"
	defaultConfigName = ".adr"
	uiInitFileExists  = `A configuration file already exists in the current directory. If you would 
like to re-initialize your adr repository please delete it first then rerun the init command.`
	uiInitInitializing = `Initializing adr repository...`
	uiInitSuccess      = `Success! You repository has been initialized and is ready to start tracking 
architecture decision records.`
	defaultTitleTemplate = "NNNN-{{ .Title }}.md"
	defaultBodyTemplate  = `# {{ .Title }}
Date: {{ .Date }}

## Status
{{ .Status }}

## Context
{{ .Context }}

## Decision
{{ .Decision }}

## Consequences
{{ .Consequences }}

`
)

var (
	dir string
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
	if !config.UsingLocalConfig {
		if config.Repository.Path == "" || cmd.Flags().Changed("repository") { // precedence to flags
			config.Repository = &Repository{
				Path: dir,
			}
		}
		f, err := os.Create(path.Join(config.WorkingDirectory, config.CfgFileName+"."+config.CfgFileExt))
		err = config.Write(f)
		cobra.CheckErr(err)
		err = f.Close()
		cobra.CheckErr(err)
		err = config.EnsureRepositoryExists()
		if err != nil {
			cmd.Println(err)
		}
		cmd.Println(uiInitSuccess)
	} else {
		// no clobbering, just print and allow exit
		cmd.Println(uiInitFileExists)
	}
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
	initCmd.Flags().StringVarP(&dir, "repository", "r", defaultRepoDir, "Change the path that ADRs will be stored in")
}
