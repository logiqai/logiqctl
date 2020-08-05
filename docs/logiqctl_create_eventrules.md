## logiqctl create eventrules

Create event rules

### Synopsis

Create event rules

```
logiqctl create eventrules [flags]
```

### Examples

```
logiqctl create eventrules -f <path to event rules file>
```

### Options

```
  -f, --file string   Path to file
  -h, --help          help for eventrules
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

