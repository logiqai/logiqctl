# `logiqctl config`

Modify logiqctl configuration options

## Synopsis


Configure  LOGIQ CLI (logiqctl) options. To get token see https://docs.logiq.ai/vewing-logs/logiqctl/obtaining-api-key

!!! tip "Note"
    The values you provide will be written to the config file located at (~/.logiqctl)


## Examples

```

View current context
	logiqctl config view

Runs an interactive prompt and let user configure
	logiqctl config init

Set default cluster
	logiqctl config set-cluster END-POINT

Set default context
	logiqctl config set-context namespace

Runs an interactive prompt and let user select namespace from the list
	logiqctl config set-context i

Set token
	logiqctl config set-token api_token

```

## Options

```
  -h, --help   help for config
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

* [logiqctl]()	 - Logiqctl - CLI for Logiq Observability stack
* [logiqctl config init](/config/logiqctl_config_init)	 - Interactive configuration command
* [logiqctl config set-cluster](/config/logiqctl_config_set-cluster)	 - Sets the logiq cluster end-point
* [logiqctl config set-context](/config/logiqctl_config_set-context)	 - Sets the default context or namespace.
* [logiqctl config set-credential](/config/logiqctl_config_set-credential)	 - Sets logiq credentials
* [logiqctl config set-token](/config/logiqctl_config_set-token)	 - Sets a logiq ui api token
* [logiqctl config view](/config/logiqctl_config_view)	 - View current defaults

