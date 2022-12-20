package converter_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/logiqai/logiqctl/converter"
	test_utils "github.com/logiqai/logiqctl/test/utils"
)

var mockResponseList = []test_utils.MockResponse{
	{
		Url:         "api/dashboards",
		HttpMethod:  http.MethodPost,
		StatusCode:  200,
		Body:        test_utils.BASE_TEST_DATA_DIR + "/mock_response/dashboard/create/201_response.json",
		Err:         nil,
		FromPayload: true,
	},
	{
		Url:        "api/dashboards",
		HttpMethod: http.MethodGet,
		StatusCode: 200,
		Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/dashboard/list/200_response.json",
		Err:        nil,
	},
	{
		Url:         "api/dashboards/1",
		HttpMethod:  http.MethodPost,
		StatusCode:  200,
		Body:        test_utils.BASE_TEST_DATA_DIR + "/mock_response/dashboard/publish/200_response.json",
		Err:         nil,
		FromPayload: true,
	},
	{
		Url:        "api/data_sources/1",
		HttpMethod: http.MethodGet,
		StatusCode: 200,
		Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/datasource/get/200_response.json",
		Err:        nil,
	},
	{
		Url:        "api/data_sources/1",
		HttpMethod: http.MethodGet,
		StatusCode: 200,
		Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/datasource/get/200_response.json",
		Err:        nil,
	},
	{
		Url:         "api/queries",
		HttpMethod:  http.MethodPost,
		StatusCode:  200,
		Body:        test_utils.BASE_TEST_DATA_DIR + "/mock_response/queries/create/201_grafana.json",
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
		FromPayload: true,
	},
	{
		Url:        "api/query_results",
		HttpMethod: http.MethodPost,
		StatusCode: 200,
		Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/query_results/create/200.json",
		Err:        nil,
	},
	{
		Url:        "api/jobs/job_1",
		HttpMethod: http.MethodGet,
		StatusCode: 200,
		Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/jobs/get/200.json",
		Err:        nil,
	},
	{
		Url:        "api/query_results/1",
		HttpMethod: http.MethodGet,
		StatusCode: 200,
		Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/query_results/get/200_response.json",
		Err:        nil,
	},
	{
		Url:        "api/queries/1",
		HttpMethod: http.MethodPost,
		StatusCode: 200,
		Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/queries/create/201_grafana.json",
		Err:        nil,
	},
}

