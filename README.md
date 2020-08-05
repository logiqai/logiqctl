## logiqctl

Logiqctl - CLI for Logiq Observability stack

### Synopsis


The LOGIQ command line toolkit, logiqctl, allows you to run commands against LOGIQ Observability stack. 
- Real-time streaming of logs
- Query historical application logs 
- Search your log data.
- View Events
- Manage Dashboards
- Create event rules
- Manage license


Find more information at: https://docs.logiq.ai/logiqctl/logiq-box



### Options

```
  -c, --cluster string       Override the default cluster set by `logiqctl set-cluster' command
  -h, --help                 help for logiqctl
  -n, --namespace string     Override the default context set by `logiqctl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. 
                             json output is not indented, use '| jq' for advanced json operations (default "table")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339. 
                             This is only applicable when the output format is table. json and yaml outputs will have time in epoch seconds. (default "relative")
```

### SEE ALSO

* [logiqctl config](docs/logiqctl_config.md)	 - Modify logiqctl configuration options
* [logiqctl create](docs/logiqctl_create.md)	 - Create a resource
* [logiqctl get](docs/logiqctl_get.md)	 - Display one or many resources
* [logiqctl license](docs/logiqctl_license.md)	 - set and get license
* [logiqctl logs](docs/logiqctl_logs.md)	 - View logs for the given namespace and application



### Quick start
The simplest way to try logiqctl is to download a pre-built binary from our release page:
https://github.com/logiqai/logiqctl/releases

Once you have the binary run the following to get started
- run `logiqctl config set-cluster CLUSTER_URL`
- run `logiqctl config set-ui-credential flash-userid password`
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

