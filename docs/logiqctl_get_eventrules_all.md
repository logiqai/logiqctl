## logiqctl get eventrules all

List all event rules

### Synopsis

List all event rules

```
logiqctl get eventrules all [flags]
```

### Examples

```
logiqctl get eventrules all
```

### Options

```
  -h, --help   help for all
```

### Options inherited from parent commands

```
  -c, --cluster string       Override the default cluster set by `logiqctl set-cluster' command
  -w, --file string          Path to file to be written to
  -g, --groups string        list of groups separated by comma
  -n, --namespace string     Override the default context set by `logiqctl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. (default "table")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339. 
                             This is only applicable when the output format is table. json and yaml outputs will have time in epoch seconds.
                             json output is not indented, use '| jq' for advanced json operations (default "relative")
```

### SEE ALSO

* [logiqctl get eventrules](logiqctl_get_eventrules.md)	 - List event rules

