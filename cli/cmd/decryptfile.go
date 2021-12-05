/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	golog "log"

	"github.com/spf13/cobra"
)

// decryptfileCmd represents the decryptfile command
var decryptfileCmd = &cobra.Command{
	Use:   "decryptfile",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("decryptfile called")
		absFilePath, err := cmd.Flags().GetString("filepath")
		if err != nil || absFilePath == "" {
			golog.Fatalln(err, ":: file path must be specified. can be specified using '--filepath' or '-f' flag.")
		}

		golog.Println("WARN: filepath must be absolute. Currently relative file-paths are not supported.")
		crypter.Init("des_input.txt")
		crypter.DecryptFile(absFilePath)
	},
}

func init() {
	rootCmd.AddCommand(decryptfileCmd)
	decryptfileCmd.Flags().StringP("filepath", "f", "", "Absolute path of the file to be decrypted")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decryptfileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// decryptfileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
