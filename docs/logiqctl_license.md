## logiqctl license

View or update license

### Synopsis


Logiq deployment comes configured with a 30-day trial license. Obtain a valid license by contacting LOGIQ at license@logiq.ai.
With this command, the user can view or apply a license. 


### Examples

```

Upload your LOGIQ deployment license
- logiqctl license set -f license.jws

View License information
 - logiqctl license get 

```

### Options

```
  -h, --help   help for license
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
* [logiqctl license get](logiqctl_license_get.md)	 - View license information
* [logiqctl license set](logiqctl_license_set.md)	 - Configure license for LOGIQ deployment

