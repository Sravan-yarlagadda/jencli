// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Sravan-yarlagadda/jencli/cli"
	"github.com/spf13/cobra"
)

var (
	url, token, user string
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a jenkins job",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("start called")

		// jencli := cli.Jencli{
		//Crumb: {UsesCrumb: false, CrumbString: "", CrumbValue: ""},
		// User:  "admin",
		// Token: "11fe33897ccc15106adca3d8110e939340",
		// }
		if len(url) != 1 && len(user) != 1 && len(token) != 1 {
			jencli := cli.Jencli{
				User:  user,
				Token: token,
			}
			jencli.Start(url, "", true)
		} else {
			fmt.Println("Usage : jencli start -l <job_url> -u <user> -t <token> -p <parameters> -m")
			os.Exit(1)
		}

		fmt.Println(strings.Join(args, " "))
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	startCmd.Flags().StringVarP(&url, "url", "l", " ", "URL to start the job")
	startCmd.Flags().BoolP("monitor", "m", false, "Monitor job")
	startCmd.Flags().StringVarP(&user, "user", "u", " ", "user")
	startCmd.Flags().StringVarP(&token, "token", "t", " ", "token")
}
