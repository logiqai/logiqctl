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
	"time"

	"github.com/logiqai/logiqctl/api/v1/query"
	"github.com/logiqai/logiqctl/utils"
	"google.golang.org/grpc"
)

const (
	OUTPUT_COLUMNS = "column"
	OUTPUT_RAW     = "raw"
	OUTPUT_JSON    = "json"
)

func Query(applicationName, searchTerm, procId string, lastSeen int64) (string, error) {
	conn, err := grpc.Dial(utils.GetClusterUrl(), grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	client := query.NewQueryServiceClient(conn)

	in := &query.QueryProperties{
		Namespace: utils.GetDefaultNamespace(),
		PageSize:  utils.GetPageSize(),
	}

	in.StartTime = utils.GetStartTime(lastSeen).Format(time.RFC3339)

	if applicationName != "" {
		in.ApplicationNames = []string{applicationName}
		//  process filter is only available for query
		if procId != "" {
			filterValuesMap := make(map[string]*query.FilterValues)
			filterValuesMap["ProcId"] = &query.FilterValues{
				Values: []string{procId},
			}
			in.Filters = filterValuesMap
		}

	}

	if searchTerm != "" {
		in.KeyWord = searchTerm
	}

	queryResponse, err := client.Query(context.Background(), in)
	if err != nil {
		return "", err
	}
	var hasData bool
	for {
		if !hasData {
			time.Sleep(time.Second)
		}
		dataResponse, err := client.GetDataNext(context.Background(), &query.GetDataRequest{
			QueryId: queryResponse.QueryId,
		})
		if err != nil {
			return "", err
		}
		if len(dataResponse.Data) > 0 {
			hasData = true
		}
		for _, entry := range dataResponse.Data {
			if utils.FlagOut == "json" {
				printSyslogMessageForType(entry, OUTPUT_JSON)
			} else {
				printSyslogMessageForType(entry, OUTPUT_RAW)
			}
			time.Sleep(20 * time.Millisecond)
		}
		if dataResponse.Remaining <= 0 && dataResponse.Status == "COMPLETE" {
			return "", err
		}
		if !utils.FlagLogsFollow {
			return queryResponse.QueryId, err
		}
	}
	return "", nil
}

func GetDataNext(queryId string) (bool, error) {
	conn, err := grpc.Dial(utils.GetClusterUrl(), grpc.WithInsecure())
	if err != nil {
		return false, err
	}
	client := query.NewQueryServiceClient(conn)
	var hasData bool
	for {
		if !hasData {
			time.Sleep(time.Second)
		}
		dataResponse, err := client.GetDataNext(context.Background(), &query.GetDataRequest{
			QueryId: queryId,
		})
		if err != nil {
			return false, err
		}
		if len(dataResponse.Data) > 0 {
			hasData = true
		}
		for _, entry := range dataResponse.Data {
			printSyslogMessageForType(entry, "raw")
			time.Sleep(20 * time.Millisecond)
		}
		if dataResponse.Remaining <= 0 && dataResponse.Status == "COMPLETE" {
			return false, err
		}
		if !utils.FlagLogsFollow {
			return true, err
		}
	}
}
