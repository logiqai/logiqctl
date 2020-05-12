package services

import (
	"context"

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
	return client.GetNamespaces(context.Background(), &namespace.NamespaceRequest{})
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

func RunSelectNamespacePrompt() (string, error) {
	namespaces, err := GetNamespacesAsStrings()
	if err != nil {
		return "", err
	}
	whatPrompt := promptui.Select{
		Label: "Select a namespace to set as default context",
		Items: namespaces,
	}
	what, _, err := whatPrompt.Run()
	if err != nil {
		return "", err
	}
	return namespaces[what], nil
}
