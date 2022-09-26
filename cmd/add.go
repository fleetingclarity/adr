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
	conf "github.com/fleetingclarity/adr/config"
	"github.com/spf13/cobra"
	"path"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:     "add \"Some title\"",
	Aliases: []string{"new", "create"},
	Args:    cobra.ExactArgs(1),
	Short:   "Add a new record to the ADR repository",
	Long: `Create a new ADR in the repository. The only argument to the Add function is the Title of
the new ADR. The Title will be used to name the file and as the heading within the markdown of the
file. Capital letters, symbols, and spaces will be converted to single dashes (-).

Example usage: adr add "Some title"
Results in: A file in your repo appropriately numbered like 'NNN-some-title.md'`,
	Run: func(cmd *cobra.Command, args []string) {
		m := make(map[string]string)
		m["Title"] = args[0]
		if verbose {
			fmt.Printf("Your title '%s' will be converted to '%s'\n", args[0], conf.Sanitize(args[0]))
		}
		cobra.CheckErr(config.New(config.Repository.Path, m))
		fmt.Printf("Success! Edit your new ADR at %s\n", path.Join(config.WorkingDirectory, config.Repository.Path))
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
