## logiqctl tail

Stream logs from LOGIQ Observability Stack

### Synopsis


'logiqctl tail' is similar to tail -f command, allows you to view the log data that is being sent to LOGIQ Observability Stack in real-time. You can see logs from the cluster at multiple levels. 'tail' without any option runs an interactive prompt and let you choose application and process in the current context. 


```
logiqctl tail [flags]
```

### Examples

```

Tail logs 
- logiqctl tail

```

### Options

```
  -f, --file string         Path to file
  -h, --help                help for tail
  -m, --max-file-size int   Max output file size (default 10)
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

* [logiqctl](logiqctl.md)	 - Logiqctl - CLI for Logiq Observability stack

