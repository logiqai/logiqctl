package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/manifoldco/promptui"

	"github.com/logiqai/logiqctl/utils"

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

func printResponse(response []*applications.ApplicationV2) {
	fmt.Println()
	if len(response) > 0 {
		tbl, err := prettytable.NewTable([]prettytable.Column{
			{Header: "Namespace"},
			{Header: "Application"},
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
			tbl.AddRow(app.Namespace, app.Name, humanize.Time(ls), humanize.Time(fs))

		}
		tbl.Print()
	}
}

func getApplicationsV2Response(all bool) (*applications.GetApplicationsResponseV2, error) {
	conn, err := grpc.Dial(utils.GetClusterUrl(), grpc.WithInsecure())
	if err != nil {
		//handleError(config, err)
		return nil, err
	}
	defer conn.Close()
	client := applications.NewApplicationsServiceClient(conn)
	request := &applications.GetApplicationsRequest{
		Page: 0,
		Size: 0,
	}
	if !all {
		request.Namespace = utils.GetDefaultNamespace()
	}
	return client.GetApplicationsV2(context.Background(), request)
}

func GetApplicationsV2(all bool) {
	response, err := getApplicationsV2Response(all)
	if err != nil {
		//handleError(config, err)
		return
	}
	printResponse(response.Applications)

}

func RunSelectApplicationForNamespacePrompt() (string, error) {
	response, err := getApplicationsV2Response(false)
	if err != nil {
		//handleError(config, err)
		return "", err
	}
	if len(response.Applications) > 0 {
		var apps []struct {
			Name    string
			Details string
		}
		for _, app := range response.Applications {
			ls := time.Unix(app.LastSeen, 0)
			apps = append(apps, struct {
				Name    string
				Details string
			}{
				Name:    app.Name,
				Details: fmt.Sprintf("Last Seen %s", humanize.Time(ls)),
			})
		}

		templates := &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "\U000000BB {{ .Name | green }} ({{ .Details | red }})",
			Inactive: "  {{ .Name | cyan }} ({{ .Details | red }})",
			Selected: "Listing processes for {{ .Name | red }}",
		}

		whatPrompt := promptui.Select{
			Label:     fmt.Sprintf("Select an application to see its processes (showing '%s' namespace)", utils.GetDefaultNamespace()),
			Items:     apps,
			Templates: templates,
			Size:      6,
		}

		what, _, err := whatPrompt.Run()
		if err != nil {
			return "", err
		}
		return apps[what].Name, nil
	}
	return "", errors.New("Cannot find processes ")
}
