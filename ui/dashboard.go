package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"io/ioutil"
	"math"
	"os"
	"strconv"

	"github.com/spf13/viper"

	"github.com/logiqai/logiqctl/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"net/http"
)

func NewListDashboardsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dashboard",
		Example: "logiqctl get dashboard|d <dashboard-slug>",
		Aliases: []string{"d"},
		Short:   "Get a dashboard",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Missing dashboard slug")
				os.Exit(-1)
			}
			exportDashboard(args)
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:     "all",
		Example: "logiqctl get dashboard all",
		Short:   "List all the available dashboards",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			listDashboards()
		},
	})

	return cmd
}

func NewDashboardCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dashboard",
		Example: "logiqctl create dashboard|d -f <path to dashboard spec>",
		Aliases: []string{"d"},
		Short:   "Create a dashboard",
		Long: `
The crowd-sourced dashboards available in https://github.com/logiqai/logiqhub can be downloaded and applied to any clusters. 
One can also export dashboards created using "logiqctl get dashboard" command and apply on different clusters.
`,
		PreRun: utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			if utils.FlagFile == "" {
				fmt.Println("Missing dashboard spec file")
				os.Exit(-1)
			} else {
				fmt.Println("Dashboard spec file :", utils.FlagFile)
				fileBytes, err := ioutil.ReadFile(utils.FlagFile)
				if err != nil {
					fmt.Println("Unable to read file ", utils.FlagFile)
					os.Exit(-1)
				}
				spec := map[string]interface{}{}
				if err = json.Unmarshal(fileBytes, &spec); err != nil {
					fmt.Println("Unable to decode dashboard spec", utils.FlagFile)
					os.Exit(-1)
				}
				createAndPublishDashboardSpec(spec)
			}
		},
	}
	cmd.Flags().StringVarP(&utils.FlagFile, "file", "f", "", "Path to file")
	return cmd
}

func exportDashboard(args []string) {
	dashboardOut := map[string]interface{}{}

	if dashboardPtr, err := GetDashboard(args); err != nil {
		fmt.Println(err.Error())
	} else {
		dashboard := *dashboardPtr

		/* json dump exercise
		fmt.Println("dashboard=<", dashboard, ">")
		v, err := json.Marshal(dashboard)
		if (err!=nil) {
			fmt.Println("error=", err)
		} else {
			fmt.Println("no error found")
		}
		fmt.Println("v=<", string(v), ">")
		gg := map[string]interface{}{}
		_ = json.Unmarshal(v, &gg)
		fmt.Println("dashboard_json=<", gg, ">")
		*/

		dashboardParams := map[string]interface{}{}
		dashboardParams["name"] = dashboard["name"]
		dashboardParams["tags"] = dashboard["tags"]
		dashboardOut["dashboard"] = dashboardParams

		widgets := dashboard["widgets"].([]interface{})
		widgetOut := []interface{}{}
		dataSources := map[int]interface{}{}

		importWidget := true
		for _, w := range widgets {
			widget := w.(map[string]interface{})

			if _, ok := widget["visualization"]; ok {

				visualization := widget["visualization"].(map[string]interface{})
				query := visualization["query"].(map[string]interface{})

				dId := (int)(query["data_source_id"].(float64))
				if _, ok := dataSources[dId]; !ok {
					if dsPtr, dsErr := getDatasource([]string{fmt.Sprintf("%d", dId)}); dsErr == nil {
						datasource := *dsPtr
						dataSources[dId] = map[string]interface{}{
							"name":    datasource["name"],
							"options": datasource["options"],
							"type":    datasource["type"],
						}
					} else {
						fmt.Printf("Data source not found: %d", dId)
						fmt.Println("This widget will not be imported")
						importWidget = false
					}
				}
				if importWidget {
					wOutEntry := map[string]interface{}{}
					wOutEntry["options"] = widget["options"]
					wOutEntry["text"] = widget["text"]
					wOutEntry["width"] = widget["width"]
					wOutEntry["visualization"] = map[string]interface{}{
						"type":    visualization["type"],
						"name":    visualization["name"],
						"options": visualization["options"],
						"query": map[string]interface{}{
							"name":           query["name"],
							"options":        query["options"],
							"description":    query["description"],
							"data_source_id": query["data_source_id"],
							"query":          query["query"],
						},
					}
					widgetOut = append(widgetOut, wOutEntry)
					dashboardOut["widgets"] = widgetOut
					dashboardOut["datasources"] = dataSources
				}
			} else {
				if importWidget {
					wOutEntry := map[string]interface{}{}
					wOutEntry["options"] = widget["options"]
					wOutEntry["text"] = widget["text"]
					wOutEntry["width"] = widget["width"]
					widgetOut = append(widgetOut, wOutEntry)
					dashboardOut["widgets"] = widgetOut
				}
			}
		}

	}

	s, _ := json.MarshalIndent(dashboardOut, "", "    ")
	fmt.Println(string(s))
}

