package ui

import (
	"crypto/tls"
	"github.com/logiqai/logiqctl/utils"
	"github.com/spf13/viper"
	"fmt"
	"net/url"
)

func addApiToken(uriNonTokenized string) string {
	u, _ := url.Parse(uriNonTokenized)
	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("api_key",viper.GetString(utils.KeyUiToken))
	u.RawQuery = q.Encode()

	return u.String()
}

func getProtocol(ipOrDns string) UriProtocol {
	if Protocol == UriUnknown {
		conf := &tls.Config{
			InsecureSkipVerify: true,
		}

		conn, err := tls.Dial("tcp", fmt.Sprintf("%s:443",ipOrDns), conf)
		if err != nil {
			return UriHttp
		}
		defer conn.Close()
		return UriHttps
	} else {
		return Protocol
	}
}

func getUrlForResource(r Resource, args...string) string {
	var uri string
	var protocolString string
	ipOrDns := viper.GetString(utils.KeyCluster)
	protocol := getProtocol(ipOrDns)

	if protocol == UriHttp {
		protocolString = "http"
	} else {
		protocolString = "https"
	}

	switch r {
	case ResourceDashboardsAll:
		uri = fmt.Sprintf("%s://%s/api/dashboards",protocolString, ipOrDns)
	case ResourceDashboardsGet:
		uri = fmt.Sprintf("%s://%s/api/dashboards/%s",protocolString, ipOrDns,args[0])
	case ResourceQueryAll:
		uri = fmt.Sprintf("%s://%s/api/queries",protocolString, ipOrDns)
	case ResourceQuery:
		uri = fmt.Sprintf("%s://%s/api/queries/%s",protocolString, ipOrDns,args[0])
	case ResourceLogin:
		uri = fmt.Sprintf("%s://%s/login",protocolString,ipOrDns)
	}

	api_key := viper.GetString(utils.KeyUiToken)
	if api_key != "" {
		uri = addApiToken(uri)
	}

	return uri
}
