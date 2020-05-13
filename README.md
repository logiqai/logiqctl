# Logiqctl 

### CLI for Logiq Observability stack

```bash

> logiqctl 

The LOGIQ command line toolkit, logiqctl, allows you to run commands against LOGIQ Observability stack.
- You can tail logs from your applications and servers
- View available namespaces/applications/processes
- Query historical data
- Search your log data.
- View Events


Find more information at: https://docs.logiq.ai/logiqctl/logiq-box

Usage:
  logiqctl [command]

Available Commands:
  config      Modify logiqctl configuration options
  get         Display one or many resources
  help        Help about any command
  logs        Print the logs for an application or process
  tail        Stream logs from logiq Observability stack

Flags:
  -c, --cluster string       Override the default cluster set by `logiqctl set-cluster' command
  -h, --help                 help for logiqctl
  -n, --namespace string     Override the default context set by `logiqctl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. (default "table")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339.
                             This is only applicable when the output format is table. json and yaml outputs will have time in epoch seconds.
                             json output is not indented, use '| jq' for advanced json operations (default "relative")
      --version              version for logiqctl

Use "logiqctl [command] --help" for more information about a command.

```


### Quick start
The simplest way to try logiqctl is to download a pre-built binary from our release page:
https://github.com/logiqai/logiqctl/releases

Once you have the binary run the following to get started
- run `logiqctl config set-cluster CLUSTER_URL`
- run `logiqctl config set-context NAMESPACE`
- run `logiqctl get namespaces` to verify


#### How to build from source

**Requirements**
- Install Go [https://golang.org/dl/]
- Install protoc [https://github.com/protocolbuffers/protobuf/releases]
    
```bash
# For MAC
PROTOC_ZIP=protoc-3.7.1-osx-x86_64.zip
curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/$PROTOC_ZIP
sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
rm -f $PROTOC_ZIP
```
   
```bash
# For Linux
PROTOC_ZIP=protoc-3.7.1-linux-x86_64.zip
curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/$PROTOC_ZIP
sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
rm -f $PROTOC_ZIP

```
- run `mkdir -p $GOPATH/src/github.com/logiqai`
- run `cd $GOPATH/src/github.com/logiqai`
- run `git clone git@github.com:logiqai/logiqctl.git`
- run `cd logiqctl`
- run `./generate_grpc.sh `
- run `go build logiqctl.go`


To know more about LOGIQ Observability stack, see https://logiq.ai/ and https://docs.logiq.ai/ 

In case of any issues either email us at cli@logiq.ai or you could create an issue in this repository.

 ### logiqctl get
 
 ```bash
> logiqctl get --help

Prints a table of the most important information about the specified resources. For example:

# List all namespaces
logiqctl get namespaces

# List all applications for the selected context
logiqctl get applications

# List all applications for all the available context
logiqctl get applications all

# List all processes
logiqctl get processes

Usage:
  logiqctl get [command]

Available Commands:
  applications List all the available applications in default namespace
  events       List all the available events for the namespace
  namespaces   List the available name spaces
  processes    List all the available processes, runs an interactive prompt to select applications

Flags:
  -h, --help   help for get

Global Flags:
  -c, --cluster string       Override the default cluster set by `logiqctl set-cluster' command
  -n, --namespace string     Override the default context set by `logiqctl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. (default "table")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339.
                             This is only applicable when the output format is table. json and yaml outputs will have time in epoch seconds.
                             json output is not indented, use '| jq' for advanced json operations (default "relative")

Use "logiqctl get [command] --help" for more information about a command. 

```

### logiqctl logs

```bash
> logiqctl logs --help
Logs expect a namespace and application to be available to return results.
Set the default namespace using 'logiqctl set-context' command or pass as '-n=NAMESPACE' flag
Application name needs to be passed as an argument to the command.
If the user is unsure of the application name, they can run an interactive prompt 
the would help them to choose filters. See examples below.

Search command searches at namespace level, flag -p is ignored.

Global flag '--time-format' is not applicable for this command.
Global flag '--output' only supports json format for this command.

Usage:
  logiqctl logs [flags]
  logiqctl logs [command]

Aliases:
  logs, log

Examples:

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



Available Commands:
  interactive Runs an interactive prompt to let user select application and filters
  search      Search for given test in logs

Flags:
  -f, --follow           Specify if the logs should be streamed.
  -h, --help             help for logs
      --page-size int    Number of log entries to return in one page (default 30)
  -p, --process string   Filter logs by  proc id
  -s, --since string     Only return logs newer than a relative duration like 2m, 3h, or 2h30m.
                         This is in relative to the last seen log time for a specified application or processes. (default "1h")

Global Flags:
  -c, --cluster string       Override the default cluster set by `logiqctl set-cluster' command
  -n, --namespace string     Override the default context set by `logiqctl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. (default "table")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339.
                             This is only applicable when the output format is table. json and yaml outputs will have time in epoch seconds.
                             json output is not indented, use '| jq' for advanced json operations (default "relative")

Use "logiqctl logs [command] --help" for more information about a command.

```