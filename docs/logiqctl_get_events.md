## logiqctl get events

List all the available events for the namespace

### Synopsis

List all the available events for the namespace

```
logiqctl get events [flags]
```

### Examples

```

List last 30 events
- logiqctl get events|e

List events by application 
- logiqctl get events -a=sshd


```

### Options

```
  -a, --application string   Filter events by application name
  -h, --help                 help for events
  -p, --process string       Filter events by application and process name
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

