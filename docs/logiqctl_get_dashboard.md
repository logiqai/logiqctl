## logiqctl get dashboard

Get a dashboard

### Synopsis

Get a dashboard

```
logiqctl get dashboard [flags]
```

### Examples

```
logiqctl get dashboard|d <dashboard-slug>
```

### Options

```
  -h, --help   help for dashboard
```

### Options inherited from parent commands

```
  -c, --cluster string       Override the default cluster set by `logiqctl set-cluster' command
  -n, --namespace string     Override the default context set by `logiqctl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. 
                             json output is not indented, use '| jq' for advanced json operations (default "table")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339. 
                             This is only applicable when the output format is table. json and yaml outputs will have time in epoch seconds. (default "relative")
```

### SEE ALSO

* [logiqctl get](logiqctl_get.md)	 - Display one or many resources
* [logiqctl get dashboard all](logiqctl_get_dashboard_all.md)	 - List all the available dashboards

