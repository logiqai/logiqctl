package ui_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	test_utils "github.com/logiqai/logiqctl/test/utils"
	"github.com/logiqai/logiqctl/types"
	"github.com/logiqai/logiqctl/ui"
	"github.com/stretchr/testify/assert"
)

func TestCreatePublishDashboardSpec(t *testing.T) {
	tearDownTestcase := test_utils.SetupTestCase(t, nil)
	defer tearDownTestcase(t)

	testcases := []test_utils.TestCase{
		{
			Name: "Create Dashboard by reusing datasources and query",
			Input: types.DashboardSpec{
				Dashboard: types.Dashboard{
					Name: "Dashboard 1",
					Tags: []string{},
				},
				Datasources: map[string]types.Datasource{
					"2": {
						Name:           "Datasource 1",
						DatasourceType: "logiq_datasource",
						Options: types.DatasourceOptions{
							Url: "localhost:1234",
						},
					},
				},
				Widgets: []types.Widget{
					{
						Options: map[string]interface{}{
							"isHidden": false,
							"parameterMappings": map[string]map[string]string{
								"duration": {
									"mapTo": "duration",
									"name":  "duration",
									"title": "",
									"type":  "dashboard-level",
									"value": "",
								},
								"step": {
									"mapTo": "step",
									"name":  "step",
									"title": "",
									"type":  "dashboard-level",
									"value": "",
								},
							},
							"position": map[string]interface{}{
								"autoHeight": false,
								"col":        3,
								"maxSizeX":   6,
								"maxSizeY":   1000,
								"minSizeX":   1,
								"minSizeY":   5,
								"row":        40,
								"sizeX":      3,
								"sizeY":      8,
							},
						},
						Width: 1,
						Visualization: &types.Visualization{
							Name: "Chart",
							Type: "CHART",
							Options: map[string]interface{}{
								"columnMapping": map[string]string{
									"timestamp": "x",
									"value":     "y",
								},
								"customCode":     "// Available variables are x, ys, element, and Plotly\n// Type console.log(x, ys); for more info about x and ys\n// To plot your graph call Plotly.plot(element, ...)\n// Plotly examples and docs: https://plot.ly/javascript/",
								"dateTimeFormat": "DD/MM/YY HH:mm",
								"defaultColumns": 3,
								"defaultRows":    8,
								"error_y": map[string]interface{}{
									"type":    "data",
									"visible": true,
								},
								"globalSeriesType": "line",
								"legend": map[string]interface{}{
									"enabled": true,
								},
								"minColumns":    1,
								"minRows":       5,
								"numberFormat":  "0,0[.]00000",
								"percentFormat": "0[.]00%",
								"series": map[string]interface{}{
									"error_y": map[string]interface{}{
										"type":    "data",
										"visible": true,
									},
									"stacking": nil,
								},
								"seriesOptions": map[string]map[string]interface{}{
									"value": {
										"index":  0,
										"type":   "line",
										"yAxis":  0,
										"zIndex": 0,
									},
								},
								"showDataLabels": false,
								"sortX":          true,
								"textFormat":     "",
								"xAxis": map[string]interface{}{
									"labels": map[string]bool{
										"enabled": true,
									},
									"type": "-",
								},
								"yAxis": []map[string]interface{}{
									{
										"type": "linear",
									},
									{
										"opposite": true,
										"type":     "linear",
									},
								},
							},
							Query: &types.Query{
								DataSourceId: 1,
								Description:  "",
								Name:         "Query with schedule",
								QueryOptions: types.QueryOptions{
									Parameters: []map[string]interface{}{
										{
											"enumOptions": "1h\n2h\n3h\n1d\n2d\n3d",
											"locals":      []string{},
											"name":        "duration",
											"title":       "Duration",
											"type":        "enum",
											"value":       "2h",
										},
										{
											"enumOptions": "10s\n30s\n60s\n1m\n5m",
											"locals":      []string{},
											"name":        "step",
											"title":       "Step",
											"type":        "enum",
											"value":       "1m",
										},
									},
								},
								Query: "query=rate(redis_keyspace_hits_total[5m])\u0026 duration={{ duration }}\u0026step={{ step }}",
							},
						},
					},
				},
			},
			MockResponses: []test_utils.MockResponse{
				{
					Url:        "api/dashboards",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/dashboard/create/201_response.json",
					Err:        nil,
				},
				{
					Url:        "api/dashboards/1",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/dashboard/publish/200_response.json",
					Err:        nil,
				},
				{
					Url:        "api/dashboards",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/dashboard/list/200_response.json",
					Err:        nil,
				},
				{
					Url:        "api/data_sources",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/datasource/list/200_response.json",
					Err:        nil,
				},
				{
					Url:        "api/queries",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/queries/list/200_success.json",
					Err:        nil,
				},
				{
					Url:         "api/widgets",
					HttpMethod:  http.MethodPost,
					StatusCode:  200,
					Body:        test_utils.BASE_TEST_DATA_DIR + "/mock_response/widget/create/201_success.json",
					Err:         nil,
					FromPayload: true,
				},
				{
					Url:         "api/visualizations",
					HttpMethod:  http.MethodPost,
					StatusCode:  200,
					Body:        test_utils.BASE_TEST_DATA_DIR + "/mock_response/visualization/create/201_success.json",
					Err:         nil,
					FromPayload: true,
				},
			},
		},
		{
			Name: "Create Dashboard by creating datasources and query",
			Input: types.DashboardSpec{
				Dashboard: types.Dashboard{
					Name: "Dashboard 1",
					Tags: []string{},
				},
				Datasources: map[string]types.Datasource{
					"2": {
						Name:           "Datasource 1",
						DatasourceType: "logiq_datasource",
						Options: types.DatasourceOptions{
							Url: "localhost:1234",
						},
					},
				},
				Widgets: []types.Widget{
					{
						Options: map[string]interface{}{
							"isHidden": false,
							"parameterMappings": map[string]map[string]string{
								"duration": {
									"mapTo": "duration",
									"name":  "duration",
									"title": "",
									"type":  "dashboard-level",
									"value": "",
								},
								"step": {
									"mapTo": "step",
									"name":  "step",
									"title": "",
									"type":  "dashboard-level",
									"value": "",
								},
							},
							"position": map[string]interface{}{
								"autoHeight": false,
								"col":        3,
								"maxSizeX":   6,
								"maxSizeY":   1000,
								"minSizeX":   1,
								"minSizeY":   5,
								"row":        40,
								"sizeX":      3,
								"sizeY":      8,
							},
						},
						Width: 1,
						Visualization: &types.Visualization{
							Name: "Chart",
							Type: "CHART",
							Options: map[string]interface{}{
								"columnMapping": map[string]string{
									"timestamp": "x",
									"value":     "y",
								},
								"customCode":     "// Available variables are x, ys, element, and Plotly\n// Type console.log(x, ys); for more info about x and ys\n// To plot your graph call Plotly.plot(element, ...)\n// Plotly examples and docs: https://plot.ly/javascript/",
								"dateTimeFormat": "DD/MM/YY HH:mm",
								"defaultColumns": 3,
								"defaultRows":    8,
								"error_y": map[string]interface{}{
									"type":    "data",
									"visible": true,
								},
								"globalSeriesType": "line",
								"legend": map[string]interface{}{
									"enabled": true,
								},
								"minColumns":    1,
								"minRows":       5,
								"numberFormat":  "0,0[.]00000",
								"percentFormat": "0[.]00%",
								"series": map[string]interface{}{
									"error_y": map[string]interface{}{
										"type":    "data",
										"visible": true,
									},
									"stacking": nil,
								},
								"seriesOptions": map[string]map[string]interface{}{
									"value": {
										"index":  0,
										"type":   "line",
										"yAxis":  0,
										"zIndex": 0,
									},
								},
								"showDataLabels": false,
								"sortX":          true,
								"textFormat":     "",
								"xAxis": map[string]interface{}{
									"labels": map[string]bool{
										"enabled": true,
									},
									"type": "-",
								},
								"yAxis": []map[string]interface{}{
									{
										"type": "linear",
									},
									{
										"opposite": true,
										"type":     "linear",
									},
								},
							},
							Query: &types.Query{
								DataSourceId: 1,
								Description:  "",
								Name:         "Query 1",
								QueryOptions: types.QueryOptions{
									Parameters: []map[string]interface{}{
										{
											"enumOptions": "1h\n2h\n3h\n1d\n2d\n3d",
											"locals":      []string{},
											"name":        "duration",
											"title":       "Duration",
											"type":        "enum",
											"value":       "2h",
										},
										{
											"enumOptions": "10s\n30s\n60s\n1m\n5m",
											"locals":      []string{},
											"name":        "step",
											"title":       "Step",
											"type":        "enum",
											"value":       "1m",
										},
									},
								},
								Query: "query=rate(redis_keyspace_hits_total[5m])\u0026 duration={{ duration }}\u0026step={{ step }}",
								QuerySchedule: &types.QuerySchedule{
									Interval: 300,
								},
							},
						},
					},
				},
			},
			MockResponses: []test_utils.MockResponse{
				{
					Url:        "api/dashboards",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/dashboard/create/201_response.json",
					Err:        nil,
				},
				{
					Url:        "api/dashboards/1",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/dashboard/publish/200_response.json",
					Err:        nil,
				},
				{
					Url:        "api/data_sources",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/datasource/list/200_response.json",
					Err:        nil,
				},
				{
					Url:        "api/queries",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/queries/list/200_success.json",
					Err:        nil,
				},
				{
					Url:         "api/queries",
					HttpMethod:  http.MethodPost,
					StatusCode:  200,
					Body:        test_utils.BASE_TEST_DATA_DIR + "/mock_response/queries/create/201_success.json",
					Err:         nil,
					FromPayload: true,
				},
				{
					Url:         "api/widgets",
					HttpMethod:  http.MethodPost,
					StatusCode:  200,
					Body:        test_utils.BASE_TEST_DATA_DIR + "/mock_response/widget/create/201_success.json",
					Err:         nil,
					FromPayload: true,
				},
				{
					Url:         "api/visualizations",
					HttpMethod:  http.MethodPost,
					StatusCode:  200,
					Body:        test_utils.BASE_TEST_DATA_DIR + "/mock_response/visualization/create/201_success.json",
					Err:         nil,
					FromPayload: false,
				},
				{
					Url:        "api/data_sources",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/datasource/create/201_success.json",
					Err:        nil,
				},
				{
					Url:        "api/alerts",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/alerts/create/success.json",
					Err:        nil,
				},
			},
		},
		{
			Name: "Create Dashboard without widgets",
			Input: types.DashboardSpec{
				Dashboard: types.Dashboard{
					Name: "Dashboard 1",
					Tags: []string{},
				},
				Datasources: map[string]types.Datasource{
					"2": {
						Name:           "Datasource 1",
						DatasourceType: "logiq_datasource",
						Options: types.DatasourceOptions{
							Url: "localhost:1234",
						},
					},
				},
			},
			MockResponses: []test_utils.MockResponse{
				{
					Url:        "api/dashboards",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/dashboard/create/201_response.json",
					Err:        nil,
				},
				{
					Url:        "api/dashboards/1",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/dashboard/publish/200_response.json",
					Err:        nil,
				},
				{
					Url:        "api/data_sources",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/datasource/list/200_response.json",
					Err:        nil,
				},
			},
		},
		{
			Name: "Create Dashboard with name which already exists",
			Input: types.DashboardSpec{
				Dashboard: types.Dashboard{
					Name: "Dashboard exist",
					Tags: []string{},
				},
				Datasources: map[string]types.Datasource{
					"2": {
						Name:           "Datasource 1",
						DatasourceType: "logiq_datasource",
						Options: types.DatasourceOptions{
							Url: "localhost:1234",
						},
					},
				},
			},
			Expected: "Dashboard with name \"Dashboard exist\" already exists",
			MockResponses: []test_utils.MockResponse{
				{
					Url:        "api/dashboards",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/dashboard/create/201_response.json",
					Err:        nil,
				},
				{
					Url:        "api/dashboards/1",
					HttpMethod: http.MethodPost,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/dashboard/publish/200_response.json",
					Err:        nil,
				},
				{
					Url:        "api/data_sources",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/datasource/list/200_response.json",
					Err:        nil,
				},
				{
					Url:        "api/data_sources",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/dashboards/list/200_response.json",
					Err:        nil,
				},
			},
		},
	}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			test_utils.MockApiResponse(&testcase.MockResponses)
			dashboardSpec := testcase.Input.(types.DashboardSpec)
			payload, err := json.Marshal(dashboardSpec)
			actual, err := ui.CreateAndPublishDashboardSpec(string(payload))
			if err != nil {
				assert.Equal(t, testcase.Expected, err.Error(), testcase.Name)
				return
			}
			actualDashboardSpec := types.DashboardSpec{}
			err = json.Unmarshal([]byte(actual), &actualDashboardSpec)
			if err != nil {
				assert.Fail(t, err.Error(), testcase.Name)
			} else {
				assert.Greater(t, actualDashboardSpec.Dashboard.Id, 0, "dashboard is not created")
				assert.Greater(t, len(actualDashboardSpec.Datasources), 0, "Datasource is empty")
				if len(dashboardSpec.Alerts) > 0 {
					assert.Equal(t, len(actualDashboardSpec.Alerts), 2, "Alerts is empty")
					for _, alert := range dashboardSpec.Alerts {
						assert.Greater(t, alert.Id, 0, testcase.Name)
					}
				}
				for _, datasource := range actualDashboardSpec.Datasources {
					assert.Greater(t, datasource.Id, 0, testcase.Name)
				}
				for _, alert := range actualDashboardSpec.Alerts {
					assert.Greater(t, alert.Id, 0, testcase.Name)
				}
				for _, widget := range actualDashboardSpec.Widgets {
					assert.Greater(t, widget.Id, 0, "widget is not created")
					assert.NotNil(t, widget.Visualization, "Visualization is not created")
					assert.Greater(t, widget.Visualization.Query.Id, 0, "Query is not created")
				}
			}
		})
	}
}

