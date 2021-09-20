// Package cmd
//Copyright © 2020 Logiq.ai <cli@logiq.ai>
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//

package cmd

import (
	"fmt"
	applications2 "github.com/logiqai/logiqctl/api/v1/applications"
	"github.com/logiqai/logiqctl/loglerpart"
	"github.com/logiqai/logiqctl/services"
	"github.com/logiqai/logiqctl/utils"
	"github.com/spf13/cobra"
)

var logsExample = `
Print logs for the LOGIQ ingest server
  % logiqctl logs -a <application_name>

Print logs in JSON format:
  % logiqctl -o=json logs -a <application_name>

In case of a Kubernetes deployment, a Stateful Set is an application, and each pod in it is a process
Print logs for logiq-flash ingest server filtered by process logiq-flash-2
The --process (-p) flag lets you view logs for the individual pod
  % logiqctl logs -p=<proc_id> -a <application_name>

Runs an interactive prompt that lets you choose filters
  % logiqctl logs interactive|i

Search logs for specific keywords or terms see help:
  % logiqctl logs search --help

More examples:  
  % logiqctl logs -a <application_name> search <searchterm>
  % logiqctl logs -a <application_name> -p <proc_id> search <searchterm>
  % logiqctl logs -a <application_name> -p <proc_id> search <searchterm> -g

If the flag --follow (-f) is specified, the logs will be streamed until the end of the log. 

One can automatically generate pattern-signature (PS) for logs using flag --psmod (-g).
Add-on executable "psmod" from logiqhub is required to run side-by-side with logiqctl. 
Enable PS generation will generate stat file ps_stat.out that computes byte and log counts and 
percentage for each pattern signature 

More examples:  
  % logiqctl config set-context <namespace>
  % logiqctl logs -a <application_name> 
  % logiqctl logs -a <application_name> -p <proc_id_name> 
  % logiqctl logs -a <application_name> -p <proc_id_name> -g

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

		/* // the condition was removed because logs take similar
	       // argument as search
		if utils.FlagBegTime != "" ||
			utils.FlagEndTime != "" ||
			utils.FlagLogsSince != "" {
			err := errors.New("Invalid arguments specified -b, -e, or -s is used\n     Default logs dump period is set at 1 hour from current\n")
			utils.HandleError(err)
			return
		} */


		if utils.FlagEnablePsmod {
			loglerpart.CheckPsmod()
			loglerpart.Init(currentReleaseVersion)
		}

	// use search gut

		fmt.Println("   BegTime: ", utils.FlagBegTime)
		fmt.Println("   EndTime: ", utils.FlagEndTime)
		fmt.Println("     Since: ", utils.FlagLogsSince)

		if utils.FlagEnablePsmod {
			loglerpart.CheckPsmod()
			loglerpart.Init(currentReleaseVersion)
		}

		var hasApp bool
		var hasMultipleApps bool
		var hasProc bool
		var applicationV2s []*applications2.ApplicationV2

		services.FindApps(&hasApp, &hasMultipleApps, &hasProc, &applicationV2s)

		// output to file if any
		go services.WriteFile()

		// if !utils.FlagEnableSerial
		if false {

			// parallelize here
			services.ParallelSearch(cmd.UsageString(), "",
				&hasApp, &hasMultipleApps, &hasProc, &applicationV2s)
		} else {
			services.SerialSearch(cmd.UsageString(), "",
				&hasApp, &hasMultipleApps, &hasProc, &applicationV2s)
		}

		/*
		hasApp := false
		hasProc := false
		if utils.FlagProcId != "" {
			hasProc = true
		}
		if utils.FlagAppName != "" {
			hasApp = true
		}
		if hasApp && hasProc {
			proc, err := services.GetProcessByApplicationAndProc(utils.FlagAppName, utils.FlagProcId)
			utils.HandleError(err)
			services.DoQuery(-1, utils.FlagAppName, "", proc.ProcID, proc.LastSeen)
			// return
		} else if hasApp {
			app, err := services.GetApplicationByName(utils.FlagAppName)
			utils.HandleError(err)
			services.DoQuery(-1, utils.FlagAppName, "", "", app.LastSeen)
		} else {
			fmt.Println(cmd.UsageString())
			return
		}

		*/

		if utils.FlagEnablePsmod {
			loglerpart.DumpCurrentPsStat("ps_stat")
		}
	},
}

var interactiveCmd = &cobra.Command{
	Use:     "interactive",
	Aliases: []string{"i"},
	Short:   `Runs an interactive prompt to display logs.`,
	Run: func(cmd *cobra.Command, args []string) {
		if utils.FlagEnablePsmod {
			loglerpart.CheckPsmod()
			loglerpart.Init(currentReleaseVersion)
		}
		app, err := services.RunSelectApplicationForNamespacePrompt(false)
		utils.HandleError(err)
		proc, err := services.RunSelectProcessesForNamespaceAndAppPrompt(app.Name, false)
		utils.HandleError(err)
		fmt.Printf("You could also run this directly `logiqctl logs -p=%s %s`\n", proc.ProcID, app.Name)
		fmt.Printf("Fetching logs for %s (namespace), %s (application) and %s (process)\n\n", utils.GetDefaultNamespace(), app.Name, proc.ProcID)
		services.DoQuery(-1, app.Name, "", proc.ProcID, proc.LastSeen)
		if utils.FlagEnablePsmod {
			loglerpart.DumpCurrentPsStat("ps_stat")
		}
	},
}

var searchExample = `
logiqctl logs search supports many time range options
  - RFC3339 and epoch timestamp formats support automatically
  - Time format in format "yyyy-MM-dd hh:mm:ss.sssss +zzzz"
  - Suffix "+zzzz" will default to UTC-to-Localtime offset
    for example, 0700 is PDT and 0000 is UTC 
    One can use option --xutc (-x) to force UTC without specifying "+zzzz"
  - Different time search range options
    * --begtime (-b) --endtime (-e) => begtime, endtime
    * --begtime (-b) and --since (-s) => begtime, begtime + duration 
    * --endtime (-e) and --since (-s) => endtime - duration, endtime
    * Single duration --since (-s) => now() - duration, now()
    * Durations --since (-s) examples are 1m, 1d, 1s, etc., default=1h
