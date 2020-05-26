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

func NewListDatasourcesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "datasource",
		Example: "logiqctl get datasource|ds <datasource-id>",
		Aliases: []string{"ds"},
		Short:   "Get a datasource",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Missing datasource id")
				os.Exit(-1)
			}
			printDataSource(args)
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:     "all",
		Example: "logiqctl get datasource all",
		Short:   "List all the available datasources",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			listDataSources()
		},
	})

	return cmd
}

func printDataSource(args []string) {
	if v, vErr := getDatasource(args); v != nil {
		s, _ := json.MarshalIndent(*v, "", "    ")
		fmt.Println(string(s))
	} else {
		fmt.Println(vErr.Error())
		os.Exit(-1)
	}
}

func getDatasource(args []string) (*map[string]interface{}, error) {
	uri := getUrlForResource(ResourceDatasource, args...)
	client := getHttpClient()

	if resp, err := client.Get(uri); err == nil {
		defer resp.Body.Close()
		var v = map[string]interface{}{}
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("Unable to fetch datasources, Error: %s", err.Error())
			}
			err = json.Unmarshal(bodyBytes, &v)
			if err != nil {
				return nil, fmt.Errorf("Unable to decode datasources, Error: %s", err.Error())
			} else {
				return &v, nil
			}
		} else {
			return nil, fmt.Errorf("Http response error, Error: %d", resp.StatusCode)
		}
	} else {
		return nil, fmt.Errorf("Unable to fetch datasources, Error: %s", err.Error())
	}
}

func createDataSource(datasourceSpec map[string]interface{}) (map[string]interface{}, error) {
	uri := getUrlForResource(ResourceDatasourceAll)
	client := getHttpClient()

	if payloadBytes, jsonMarshallError := json.Marshal(datasourceSpec); jsonMarshallError != nil {
		return nil, jsonMarshallError
	} else {
		if resp, err := client.Post(uri, "application/json", bytes.NewBuffer(payloadBytes)); err == nil {
			jsonStr, _ := json.MarshalIndent(datasourceSpec, "", "    ")
			fmt.Printf("Successfully created datasource : %s", jsonStr)

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("Unable to read create datasource response, Error: %s", err.Error())
			}
			respDict := map[string]interface{}{}
			if errUnmarshall := json.Unmarshal(bodyBytes, &respDict); errUnmarshall != nil {
				return nil, fmt.Errorf("Unable to decode create datasource response")
			}

			return respDict, nil
		} else {
			return nil, err
		}
	}
}

func getDataSourceByName(name string) map[string]interface{} {
	if v, err := getDatasources(); err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	} else {
		for _, ds := range v {
			if ds["name"] == name {
				return ds
			}
		}
	}
	return nil
}

func listDataSources() {
	if v, err := getDatasources(); err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	} else {
		count := len(v)
		fmt.Println("(", count, ") datasources found")
		if count > 0 {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Type", "Id"})
			for _, datasource := range v {
				name := datasource["name"].(string)
				typ := datasource["type"].(string)
				id := (int)(datasource["id"].(float64))
				table.Append([]string{name, typ,
					fmt.Sprintf("%d", id)})
			}
			table.Render()
		}
	}
}

func getDatasources() ([]map[string]interface{}, error) {
	uri := getUrlForResource(ResourceDatasourceAll)
	client := getHttpClient()

	if resp, err := client.Get(uri); err == nil {
		defer resp.Body.Close()
		var v = []map[string]interface{}{}
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("Unable to fetch datasources, Error: %s", err.Error())
			}
			err = json.Unmarshal(bodyBytes, &v)
			if err != nil {
				return nil, fmt.Errorf("Unable to decode datasources, Error: %s", err.Error())
			} else {
				return v, nil
			}
		} else {
			return nil, fmt.Errorf("Http response error, Error: %d", resp.StatusCode)
		}
	} else {
		return nil, fmt.Errorf("Unable to fetch datasources, Error: %s", err.Error())
	}
}
