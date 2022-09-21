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
	"strconv"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Args:  cobra.ExactArgs(2),
	Short: "Update status of an existing ADR",
	Long:  `Update an existing ADR Status by providing the ADR number and the new Status`,
	Run: func(cmd *cobra.Command, args []string) {
		n, err := strconv.Atoi(args[0])
		s := args[1]
		cobra.CheckErr(err)
		if verbose {
			cmd.Println("Updato potato")
		}
		a, err := config.Find(config.Repository.Path, n)
		cobra.CheckErr(err)
		err = config.UpdateStatus(a, s)
		cobra.CheckErr(err)
		cmd.Printf("%s status updated to %s\n", a, s)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
