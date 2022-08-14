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
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	ignore  bool
	config  *Config
	verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "adr",
	Short: "A tool for managing Architectural Decision Records (ADRs), written in Go",
	Long: `adr is a CLI tool to ease the management (create, update, link, supersede, lint) of
ADR repositories. It will generate new files based on templates that you can specify
and it will ensure that all files in the repository match that template when used
in its linting capacity.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.PersistentFlags().BoolVarP(&ignore, "ignore", "i", false, "Ignore config files and only use defaults + options (other than config file)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output useful for debugging")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	/*
		todo: redo flow for loading configs
		0. [x] if flag for ignore and cfg file are both used then exit with error msg
		1. [ ] if ignore flag set then skip all files and only load defaults
		2. [ ] if flag for cfg file is used then load it with non-colliding defaults, skip all other
		4. [ ] if not ignore and not file flag and local file exists then load it
		3. [ ] else if not ignore and not file flag and $HOME file exists then load it to &config
		5. [ ] using flags to set individual configs not currently supported
		6. [ ] using env vars to set/override configs may be possible but not a supported feature yet
	*/
	config = &Config{} // necessary for test suites, hoping this won't affect production but idk cobra/viper very well
	if ignore && cfgFile != "" {
		fmt.Println("You cannot ignore config files and use a specific file, please choose one option or the other")
		os.Exit(1)
	}
	wd, err := os.Getwd()
	home := os.Getenv("HOME")

	cobra.CheckErr(err)
	if !ignore {
		if cfgFile != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfgFile)
		} else {
			viper.SetConfigType(defaultConfigExt)
			fileName := defaultConfigName + "." + defaultConfigExt
			wdFile := path.Join(wd, fileName)
			homeFile := path.Join(home, defaultConfigName, fileName)
			// always use local if it exists, then try to use home
			if _, err := os.Stat(wdFile); errors.Is(err, os.ErrNotExist) {
				if _, err = os.Stat(homeFile); errors.Is(err, os.ErrNotExist) {
					cfgFile = ""
				} else {
					cfgFile = homeFile
				}
			} else {
				cfgFile = wdFile
				config.UsingLocalConfig = true
			}
			viper.SetConfigFile(cfgFile)
		}
	}

	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
		err = viper.Unmarshal(&config)
		cobra.CheckErr(err)
	}
	if config.CfgFile == "" {
		config.CfgFile = viper.ConfigFileUsed()
	}
	if config.WorkingDirectory == "" {
		config.WorkingDirectory = wd
	}
	if config.CfgFileName == "" {
		config.CfgFileName = defaultConfigName
	}
	if config.CfgFileExt == "" {
		config.CfgFileExt = defaultConfigExt
	}
	if config.UserHome == "" {
		config.UserHome = home
	}
	if config.Repository == nil {
		config.Repository = &Repository{
			Path: "",
		}
	}
}
