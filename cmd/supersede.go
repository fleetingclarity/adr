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

// supersedeCmd represents the supersede command
var supersedeCmd = &cobra.Command{
	Use:   "supersede",
	Args:  cobra.ExactArgs(4),
	Short: "Mark an ADR as Superseded",
	Long: `Update two existing ADRs for superseding. The Source is updated as Superseded and a link is created to
the Target (Superseding) ADR. The Target is updated with a backlink to the Source (Superseded).

Expected usage: adr supersede <Source#> <Msg> <Target#> <BackMsg>
Example: adr supersede 1 "some note" 2 "" # empty quotes if you don't want a message'`,
	Run: func(cmd *cobra.Command, args []string) {
		sNum, err := strconv.Atoi(args[0])
		cobra.CheckErr(err)
		tNum, err := strconv.Atoi(args[2])
		cobra.CheckErr(err)
		lp := &conf.LinkPair{
			SourceNum: sNum,
			TargetNum: tNum,
			SourceMsg: args[1],
			BackMsg:   args[3],
			RepoDir:   config.Repository.Path,
		}
		cobra.CheckErr(conf.Supersede(lp))
	},
}

func init() {
	rootCmd.AddCommand(supersedeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// supersedeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// supersedeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
