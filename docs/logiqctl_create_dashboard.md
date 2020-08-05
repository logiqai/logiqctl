## logiqctl create dashboard

Create a dashboard

### Synopsis


The crowd-sourced dashboards available in https://github.com/logiqai/logiqhub can be downloaded and applied to any clusters. 
One can also export dashboards created using "logiqctl get dashboard" command and apply on different clusters.


```
logiqctl create dashboard [flags]
```

### Examples

```
logiqctl create dashboard|d -f <path to dashboard spec>
```

### Options

```
  -f, --file string   Path to file
  -h, --help          help for dashboard
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

* [logiqctl create](logiqctl_create.md)	 - Create a resource

