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
	"github.com/logiqai/logiqctl/services"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get namespaces(n)|applications(a)|processes(p)",
	Short: "Display one or many resources",
	Long: `Prints a table of the most important information about the specified resources. For example:

# List all namespaces
logiqctl get namespaces

# List all applications for the selected context
logiqctl get applications

# List all applications for all the available context
logiqctl get applications all

# List all processes
logiqctl get processes

`,
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(NewListNameSpaceCommand())
	getCmd.AddCommand(NewListApplicationsCommand())
	getCmd.AddCommand(NewListProcessesCommand())
	getCmd.AddCommand(NewListEventsCommand())
}

func NewListNameSpaceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "namespaces",
		Example: "logiqctl get namespaces|ns|n",
		Aliases: []string{"n", "ns"},
		Short:   "List the available name spaces",
		PreRun:  preRun,
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
		Short:   "List all the available applications in default namespace",
		PreRun:  preRunWithNs,
		Run: func(cmd *cobra.Command, args []string) {
			services.GetApplicationsV2(false)
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:     "all",
		Example: "logiqctl get applications all",
		Short:   "List all the available applications across namespaces",
		PreRun:  preRun,
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
# logiqctl get events|e

#List events by application 
logiqctl get events -a=sshd

`,
		Aliases: []string{"e"},
		Short:   "List all the available events for the namespace",
		PreRun:  preRunWithNs,
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
		Short:   "List all the available processes, runs an interactive prompt to select applications",
		PreRun:  preRunWithNs,
		Run: func(cmd *cobra.Command, args []string) {
			services.ListProcesses()
		},
	}
}
