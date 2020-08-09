---
description: >-
  The LOGIQ command line toolkit, logiqctl, allows you to run commands against
  LOGIQ clusters. You can tail logs from your applications and servers, query
  historical data and search your log data.
---

# Overview

## Installing logiqctl - 1.0.0/1.0.1

Download Logiqctl binary based on your workstation OS and architecture from [here](https://github.com/logiqai/logiqctl/releases).

## Configuring logiqctl

Before using logiqctl can be used , it is important to first configure the logiq cluster endpoint.

```text
$ ./logiqctl c
Enter the Logiq cluster URL: logiq.my-domain.com
Created following profile at /Users/johndoe/.logiqctl/config.toml
====================================
Version = ""

[[Profile]]
  Name = "logiq.my-domain.com"
  ClusterURL = "logiq.my-domain.com:8081"
  Default = true

====================================
```

## Help

```text
$ ./logiqctl h
NAME:
   Logiqctl - LOGIQ command line toolkit

USAGE:
   logiqctl [global options] command [command options] [arguments...]

VERSION:
   1.0.0-rc

AUTHOR:
   logiq.ai <cli@logiq.ai>

COMMANDS:
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

