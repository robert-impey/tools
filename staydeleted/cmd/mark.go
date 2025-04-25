package cmd

// Copyright © 2018 - 2025 Robert Impey robert-impey@users.noreply.github.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"log"

	"github.com/robert-impey/tools/staydeleted/sdlib"
	"github.com/spf13/cobra"
)

var Keep bool

// markCmd represents the mark command
var markCmd = &cobra.Command{
	Use:   "mark",
	Short: "Mark a file for deletion or keeping",
	Long:  `Files marked for deletion or keeping will be taken care of by the sweep command.`,
	Run: func(cmd *cobra.Command, args []string) {
		action := sdlib.GetActionForBool(Keep)

		for _, arg := range args {
			err := sdlib.SetActionForFile(arg, action)
			if err != nil {
				log.Printf("couldn't set action for file '%s'\n", arg)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(markCmd)

	markCmd.Flags().BoolVarP(&Keep, "keep", "k", false, "Keep this file.")
}
