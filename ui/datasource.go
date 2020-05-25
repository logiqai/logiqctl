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
"github.com/TylerBrock/colorjson"
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
			listDatasources()
		},
	})

	return cmd
}

func printDataSource(args []string) {
	if v, vErr := getDatasource(args); v != nil {
		f := colorjson.NewFormatter()
		f.Indent = 2
		s, _ := f.Marshal(*v)
		fmt.Println(string(s))
	} else {
		fmt.Println(vErr.Error())
		os.Exit(-1)
	}
}

func getDatasource(args []string) (*map[string]interface{}, error) {
	uri := getUrlForResource(ResourceDatasource,args...)
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

func listDatasources() {
	uri := getUrlForResource(ResourceDatasourceAll)
	client := getHttpClient()

	if resp, err := client.Get(uri); err == nil {
		defer resp.Body.Close()
		var v = []map[string]interface{}{}
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Unable to fetch datasources, Error: %s", err.Error())
				os.Exit(-1)
			}
			err = json.Unmarshal(bodyBytes, &v)
			if err != nil {
				fmt.Printf("Unable to decode datasources, Error: %s", err.Error())
			} else {
				count := len(v)
				fmt.Println("(",count, ") datasources found")
				if ( count > 0 ) {
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"Name", "Type", "Id"})
					for _, datasource := range v {
						name := datasource["name"].(string)
						typ := datasource["type"].(string)
						id := (int)(datasource["id"].(float64))
						table.Append([]string{name, typ,
							fmt.Sprintf("%d",id),})
					}

					table.Render()
				}
			}
		}
	} else {
		fmt.Printf("Unable to fetch datasources, Error: %s", err.Error())
		os.Exit(-1)
	}
}
