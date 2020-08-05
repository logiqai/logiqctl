## logiqctl config set-context interactive

Runs an interactive prompt and let user select namespace from the list

### Synopsis

Runs an interactive prompt and let user select namespace from the list

```
logiqctl config set-context interactive [flags]
```

### Options

```
  -h, --help   help for interactive
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

* [logiqctl config set-context](logiqctl_config_set-context.md)	 - Sets the default context or namespace.

