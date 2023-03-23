package utils

import (
	"crypto/tls"
	"fmt"
	"github.com/logiqai/logiqctl/defines"
	"net"
	"time"
)

var (
	protocol defines.UriProtocol = defines.UriUnknown
)

func getProtocol(ipOrDns string) defines.UriProtocol {
	if protocol != defines.UriUnknown {
		return protocol
	}
	if defines.Protocol == defines.UriUnknown {
		conf := &tls.Config{
			InsecureSkipVerify: true,
		}
		dialer := net.Dialer{Timeout: time.Duration(time.Second)}
		conn, err := tls.DialWithDialer(&dialer, "tcp", fmt.Sprintf("%s:443", ipOrDns), conf)
		if err != nil {
			//fmt.Println("Uri is http")
			protocol = defines.UriHttp
			return defines.UriHttp
		}
		defer conn.Close()
		protocol = defines.UriHttps
		return defines.UriHttps
	} else {
		return defines.Protocol
	}
}

func GetUrlForResource(r defines.Resource, args ...string) string {
	var uri string
	switch r {
	case defines.ResourceDashboardsAll:
		uri = fmt.Sprintf("api/dashboards")
	case defines.ResourceDashboardsGet:
		uri = fmt.Sprintf("api/dashboards/%s", args[0])
	case defines.ResourceQueryAll:
		uri = fmt.Sprintf("api/queries")
	case defines.ResourceQuery:
		uri = fmt.Sprintf("api/queries/%s", args[0])
	case defines.ResourceQueryResultGet:
		uri = fmt.Sprintf("api/query_results/%s", args[0])
	case defines.ResourceQueryResult:
		uri = fmt.Sprintf("api/query_results")
	case defines.ResourceJobGet:
		uri = fmt.Sprintf("api/jobs/%s", args[0])
	case defines.ResourceDatasourceAll:
		uri = fmt.Sprintf("api/data_sources")
	case defines.ResourceDatasource:
		uri = fmt.Sprintf("api/data_sources/%s", args[0])
	case defines.ResourceVisualizationAll:
		uri = fmt.Sprintf("api/visualizations")
	case defines.ResourceVizualization:
		uri = fmt.Sprintf("api/visualizations/%s", args[0])
	case defines.ResourceWidgetAll:
		uri = fmt.Sprintf("api/widgets")
	case defines.ResourceWidget:
		uri = fmt.Sprintf("api/widgets/%s", args[0])
	case defines.ResourceLogin:
		uri = fmt.Sprintf("login")
	case defines.ResourceJWTToken:
		uri = fmt.Sprintf("token")
	case defines.ResourcePrometheusProxy:
		uri = fmt.Sprintf("api/logiq_proxy")
	case defines.ResourceAlertsAll:
		uri = fmt.Sprintf("api/alerts")
	case defines.ResourceAlert:
		uri = fmt.Sprintf("api/alerts/%s", args[0])
	case defines.ResourceForwardersAll:
		uri = fmt.Sprintf("v1/forwards")
	}

	return uri
}