func GetDashboard(args []string) (*map[string]interface{}, error) {

	uri := GetUrlForResource(ResourceDashboardsGet, args...)
	client := getHttpClient()
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Println("Unable to get dashboards ", err.Error())
		os.Exit(-1)
	}

	req = utils.AddNetTrace(req)

	if api_key := viper.GetString(utils.AuthToken); api_key != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Key %s", api_key))
	}

	if resp, err := client.Do(req); err == nil {
		defer resp.Body.Close()
		var v = map[string]interface{}{}
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("Unable to fetch dashboard, Error: %s", err.Error())
			}
			err = json.Unmarshal(bodyBytes, &v)
			if err != nil {
				return nil, fmt.Errorf("Unable to decode dashboard, Error: %s", err.Error())
			} else {
				return &v, nil
			}
		} else {
			return nil, fmt.Errorf("Http response error, Error: %d", resp.StatusCode)
		}
	} else {
		return nil, fmt.Errorf("Unable to fetch dashboard, Error: %s", err.Error())
	}
}

func createAndPublishDashboard(name string) (map[string]interface{}, error) {
	dashboardParams := map[string]interface{}{
		"name": name,
	}

	if payloadBytes, jsonMarshallError := json.Marshal(dashboardParams); jsonMarshallError != nil {
		return nil, jsonMarshallError
	} else {
		// Create dashboard
		uri := GetUrlForResource(ResourceDashboardsAll)
		client := getHttpClient()
		req, err := http.NewRequest("POST", uri, bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Println("Unable to get dashboards ", err.Error())
			os.Exit(-1)
		}
		if api_key := viper.GetString(utils.AuthToken); api_key != "" {
			req.Header.Add("Authorization", fmt.Sprintf("Key %s", api_key))
			req.Header.Add("Content-Type", "application/json")
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("createAndPublishDashboard1, Http response error while creating dashboard, Error: %d", resp.StatusCode)
		}
		defer resp.Body.Close()
		/*
			var v = map[string]interface{}{}
			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("createAndPublishDashboard2, Http response error while creating dashboard, Error: %d", resp.StatusCode)
			}
		*/

		// Decode create response
		var v = map[string]interface{}{}
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Unable to create dashboard, Read Error: %s", err.Error())
		}
		err = json.Unmarshal(bodyBytes, &v)
		if err != nil {
			return nil, fmt.Errorf("Unable to decode create dashboard response, Error: %s", err.Error())
		}

		// check for server error

		utils.CheckMesgErr(v, "createAndPublishDashboard")

		// Create publish payload
		payloadPublish := map[string]interface{}{
			"is_draft": false, "name": name, "slug": v["id"],
		}
		payloadBytes, jsonMarshallError = json.Marshal(payloadPublish)
		if jsonMarshallError != nil {
			return nil, jsonMarshallError
		}
		args := []string{fmt.Sprintf("%v", v["id"])}

		// Publish dashboard
		uri = GetUrlForResource(ResourceDashboardsGet, args...)
		resp, err = client.Post(uri, "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("createAndPublishDashboard2, Http response error while publishing dashboard, Error: %d", resp.StatusCode)
		}
		defer resp.Body.Close()

		// Decode publish response
		bodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Unable to publish dashboard, Read Error: %s", err.Error())
		}
		err = json.Unmarshal(bodyBytes, &v)
		if err != nil {
			return nil, fmt.Errorf("Unable to decode publish dashboard response, Error: %s", err.Error())
		}
		return v, nil
	}
}

