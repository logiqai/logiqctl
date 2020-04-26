package services

import (
	"context"
	"fmt"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/tatsushid/go-prettytable"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/logiqai/logiqctl/api/v1/applications"
	"github.com/logiqai/logiqctl/cfg"
	"google.golang.org/grpc"
)

func GetApplications(config *cfg.Config, listNamespaces bool, namespaces []string) {
	namespaceMap := map[string]bool{}
	if listNamespaces {
		for _, k := range namespaces {
			namespaceMap[k] = true
		}
	}
	conn, err := grpc.Dial(config.Cluster, grpc.WithInsecure())
	if err != nil {
		handleError(config, err)
		return
	}
	defer conn.Close()
	client := applications.NewApplicationsServiceClient(conn)
	response, err := client.GetApplications(context.Background(), &empty.Empty{})
	if err != nil {
		handleError(config, err)
		return
	}

	if response != nil && response.Response != nil && response.Response.ApplicationsList != nil {
		tbl, err := prettytable.NewTable([]prettytable.Column{
			{Header: "Namespace"},
			{Header: "Application"},
			{Header: "ProcId"},
			{Header: "Last Seen"},
			{Header: "First Seen"},
		}...)
		if err != nil {
			panic(err)
		}
		tbl.Separator = " | "
		for _, app := range response.Response.ApplicationsList {
			if listNamespaces {
				if _, ok := namespaceMap[app.Namespace]; ok {
					fs := time.Unix(app.FirstSeen, 0)
					ls := time.Unix(app.LastSeen, 0)
					tbl.AddRow(app.Namespace, app.Name, app.Procid, humanize.Time(ls), humanize.Time(fs))
				}
			} else {
				fs := time.Unix(app.FirstSeen, 0)
				ls := time.Unix(app.LastSeen, 0)
				tbl.AddRow(app.Namespace, app.Name, app.Procid, humanize.Time(ls), humanize.Time(fs))
			}
		}
		tbl.Print()
	}
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	if h == 0 {
		return fmt.Sprintf("%02dm ago", m)
	}
	return fmt.Sprintf("%02dh:%02dm ago", h, m)
}