logiqctl logs search supports search into multiple applications using the same -a option
    * -a <app1>,<app2>,<app3>,...
logiqctl advanced search supports nested expression and regex search 
    * -r 

Examples:
  % logiqctl -a app1,app2,app3 -p pid134 logs search "https"
  %	logiqctl -a app2 logs search "https" -b "2021-07-04 23:30:00.1234 0000" -s 5m
  %	logiqctl -a app3 logs search "error" -b "2021-07-04 23:30:00.1234" -e "2021-07-04 23:35:00.1234"
  %	logiqctl -a app3 logs search "message =~ 'code|piping' -r -b "2021-07-04 23:30:00.1234" -e "2021-07-04 23:35:00.1234"
  %	logiqctl -a app3 logs search "message =~ 'code' || message =~ 'piping'" -r -b "2021-07-04 23:30:00.1234" -e "2021-07-04 23:35:00.1234"
`

var searchCmd = &cobra.Command{
	Use:     "search [SearchString]",
	Aliases: []string{"s"},
	Example: searchExample,
	Short:   `Search logs for specific keywords or terms.`,
	Long:    `Search for specific keywords or terms in logs within a namespace, app, proc`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("SearchArgs: ", args[0])
		fmt.Println("   BegTime: ", utils.FlagBegTime)
		fmt.Println("   EndTime: ", utils.FlagEndTime)
		fmt.Println("     Since: ", utils.FlagLogsSince)

		if utils.FlagEnablePsmod {
			loglerpart.CheckPsmod()
			loglerpart.Init(currentReleaseVersion)
		}

		if len(args) != 1 {
			fmt.Println(cmd.Usage())
			return
		}

		var hasApp bool
		var hasMultipleApps bool
		var hasProc bool
		var applicationV2s []*applications2.ApplicationV2

		services.FindApps(&hasApp, &hasMultipleApps, &hasProc, &applicationV2s)

		// output to file if any
		go services.WriteFile()

		if !utils.FlagEnableSerial {

			// parallelize here
			services.ParallelSearch(cmd.UsageString(), args[0],
				&hasApp, &hasMultipleApps, &hasProc, &applicationV2s)
		} else {
			services.SerialSearch(cmd.UsageString(), args[0],
				&hasApp, &hasMultipleApps, &hasProc, &applicationV2s)
		}

		if utils.FlagEnablePsmod {
			loglerpart.DumpCurrentPsStat("ps_stat")
		}
	},
}




func init() {
	logsCmd.PersistentFlags().StringVarP(&utils.FlagLogsSince, "since", "s", "",
		`Only return logs newer than a relative duration. This is in relative to the last
