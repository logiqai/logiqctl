package services

import (
	"context"
	"fmt"

	"github.com/tatsushid/go-prettytable"

	"github.com/logiqai/logiqctl/api/v1/events"
	"github.com/logiqai/logiqctl/utils"
	"google.golang.org/grpc"
)

func GetEvents(applicationName, process string) error {
	conn, err := grpc.Dial(utils.GetClusterUrl(), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	in := &events.EventsPostRequest{
		Count: 30,
	}

	inFilter := &events.EventsFilter{
		Namespace: []string{utils.GetDefaultNamespace()},
	}
	if applicationName != "" {
		inFilter.AppName = []string{applicationName}
	}
	if process != "" {
		inFilter.ProcId = []string{process}
	}
	in.Filter = inFilter

	client := events.NewEventsServiceClient(conn)
	events, err := client.GetEvents(context.Background(), in)
	if err != nil {
		return err
	}
	if len(events.Events) > 0 {
		if !utils.PrintResponse(events) {
			tbl, err := prettytable.NewTable([]prettytable.Column{
				{Header: "Name"},
				{Header: "Event Time"},
				{Header: "Level"},
				{Header: "Namespace"},
				{Header: "Application"},
				{Header: "Process"},
				{Header: "Message"},
			}...)
			if err != nil {
				panic(err)
			}
			tbl.Separator = " | "
			for _, e := range events.Events {
				tbl.AddRow(e.Name, utils.GetTimeAsString(e.TimestampInt), e.Level, e.Namespace, e.AppName, e.Sender, e.Message)
			}
			tbl.Print()
		}

	} else {
		display := fmt.Sprintf("No events found for %s (namespace)", utils.GetDefaultNamespace())
		if applicationName != "" {
			display = fmt.Sprintf("%s, %s (application)", display, applicationName)
		}
		if process != "" {
			display = fmt.Sprintf("%s, and %s (process)", display, process)
		}
		fmt.Println(display)
	}
	return nil
}
