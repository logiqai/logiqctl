package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/logiqai/logiqctl/defines"

	"github.com/spf13/viper"

	"github.com/logiqai/logiqctl/types"
	"github.com/logiqai/logiqctl/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"net/http"
)

func getDatasourceFromWidgets(dashboardSpec *types.DashboardSpec) {
	datasourceMap := map[string]types.Datasource{}
	validWidgets := []types.Widget{}

	for _, widget := range dashboardSpec.Widgets {
		if widget.Visualization != nil {
			visualization := *widget.Visualization
			query := visualization.Query
			datasourceId := query.DataSourceId
			datasource, err := GetDatasource(fmt.Sprintf("%d", datasourceId))
			if err != nil {
				fmt.Printf("Datasource with %d not found..\n", datasourceId)
			} else {
				datasourceMap[fmt.Sprintf("%d", datasourceId)] = *datasource
				validWidgets = append(validWidgets, widget)
			}
		}
	}
	dashboardSpec.Datasources = datasourceMap
	dashboardSpec.Widgets = validWidgets
}

func GetDashboard(args []string) (string, error) {

	uri := utils.GetUrlForResource(defines.ResourceDashboardsGet, args...)
	client := utils.GetApiClient()

	resp, err := client.MakeApiCall(http.MethodGet, uri, nil)
	if err == nil {
		defer resp.Body.Close()
		var v = types.DashboardSpec{}
		var d = types.Dashboard{}
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", fmt.Errorf("Unable to fetch dashboard, Error: %s", err.Error())
			}
			json.Unmarshal(bodyBytes, &d)
			err = json.Unmarshal(bodyBytes, &v)
			if err != nil {
				return "", fmt.Errorf("Unable to decode dashboard, Error: %s", err.Error())
			} else {
				v.Dashboard = d
				getDatasourceFromWidgets(&v)
				response, err := json.MarshalIndent(v, "", " ")
				if err != nil {
					return "", fmt.Errorf("Error: Unable to convert dashboard to json, %s", err.Error())
				}
				return string(response), nil
			}
		} else {
			return "", fmt.Errorf("Http response error, Error: %d", resp.StatusCode)
		}
	} else {
		return "", fmt.Errorf("Unable to fetch dashboard, Error: %s", err.Error())
	}
}

func CreateAndPublishDashboard(name string, tags []string) (types.Dashboard, error) {
	dashboardParams := map[string]interface{}{
		"name": name,
	}
	if len(tags) > 0 {
		dashboardParams["tags"] = tags
	}

	if payloadBytes, jsonMarshallError := json.Marshal(dashboardParams); jsonMarshallError != nil {
		return types.Dashboard{}, jsonMarshallError
	} else {
		// Create dashboard
		uri := utils.GetUrlForResource(defines.ResourceDashboardsAll)
		client := utils.GetApiClient()
		resp, err := client.MakeApiCall(http.MethodPost, uri, bytes.NewBuffer(payloadBytes))

		if err != nil {
			return types.Dashboard{}, err
		}

		if resp.StatusCode != http.StatusOK {
			return types.Dashboard{}, fmt.Errorf("createAndPublishDashboard1, Http response error while creating dashboard, Error: %d", resp.StatusCode)
		}
		defer resp.Body.Close()
		/*
			var v = map[string]interface{}{}
			if resp.StatusCode != http.StatusOK {
				return types.Dashboard{}, fmt.Errorf("createAndPublishDashboard2, Http response error while creating dashboard, Error: %d", resp.StatusCode)
			}
		*/

		// Decode create response
		var dashboard types.Dashboard
		var v map[string]interface{}
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return types.Dashboard{}, fmt.Errorf("Unable to create dashboard, Read Error: %s", err.Error())
		}
		err = json.Unmarshal(bodyBytes, &v)
		if err != nil {
			return types.Dashboard{}, fmt.Errorf("Unable to decode create dashboard response, Error: %s", err.Error())
		}

		// check for server error

		utils.CheckMesgErr(v, "createAndPublishDashboard")

		// Create publish payload
		payloadPublish := map[string]interface{}{
			"is_draft": false, "name": name, "slug": v["id"],
		}
		payloadBytes, jsonMarshallError = json.Marshal(payloadPublish)
		if jsonMarshallError != nil {
			return types.Dashboard{}, jsonMarshallError
		}
		args := []string{fmt.Sprintf("%v", v["id"])}

		// Publish dashboard
		uri = utils.GetUrlForResource(defines.ResourceDashboardsGet, args...)
		resp, err = client.MakeApiCall(http.MethodPost, uri, bytes.NewBuffer(payloadBytes))

		if err != nil {
			return types.Dashboard{}, err
		}
		if resp.StatusCode != http.StatusOK {
			return types.Dashboard{}, fmt.Errorf("createAndPublishDashboard2, Http response error while publishing dashboard, Error: %d", resp.StatusCode)
		}
		defer resp.Body.Close()

		// Decode publish response
		bodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return types.Dashboard{}, fmt.Errorf("Unable to publish dashboard, Read Error: %s", err.Error())
		}
		err = json.Unmarshal(bodyBytes, &dashboard)
		if err != nil {
			return types.Dashboard{}, fmt.Errorf("Unable to decode publish dashboard response, Error: %s", err.Error())
		}
		return dashboard, nil
	}
}

