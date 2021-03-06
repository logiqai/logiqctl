# `logiqctl get`

Display one or more of your LOGIQ resources

## Synopsis

Prints a table that displays the most important information about the LOGIQ resource you specify.

## Examples

```

List all applications for the selected context
logiqctl get applications

List all applications for all available contexts
logiqctl get applications all

List all dashboards
logiqctl get dashboards all

Get dashboard
logiqctl get dashboard dashboard-slug

List events for the available namespace
logiqctl get events

List events for all available namespaces
logiqctl get events all

List or export eventrules
logiqctl get eventrules

List all namespaces
logiqctl get namespaces

List all processes
logiqctl get processes

List all queries
logiqctl get queries all

Get httpingestkey
logiqctl get httpingestkey

Get query
logiqctl get query query-slug

```

## Options

```
  -h, --help   help for get
```

## Options inherited from parent commands

```
  -c, --cluster string       Override the default cluster set by `logiqctl set-cluster' command
  -n, --namespace string     Override the default context set by `logiqctl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. 
                             JSON output is not indented, use '| jq' for advanced JSON operations (default "table")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339. 
                             This is only applicable when the output format is table. JSON and YAML outputs will have time in epoch seconds. (default "relative")
```

## SEE ALSO

* [logiqctl](/)	 - Logiqctl - CLI for Logiq Observability stack
* [logiqctl get applications](/get/logiqctl_get_applications)	 - List all available applications within the default namespace
* [logiqctl get dashboard](/get/logiqctl_get_dashboard)	 - Get a dashboard
* [logiqctl get datasource](/get/logiqctl_get_datasource)	 - Get a datasource
* [logiqctl get eventrules](/get/logiqctl_get_eventrules)	 - List event rules
* [logiqctl get events](/get/logiqctl_get_events)	 - List all available events for the namespace
* [logiqctl get httpingestkey](/get/logiqctl_get_httpingestkey)	 - Get httpingestkey
* [logiqctl get namespaces](/get/logiqctl_get_namespaces)	 - List all available namespaces
* [logiqctl get processes](/get/logiqctl_get_processes)	 - List all available processes. This command runs an interactive prompt that lets you choose an application from a list of available applications.
* [logiqctl get query](/get/logiqctl_get_query)	 - Get a query

