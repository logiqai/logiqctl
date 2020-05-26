package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func createVisualization(vFromDashSpec map[string]interface{}, qId interface{}) (map[string]interface{}, error) {
	uri := getUrlForResource(ResourceVisualizationAll)
	client := getHttpClient()
	vSpec := map[string]interface{}{}
	vSpec["name"] = vFromDashSpec["name"]
	vSpec["options"] = vFromDashSpec["options"]
	vSpec["description"] = vFromDashSpec["description"]
	vSpec["type"] = vFromDashSpec["type"]
	vSpec["query_id"] = qId

	if payloadBytes, jsonMarshallError := json.Marshal(vSpec); jsonMarshallError != nil {
		return nil, jsonMarshallError
	} else {
		if resp, err := client.Post(uri, "application/json", bytes.NewBuffer(payloadBytes)); err == nil {
			jsonStr, _ := json.MarshalIndent(vSpec, "", "    ")
			fmt.Printf("Successfully created visualization : %s", jsonStr)

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("Unable to read create visualization response, Error: %s", err.Error())
			}
			respDict := map[string]interface{}{}
			if errUnmarshall := json.Unmarshal(bodyBytes, &respDict); errUnmarshall != nil {
				return nil, fmt.Errorf("Unable to decode create visualization response")
			}

			return respDict, nil
		} else {
			return nil, err
		}
	}
}
