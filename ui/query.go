package ui

import (
	"encoding/json"
	"github.com/logiqai/logiqctl/utils"
	"github.com/spf13/cobra"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"github.com/olekukonko/tablewriter"
)

func NewListQueriesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query",
		Example: "logiqctl get query|q <query-id>",
		Aliases: []string{"q"},
		Short:   "Get a query",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Missing query id")
				os.Exit(-1)
			}
			getQuery(args)
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:     "all",
		Example: "logiqctl get query all",
		Short:   "List all the available queries",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			listQueries()
		},
	})

	return cmd
}

func getQuery(args []string) {
	uri := getUrlForResource(ResourceQuery,args...)
	client := getHttpClient()

	if resp, err := client.Get(uri); err == nil {
		defer resp.Body.Close()
		var v = map[string]interface{}{}
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Unable to fetch queries, Error: %s", err.Error())
				os.Exit(-1)
			}
			err = json.Unmarshal(bodyBytes, &v)
			if err != nil {
				fmt.Printf("Unable to decode queries, Error: %s", err.Error())
			} else {
				fmt.Println(string(bodyBytes))
			}
		}
	} else {
		fmt.Printf("Unable to fetch queries, Error: %s", err.Error())
		os.Exit(-1)
	}
}

func listQueries() {
	uri := getUrlForResource(ResourceQueryAll)
	client := getHttpClient()

	if resp, err := client.Get(uri); err == nil {
		defer resp.Body.Close()
		var v = map[string]interface{}{}
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Unable to fetch queries, Error: %s", err.Error())
				os.Exit(-1)
			}
			err = json.Unmarshal(bodyBytes, &v)
			if err != nil {
				fmt.Printf("Unable to decode quesries, Error: %s", err.Error())
			} else {
				count := (int)(v["count"].(float64))
				queries := v["results"].([]interface{})
				fmt.Println("(",count, ") queries found")
				if ( count > 0 ) {
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"Name", "Data source Id", "Id","Archived","Draft"})
					for _, q := range queries {
						query := q.(map[string]interface{})
						dataSourceId := (int)(query["data_source_id"].(float64))
						name := query["name"].(string)
						isArchived := query["is_archived"].(bool)
						isDraft := query["is_draft"].(bool)
						id := (int)(query["id"].(float64))
						table.Append([]string{name, fmt.Sprintf("%d",dataSourceId),
							fmt.Sprintf("%d",id), fmt.Sprintf("%v",isArchived),
							fmt.Sprintf("%v",isDraft),
						})
					}

					table.Render()
				}
			}
		}
	} else {
		fmt.Printf("Unable to fetch queries, Error: %s", err.Error())
		os.Exit(-1)
	}
}