func createAndPublishDashboardSpec(dashboardSpec map[string]interface{}) {
	// First create the datasources and add id's for the data sources
	var dataSources = map[string]interface{}{}

	if _, ok := dashboardSpec["datasources"]; ok {

		dsKeyValue := dashboardSpec["datasources"]
		dataSources = dsKeyValue.(map[string]interface{})
		for _, ds := range dataSources {
			datasource := ds.(map[string]interface{})
			if existingDs := getDataSourceByName(datasource["name"].(string)); existingDs != nil {
				datasource["id"] = existingDs["id"]
			} else {
				if respDict, err := createDataSource(datasource); err != nil {
					fmt.Println("Error creating data source ", err.Error())
					os.Exit(-1)
				} else {
					datasource["id"] = respDict["id"]
				}
			}
		}

	}
	// We now create the dashboard
	if _, dashboardExists := dashboardSpec["dashboard"]; dashboardExists {
		dashboardParams := dashboardSpec["dashboard"].(map[string]interface{})

		if existingDashboard := getDashboardByName(dashboardParams["name"].(string)); existingDashboard != nil {
			fmt.Println("Dashboard already exists ", dashboardParams["name"])
			os.Exit(0)
			dashboardParams["id"] = existingDashboard["id"]
		} else {

			respDict, err := createAndPublishDashboard(dashboardParams["name"].(string))
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(-1)
			}
			dashboardParams["id"] = respDict["id"]

		}

		// We will now create the queries and visualizations
		widgets := dashboardSpec["widgets"]
		for _, w := range widgets.([]interface{}) {
			widget := w.(map[string]interface{})

			var vflag = false
			if _, ok := widget["visualization"]; ok {
				vflag = true
			}
			if vflag {
				visualization := widget["visualization"].(map[string]interface{})
				q := visualization["query"].(map[string]interface{})
				dsIdLookupKey := fmt.Sprintf("%d", (int)(q["data_source_id"].(float64)))
				dsIdForQuery := dataSources[dsIdLookupKey].(map[string]interface{})["id"]
				q["data_source_id"] = dsIdForQuery

				var isDraft bool

				if existingQuery := getQueryByName(q["name"].(string)); existingQuery == nil {
					if respDict, err := createQuery(q); err != nil {
						fmt.Println("Failed creating query,", q)
						os.Exit(-1)
					} else {
						q["id"] = respDict["id"]
						q["version"] = respDict["version"]
						isDraft = respDict["is_draft"].(bool)
					}
				} else {
					q["id"] = existingQuery["id"]
					q["version"] = existingQuery["version"]
					isDraft = existingQuery["is_draft"].(bool)
				}

				if isDraft {
					publishArgs := []string{fmt.Sprintf("%v", q["id"]), fmt.Sprintf("%v", q["version"])}
					//fmt.Println(publisArgs)
					publishQuery(publishArgs)
				}

				// Create the visualization for the widget/query
				// The visualization is attached to
				if respDict, err := createVisualization(visualization, q["id"]); err != nil {
					fmt.Println(err.Error())
					os.Exit(-1)
				} else {
					q["visualization_id"] = respDict["id"]

					// The visualization is ready, lets add to the dashboard
					if _, err := createWidget(widget, q["visualization_id"], dashboardParams["id"]); err != nil {
						fmt.Println(err.Error())
						os.Exit(-1)
					}
					fmt.Println("Added visualize widgets to dashboards ", visualization["name"])

				}
			} else {
				if _, err := createWidget(widget, nil, dashboardParams["id"]); err != nil {
					fmt.Println(err.Error())
					os.Exit(-1)
				}
				fmt.Println("Added none-visualization widget to dashboards")
			}
		}
	}

	//b, _ := json.MarshalIndent(dashboardSpec,"","    ")
	//fmt.Println((string)(b))
}

func getDashboardByName(name string) map[string]interface{} {
	if v, err := GetDashboards(); err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	} else {
		dashboards := v["results"].([]interface{})
		for _, dash := range dashboards {
			dashboard := dash.(map[string]interface{})
			if dashboard["name"] == name {
				return dashboard
			}
		}
	}
	return nil
}

func GetDashboards() (map[string]interface{}, error) {
	uri := GetUrlForResource(ResourceDashboardsAll)
	client := getHttpClient()

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Println("Unable to get dashboards ", err.Error())
		os.Exit(-1)
	}

	req = utils.AddNetTrace(req)

	if api_key := viper.GetString(utils.AuthToken); api_key != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Key %s", api_key))
	}

	if resp, err := client.Do(req); err == nil {
		defer resp.Body.Close()
		var v = map[string]interface{}{}
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("Unable to fetch dashboards, Error: %s", err.Error())
			}
			err = json.Unmarshal(bodyBytes, &v)
			if err != nil {
				return nil, fmt.Errorf("Unable to decode dashboards, Error: %s", err.Error())
			} else {
				return v, nil
			}
		} else {
			return nil, fmt.Errorf("Http response error with get dashboards, Error: %d", resp.StatusCode)
		}
	} else {
		return nil, fmt.Errorf("Unable to fetch dashboards, Error: %s", err.Error())
	}
}

func listDashboards() {
	if v, err := GetDashboards(); err == nil {
		count := (int)(v["count"].(float64))
		dashboards := v["results"].([]interface{})
		fmt.Println("(", count, ") dashboards found")
		if count > 0 {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Slug", "Id"})
			for _, d := range dashboards {
				dash := d.(map[string]interface{})
				slug := dash["slug"].(string)
				name := dash["name"].(string)
				id := (int)(dash["id"].(float64))
				table.Append([]string{name, slug, fmt.Sprintf("%d", id)})
			}

			table.Render()
		}
	} else {
		fmt.Println("Unable to get dashboards ", err.Error())
		os.Exit(-1)
	}
}

