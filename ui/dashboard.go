package ui

import (
	"bytes"
	"encoding/json"
	"github.com/logiqai/logiqctl/utils"
	"github.com/spf13/cobra"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
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
			getDashboard(args)
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:     "all",
		Example: "logiqctl get dashboards all",
		Short:   "List all the available dashboards",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			listDashboards()
		},
	})

	return cmd
}
func getHttpClient() *http.Client {
	var client *http.Client
	api_key := viper.GetString(utils.KeyUiToken)

	if  api_key != "" {
		client = &http.Client{}
	} else {
		user := viper.GetString(utils.KeyUiUser)
		password := viper.GetString(utils.KeyUiPassword)

		if user != "" && password != "" {
			cookieJar, _ := cookiejar.New(nil)

			client = &http.Client{
				Jar: cookieJar,
			}
			loginUrl := getUrlForResource(ResourceLogin)
			u, _ := url.Parse(loginUrl)
			q, _ := url.ParseQuery(u.RawQuery)
			q.Add("remember","on")
			q.Add("email",user)
			q.Add("password",password)
			u.RawQuery = q.Encode()

			if resp, err := client.Post(u.String(),"application/x-www-form-urlencoded",bytes.NewReader(([]byte)(q.Encode()))); err != nil {
				fmt.Println("Error login with provided credentials, Error:", err.Error())
				os.Exit(-1)
			} else {
				defer resp.Body.Close()
			}
		} else {
			fmt.Println("api token or ui credentials must be set. See \"logiqctl config help\" for more details")
			os.Exit(-1)
		}
	}

	return client
}

func getDashboard(args []string) {
	uri := getUrlForResource(ResourceDashboardsGet,args...)
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
				fmt.Println(string(bodyBytes))
			}
		}
	} else {
		fmt.Printf("Unable to fetch dashboards, Error: %s", err.Error())
		os.Exit(-1)
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
