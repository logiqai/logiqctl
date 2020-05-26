package ui

import (
	"bytes"
	"fmt"
	"github.com/logiqai/logiqctl/utils"
	"github.com/spf13/viper"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

var (
	client *http.Client = nil
)

func getHttpClient() *http.Client {
	if client != nil {
		return client
	}

	api_key := viper.GetString(utils.KeyUiToken)

	if api_key != "" {
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
			q.Add("remember", "on")
			q.Add("email", user)
			q.Add("password", password)
			u.RawQuery = q.Encode()

			if resp, err := client.Post(u.String(), "application/x-www-form-urlencoded", bytes.NewReader(([]byte)(q.Encode()))); err != nil {
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