func GetLogEvents(numDays int) error {
	response, err := ExecutePrometheusQuery(fmt.Sprintf("round(sum(increase(logiq_message_count[%dh])))", numDays))
	if err != nil {
		return err
	}
	if data, ok := response["data"]; ok {
		if data, ok := data.(map[string]interface{}); ok {
			if result, ok := data["result"]; ok {
				if result, ok := result.([]interface{}); ok {
					for _, v := range result {
						if v, ok := v.(map[string]interface{}); ok {
							for k, vv := range v {
								if k == "value" {
									if vv, ok := vv.([]interface{}); ok {
										fmt.Printf("Total Log Events for %d days in %s: %s\n", numDays, viper.GetString(utils.KeyCluster), vv[1])
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func GetLogVolume(numDays int, raw bool) error {
	response, err := ExecutePrometheusQuery(fmt.Sprintf("round(sum(increase(logiq_data_received_bytes[%dd])))", numDays))
	if err != nil {
		return err
	}
	if data, ok := response["data"]; ok {
		if data, ok := data.(map[string]interface{}); ok {
			if result, ok := data["result"]; ok {
				if result, ok := result.([]interface{}); ok {
					if len(result) == 0 {
						if raw {
							fmt.Println(0)
						} else {
							fmt.Printf("Total Log volume for %d day(s) in %s: 0B\n", numDays, viper.GetString(utils.KeyCluster))
							return nil
						}
					}
					for _, v := range result {
						if v, ok := v.(map[string]interface{}); ok {
							for k, vv := range v {
								if k == "value" {
									if vv, ok := vv.([]interface{}); ok {
										if len(vv) > 1 {
											if vvv, ok := vv[1].(string); ok {
												if raw {
													vInt, err := strconv.ParseFloat(vvv, 64)
													if err != nil {
														fmt.Println(err)
													} else {
														fmt.Println(math.Round(vInt / (1 << 20)))
													}
												} else {
													vInt, err := strconv.ParseInt(vvv, 10, 64)
													if err != nil {
														fmt.Println(err)
													} else {
														fmt.Printf("Total Log volume for %d day(s) in %s: %s\n", numDays, viper.GetString(utils.KeyCluster), humanize.Bytes(uint64(vInt)))
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func ExecutePrometheusQuery(query string) (map[string]interface{}, error) {
	uri := GetUrlForResource(ResourcePrometheusProxy)
	client := getHttpClient()
	payload := fmt.Sprintf(`{"query":"%s","type":"query"}`, query)
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		fmt.Println("Unable to create widget", err.Error())
		os.Exit(-1)
	}
	if api_key := viper.GetString(utils.AuthToken); api_key != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Key %s", api_key))
	}
	if resp, err := client.Do(req); err == nil {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("unable to execute Query %s", err.Error())
		}
		respDict := map[string]interface{}{}
		if errUnmarshall := json.Unmarshal(bodyBytes, &respDict); errUnmarshall != nil {
			return nil, fmt.Errorf("unable to decode response")
		}
		return respDict, nil
	} else {
		fmt.Println("ExecutePrometheusQuery err=<", err, ">")
		return nil, err
	}
}

func NewGetLogEvents() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "log-events",
		Example: "logiqctl get log-events 7",
		Aliases: []string{"t"},
		Short:   "Get total log events for the duration in days",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			var days int = 1
			if len(args) == 1 {
				d, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					fmt.Printf("Unable to parse input, [%v], expects an integer\n", args[0])
					os.Exit(-1)
				}
				days = int(d)
			}
			GetLogEvents(days)
		},
	}
	return cmd
}

func NewGetLogVolume() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "log-volume",
		Example: "logiqctl get log-volume 7",
		Aliases: []string{"v"},
		Short:   "Get total log volume for the duration in days",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			var days int = 1
			if len(args) == 1 {
				d, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					fmt.Printf("Unable to parse input, [%v], expects an integer\n", args[0])
					os.Exit(-1)
				}
				days = int(d)
			}
			GetLogVolume(days, false)
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:     "raw",
		Example: "logiqctl get log-volume raw",
		Short:   "Get volume in MB",
		PreRun:  utils.PreRun,
		Run: func(cmd *cobra.Command, args []string) {
			var days int = 1
			if len(args) == 2 {
				d, err := strconv.ParseInt(args[1], 10, 64)
				if err != nil {
					fmt.Printf("Unable to parse input, [%v], expects an integer\n", args[0])
					os.Exit(-1)
				}
				days = int(d)
			}
			GetLogVolume(days, true)
		},
	})
	return cmd
}
