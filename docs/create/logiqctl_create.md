# `logiqctl create`

Create a resource on LOGIQ. 

## Synopsis

Creates a resource from a resource specification

## Examples

```text
Create a dashboard
logiqctl create dashboard -f <path to dashboard_spec_file.json>

Create eventrules
logiqctl create eventrules -f <path to eventrules_file.json>
```

## Options

```text
  -h, --help   help for create
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
* [logiqctl create dashboard](/create/logiqctl_create_dashboard)     - Create a dashboard
* [logiqctl create eventrules](/create/logiqctl_create_eventrules)     - Create event rules

