package ui_test

import (
	"testing"

	"net/http"

	"github.com/jarcoal/httpmock"
	test_utils "github.com/logiqai/logiqctl/test/utils"
	"github.com/logiqai/logiqctl/ui"
	"github.com/stretchr/testify/assert"
)

var mockResponseList = []test_utils.MockResponse{
	{
		Url:        "api/dashboards/dashboard-1",
		HttpMethod: http.MethodDelete,
		StatusCode: 200,
		Body:       test_utils.BASE_TEST_DATA_DIR + "/mock_response/dashboard/delete/200_response.json",
		Err:        nil,
	},
}

func TestDeleteDashboard(t *testing.T) {
	tearDownTestcase := test_utils.SetupTestCase(t, &mockResponseList)
	defer tearDownTestcase(t)

	testcases := []test_utils.TestCase{
		{
			Name: "Archive dashboard with all the widgets given the dashboard slug",
			Input: map[string]interface{}{
				"slug":    "dashboard-1",
				"archive": true,
				"query":   false,
			},
			Expected: true,
		},
		{
			Name: "Delete dashboard with all the widgets given the dashboard slug including queries",
			Input: map[string]interface{}{
				"slug":    "dashboard-1",
				"archive": false,
				"query":   true,
			},
			Expected: true,
		},
	}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, testCase := range testcases {
		t.Run(testCase.Name, func(t *testing.T) {
			test_utils.SetupOutputFormat(testCase.OutputFormat)
			input := testCase.Input.(map[string]interface{})
			res, err := ui.DeleteDashboardBySlug(input["slug"].(string), input["archive"].(bool), input["query"].(bool))
			if err != nil {
				if testCase.Expected != nil {
					assert.Equal(t, testCase.Expected.(string), err.Error())
				}
				assert.Fail(t, err.Error())
				return
			}
			assert.Equal(t, testCase.Expected.(bool), res)
		})
	}
}
