# `logiqctl config set-token`

Sets a logiq ui api token

## Synopsis

Sets the cluster UI api token, a valid logiq cluster end point is also required for all the operations

```text
logiqctl config set-token [flags]
```

## Examples

```text
logiqctl set-token api_token
```

## Options

```text
  -h, --help   help for set-token
```

## Options inherited from parent commands

```text
  -c, --cluster string       Override the default cluster set by `logiqctl set-cluster' command
  -n, --namespace string     Override the default context set by `logiqctl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. 
                             json output is not indented, use '| jq' for advanced json operations (default "table")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339. 
                             This is only applicable when the output format is table. json and yaml outputs will have time in epoch seconds. (default "relative")
```

## SEE ALSO

* [logiqctl config](/config/logiqctl_config)     - Modify logiqctl configuration options

