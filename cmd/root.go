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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	ignore  bool
	config  Config
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

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	wd, err := os.Getwd()
	cobra.CheckErr(err)
	if !ignore {
		if cfgFile != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfgFile)
		} else {
			// Find home directory.
			// todo: check if only use local config file?
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)

			// Search config in home directory with name ".adr" (without extension).
			viper.AddConfigPath(home + "/.adr")
			viper.AddConfigPath(wd)

			viper.SetConfigType(defaultConfigExt)
			viper.SetConfigName(defaultConfigName)
		}
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		err = viper.Unmarshal(&config)
		cobra.CheckErr(err)
	}
	viper.AutomaticEnv() // read in environment variables that match
	if config.WorkingDirectory == "" {
		config.WorkingDirectory = wd
	}
	if config.CfgFile == "" {
		config.CfgFile = viper.ConfigFileUsed()
	}
}
