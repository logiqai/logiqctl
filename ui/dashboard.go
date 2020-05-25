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

func exportDashboard(args []string) {
	dashboardOut := map[string]interface{}{}

	if dashboardPtr, err := getDashboard(args); err != nil {
		fmt.Println(err.Error())
	} else {
		dashboard := *dashboardPtr
		dashboardParams := map[string]interface{}{}
		dashboardParams["name"]=dashboard["name"]
		dashboardParams["tags"]=dashboard["tags"]
		dashboardOut["dashboard"] = dashboardParams

		widgets := dashboard["widgets"].([]interface{})
		widgetOut := []interface{}{}
		dataSources := map[int]interface{}{}

		for _, w := range widgets {
			widget := w.(map[string]interface{})
			visualization := widget["visualization"].(map[string]interface{})
			query := visualization["query"].(map[string]interface{})
			wOutEntry := map[string]interface{}{}
			wOutEntry["options"] = widget["options"]
			wOutEntry["text"] = widget["text"]
			wOutEntry["visualization"] = map[string]interface{}{
				"type": visualization["type"],
				"name": visualization["name"],
				"options": visualization["options"],
				"query": map[string]interface{}{
					"name":query["name"],
					"options":query["options"],
					"description":query["description"],
					"data_source_id":query["data_source_id"],
				},
			}
			dId := (int)(query["data_source_id"].(float64))
			if _, ok := dataSources[dId]; !ok {
				if dsPtr, dsErr := getDatasource([]string{fmt.Sprintf("%d",dId)});dsErr == nil {
					datasource := *dsPtr
					dataSources[dId] = map[string]interface{}{
						"name": datasource["name"],
						"options": datasource["options"],
						"type": datasource["type"],
					}
				} else {
					fmt.Printf("Data source not found: %d",dId)
					fmt.Println("This widget will not be imported")
				}
			} else {
				fmt.Println("Here", dataSources[dId])
			}
			widgetOut = append(widgetOut, wOutEntry)
		}

		dashboardOut["widgets"] = widgetOut
		dashboardOut["datasources"] = dataSources
	}

	s, _ := json.MarshalIndent(dashboardOut,"","    ")
	fmt.Println(string(s))
}

func getDashboard(args []string) (*map[string]interface{}, error){
	uri := getUrlForResource(ResourceDashboardsGet,args...)
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

func listDashboards() {
	uri := getUrlForResource(ResourceDashboardsAll)
	client := getHttpClient()

	if resp, err := client.Get(uri); err == nil {
		defer resp.Body.Close()
		var v = map[string]interface{}{}
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Unable to fetch dashboards, Error: %s", err.Error())
				os.Exit(-1)
			}
			err = json.Unmarshal(bodyBytes, &v)
			if err != nil {
				fmt.Printf("Unable to decode dashboards, Error: %s", err.Error())
			} else {
				count := (int)(v["count"].(float64))
				dashboards := v["results"].([]interface{})
				fmt.Println("(",count, ") dashboards found")
				if ( count > 0 ) {
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"Name", "Slug", "Id"})
					for _, d := range dashboards {
						dash := d.(map[string]interface{})
						slug := dash["slug"].(string)
						name := dash["name"].(string)
						id := (int)(dash["id"].(float64))
						table.Append([]string{name, slug, fmt.Sprintf("%d",id)})
					}

					table.Render()
				}
			}
		}
	} else {
		fmt.Printf("Unable to fetch dashboards, Error: %s", err.Error())
		os.Exit(-1)
	}
}
