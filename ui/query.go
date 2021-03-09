package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/logiqai/easymap"
	"github.com/logiqai/logiqctl/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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
			printQuery(args)
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:     "result",
		Example: "logiqctl get query result|q <query-result-id>",
		Aliases: []string{"q"},
		Short:   "Get a query result",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Missing query result id")
				os.Exit(-1)
			}
			printQueryResult(args)
		},
	})
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

func createQuery(qyerySpec map[string]interface{}) (map[string]interface{}, error) {
	uri := GetUrlForResource(ResourceQueryAll)
	client := getHttpClient()

	if payloadBytes, jsonMarshallError := json.Marshal(qyerySpec); jsonMarshallError != nil {
		return nil, jsonMarshallError
	} else {
		req, err := http.NewRequest("POST",uri,bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Println("Unable to create query ", err.Error())
			os.Exit(-1)
		}
		if api_key := viper.GetString(utils.AuthToken); api_key != "" {
			req.Header.Add("Authorization", fmt.Sprintf("Key %s", api_key))
		}
		if resp, err := client.Do(req); err == nil {
			jsonStr, _ := json.MarshalIndent(qyerySpec, "", "    ")
			fmt.Printf("Successfully created query : %s", jsonStr)
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("Unable to read create query response, Error: %s", err.Error())
			}
			respDict := map[string]interface{}{}
			if errUnmarshall := json.Unmarshal(bodyBytes, &respDict); errUnmarshall != nil {
				return nil, fmt.Errorf("Unable to decode create query response")
			}

			return respDict, nil
		} else {
			return nil, err
		}
	}
}

func printQuery(args []string) {
	if v, vErr := getQuery(args); v != nil {
		query := (*v)
		s, _ := json.MarshalIndent(*v, "", "    ")
		if utils.FlagOut == "json" {
			fmt.Println(string(s))
		} else if utils.FlagOut == "table" {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Data source Id", "Id", "Latest query result Id", "Archived", "Draft"})
			dataSourceId := (int)(query["data_source_id"].(float64))
			name := query["name"].(string)
			isArchived := query["is_archived"].(bool)
			isDraft := query["is_draft"].(bool)
			id := (int)(query["id"].(float64))
			lastRun := ""
			if lrId, lrIdOk := query["latest_query_data_id"]; lrIdOk {
				lastRun = fmt.Sprintf("%d",(int)(lrId.(float64)))
			}
			table.Append([]string{name, fmt.Sprintf("%d", dataSourceId),
				fmt.Sprintf("%d", id), fmt.Sprintf("%s",lastRun),
				fmt.Sprintf("%v", isArchived),
				fmt.Sprintf("%v", isDraft),
			})
			table.Render()
		} else if utils.FlagOut == "yaml" {
			a, yamlErr := yaml.Marshal(v)
			if yamlErr != nil {
				fmt.Errorf("Error converting to YAML")
			} else {
				fmt.Println(string(a))
			}
		}
	} else {
		fmt.Println(vErr.Error())
		os.Exit(-1)
	}
}

func getQueryResult(args []string) (*map[string]interface{}, error) {
	uri := GetUrlForResource(ResourceQueryResult, args...)
	client := getHttpClient()
	req, err := http.NewRequest("GET",uri, nil)
	if err != nil {
		fmt.Println("Unable to get query ", err.Error())
		os.Exit(-1)
	}
	if api_key := viper.GetString(utils.AuthToken); api_key != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Key %s", api_key))
	}

	if resp, err := client.Do(req); err == nil {
		defer resp.Body.Close()
		var v = map[string]interface{}{}
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("Unable to fetch queries, Error: %s", err.Error())
			}
			err = json.Unmarshal(bodyBytes, &v)
			if err != nil {
				return nil, fmt.Errorf("Unable to decode queries, Error: %s", err.Error())
			} else {
				return &v, nil
			}
		} else {
			return nil, fmt.Errorf("Http response error, Error: %d", resp.StatusCode)
		}
	} else {
		return nil, fmt.Errorf("Unable to fetch queries, Error: %s", err.Error())
	}
}

func printQueryResult(args []string) {
	if v, vErr := getQueryResult(args); v != nil {
		em := (easymap.EasyMap)(*v)
		rows := em.Get("rows")
		columns := em.Get("columns")
		s, _ := json.MarshalIndent(*v, "", "    ")
		if utils.FlagOut == "json" {
			fmt.Println(string(s))
		} else if utils.FlagOut == "table" {
			table := tablewriter.NewWriter(os.Stdout)
			th := []string{}
			for _, colDetails := range columns["query_result.data.columns"].([]interface{}) {
				c := colDetails.(map[string]interface{})
				th = append(th, c["friendly_name"].(string))
			}
			table.SetHeader(th)
			for _, rowDetails := range rows["query_result.data.rows"].([]interface{}) {
				r := rowDetails.(map[string]interface{})
				rvals := []string{}
				for _, key := range th {
					if strings.ToLower(key) == "timestamp" {
						if v, vOk := r[key].(float64); vOk {
							t := time.Unix((int64)(v/1000),0)
							rvals = append(rvals, t.String() )
						}
					} else {
						rvals = append(rvals, fmt.Sprintf("%v", r[key]))
					}
				}
				table.Append(rvals)
			}
			table.Render()
		} else if utils.FlagOut == "yaml" {
			a, yamlErr := yaml.Marshal(v)
			if yamlErr != nil {
				fmt.Errorf("Error converting to YAML")
			} else {
				fmt.Println(string(a))
			}
		}
	} else {
		fmt.Println(vErr.Error())
		os.Exit(-1)
	}
}

