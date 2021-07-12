# `logiqctl get`

Display one or many resources

## Synopsis

Prints a table of the most important information about the specified resources

## Examples

```text
List all applications for the selected context
logiqctl get applications

List all applications for all the available context
logiqctl get applications all

List all dashboards
logiqctl get dashboards all

Get dashboard
logiqctl get dashboard dashboard-slug

List events for the available namespace
logiqctl get events

List events for all the available namespaces
logiqctl get events all

List or export eventrules
logiqctl get eventrules

List all namespaces
logiqctl get namespaces

List all processes
logiqctl get processes

List all queries
logiqctl get queries all

Get query
logiqctl get query query-slug
```

## Options

```text
  -h, --help   help for get
```

## Options inherited from parent commands

```text
  -c, --cluster string       Override the default cluster set by `logiqctl set-cluster' command
  -n, --namespace string     Override the default context set by `logiqctl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. 
                             json output is not indented, use '| jq' for advanced json operations (default "table")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339. 
                             This is only applicable when the output format is table. json and yaml outputs will have time in epoch seconds. (default "relative")
```

## SEE ALSO

* [logiqctl](/)     - Logiqctl - CLI for Logiq Observability stack
* [logiqctl get applications](/get/logiqctl_get_applications)     - List all the available applications in default namespace
* [logiqctl get dashboard](/get/logiqctl_get_dashboard)     - Get a dashboard
* [logiqctl get datasource](/get/logiqctl_get_datasource)     - Get a datasource
* [logiqctl get eventrules](/get/logiqctl_get_eventrules)     - List event rules
* [logiqctl get events](/get/logiqctl_get_events)     - List all the available events for the namespace
* [logiqctl get namespaces](/get/logiqctl_get_namespaces)     - List the available namespaces
* [logiqctl get processes](/get/logiqctl_get_processes)     - List all the available processes, runs an interactive prompt to select applications
* [logiqctl get query](/get/logiqctl_get_query)     - Get a query