seen log time for a specified application or processes within the namespace.
A duration string is a possibly signed sequence of decimal numbers, each with optional
fraction and a unit suffix, such as "3h34m", "1.5h" or "24h". Valid time units are "s", "m", "h"`)
	logsCmd.PersistentFlags().Uint32Var(&utils.FlagLogsPageSize, "page-size", 30, `Number of log entries to return in one page`)
	logsCmd.PersistentFlags().BoolVarP(&utils.FlagLogsFollow, "follow", "f", false, `Specify if the logs should be streamed.`)
	logsCmd.PersistentFlags().StringVarP(&utils.FlagProcId, "process", "p", "", `Filter logs by  proc id`)
	logsCmd.PersistentFlags().StringVarP(&utils.FlagAppName, "application", "a", "", `Filter logs by application`)
	logsCmd.PersistentFlags().StringVarP(&utils.FlagBegTime, "begtime", "b", "",
		`Search begin time range format "yyyy-MM-dd hh:mm:ss +0000". 
"+0000" suffix is required for search using UTC time.  
Localtime time search is assumed WITHOUT specifying "+0000."`)
	logsCmd.PersistentFlags().StringVarP(&utils.FlagEndTime, "endtime", "e", "",
		`Search end time range format "yyyy-MM-dd hh:mm:ss +0000". 
"+0000" suffix is required for search using UTC time.  
Localtime time search is assumed WITHOUT specifying "+0000."`)
	logsCmd.PersistentFlags().BoolVarP(&utils.FlagRegex, "advanced", "r", false, `Advanced expression and regex search`)
	logsCmd.PersistentFlags().BoolVarP(&utils.FlagSubSecond, "xutc", "x", false, `Force UTC date-time`)
	logsCmd.PersistentFlags().BoolVarP(&utils.FlagEnablePsmod, "psmod", "g", false, `Enable pattern signature generation module`)
	rootCmd.AddCommand(logsCmd)
	logsCmd.AddCommand(interactiveCmd)
	logsCmd.AddCommand(searchCmd)
	logsCmd.PersistentFlags().StringVarP(&utils.FlagFile, "write-to-file", "w", "", "Path to file")
	logsCmd.PersistentFlags().IntVarP(&utils.FlagMaxLogLines, "max-log-lines", "m", 200000, "Max log record output size")
	logsCmd.PersistentFlags().IntVarP(&utils.FlagParPeriod, "if1", "", 3, "Internal flag #1")
	logsCmd.PersistentFlags().IntVarP(&utils.FlagParCopies, "if2", "", 5, "Internal flag #2")
	logsCmd.PersistentFlags().BoolVarP(&utils.FlagEnableSerial, "if3", "", false, `Internal flag #3`)
}
