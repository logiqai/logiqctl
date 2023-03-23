package ui

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/logiqai/logiqctl/defines"
	"github.com/logiqai/logiqctl/types"
	"github.com/logiqai/logiqctl/utils"
)

type Forwarders struct {
	Configs []types.Forwarder `json:"configs"`
}

func GetForwarders() ([]types.Forwarder, error) {
	client := utils.GetApiClient()
	uri := utils.GetUrlForResource(defines.ResourceForwardersAll)
	res, err := client.MakeApiCall(http.MethodGet, uri, nil)
	if err != nil {
		return []types.Forwarder{}, err
	}

	defer res.Body.Close()
	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return []types.Forwarder{}, err
	}

	if res.StatusCode != 200 {
		var errorResponse map[string]string
		json.Unmarshal(responseData, &errorResponse)
		return nil, fmt.Errorf("Error: %s", errorResponse["message"])
	}

	var forwarderList Forwarders
	json.Unmarshal(responseData, &forwarderList)

	return forwarderList.Configs, nil
}
