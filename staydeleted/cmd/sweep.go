package cmd

// Copyright © 2018 Robert Impey robert-impey@users.noreply.github.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"log"
	"os"

	"github.com/robert-impey/tools/staydeleted/sdlib"

	"github.com/spf13/cobra"
)

var LogsDir string

var ExpiryMonths int
var Verbose bool

// sweepCmd represents the sweep command
var sweepCmd = &cobra.Command{
	Use:   "sweep",
	Short: "Sweep directories of files marked for deletion.",
	Long: `Walk through the directories given in the command line args
looking for files that have been marked for deletion.
`,
	Run: func(cmd *cobra.Command, args []string) {
		sweep(args)
	},
}

func init() {
	rootCmd.AddCommand(sweepCmd)
	sweepCmd.Flags().StringVarP(&LogsDir, "logs", "l", "",
		"The logs directory.")
	sweepCmd.Flags().IntVarP(&ExpiryMonths, "expiry", "e", 12,
		"The number of months before SD files expire.")
	sweepCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "Print verbosely.")
}

func sweep(paths []string) {
	for _, path := range paths {
		stat, err := os.Stat(path)
		if err != nil {
			log.Println(err)
			continue
		}

		if stat.IsDir() {
			err := sdlib.SweepDirectory(path, ExpiryMonths, Verbose)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			log.Printf("%v\n is not a directory!", path)
		}
	}
}
