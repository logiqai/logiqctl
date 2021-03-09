package ui

import (
	"crypto/tls"
	"fmt"
	"github.com/logiqai/logiqctl/utils"
	"github.com/spf13/viper"
	"net"
	"time"
)

var (
	protocol UriProtocol = UriUnknown
)

func getProtocol(ipOrDns string) UriProtocol {
	if protocol != UriUnknown {
		return protocol
	}
	if Protocol == UriUnknown {
		conf := &tls.Config{
			InsecureSkipVerify: true,
		}
		dialer := net.Dialer{Timeout: time.Duration(time.Second)}
		conn, err := tls.DialWithDialer(&dialer, "tcp", fmt.Sprintf("%s:443", ipOrDns), conf)
		if err != nil {
			//fmt.Println("Uri is http")
			protocol = UriHttp
			return UriHttp
		}
		defer conn.Close()
		protocol = UriHttps
		return UriHttps
	} else {
		return Protocol
	}
}

func GetUrlForResource(r Resource, args ...string) string {
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
		uri = fmt.Sprintf("%s://%s/api/dashboards", protocolString, ipOrDns)
	case ResourceDashboardsGet:
		uri = fmt.Sprintf("%s://%s/api/dashboards/%s", protocolString, ipOrDns, args[0])
	case ResourceQueryAll:
		uri = fmt.Sprintf("%s://%s/api/queries", protocolString, ipOrDns)
	case ResourceQuery:
		uri = fmt.Sprintf("%s://%s/api/queries/%s", protocolString, ipOrDns, args[0])
	case ResourceQueryResult:
		uri = fmt.Sprintf("%s://%s/api/query_results/%s", protocolString, ipOrDns, args[0])
	case ResourceDatasourceAll:
		uri = fmt.Sprintf("%s://%s/api/data_sources", protocolString, ipOrDns)
	case ResourceDatasource:
		uri = fmt.Sprintf("%s://%s/api/data_sources/%s", protocolString, ipOrDns, args[0])
	case ResourceVisualizationAll:
		uri = fmt.Sprintf("%s://%s/api/visualizations", protocolString, ipOrDns)
	case ResourceVizualization:
		uri = fmt.Sprintf("%s://%s/api/visualizations/%s", protocolString, ipOrDns, args[0])
	case ResourceWidgetAll:
		uri = fmt.Sprintf("%s://%s/api/widgets", protocolString, ipOrDns)
	case ResourceWidget:
		uri = fmt.Sprintf("%s://%s/api/widgets/%s", protocolString, ipOrDns, args[0])
	case ResourceLogin:
		uri = fmt.Sprintf("%s://%s/login", protocolString, ipOrDns)
	case ResourceJWTToken:
		uri = fmt.Sprintf("%s://%s/token", protocolString, ipOrDns)
	}

	return uri
}
