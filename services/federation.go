package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/logiqai/logiqctl/api/v1/federation"
	"github.com/logiqai/logiqctl/utils"
	"google.golang.org/grpc"
)

func printError() {
	fmt.Printf("Unable to PING cluster endpoint %s, please verify if the endpoint is accessible.\n", utils.GetClusterUrl())
}

func Ping() error {
	conn, err := grpc.Dial(utils.GetClusterUrl(), grpc.WithInsecure())
	if err != nil {
		printError()
		return err
	}
	client := federation.NewFederationClient(conn)
	p, err := client.Ping(context.Background(), &empty.Empty{})
	if err != nil {
		printError()
		return err
	}
	if p.Data != "PONG" {
		printError()
		return errors.New("")
	}
	fmt.Printf("Endpoint %s configured succesfully.\n\n", utils.GetClusterUrl())
	return nil
}
