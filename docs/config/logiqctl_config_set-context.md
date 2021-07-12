# `logiqctl config set-context`

Sets the default context or namespace.

## Synopsis


A context or namespace is required for all the logiqctl operations. To override the default namespace for individual command use flag `-n`. 
		

```
logiqctl config set-context [flags]
```

## Examples

```
set-context <namespace name>
```

## Options

```
  -h, --help   help for set-context
```

## Options inherited from parent commands

```
  -c, --cluster string       Override the default cluster set by `logiqctl set-cluster' command
  -n, --namespace string     Override the default context set by `logiqctl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. 
                             json output is not indented, use '| jq' for advanced json operations (default "table")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339. 
                             This is only applicable when the output format is table. json and yaml outputs will have time in epoch seconds. (default "relative")
```

## SEE ALSO

* [logiqctl config](/config/logiqctl_config)	 - Modify logiqctl configuration options
* [logiqctl config set-context interactive](/config/logiqctl_config_set-context_interactive)	 - Runs an interactive prompt and let user select namespace from the list