func CreateAndPublishDashboardSpec(dashboardSpecJson string) (string, error) {
	var dashboardSpec types.DashboardSpec
	var dashboard types.Dashboard
	var err error
	var responseSpec types.DashboardSpec
	var response []byte
	err = json.Unmarshal([]byte(dashboardSpecJson), &dashboardSpec)

	if err != nil {
		return "", fmt.Errorf("Error: %s", err.Error())
	}

	if dashboardSpec.Dashboard.Name != "" {
		existingDashboard, err := GetDashboardByName(dashboardSpec.Dashboard.Name)
		if err != nil {
			return "", err
		}
		if existingDashboard != nil {
			return "", fmt.Errorf("Dashboard with name \"%s\" already exists", dashboardSpec.Dashboard.Name)
		}
		dashboard, err = CreateAndPublishDashboard(dashboardSpec.Dashboard.Name, []string{})
		if err != nil {
			return "", fmt.Errorf("Error: %s", err.Error())
		}
		responseSpec.Dashboard = dashboard
	} else {
		return "", fmt.Errorf("Error: %s", "Dashboard name is missing in spec")
	}

	if len(dashboardSpec.Datasources) > 0 {
		responseSpec.Datasources = make(map[string]types.Datasource)
		for _, datasource := range dashboardSpec.Datasources {
			if existingDatasource, err := getDataSourceByName(datasource.Name); existingDatasource.Id != 0 {
				datasource.Id = existingDatasource.Id
				responseSpec.Datasources[strconv.Itoa(existingDatasource.Id)] = existingDatasource
			} else if err != nil {
				return "", fmt.Errorf("Error: %s", err.Error())
			} else {
				datasource, err = createDataSource(datasource)
				if err != nil {
					return "", fmt.Errorf("Error: %s", err.Error())
				}
				responseSpec.Datasources[strconv.Itoa(datasource.Id)] = datasource
			}
		}
	}
	if len(dashboardSpec.Widgets) > 0 {
		responseSpec.Widgets = []types.Widget{}
		for _, widget := range dashboardSpec.Widgets {
			if widget.Visualization.Type != "" {
				query := widget.Visualization.Query
				query = getQueryByName(query.Name)

				if query == nil {
					tempQuery, _ := json.Marshal(widget.Visualization.Query)
					queryPayload := types.CreateQueryPayload{}
					err = json.Unmarshal(tempQuery, &queryPayload)
					if err != nil {
						return "", err
					}
					queryResponse, err := CreateQuery(queryPayload)
					if err != nil {
						return "", err
					}
					query = &queryResponse
				}
				if query.IsDraft {
					publishArgs := []string{fmt.Sprintf("%d", query.Id), fmt.Sprintf("%d", query.Version)}
					PublishQuery(publishArgs)
				}
				visualization, err := CreateVisualization(widget.Visualization, query.Id)
				if err != nil {
					return "", fmt.Errorf("Error: %s", err.Error())
				}
				newWidget, err := CreateWidget(widget, visualization.Id, responseSpec.Dashboard.Id)
				visualization.Query = query
				newWidget.Visualization = visualization
				responseSpec.Widgets = append(responseSpec.Widgets, newWidget)
			}
		}

		if len(dashboardSpec.Alerts) > 0 {
			alertsJson, _ := json.Marshal(dashboardSpec.Alerts)
			alertPayload := []types.CreateAlertPayload{}
			err := json.Unmarshal(alertsJson, &alertPayload)
			if err != nil {
				return "", fmt.Errorf("Error: %s", err.Error())
			}
			payload, err := json.Marshal(alertPayload)
			if err != nil {
				return "", fmt.Errorf("Error: %s", err.Error())
			}
			alertResponse, err := CreateAlert(string(payload))
			if err != nil {
				return "", fmt.Errorf("Error: %s", err.Error())
			}
			alertList := []types.Alert{}
			json.Unmarshal([]byte(alertResponse), &alertList)
			responseSpec.Alerts = alertList
		}
	}

	response, _ = json.MarshalIndent(&responseSpec, "", " ")

	return string(response), nil
}

func GetDashboardByName(name string) (map[string]interface{}, error) {
	if v, err := GetDashboards(); err != nil {
		return nil, err
	} else {
		dashboards := v["results"].([]interface{})
		for _, dash := range dashboards {
			dashboard := dash.(map[string]interface{})
			if dashboard["name"] == name {
				return dashboard, nil
			}
		}
	}
	return nil, nil
}

func GetDashboards() (map[string]interface{}, error) {
	uri := utils.GetUrlForResource(defines.ResourceDashboardsAll)
	client := utils.GetApiClient()

	resp, err := client.MakeApiCall(http.MethodGet, uri, nil)
	if err == nil {
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

func ListDashboards() {
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

func DeleteDashboard(slug string) error {
	uri := utils.GetUrlForResource(defines.ResourceDashboardsGet, slug)
	client := utils.GetApiClient()

	res, err := client.MakeApiCall(http.MethodDelete, uri, nil)
	if err != nil {
		return err
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var errorResponse map[string]string
		json.Unmarshal(responseData, &errorResponse)

		return fmt.Errorf("Error: %s", errorResponse["message"])
	}

	return nil
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

func ExecutePrometheusQuery(query string) (map[string]interface{}, error) {
	uri := utils.GetUrlForResource(defines.ResourcePrometheusProxy)
	client := utils.GetApiClient()
	payload := fmt.Sprintf(`{"query":"%s","type":"query"}`, query)
	resp, err := client.MakeApiCall(http.MethodPost, uri, bytes.NewBuffer([]byte(payload)))
	if err == nil {
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
