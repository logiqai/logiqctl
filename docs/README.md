## logiqctl

Logiqctl - CLI for Logiq Observability stack

### Synopsis


The LOGIQ command line toolkit, logiqctl, allows you to run commands against LOGIQ Observability stack. 
- Real-time streaming of logs
- Query historical application logs 
- Search your log data.
- View Events
- Manage Dashboards
- Create event rules
- Manage license


Find more information at: https://docs.logiq.ai/logiqctl/logiq-box



### Options

```
  -c, --cluster string       Override the default cluster set by `logiqctl set-cluster' command
  -h, --help                 help for logiqctl
  -n, --namespace string     Override the default context set by `logiqctl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. 
                             json output is not indented, use '| jq' for advanced json operations (default "table")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339. 
                             This is only applicable when the output format is table. json and yaml outputs will have time in epoch seconds. (default "relative")
```

### SEE ALSO

* [logiqctl config](logiqctl_config.md)	 - Modify logiqctl configuration options
* [logiqctl create](logiqctl_create.md)	 - Create a resource
* [logiqctl get](logiqctl_get.md)	 - Display one or many resources
* [logiqctl license](logiqctl_license.md)	 - View or update license
* [logiqctl logs](logiqctl_logs.md)	 - View logs for the given namespace and application