func getQueryByName(name string) map[string]interface{} {
	if v, err := getQueries(); err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	} else {
		queries := v["results"].([]interface{})
		for _, q := range queries {
			query := q.(map[string]interface{})
			if query["name"] == name {
				return query
			}
		}
	}
	return nil
}

func getQuery(args []string) (*map[string]interface{}, error) {
	uri := GetUrlForResource(ResourceQuery, args...)
	client := getHttpClient()
	req, err := http.NewRequest("GET",uri, nil)
	if err != nil {
		fmt.Println("Unable to get query ", err.Error())
		os.Exit(-1)
	}
	if api_key := viper.GetString(utils.AuthToken); api_key != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Key %s", api_key))
	}

	if resp, err := client.Do(req); err == nil {
		defer resp.Body.Close()
		var v = map[string]interface{}{}
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("Unable to fetch queries, Error: %s", err.Error())
			}
			err = json.Unmarshal(bodyBytes, &v)
			if err != nil {
				return nil, fmt.Errorf("Unable to decode queries, Error: %s", err.Error())
			} else {
				return &v, nil
			}
		} else {
			return nil, fmt.Errorf("Http response error, Error: %d", resp.StatusCode)
		}
	} else {
		return nil, fmt.Errorf("Unable to fetch queries, Error: %s", err.Error())
	}
}

func publishQuery(args []string) (*map[string]interface{}, error) {
	uri := GetUrlForResource(ResourceQuery, args...)
	client := getHttpClient()
	id, _ := strconv.Atoi(args[0])
	version, _ := strconv.Atoi(args[1])
	queryPublishSpec := map[string]interface{}{
		"is_draft": false,
		"id":       id,
		"version":  version,
	}

	if payloadBytes, jsonMarshallError := json.Marshal(queryPublishSpec); jsonMarshallError != nil {
		return nil, jsonMarshallError
	} else {
		req, err := http.NewRequest("POST",uri, bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Println("Unable to publish query ", err.Error())
			os.Exit(-1)
		}
		if api_key := viper.GetString(utils.AuthToken); api_key != "" {
			req.Header.Add("Authorization", fmt.Sprintf("Key %s", api_key))
		}
		if resp, err := client.Do(req); err == nil {
			defer resp.Body.Close()
			var v = map[string]interface{}{}
			if resp.StatusCode == http.StatusOK {
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("Unable to publish query, Read Error: %s", err.Error())
				}
				err = json.Unmarshal(bodyBytes, &v)
				if err != nil {
					return nil, fmt.Errorf("Unable to decode publish query response, Error: %s", err.Error())
				} else {
					return &v, nil
				}
			} else {
				return nil, fmt.Errorf("Http response error while publishing query, Error: %d", resp.StatusCode)
			}
		} else {
			return nil, fmt.Errorf("Unable to publish query, Error: %s", err.Error())
		}
	}
}

func listQueries() {
	if v, err := getQueries(); err == nil {
		count := (int)(v["count"].(float64))
		queries := v["results"].([]interface{})
		fmt.Println("(", count, ") queries found")
		if count > 0 {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Data source Id", "Id", "Archived", "Draft"})
			for _, q := range queries {
				query := q.(map[string]interface{})
				dataSourceId := (int)(query["data_source_id"].(float64))
				name := query["name"].(string)
				isArchived := query["is_archived"].(bool)
				isDraft := query["is_draft"].(bool)
				id := (int)(query["id"].(float64))
				table.Append([]string{name, fmt.Sprintf("%d", dataSourceId),
					fmt.Sprintf("%d", id), fmt.Sprintf("%v", isArchived),
					fmt.Sprintf("%v", isDraft),
				})
			}

			table.Render()
		}
	} else {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
}

func getQueries() (map[string]interface{}, error) {
	uri := GetUrlForResource(ResourceQueryAll)
	client := getHttpClient()
	req, err := http.NewRequest("GET",uri,nil)
	if err != nil {
		fmt.Println("Unable to get queries ", err.Error())
		os.Exit(-1)
	}
	if api_key := viper.GetString(utils.AuthToken); api_key != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Key %s", api_key))
	}

	if resp, err := client.Do(req); err == nil {
		defer resp.Body.Close()
		var v = map[string]interface{}{}
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("Unable to fetch queries, Error: %s", err.Error())
				os.Exit(-1)
			}
			err = json.Unmarshal(bodyBytes, &v)
			if err != nil {
				return nil, fmt.Errorf("Unable to decode quesries, Error: %s", err.Error())
			} else {
				return v, nil
			}
		} else {
			return nil, fmt.Errorf("Http response error, Error: %d", resp.StatusCode)
		}
	} else {
		return nil, fmt.Errorf("Unable to fetch queries, Error: %s", err.Error())
	}
}
