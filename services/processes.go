package services

import (
	"context"
	"time"

	"github.com/logiqai/logiqctl/api/v1/processes"
	"github.com/logiqai/logiqctl/cfg"
	"google.golang.org/grpc"
)

func GetProcesses(ns, appName string) []*processes.Process {
	config := cfg.CONFIG
	conn, err := grpc.Dial(config.Cluster, grpc.WithInsecure())
	if err != nil {
		handleError(config, err)
		return nil
	}
	defer conn.Close()
	client := processes.NewProcessDetailsServiceClient(conn)
	ctx := context.Background()
	response, err := client.GetProcesses(ctx, &processes.ProcessesRequest{
		Namespace:       ns,
		ApplicationName: appName,
		Page:            0,
		Size:            0,
	})
	if err != nil {
		handleError(config, err)
		return nil
	}
	return response.Processes

}

func GetProcessesMock(ns, appName string) []*processes.Process {
	time.Sleep(3 * time.Second)
	return []*processes.Process{
		{
			ProcID:    "45682",
			LastSeen:  time.Now().Unix() - 60,
			FirstSeen: time.Now().Unix() - 1000,
		}, {
			ProcID:    "45688",
			LastSeen:  time.Now().Unix() - 60,
			FirstSeen: time.Now().Unix() - 1000,
		},
	}
}
