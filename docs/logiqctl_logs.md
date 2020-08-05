## logiqctl logs

Print the logs for an application or process

### Synopsis

Logs expect a namespace and application to be available to return results.
Set the default namespace using 'logiqctl set-context' command or pass as '-n=NAMESPACE' flag
Application name needs to be passed as an argument to the command. 
If the user is unsure of the application name, they can run an interactive prompt the would help them to choose filters.  
See examples below. 

Search command searches at namespace level, flags -p is ignored. 

Global flag '--time-format' is not applicable for this command.
Global flag '--output' only supports json format for this command.

```
logiqctl logs [flags]
```

### Examples

```

Print logs for logiq-flash ingest server
# logiqctl logs logiq-flash

Print logs in json format
# logiqctl -o=json logs logiq-flash

Print logs for logiq-flash ingest server filtered by process logiq-flash-2
In case of Kubernetes deployment a Stateful Set is an application, and each pods in it is a process
The --process (-p) flag lets you view logs for the individual pod
# logiqctl logs -p=logiq-flash-2 logiq-flash

Runs an interactive prompt to let user choose filters
# logiqctl logs i

Search logs for the given text
# logiqctl logs search "your search term"   

If the flag --follow (-f) is specified the logs will be streamed till it over. 


```

### Options

```
  -f, --follow             Specify if the logs should be streamed.
  -h, --help               help for logs
      --page-size uint32   Number of log entries to return in one page (default 30)
  -p, --process string     Filter logs by  proc id
  -s, --since string       Only return logs newer than a relative duration like 2m, 3h, or 2h30m.
                           This is in relative to the last seen log time for a specified application or processes. (default "1h")
```

### Options inherited from parent commands

```
  -c, --cluster string       Override the default cluster set by `logiqctl set-cluster' command
  -n, --namespace string     Override the default context set by `logiqctl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. (default "table")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339. 
                             This is only applicable when the output format is table. json and yaml outputs will have time in epoch seconds.
                             json output is not indented, use '| jq' for advanced json operations (default "relative")
```

### SEE ALSO

* [logiqctl](logiqctl.md)	 - Logiqctl - CLI for Logiq Observability stack
* [logiqctl logs interactive](logiqctl_logs_interactive.md)	 - Runs an interactive prompt to let user select application and filters
* [logiqctl logs search](logiqctl_logs_search.md)	 - Search for given text in logs

