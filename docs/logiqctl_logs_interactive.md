## logiqctl logs interactive

Runs an interactive prompt to let the user select application and filters

### Synopsis

Runs an interactive prompt to let the user select application and filters

```
logiqctl logs interactive [flags]
```

### Options

```
  -h, --help   help for interactive
```

### Options inherited from parent commands

```
  -c, --cluster string       Override the default cluster set by `logiqctl set-cluster' command
  -f, --follow               Specify if the logs should be streamed.
  -n, --namespace string     Override the default context set by `logiqctl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. 
                             json output is not indented, use '| jq' for advanced json operations (default "table")
      --page-size uint32     Number of log entries to return in one page (default 30)
  -s, --since string         Only return logs newer than a relative duration. This is in relative to the last
                             seen log time for a specified application or processes within the namespace.
                             A duration string is a possibly signed sequence of decimal numbers, each with optional
                             fraction and a unit suffix, such as "3h34m", "1.5h" or "24h". Valid time units are "s", "m", "h" (default "1h")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339. 
                             This is only applicable when the output format is table. json and yaml outputs will have time in epoch seconds. (default "relative")
```

### SEE ALSO

* [logiqctl logs](logiqctl_logs.md)	 - View logs for the given namespace and application

