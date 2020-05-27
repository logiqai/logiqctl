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

	"github.com/logiqai/logiqctl/utils"

	"github.com/tatsushid/go-prettytable"

	"github.com/logiqai/logiqctl/api/v1/applications"
	"google.golang.org/grpc"
)

func printResponse(response []*applications.ApplicationV2) {
	fmt.Println()
	if len(response) > 0 {
		if !utils.PrintResponse(response) {
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
				tbl.AddRow(app.Namespace, app.Name, utils.GetTimeAsString(app.LastSeen), utils.GetTimeAsString(app.FirstSeen))
			}
			tbl.Print()
		}
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

func GetApplicationByName(application string) (*applications.ApplicationV2, error) {
	response, err := getApplicationsV2Response(false)
	if err != nil {
		return nil, err
	}
	if len(response.Applications) > 0 {
		for _, app := range response.Applications {
			if app.Name == application {
				return app, nil
			}
		}
	}
	return nil, errors.New("Cannot find application ")
}

func RunSelectApplicationForNamespacePrompt(all bool) (*applications.ApplicationV2, error) {
	response, err := getApplicationsV2Response(false)
	if err != nil {
		return nil, err
	}
	if len(response.Applications) > 0 {
		var apps []SelectDisplay
		for _, app := range response.Applications {
			apps = append(apps, SelectDisplay{
				Name:    app.Name,
				Details: fmt.Sprintf("Last Seen %s", utils.GetTimeAsString(app.LastSeen)),
			})
		}
		if all {
			apps = append([]SelectDisplay{
				{
					Name:    "*",
					Details: "Select All",
				},
			}, apps...)
		}

		whatPrompt := promptui.Select{
			Label:     fmt.Sprintf("Select an application (showing '%s' namespace)", utils.GetDefaultNamespace()),
			Items:     apps,
			Templates: GetTemplateForType(templateTypeApplication),
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
			return response.Applications[what-1], nil
		}
		return response.Applications[what], nil
	}
	return nil, errors.New("Cannot find application ")
}