func TestGetDashboard(t *testing.T) {
	tearDownTestcase := test_utils.SetupTestCase(t, nil)
	defer tearDownTestcase(t)

	testcases := []test_utils.TestCase{
		{
			Name:  "Get Dashboard by slug",
			Input: "dashboard-1",
			Expected: types.DashboardSpec{
				Dashboard: types.Dashboard{
					Id:   1,
					Name: "Dashboard 1",
					Tags: []string{},
					Slug: "dashboard-1",
				},
				Datasources: map[string]types.Datasource{
					"1": {
						Id:             1,
						Name:           "Datasource 1",
						DatasourceType: "logiq_datasource",
						Options: types.DatasourceOptions{
							Url: "localhost:1234",
						},
					},
				},
				Widgets: []types.Widget{
					{
						Id: 1,
						Options: map[string]interface{}{
							"isHidden":          false,
							"parameterMappings": map[string]map[string]string{},
							"position": map[string]interface{}{
								"autoHeight": false,
								"sizeX":      2,
								"sizeY":      5,
								"minSizeX":   1,
								"maxSizeX":   6,
								"minSizeY":   1,
								"maxSizeY":   1000,
								"col":        0,
								"row":        0,
							},
						},
						Width: 1,
						Visualization: &types.Visualization{
							Id:   2,
							Name: "Go Version",
							Type: "COUNTER",
							Options: map[string]interface{}{
								"counterColName":  "value",
								"rowNumber":       1,
								"targetRowNumber": 1,
								"stringDecimal":   0,
								"stringDecChar":   ".",
								"stringThouSep":   ",",
								"defaultColumns":  2,
								"defaultRows":     5,
								"bgColor":         nil,
								"textColor":       "#049235",
							},
							Query: &types.Query{
								Id:                1,
								LatestQueryDataId: 9,
								DataSourceId:      1,
								Description:       "",
								Name:              "Go Version",
								QueryOptions: types.QueryOptions{
									Parameters: []map[string]interface{}{},
								},
								Query:   "go_info{job=\"flash\"}",
								Tags:    []string{},
								Version: 1,
								QuerySchedule: &types.QuerySchedule{
									Interval: 300,
								},
							},
						},
					},
				},
			},
			MockResponses: []test_utils.MockResponse{
				{
					Url:        "api/dashboards/dashboard-1",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/dashboard/get/200_response.json",
					Err:        nil,
				},
				{
					Url:        "api/data_sources/1",
					HttpMethod: http.MethodGet,
					StatusCode: 200,
					Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/datasource/get/200_response.json",
					Err:        nil,
				},
			},
		},
	}
	httpmock.Activate()
	defer httpmock.Deactivate()
	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			test_utils.MockApiResponse(&testcase.MockResponses)
			dashboardSlug := testcase.Input.(string)
			actual, err := ui.GetDashboard([]string{dashboardSlug})
			if err != nil {
				assert.Fail(t, err.Error(), testcase.Name)
			}
			expected, err := json.MarshalIndent(testcase.Expected, "", " ")
			if err != nil {
				fmt.Printf("Error decoding expected response: %s\n", err.Error())
			}
			assert.Equal(t, string(expected), actual, testcase.Name)
		})
	}
}
