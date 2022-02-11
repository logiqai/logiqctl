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

	"github.com/logiqai/logiqctl/grpc_utils"
	"github.com/logiqai/logiqctl/services"
	"github.com/logiqai/logiqctl/ui"
	"github.com/logiqai/logiqctl/utils"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get <resource_name>",
	Short: "Display one or more of your LOGIQ resources",
	Long:  `Prints a table that displays the most important information about the LOGIQ resource you specify.`,
	Example: `
List all applications for the selected context
logiqctl get applications

List all applications for all available contexts
logiqctl get applications all

List all dashboards
logiqctl get dashboards all

Get dashboard
logiqctl get dashboard dashboard-slug

List events for the available namespace
logiqctl get events

List events for all available namespaces
logiqctl get events all

List or export eventrules
logiqctl get eventrules

List all namespaces
logiqctl get namespaces

List all processes
logiqctl get processes

List all queries
logiqctl get queries all

Get httpingestkey
logiqctl get httpingestkey

Get query
logiqctl get query query-slug
`,
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(NewListNameSpaceCommand())
	getCmd.AddCommand(NewListApplicationsCommand())
	getCmd.AddCommand(NewListProcessesCommand())
	getCmd.AddCommand(NewListEventsCommand())
	getCmd.AddCommand(NewGetEventRulesCommand())
	getCmd.AddCommand(ui.NewListDashboardsCommand())
	getCmd.AddCommand(ui.NewListQueriesCommand())
	getCmd.AddCommand(ui.NewListDatasourcesCommand())
	getCmd.AddCommand(getForwardsCommand())
	getCmd.AddCommand(getMappersCommand())
	getCmd.AddCommand(getHttpingestkeyCommand())
	getCmd.AddCommand(ui.NewGetLogEvents())
}

func getMappersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mappers",
		Example: "logiqctl get mappers",
		Aliases: []string{"mappers"},
		Short:   "Get logflow log mappers",
		PreRun:  utils.PreRun,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("getMappersCommand - coming soon")
		},
	}
	return cmd
}
func getForwardsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "forwards",
		Example: "logiqctl get forwards",
		Aliases: []string{"forwards"},
		Short:   "Get logflow log forwards",
		PreRun:  utils.PreRun,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("getforwardsCommand - comming soon")
		},
	}
	return cmd
}

func getHttpingestkeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "httpingestkey",
		Example: "logiqctl get httpingestkey",
		Aliases: []string{"ingestkey"},
		Short:   "Get httpingestkey",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			if u, cookieJar, err := grpc_utils.GetCookies(); err != nil {
				fmt.Println("Error getting httpingestkey: ", err.Error())
			} else {
				tokFound := false
				for _, c := range cookieJar.Cookies(u) {
					if "x-api-key" == c.Name {
						tokFound = true
						fmt.Println(c.Value)
					}
				}
				if !tokFound {
					fmt.Println("Error getting the httpingestkey")
				}
			}
		},
	}
	return cmd
}

func NewListNameSpaceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "namespaces",
		Example: "logiqctl get namespaces|ns|n",
		Aliases: []string{"n", "ns"},
		Short:   "List all available namespaces",
		PreRun:  utils.PreRun,
		Run: func(cmd *cobra.Command, args []string) {
			services.ListNamespaces()
		},
	}

	return cmd
}

func NewListApplicationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "applications",
		Example: "logiqctl get applications|apps|a",
		Aliases: []string{"a", "apps"},
		Short:   "List all available applications within the default namespace",
		PreRun:  utils.PreRunWithNs,
		Run: func(cmd *cobra.Command, args []string) {
			services.GetApplicationsV2(false)
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:     "all",
		Example: "logiqctl get applications all",
		Short:   "List all available applications across namespaces",
		PreRun:  utils.PreRun,
		Run: func(cmd *cobra.Command, args []string) {
			services.GetApplicationsV2(true)
		},
	})
	return cmd
}

func NewListEventsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "events",
		Example: `
List last 30 events
- logiqctl get events|e

List events by application 
- logiqctl get events -a=sshd

`,
		Aliases: []string{"ev"},
		Short:   "List all available events for the namespace",
		PreRun:  utils.PreRunWithNs,
		Run: func(cmd *cobra.Command, args []string) {
			services.GetEvents(application, process)
		},
	}
	cmd.Flags().StringVarP(&application, "application", "a", "", `Filter events by application name`)
	cmd.Flags().StringVarP(&process, "process", "p", "", `Filter events by application and process name`)
	return cmd
}

func NewListProcessesCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "processes",
		Example: "logiqctl get processes|proc|p",
		Aliases: []string{"p", "proc"},
		Short:   "List all available processes. This command runs an interactive prompt that lets you choose an application from a list of available applications.",
		PreRun:  utils.PreRunWithNs,
		Run: func(cmd *cobra.Command, args []string) {
			services.ListProcesses()
		},
	}
}

func NewGetEventRulesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "eventrules",
		Example: "logiqctl get eventrule",
		Aliases: []string{"eventrule", "er"},
		Short:   "List event rules",
		PreRun:  utils.PreRunWithNs,
		Run: func(cmd *cobra.Command, args []string) {
			help := `Usage:
logiqctl get eventrules all
logiqctl get eventrules all -w filename.json
logiqctl get eventrules groups
logiqctl get eventrules groups -g=group1,group2,...
logiqctl get eventrules groups -g=group1,group2,... -w filename.json
`
			fmt.Print(help)
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:     "all",
		Example: "logiqctl get eventrules all",
		Aliases: []string{},
		Short:   "List all event rules",
		PreRun:  utils.PreRunWithNs,
		Run: func(cmd *cobra.Command, args []string) {
			services.GetEventRules(args, nil)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:     "groups",
		Example: "logiqctl get eventrules all",
		Aliases: []string{},
		Short:   "List all event rules",
		PreRun:  utils.PreRunWithNs,
		Run: func(cmd *cobra.Command, args []string) {
			services.GetEventRuleGroups(args)
		},
	})
	cmd.Flags().StringVarP(&utils.EventRuleGroupsFlag, "groups", "g", "", "list of groups separated by comma")
	cmd.Flags().StringVarP(&utils.FlagFile, "file", "w", "", "Path to file to be written to")
	return cmd
}