func TestGraphanaConvert(t *testing.T) {
	tearDownTestcase := test_utils.SetupTestCase(t, &mockResponseList)
	defer tearDownTestcase(t)

	testcases := []test_utils.TestCase{
		// {
		// 	Name: "Convert graphana dasboard with line charts to logiq dashboard",
		// 	Input: map[string]interface{}{
		// 		"grafanaJson":   test_utils.BASE_TEST_DATA_DIR + "/data/converter/line-charts/input.json",
		// 		"dashboardName": "Dashboard 1",
		// 		"inputMap":      test_utils.BASE_TEST_DATA_DIR + "/data/converter/input_map.json",
		// 	},
		// 	Expected: test_utils.BASE_TEST_DATA_DIR + "/data/converter/line-charts/output.json",
		// },
		// {
		// 	Name: "Convert graphana dasboard with bar charts to logiq dashboard",
		// 	Input: map[string]interface{}{
		// 		"grafanaJson":   test_utils.BASE_TEST_DATA_DIR + "/data/converter/bar-charts/input.json",
		// 		"dashboardName": "Dashboard 1",
		// 		"inputMap":      test_utils.BASE_TEST_DATA_DIR + "/data/converter/input_map.json",
		// 	},
		// 	Expected: test_utils.BASE_TEST_DATA_DIR + "/data/converter/bar-charts/output.json",
		// },
		// {
		// 	Name: "Convert graphana dasboard with gauge to logiq dashboard",
		// 	Input: map[string]interface{}{
		// 		"grafanaJson":   test_utils.BASE_TEST_DATA_DIR + "/data/converter/gauge/input.json",
		// 		"dashboardName": "Dashboard 1",
		// 		"inputMap":      test_utils.BASE_TEST_DATA_DIR + "/data/converter/input_map.json",
		// 	},
		// 	Expected: test_utils.BASE_TEST_DATA_DIR + "/data/converter/gauge/output.json",
		// },
		// {
		// 	Name: "Convert graphana dasboard with heatmap to logiq dashboard",
		// 	Input: map[string]interface{}{
		// 		"grafanaJson":   test_utils.BASE_TEST_DATA_DIR + "/data/converter/heatmap/input.json",
		// 		"dashboardName": "Dashboard 1",
		// 		"inputMap":      test_utils.BASE_TEST_DATA_DIR + "/data/converter/input_map.json",
		// 	},
		// 	Expected: test_utils.BASE_TEST_DATA_DIR + "/data/converter/heatmap/output.json",
		// },
		// {
		// 	Name: "Convert graphana dasboard with piechart to logiq dashboard",
		// 	Input: map[string]interface{}{
		// 		"grafanaJson":   test_utils.BASE_TEST_DATA_DIR + "/data/converter/pie-chart/input.json",
		// 		"dashboardName": "Dashboard 1",
		// 		"inputMap":      test_utils.BASE_TEST_DATA_DIR + "/data/converter/input_map.json",
		// 	},
		// 	Expected: test_utils.BASE_TEST_DATA_DIR + "/data/converter/pie-chart/output.json",
		// },
		// {
		// 	Name: "Convert graphana dasboard with table to logiq dashboard",
		// 	Input: map[string]interface{}{
		// 		"grafanaJson":   test_utils.BASE_TEST_DATA_DIR + "/data/converter/table/input.json",
		// 		"dashboardName": "Dashboard 1",
		// 		"inputMap":      test_utils.BASE_TEST_DATA_DIR + "/data/converter/input_map.json",
		// 	},
		// 	Expected: test_utils.BASE_TEST_DATA_DIR + "/data/converter/table/output.json",
		// },
		// {
		// 	Name: "Convert graphana dasboard with counter to logiq dashboard",
		// 	Input: map[string]interface{}{
		// 		"grafanaJson":   test_utils.BASE_TEST_DATA_DIR + "/data/converter/counter/input.json",
		// 		"dashboardName": "Dashboard 1",
		// 		"inputMap":      test_utils.BASE_TEST_DATA_DIR + "/data/converter/input_map.json",
		// 	},
		// 	Expected: test_utils.BASE_TEST_DATA_DIR + "/data/converter/counter/output.json",
		// },
		// {
		// 	Name: "Convert graphana dasboard with template label query which givens empty values",
		// 	Input: map[string]interface{}{
		// 		"grafanaJson":   test_utils.BASE_TEST_DATA_DIR + "/data/converter/edge_cases/empty_label_values/input.json",
		// 		"dashboardName": "Dashboard 1",
		// 		"inputMap":      test_utils.BASE_TEST_DATA_DIR + "/data/converter/input_map.json",
		// 	},
		// 	Expected: test_utils.BASE_TEST_DATA_DIR + "/data/converter/edge_cases/empty_label_values/output.json",
		// },
		{
			Name: "Convert graphana dasboard with formatted result",
			Input: map[string]interface{}{
				"grafanaJson":   test_utils.BASE_TEST_DATA_DIR + "/data/converter/edge_cases/format_result/input.json",
				"dashboardName": "Dashboard 1",
				"inputMap":      test_utils.BASE_TEST_DATA_DIR + "/data/converter/input_map.json",
			},
			Expected: test_utils.BASE_TEST_DATA_DIR + "/data/converter/edge_cases/format_result/output.json",
		},
		// {
		// 	Name: "Convert graphana dasboard by parsing query with template variables",
		// 	Input: map[string]interface{}{
		// 		"grafanaJson":   test_utils.BASE_TEST_DATA_DIR + "/data/converter/edge_cases/parse_query/input.json",
		// 		"dashboardName": "Dashboard 1",
		// 		"inputMap":      test_utils.BASE_TEST_DATA_DIR + "/data/converter/input_map.json",
		// 	},
		// 	Expected: test_utils.BASE_TEST_DATA_DIR + "/data/converter/edge_cases/parse_query/output.json",
		// },
		// {
		// 	Name: "Convert graphana dasboard by parsing teamplates with different formats and types",
		// 	Input: map[string]interface{}{
		// 		"grafanaJson":   test_utils.BASE_TEST_DATA_DIR + "/data/converter/edge_cases/parse_template/input.json",
		// 		"dashboardName": "Dashboard 1",
		// 		"inputMap":      test_utils.BASE_TEST_DATA_DIR + "/data/converter/input_map.json",
		// 	},
		// 	Expected: test_utils.BASE_TEST_DATA_DIR + "/data/converter/edge_cases/parse_template/output.json",
		// },
		// {
		// 	Name: "Convert graphana dasboard by parsing datasource with different schema",
		// 	Input: map[string]interface{}{
		// 		"grafanaJson":   test_utils.BASE_TEST_DATA_DIR + "/data/converter/edge_cases/parse_datasource/input.json",
		// 		"dashboardName": "Dashboard 1",
		// 		"inputMap":      test_utils.BASE_TEST_DATA_DIR + "/data/converter/input_map.json",
		// 	},
		// 	Expected: test_utils.BASE_TEST_DATA_DIR + "/data/converter/edge_cases/parse_datasource/output.json",
		// },
		// {
		// 	Name: "Import grafana dashboard with the dashboard name which already exist",
		// 	Input: map[string]interface{}{
		// 		"grafanaJson":   test_utils.BASE_TEST_DATA_DIR + "/data/converter/edge_cases/parse_template/input.json",
		// 		"dashboardName": "Dashboard exist",
		// 		"inputMap":      test_utils.BASE_TEST_DATA_DIR + "/data/converter/input_map.json",
		// 	},
		// 	Expected: "error: Dashboard with name 'Dashboard exist' already exist",
		// },
		// {
		// 	Name: "Import grafana dashboard without the datasource id",
		// 	Input: map[string]interface{}{
		// 		"grafanaJson":   test_utils.BASE_TEST_DATA_DIR + "/data/converter/edge_cases/parse_template/input.json",
		// 		"dashboardName": "Dashboard 1",
		// 		"inputMap":      test_utils.BASE_TEST_DATA_DIR + "/data/converter/edge_cases/nil_datasource/input_map.json",
		// 	},
		// 	Expected: "error: Your template is missing with input for datasource. Datasource Id is required to import this dashboard. Please update your template with datasource input",
		// },
		// {
		// 	Name: "Convert graphana dasboard with widget groups",
		// 	Input: map[string]interface{}{
		// 		"grafanaJson":   test_utils.BASE_TEST_DATA_DIR + "/data/converter/widget-group/simple-group/input.json",
		// 		"dashboardName": "Dashboard 1",
		// 		"inputMap":      test_utils.BASE_TEST_DATA_DIR + "/data/converter/input_map.json",
		// 	},
		// 	Expected: test_utils.BASE_TEST_DATA_DIR + "/data/converter/widget-group/simple-group/output.json",
		// },
		// {
		// 	Name: "Convert graphana dasboard with widget groups which have panels embbeded in the panel row",
		// 	Input: map[string]interface{}{
		// 		"grafanaJson":   test_utils.BASE_TEST_DATA_DIR + "/data/converter/widget-group/embedded-panels/input.json",
		// 		"dashboardName": "Dashboard 1",
		// 		"inputMap":      test_utils.BASE_TEST_DATA_DIR + "/data/converter/input_map.json",
		// 	},
		// 	Expected: test_utils.BASE_TEST_DATA_DIR + "/data/converter/widget-group/embedded-panels/output.json",
		// },
		// {
		// 	Name: "Convert graphana dasboard with widget groups which have rows embedded with panels",
		// 	Input: map[string]interface{}{
		// 		"grafanaJson":   test_utils.BASE_TEST_DATA_DIR + "/data/converter/widget-group/rows/input.json",
		// 		"dashboardName": "Dashboard 1",
		// 		"inputMap":      test_utils.BASE_TEST_DATA_DIR + "/data/converter/input_map.json",
		// 	},
		// 	Expected: test_utils.BASE_TEST_DATA_DIR + "/data/converter/widget-group/rows/output.json",
		// },
	}
	t.Parallel()
	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			input := testcase.Input.(map[string]interface{})
			dashboardName := input["dashboardName"].(string)
			grafanaJson, err := ioutil.ReadFile(input["grafanaJson"].(string))
			inputMap, err := ioutil.ReadFile(input["inputMap"].(string))
			if err != nil {
				fmt.Println("error reading input map, " + err.Error())
				os.Exit(-1)
			}
			if len(testcase.MockResponses) > 0 {
				test_utils.MockApiResponse(&testcase.MockResponses)
			}
			if err != nil {
				fmt.Println("Error reading payload:", err.Error())
				os.Exit(-1)
			}

			numRoutines := runtime.NumGoroutine()
			actual, err := converter.ConvertToLogiqDashboard(string(grafanaJson), "grafana", dashboardName, &inputMap)
			assert.Equal(t, numRoutines, runtime.NumGoroutine(), "verify goroutines leak")
			if err != nil {
				assert.Equal(t, testcase.Expected, err.Error(), testcase.Name)
				return
			}

			expected, err := ioutil.ReadFile(testcase.Expected.(string))
			if err != nil {
				fmt.Println("Error reading output file:", err.Error())
				os.Exit(-1)
			}
			assert.JSONEq(t, string(expected), actual, testcase.Name)
		})
	}
}
