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
	"strings"

	"github.com/manifoldco/promptui"

	"github.com/logiqai/logiqctl/utils"

	"github.com/logiqai/logiqctl/services"

	"github.com/spf13/cobra"
)

var logsExample = `
Print logs for the LOGIQ ingest server
- logiqctl logs logiq-flash

Print logs in JSON format
- logiqctl -o=json logs logiq-flash

In case of a Kubernetes deployment, a Stateful Set is an application, and each pod in it is a process
Print logs for logiq-flash ingest server filtered by process logiq-flash-2
The --process (-p) flag lets you view logs for the individual pod
- logiqctl logs -p=logiq-flash-2 logiq-flash

Runs an interactive prompt that lets you choose filters
- logiqctl logs interactive|i

Search logs for specific keywords or terms
- logiqctl logs search "your search term"   

If the flag --follow (-f) is specified, the logs will be streamed until the end of the log. 

`

var logsLong = `
The 'logs' command is used to view historical logs. This command expects a namespace and an application to be available to return results. You can set the default namespace using the 'logiqctl set-context' command or pass the namespace as '-n=NAMESPACE' flag. The application name also needs to be passed as an argument to the command. You can also use the 'interactive' command to choose from the list of available applications and processes.   

**Note:**
- The global flag '--time-format' is not applicable for this command.
- The global flag '--output' only supports JSON format for this command.`

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:     "logs",
	Example: logsExample,
	Aliases: []string{"log"},
	Short:   "View logs for the given namespace and application",
	Long:    logsLong,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println(cmd.UsageString())
			return
		}
		if utils.FlagProcId != "" {
			proc, err := services.GetProcessByApplicationAndProc(args[0], utils.FlagProcId)
			handleError(err)
			query(args[0], "", proc.ProcID, proc.LastSeen)
			return
		}
		app, err := services.GetApplicationByName(args[0])
		handleError(err)
		query(args[0], "", "", app.LastSeen)
	},
}

func query(appName, searchTerm, procId string, lastSeen int64) {
	queryId, err := services.Query(appName, searchTerm, procId, lastSeen)
	handleError(err)
	if queryId != "" {
		for {
			prompt := promptui.Prompt{
				Label:     "View More ",
				IsConfirm: true,
			}

			result, _ := prompt.Run()
			if strings.ToLower(result) == "y" {
				hasMoreData, err := services.GetDataNext(queryId)
				handleError(err)
				if !hasMoreData {
					break
				}
			} else {
				break
			}
		}
	}
}

var interactiveCmd = &cobra.Command{
	Use:     "interactive",
	Aliases: []string{"i"},
	Short:   `Runs an interactive prompt to display logs.`,
	Run: func(cmd *cobra.Command, args []string) {
		app, err := services.RunSelectApplicationForNamespacePrompt(false)
		handleError(err)
		proc, err := services.RunSelectProcessesForNamespaceAndAppPrompt(app.Name, false)
		handleError(err)
		fmt.Printf("You could also run this directly `logiqctl logs -p=%s %s`\n", proc.ProcID, app.Name)
		fmt.Printf("Fetching logs for %s (namespace), %s (application) and %s (process)\n\n", utils.GetDefaultNamespace(), app.Name, proc.ProcID)
		query(app.Name, "", proc.ProcID, proc.LastSeen)
	},
}

var searchCmd = &cobra.Command{
	Use:     "search",
	Aliases: []string{"s"},
	Example: ``,
	Short:   `Search logs for specific keywords or terms.`,
	Long:    `Search for specific keywords or terms in logs within a namespace.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println(cmd.Usage())
			return
		}
		app, err := services.RunSelectApplicationForNamespacePrompt(false)
		handleError(err)
		query(app.Name, args[0], "", -1)
	},
}

func init() {
	logsCmd.PersistentFlags().StringVarP(&utils.FlagLogsSince, "since", "s", "1h", `Only return logs newer than a relative duration. This is in relative to the last
seen log time for a specified application or processes within the namespace.
A duration string is a possibly signed sequence of decimal numbers, each with optional
fraction and a unit suffix, such as "3h34m", "1.5h" or "24h". Valid time units are "s", "m", "h"`)
	logsCmd.PersistentFlags().Uint32Var(&utils.FlagLogsPageSize, "page-size", 30, `Number of log entries to return in one page`)
	logsCmd.PersistentFlags().BoolVarP(&utils.FlagLogsFollow, "follow", "f", false, `Specify if the logs should be streamed.`)
	logsCmd.Flags().StringVarP(&utils.FlagProcId, "process", "p", "", `Filter logs by  proc id`)
	rootCmd.AddCommand(logsCmd)
	logsCmd.AddCommand(interactiveCmd)
	logsCmd.AddCommand(searchCmd)
}
