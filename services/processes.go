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

package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/manifoldco/promptui"

	"github.com/tatsushid/go-prettytable"

	"github.com/logiqai/logiqctl/api/v1/processes"
	"github.com/logiqai/logiqctl/utils"
	"google.golang.org/grpc"
)

func getProcessesResponse(application string) (*processes.ProcessesResponse, error) {
	conn, err := grpc.Dial(utils.GetClusterUrl(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := processes.NewProcessDetailsServiceClient(conn)
	return client.GetProcesses(context.Background(), &processes.ProcessesRequest{
		Namespace:       utils.GetDefaultNamespace(),
		ApplicationName: application,
	})
}

func ListProcesses() {
	application, err := RunSelectApplicationForNamespacePrompt(false)
	if err != nil {
		return
	}
	response, err := getProcessesResponse(application.Name)
	if err != nil {
		return
	}
	printProcessesResponse(response.Processes)
}

func printProcessesResponse(response []*processes.Process) {
	fmt.Println()
	if len(response) > 0 {
		if !utils.PrintResponse(response) {
			tbl, err := prettytable.NewTable([]prettytable.Column{
				{Header: "Namespace"},
				{Header: "ProcID"},
				{Header: "Last Seen"},
				{Header: "First Seen"},
			}...)
			if err != nil {
				panic(err)
			}
			tbl.Separator = " | "
			for _, app := range response {
				tbl.AddRow(utils.GetDefaultNamespace(), app.ProcID, utils.GetTimeAsString(app.LastSeen), utils.GetTimeAsString(app.FirstSeen))
			}
			tbl.Print()
		}
	}
}

func GetProcessByApplicationAndProc(application, procId string) (*processes.Process, error) {
	response, err := getProcessesResponse(application)
	if err != nil {
		return nil, err
	}
	if len(response.Processes) > 0 {
		for _, proc := range response.Processes {
			if proc.ProcID == procId {
				return proc, nil
			}
		}
	}
	return nil, errors.New("NOT_FOUND")
}

func RunSelectProcessesForNamespaceAndAppPrompt(application string, all bool) (*processes.Process, error) {
	response, err := getProcessesResponse(application)
	if err != nil {
		return nil, err
	}
	//if len(response.Processes) == 1 {
	//	return response.Processes[0], nil
	//}
	if len(response.Processes) > 0 {
		//TODO clean this up
		var procForDisplay []SelectDisplay
		for _, proc := range response.Processes {
			procForDisplay = append(procForDisplay, SelectDisplay{
				Name:    proc.ProcID,
				Details: fmt.Sprintf("Last Seen %s", utils.GetTimeAsString(proc.LastSeen)),
			})
		}
		if all {
			procForDisplay = append([]SelectDisplay{
				{
					Name:    "*",
					Details: "Select All",
				},
			}, procForDisplay...)
		}

		whatPrompt := promptui.Select{
			Label:     fmt.Sprintf("Select an process (showing '%s' namespace and '%s' application", utils.GetDefaultNamespace(), application),
			Items:     procForDisplay,
			Templates: GetTemplateForType(templateTypeProcess),
			Size:      6,
		}

		what, _, err := whatPrompt.Run()
		if err != nil {
			return nil, err
		}
		if all && what == 0 {
			//select all
			return nil, nil
		}
		if all {
			return response.Processes[what-1], nil
		}
		return response.Processes[what], nil
	}
	return nil, errors.New("NOT_FOUND")
}
