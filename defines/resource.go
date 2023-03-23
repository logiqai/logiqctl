package defines

type Resource int

var (
	Protocol UriProtocol = UriUnknown
)

const (
	ResourceDashboardsAll    Resource = iota
	ResourceDashboardsGet    Resource = iota
	ResourceLogin                     = iota
	ResourceQuery            Resource = iota
	ResourceQueryAll         Resource = iota
	ResourceDatasource       Resource = iota
	ResourceDatasourceAll    Resource = iota
	ResourceVizualization    Resource = iota
	ResourceVisualizationAll Resource = iota
	ResourceWidget           Resource = iota
	ResourceWidgetAll        Resource = iota
	ResourceJWTToken         Resource = iota
	ResourceQueryResultGet   Resource = iota
	ResourcePrometheusProxy  Resource = iota
	ResourceAlertsAll        Resource = iota
	ResourceAlert            Resource = iota
	ResourceQueryResult      Resource = iota
	ResourceJobGet           Resource = iota
	ResourceForwardersAll    Resource = iota
)

type UriProtocol int

const (
	UriUnknown             = iota
	UriHttp    UriProtocol = iota
	UriHttps   UriProtocol = iota
)
