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
	"github.com/logiqai/logiqctl/grpc_utils"
	"github.com/manifoldco/promptui"

	"github.com/tatsushid/go-prettytable"

	"github.com/logiqai/logiqctl/api/v1/namespace"
	"github.com/logiqai/logiqctl/utils"
	"google.golang.org/grpc"
)

func getNamespaces() (*namespace.NamespaceResponse, error) {
	conn, err := grpc.Dial(utils.GetClusterUrl(), grpc.WithInsecure())
	if err != nil {
		//handleError(config, err)
		return nil, err
	}
	defer conn.Close()
	client := namespace.NewNamespaceServiceClient(conn)
	return client.GetNamespaces(grpc_utils.GetGrpcContext(), &namespace.NamespaceRequest{})
}

func GetNamespacesAsStrings() ([]string, error) {
	response, err := getNamespaces()
	if err != nil {
		//handleError(config, err)
		return nil, err
	}
	var namespaces []string
	for _, ns := range response.Namespaces {
		namespaces = append(namespaces, ns.Namespace)
	}
	return namespaces, nil
}

func ListNamespaces() {
	response, err := getNamespaces()
	if err != nil {
		//handleError(config, err)
		return
	}
	if response != nil && len(response.Namespaces) > 0 {
		if !utils.PrintResponse(response) {
			tbl, err := prettytable.NewTable([]prettytable.Column{
				{Header: "Namespace"},
				{Header: "Type"},
				{Header: "Last Seen"},
				{Header: "First Seen"},
			}...)
			if err != nil {
				panic(err)
			}
			tbl.Separator = " | "
			for _, ns := range response.Namespaces {
				readableType := "Namespace"
				if ns.Type == "H" {
					readableType = "Host"
				}
				tbl.AddRow(ns.Namespace, readableType, utils.GetTimeAsString(ns.LastSeen), utils.GetTimeAsString(ns.FirstSeen))
			}
			tbl.Print()
		}
	}
}

func RunSelectNamespacePrompt(all bool) (string, error) {
	namespaces, err := GetNamespacesAsStrings()
	if err != nil {
		return "", err
	}
	if all {
		namespaces = append(namespaces, "*")
	}
	whatPrompt := promptui.Select{
		Label: "Select a namespace",
		Items: namespaces,
	}
	what, _, err := whatPrompt.Run()
	if err != nil {
		return "", err
	}
	return namespaces[what], nil
}
