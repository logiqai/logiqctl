package services

import (
	"context"
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
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
			fs := time.Unix(app.FirstSeen, 0)
			ls := time.Unix(app.LastSeen, 0)
			tbl.AddRow(utils.GetDefaultNamespace(), app.ProcID, humanize.Time(ls), humanize.Time(fs))

		}
		tbl.Print()
	}
}
