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
	conf "github.com/fleetingclarity/adr/config"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	config  *conf.Config
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

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output useful for debugging")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config = &conf.Config{} // necessary for test suites, hoping this won't affect production but IDK cobra/viper very well
	wd, err := os.Getwd()
	cobra.CheckErr(err)

	viper.SetConfigType(conf.DefaultConfigExt)
	fileName := conf.DefaultConfigName + "." + conf.DefaultConfigExt
	wdFile := path.Join(wd, fileName)
	if _, err := os.Stat(wdFile); errors.Is(err, os.ErrNotExist) {
		cfgFile = ""
	} else {
		cfgFile = wdFile
	}
	viper.SetConfigFile(cfgFile)

	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
		err = viper.Unmarshal(&config)
		cobra.CheckErr(err)
	} else { // otherwise use default values
		config = conf.NewDefaultConfig()
	}
}
