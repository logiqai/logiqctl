# Logiqbox 
### CLI for Logiq Platform

- Tail logs in realtime
- Query historic data
- Do Text Search on data 


#### How to run

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
- run `git clone git@github.com:logiqai/logiqbox.git`
- run `cd logiqbox`
- run `./generate_grpc.sh `
- run `go build logiqbox.go`

```bash

> ./logiqbox 
               
NAME:
   Logiq-box - Logiq CLI Tool

USAGE:
   logiqbox [global options] command [command options] [arguments...]

VERSION:
   1.0.0

AUTHOR:
   logiq.ai <cli@logiq.ai>

COMMANDS:
   configure, c  Configure Logiq-box
   list, ls      List of applications that you can tail
   tail, t       tail logs filtered by namespace, application, labels or process / pod name
   next, n       query n
   query, q      query "sudo cron" 2h
   search, s     search sudo
   help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --st value     Relative start time. (default: "10h")
   --et value     Relative end time. (default: "10h")
   --debug value  --debug true (default: "false")
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)


```

To know more about Logiq Platform, see https://logiq.ai/ and https://docs.logiq.ai/ 