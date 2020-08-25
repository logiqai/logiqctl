package grpc_utils

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"

	"github.com/logiqai/logiqctl/ui"
	"github.com/logiqai/logiqctl/utils"
	"github.com/spf13/viper"
	"google.golang.org/grpc/metadata"
)

var (
	client      *http.Client    = nil
	grpcContext context.Context = nil
)

func GetGrpcContext() context.Context {
	if grpcContext != nil {
		return grpcContext
	}

	if url, cookieJar, err := GetCookies(); err != nil {
		fmt.Println("api token or ui credentials must be set. See \"logiqctl config help\" for more details")
		os.Exit(-1)
	} else {
		var cookieStr string
		for _, c := range cookieJar.Cookies(url) {
			if c.Name == "x-api-key" {
				cookieStr = fmt.Sprintf("%s=%s", c.Name, c.Value)
				break
			}
		}
		md := metadata.Pairs("grpcgateway-cookie", cookieStr)
		grpcContext = metadata.NewOutgoingContext(context.Background(), md)
	}
	return grpcContext
}

func GetCookies() (*url.URL, *cookiejar.Jar, error) {
	api_key := viper.GetString(utils.AuthToken)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	loginUrl := ui.GetUrlForResource(ui.ResourceLogin)
	u, _ := url.Parse(loginUrl)

	if api_key != "" {
		cookieJar, _ := cookiejar.New(nil)
		client = &http.Client{
			Jar: cookieJar,
		}
		req, err := http.NewRequest("GET", loginUrl, nil)
		if err != nil {
			fmt.Println("Unable to create login Request: ", err.Error())
			os.Exit(-1)
		}
		if api_key := viper.GetString(utils.AuthToken); api_key != "" {
			req.Header.Add("Authorization", fmt.Sprintf("Key %s", api_key))
		}
		if _, err := client.Do(req); err != nil {
			fmt.Println("Error login with provided token, Error:", err.Error())
			os.Exit(-1)
		} else {
			u, _ := url.Parse(loginUrl)
			return u, cookieJar, nil
		}
	} else {
		user := utils.GetUIUser()
		password := utils.GetUIPass()

		if user != "" && password != "" {
			cookieJar, _ := cookiejar.New(nil)

			client = &http.Client{
				Jar: cookieJar,
			}
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
				b, _ := ioutil.ReadAll(resp.Body)
				if match, _ := regexp.Match("Wrong email or password", b); match {
					fmt.Println("Error credentials")
					os.Exit(-1)
				}
				u, _ := url.Parse(loginUrl)
				return u, cookieJar, nil
			}
		} else {
			fmt.Println("api token or ui credentials must be set. See \"logiqctl config help\" for more details")
			os.Exit(-1)
		}
	}
	return nil, nil, fmt.Errorf("Error getting the token")
}
