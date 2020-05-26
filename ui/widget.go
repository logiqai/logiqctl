package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func createWidget(vFromDashSpec map[string]interface{}, vId interface{}, dId interface{}) (map[string]interface{}, error) {
	uri := getUrlForResource(ResourceWidgetAll)
	client := getHttpClient()
	vSpec := map[string]interface{}{}
	vSpec["dashboard_id"] = dId
	vSpec["visualization_id"] = vId
	vSpec["options"] = vFromDashSpec["options"]
	vSpec["width"] = vFromDashSpec["width"]
	vSpec["text"] = vFromDashSpec["text"]

	if payloadBytes, jsonMarshallError := json.Marshal(vSpec); jsonMarshallError != nil {
		return nil, jsonMarshallError
	} else {
		if resp, err := client.Post(uri, "application/json", bytes.NewBuffer(payloadBytes)); err == nil {
			jsonStr, _ := json.MarshalIndent(vSpec, "", "    ")
			fmt.Printf("Successfully added widget to dashboard : %s", jsonStr)

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("Unable to read added widget to dashboard response, Error: %s", err.Error())
			}
			respDict := map[string]interface{}{}
			if errUnmarshall := json.Unmarshal(bodyBytes, &respDict); errUnmarshall != nil {
				return nil, fmt.Errorf("Unable to decode added widget to dashboard response")
			}

			return respDict, nil
		} else {
			return nil, err
		}
	}
}
