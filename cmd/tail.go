/*
Copyright © 2020 Logiq.ai <cli@logiq.ai>

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
	"os"

	"github.com/logiqai/logiqctl/utils"

	"github.com/logiqai/logiqctl/services"

	"github.com/spf13/cobra"
)

var application, process, labels string

var tailExample = `
Tail logs 
- logiqctl tail
`

// tailCmd represents the tail command
var tailCmd = &cobra.Command{
	Use:     "tail",
	Aliases: []string{"t"},
	Example: tailExample,
	Short:   "Stream logs sent to your LOGIQ Observability platform in real-time.",
	Long: `
The 'logiqctl tail' command is similar to the 'tail -f' command. It allows you to stream the log data that is being sent to your LOGIQ Observability platform in real-time. You can see logs from the cluster at multiple levels. Running the command 'tail' without any options brings up an interactive prompt that lets you choose an application and process in the current context. 
`,
	Run: func(cmd *cobra.Command, args []string) {
		var labelsArray []string
		var appName string
		var procId string
		app, err := services.RunSelectApplicationForNamespacePrompt(true)
		handleError(err)
		if app != nil {
			appName = app.Name
			process, err := services.RunSelectProcessesForNamespaceAndAppPrompt(app.Name, true)
			handleError(err)
			if process != nil {
				procId = process.ProcID
			} else {
				procId = "*"
			}
		} else {
			appName = "*"
		}
		if utils.FlagFile != "" {
			if _, err := os.Stat(utils.FlagFile); err == nil {
				fmt.Printf("File %s exists\n", utils.FlagFile)
				os.Exit(1)
			}
			if f, err := os.Create(utils.FlagFile); err != nil {
				fmt.Printf("Cannot create file: %s\n", utils.FlagFile)
				os.Exit(1)
			} else {
				f.Close()
			}
		}
		services.Tail(appName, procId, labelsArray)
		return
	},
}

func init() {

	rootCmd.AddCommand(tailCmd)
	tailCmd.Flags().StringVarP(&utils.FlagFile, "write-to-file", "w", "", "Path to file")
	tailCmd.Flags().IntVarP(&utils.FlagMaxFileSize, "max-file-size", "m", 10, "Max output file size")
	//tailCmd.Flags().StringVarP(&process, "process", "p", "", `Filter logs by process id`)
	//tailCmd.Flags().StringVarP(&labels, "labels", "l", "", `Filter logs by label`)
}
