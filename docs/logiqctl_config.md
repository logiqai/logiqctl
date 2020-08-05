## logiqctl config

Modify logiqctl configuration options

### Synopsis


Configure  LOGIQ CLI (logiqctl) options. 
Note: The values you provide will be written to the config file located at (~/.logiqctl)


### Examples

```

View current context
	logiqctl config view

Set default cluster
	logiqctl config set-cluster END-POINT

Set default context
	logiqctl config set-context namespace

Runs an interactive prompt and let user select namespace from the list
	logiqctl config set-context i

Set Credential
	logiqctl config set-ui-credential user password

```

### Options

```
  -h, --help   help for config
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
* [logiqctl config init](logiqctl_config_init.md)	 - Interactive configuration command
* [logiqctl config set-cluster](logiqctl_config_set-cluster.md)	 - Sets the logiq cluster end-point
* [logiqctl config set-context](logiqctl_config_set-context.md)	 - Sets the default context or namespace.
* [logiqctl config set-credential](logiqctl_config_set-credential.md)	 - Sets logiq credentials
* [logiqctl config view](logiqctl_config_view.md)	 - View current defaults

