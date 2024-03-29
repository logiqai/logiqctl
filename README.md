# logiqctl

`logiqctl` is LOGIQ's inbuilt command-line toolkit that lets you interact with the LOGIQ Observability platform without logging into the UI. Using `logiqctl`, you can:

- Stream logs in real-time
- Query historical application logs
- Search within logs across namespaces
- Query and view events across your LOGIQ stack
- View and create event rules
- Create and manage dashboards
- Query and view all your resources on LOGIQ such as applications, dashboards, namespaces, processes, and queries
- Manage LOGIQ licenses

# Quickstart

The quickest way to start using `logiqctl` is to download a pre-built binary from our [release page on GitHub](https://github.com/logiqai/logiqctl/releases). 

## Configuring `logiqctl`

Once you've downloaded the binary, you can configure `logiqctl` to interact with your LOGIQ instance by doing the following:
1. Set your cluster URL:
    ```
    logiqctl config set-cluster <CLUSTER_URL>
    ```
1. Set the API token:
    ```
    logiqctl config set-token <LOGIQ_API_KEY>
    ```
    **Note:** If you don't have a LOGIQ API key, read [Obtaining API Key](https://docs.logiq.ai/vewing-logs/logiqctl/obtaining-api-key) to learn how to obtain one. 
1. Set your LOGIQ credentials:
    ```
    logiqctl config set-ui-credential flash-userid password
    ```
1. Set your default namespace:
    ```
    logiqctl config set-context NAMESPACE
    ```
1. Verify your `logiqctl` configuration:
    ```
    logiqctl get namespaces
    ```
This completes the installation of `logiqctl`. You can now use `logiqctl` to interact with your LOGIQ instance right from your terminal.


# Pattern-signature generation
`Logiqctl` is equipped with log Pattern-Signature (PS) generation and post PS statistics analysis. All the logs dumped by `logiqctl` client can be automatically calcaulated common text patterns using the flag (-g).  This feature supports log dumping functions 'logiqctl logs', 'logiqctl logs search', and 'logiqctl tail'.  

PS generation is processed in binary [psmod](https://github.com/logiqai/logiqctl/releases/tag/2.1.2) executable.  
- running with ps gen requires psmod be at the same location as logiqctl.
- From the downloaded releases zip file, copy both the psmod and logiqctl binaries for your architecture/os before running e.g. if your architecture is darwin_amd64, copy logiqctl_darwin_amd64 and psmod_darwin_amd64 to a folder. Rename psmod_darwin_amd64 to psmod before running logiqctl
- Once pattern signatures are generated, see the signatures extracted in the ps_stat.out file.

# Building `logiqctl` from source

Another way of installing `logiqctl` is by building it from the source code. Building `logiqctl` from its source code involves two steps:
- Installing dependencies
- Downloading and building the `logiqctl` binary

## Installing dependencies

`logiqctl` has the following dependencies:
- Go: You can install Go by following the instructions listed on [https://golang.org/dl/]
- Protocol Buffers: Download the binary and set it up by running the following commands:

On macOS:

```bash
PROTOC_ZIP=protoc-3.15.6-osx-x86_64.zip
curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.15.6/$PROTOC_ZIP
sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
rm -f $PROTOC_ZIP
```

On Linux OS:
   
```bash
PROTOC_ZIP=protoc-3.15.6-linux-x86_64.zip
curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.15.6/$PROTOC_ZIP
sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
rm -f $PROTOC_ZIP

```

This completes the installation of all `logiqctl` dependencies. 

## Building `logiqctl`

Run the following commands to build `logiqctl` from the source code:
1. Create a directory inside your workspace in which to keep source code:
    ```
    mkdir -p $GOPATH/src/github.com/logiqai
    ```
1. Accesss the source code directory:
    ```
    cd $GOPATH/src/github.com/logiqai
    ```
1. Clone the `logiqctl` GitHub repository into this folder:
    ```
    git clone git@github.com:logiqai/logiqctl.git
    ```
1. Access the repository you just cloned:
    ```
    cd logiqctl
    ```
1. Build `logiqctl`:
    ```
    go build logiqctl.go
    ```
1. Make the binary `logiqctl` executable:
    ```
    chmod +x ./logiqctl
    ```
1. Verify the build:
    ```
    logiqctl -h
    ```

`logiqctl` is now built and ready for configuration and use. To configure `logiqctl`, refer to the configuration instructions listed under [Configuring `logiqctl`](#configuring-logiqctl). 

# Available `logiqctl` commands

| Command | Operation |
|---|---|
| [`logiqctl config`](docs/logiqctl_config.md) | Configure `logiqctl` or modify existing `logiqctl` configuration |
| [`logiqctl tail`](docs/logiqctl_tail.md) | Stream logs from your LOGIQ instance in real-time |
| [`logiqctl create`](docs/logiqctl_create.md) | Create LOGIQ resources such as dashboards and event rules |
| [`logiqctl get`](docs/logiqctl_get.md) | Display one or more LOGIQ resources |
| [`logiqctl license`](docs/logiqctl_license.md) | View and manage your LOGIQ license |
| [`logiqctl logs`](docs/logiqctl_logs.md) | View logs for the given namespace and application |

# Options

```
  -c, --cluster string       Override the default cluster set by `logiqctl set-cluster' command
  -h, --help                 help for logiqctl
  -n, --namespace string     Override the default context set by `logiqctl set-context' command
  -o, --output string        Output format. One of: table|json|yaml. 
                             json output is not indented, use '| jq' for advanced json operations (default "table")
  -t, --time-format string   Time formatting options. One of: relative|epoch|RFC3339. 
                             This is only applicable when the output format is table. json and yaml outputs will have time in epoch seconds. (default "relative")
```

### Release Note:
```
Thu Jul  8 15:04:59 PDT 2021 (2.1.0)
  - Enhance search operation with time-ranges
  - PS enhancement using addon binary module psmod
Mon Jul  5 21:24:25 PDT 2021
  - Enhance with log pattern-signature (PS) generation support
Wed Jul 14 17:13:48 PDT 2021 (2.1.1)
  - Multiple application searches
Thu Jul 15 09:02:23 PDT 2021 (2.1.2)
  - Inconsistent multi-apps display fixes
Tue Aug  3 07:54:39 PDT 2021
  - Enhance search capability
Wed Sep 15 14:10:46 PDT 2021
  - Enable advanced search feature
Fri Nov  5 16:05:50 PDT 2021
  - Fixed dashboard text import/export, add http network trace
```




To know more about the LOGIQ Observability stack, see https://logiq.ai/ and https://docs.logiq.ai/. 

In case of issues or questions, do reach out to us at [cli@logiq.ai]. You can also [log an issue](https://github.com/logiqai/logiqctl/issues/new) in our `logiqctl` source code repository on GitHub. 
