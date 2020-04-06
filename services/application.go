package services

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/logiqai/logiqbox/api/v1/applications"
	"github.com/logiqai/logiqbox/cfg"
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
		fmt.Printf("%-16s| %-16s | %-16s\n", "Namespace", "Application", "ProcId")
		for _, app := range response.Response.ApplicationsList {
			if listNamespaces {
				if _, ok := namespaceMap[app.Namespace]; ok {
					fmt.Printf("%-16s | %-16s | %-16s\n", app.Namespace, app.Name, app.Procid)
				}
			} else {
				fmt.Printf("%-16s | %-16s | %-16s\n", app.Namespace, app.Name, app.Procid)
			}
		}
	}
}
