## logiqctl logs search

Search for given text in logs

### Synopsis

Search for given text in logs

```
logiqctl logs search [flags]
```

### Options

```
  -h, --help   help for search
```

### Options inherited from parent commands

```
  -c, --cluster string       Override the default cluster set by `logiqctl set-cluster' command
  -f, --follow               Specify if the logs should be streamed.
  -n, --namespace string     Override the default context set by `logiqctl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. (default "table")
      --page-size uint32     Number of log entries to return in one page (default 30)
  -s, --since string         Only return logs newer than a relative duration like 2m, 3h, or 2h30m.
                             This is in relative to the last seen log time for a specified application or processes. (default "1h")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339. 
                             This is only applicable when the output format is table. json and yaml outputs will have time in epoch seconds.
                             json output is not indented, use '| jq' for advanced json operations (default "relative")
```

### SEE ALSO

* [logiqctl logs](logiqctl_logs.md)	 - Print the logs for an application or process

