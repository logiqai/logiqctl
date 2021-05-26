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
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/logiqai/logiqctl/grpc_utils"

	"github.com/logiqai/logiqctl/api/v1/query"
	"github.com/logiqai/logiqctl/utils"
	"google.golang.org/grpc"
)

const (
	OUTPUT_COLUMNS = "column"
	OUTPUT_RAW     = "raw"
	OUTPUT_JSON    = "json"
)

func postQuery(applicationName, searchTerm, procId string, lastSeen int64) (string, query.QueryServiceClient, error) {
	conn, err := grpc.Dial(utils.GetClusterUrl(), grpc.WithInsecure())
	if err != nil {
		return "", nil, err
	}
	client := query.NewQueryServiceClient(conn)

	in := &query.QueryProperties{
		Namespace: utils.GetDefaultNamespace(),
		PageSize:  utils.GetPageSize(),
		QType:     query.QueryType_Fetch,
	}
	var st time.Time
	if lastSeen > 0 {
		st = utils.GetStartTime(lastSeen)
	} else {
		st = time.Now().UTC()
	}
	in.StartTime = st.Format(time.RFC3339)
	if utils.FlagLogsSince != "" {
		d, err := time.ParseDuration(utils.FlagLogsSince)
		if err != nil {
			fmt.Printf("Unable to parse duration %s\n", utils.FlagLogsSince)
			os.Exit(1)
		}
		in.EndTime = st.Add(-1 * d).Format(time.RFC3339)
	}

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
		in.QType = query.QueryType_Search
	}

	queryResponse, err := client.Query(grpc_utils.GetGrpcContext(), in)
	if err != nil {
		return "", nil, err
	}
	return queryResponse.QueryId, client, nil
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("Error Occured: %s", err.Error())
		os.Exit(-1)
	}
}

func DoQuery(appName, searchTerm, procId string, lastSeen int64) {
	queryId, client, err := postQuery(appName, searchTerm, procId, lastSeen)
	handleError(err)
	if queryId != "" {
		var f *os.File
		var writeToFile bool
		if utils.FlagFile != "" {
			once.Do(func() {
				writeToFile = true
				if fTmp, err := os.OpenFile(utils.FlagFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600); err != nil {
					fmt.Printf("1 Unable to write to file: %s \n", err.Error())
					os.Exit(1)
				} else {
					f = fTmp
					fmt.Printf("Writing output to %s\n", utils.FlagFile)
				}
			})
			defer f.Close()
		}
		for {
			response, err := client.GetDataNext(grpc_utils.GetGrpcContext(), &query.GetDataRequest{
				QueryId: queryId,
			})
			if err != nil {
				handleError(err)
			}
			if len(response.Data) > 0 {
				for _, entry := range response.Data {
					if writeToFile {
						line := fmt.Sprintf("%s %s %s %s - %s",
							entry.Timestamp,
							entry.SeverityString,
							entry.FacilityString,
							entry.ProcID,
							entry.Message,
						)
						if strings.HasSuffix(line, "\n") {
							line = strings.ReplaceAll(line, "\n", "")
						}
						line = fmt.Sprintf("%s\n", line)
						if _, err := f.WriteString(line); err != nil {
							fmt.Printf("Cannot write file: %s\n", err.Error())
							os.Exit(1)
						}
						if stat, err := os.Stat(utils.FlagFile); err == nil {
							if stat.Size() > int64(utils.FlagMaxFileSize*1048576) {
								fmt.Printf("Max file size reached. Control file size using -m\n")
								os.Exit(1)
							}
						}
					} else {
						PrintSyslogMessageForType(entry, "raw")
						time.Sleep(20 * time.Millisecond)
					}
				}
			} else {
				if response.Remaining <= 0 && response.Status == "COMPLETE" {
					os.Exit(0)
				}
				time.Sleep(2 * time.Second)
			}
		}
	}
}
