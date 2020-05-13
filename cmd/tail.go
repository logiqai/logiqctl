/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"strings"

	"github.com/logiqai/logiqctl/services"
	"github.com/logiqai/logiqctl/utils"

	"github.com/spf13/cobra"
)

var procs, labels string

// tailCmd represents the tail command
var tailCmd = &cobra.Command{
	Use:     "tail",
	Aliases: []string{"t"},
	Short:   "Stream logs from logiq Observability stack",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var procsArray []string
		var labelsArray []string
		var namespaces = []string{utils.GetDefaultNamespace()}
		if procs != "" {
			procsArray = strings.Split(procs, ",")
		}
		if labels != "" {
			labelsArray = strings.Split(labels, ",")
		}

		fmt.Println("Crunching data for you...")
		services.Tail(namespaces, labelsArray, args, procsArray, nil)

		return
	},
}

func init() {

	rootCmd.AddCommand(tailCmd)
	tailCmd.Flags().StringVarP(&procs, "process", "p", "", `Filter logs by process id`)
	tailCmd.Flags().StringVarP(&labels, "labels", "l", "", `Filter logs by  label`)

}
