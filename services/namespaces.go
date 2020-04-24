package services

import (
	"context"

	"github.com/logiqai/logiqctl/cfg"

	"github.com/logiqai/logiqctl/api/v1/namespace"
	"google.golang.org/grpc"
)

func GetNamespaces() []*namespace.Namespace {
	config := cfg.CONFIG
	conn, err := grpc.Dial(config.Cluster, grpc.WithInsecure())
	if err != nil {
		handleError(config, err)
		return nil
	}
	defer conn.Close()
	client := namespace.NewNamespaceServiceClient(conn)
	ctx := context.Background()
	response, err := client.GetNamespaces(ctx, &namespace.NamespaceRequest{})
	if err != nil {
		handleError(config, err)
		return nil
	}
	return response.Namespaces

}
