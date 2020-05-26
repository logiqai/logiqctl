package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/logiqai/logiqctl/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
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
		PreRun:  utils.PreRunUiTokenOrCredentials,
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

	return cmd
}

func exportDashboard(args []string) {
	dashboardOut := map[string]interface{}{}

	if dashboardPtr, err := getDashboard(args); err != nil {
		fmt.Println(err.Error())
	} else {
		dashboard := *dashboardPtr
		dashboardParams := map[string]interface{}{}
		dashboardParams["name"] = dashboard["name"]
		dashboardParams["tags"] = dashboard["tags"]
		dashboardOut["dashboard"] = dashboardParams

		widgets := dashboard["widgets"].([]interface{})
		widgetOut := []interface{}{}
		dataSources := map[int]interface{}{}

		for _, w := range widgets {
			widget := w.(map[string]interface{})
			visualization := widget["visualization"].(map[string]interface{})
			query := visualization["query"].(map[string]interface{})

			dId := (int)(query["data_source_id"].(float64))
			importWidget := true
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
			}
		}

		dashboardOut["widgets"] = widgetOut
		dashboardOut["datasources"] = dataSources
	}

	s, _ := json.MarshalIndent(dashboardOut, "", "    ")
	fmt.Println(string(s))
}

func getDashboard(args []string) (*map[string]interface{}, error) {
	uri := getUrlForResource(ResourceDashboardsGet, args...)
	client := getHttpClient()

	if resp, err := client.Get(uri); err == nil {
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
		uri := getUrlForResource(ResourceDashboardsAll)
		client := getHttpClient()
		resp, err := client.Post(uri, "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Http response error while creating dashboard, Error: %d", resp.StatusCode)
		}
		defer resp.Body.Close()
		var v = map[string]interface{}{}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Http response error while creating dashboard, Error: %d", resp.StatusCode)
		}

		// Decode create response
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Unable to create dashboard, Read Error: %s", err.Error())
		}
		err = json.Unmarshal(bodyBytes, &v)
		if err != nil {
			return nil, fmt.Errorf("Unable to decode create dashboard response, Error: %s", err.Error())
		}

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
		uri = getUrlForResource(ResourceDashboardsGet, args...)
		resp, err = client.Post(uri, "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Http response error while publishing dashboard, Error: %d", resp.StatusCode)
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
	dsKeyValue := dashboardSpec["datasources"]
	dataSources := dsKeyValue.(map[string]interface{})
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

	// We now create the dashboard
	dashboardParams := dashboardSpec["dashboard"].(map[string]interface{})

	if existingDashboard := getDashboardByName(dashboardParams["name"].(string)); existingDashboard != nil {
		fmt.Println("Dashboard already exists ", dashboardParams["name"])
		os.Exit(-1)
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
			fmt.Println("Query with name already exists ", q["name"])
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
			} else {
				fmt.Println("Added visualization to dashboards ", visualization["name"])
			}
		}
	}

	//b, _ := json.MarshalIndent(dashboardSpec,"","    ")
	//fmt.Println((string)(b))
}

func getDashboardByName(name string) map[string]interface{} {
	if v, err := getDashboards(); err != nil {
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

func getDashboards() (map[string]interface{}, error) {
	uri := getUrlForResource(ResourceDashboardsAll)
	client := getHttpClient()

	if resp, err := client.Get(uri); err == nil {
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
	if v, err := getDashboards(); err == nil {
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
