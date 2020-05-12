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
	"fmt"

	"github.com/tatsushid/go-prettytable"

	"github.com/logiqai/logiqctl/api/v1/processes"
	"github.com/logiqai/logiqctl/utils"
	"google.golang.org/grpc"
)

func ListProcesses() {
	conn, err := grpc.Dial(utils.GetClusterUrl(), grpc.WithInsecure())
	if err != nil {
		return
	}

	client := processes.NewProcessDetailsServiceClient(conn)
	application, err := RunSelectApplicationForNamespacePrompt()
	if err != nil {
		return
	}
	response, err := client.GetProcesses(context.Background(), &processes.ProcessesRequest{
		Namespace:       utils.GetDefaultNamespace(),
		ApplicationName: application,
	})
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
