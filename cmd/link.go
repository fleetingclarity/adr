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
	conf "github.com/fleetingclarity/adr/config"
	"github.com/spf13/cobra"
	"strconv"
)

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link",
	Args:  cobra.ExactArgs(4),
	Short: "Link two ADRs",
	Long: `Create a link between two ADRs, such as when an ADR is amended or would otherwise be
altered but not fully superseded.

Expected usage: adr link <linker#> <link-message> <linked#> <back-link-message>

Example: adr link 182 "Amends some important thing" 10 "Important thing is amended"`,
	Run: func(cmd *cobra.Command, args []string) {
		lnum, err := strconv.Atoi(args[0])
		cobra.CheckErr(err)
		rnum, err := strconv.Atoi(args[2])
		cobra.CheckErr(err)
		lp := &conf.LinkPair{
			SourceNum: lnum,
			TargetNum: rnum,
			SourceMsg: args[1],
			BackMsg:   args[3],
			RepoDir:   config.Repository.Path,
		}
		err = config.Link(lp)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// linkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// linkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
