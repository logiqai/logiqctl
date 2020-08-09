---
description: This page documents command-line usage for the List command
---

# List

## Help

```bash
#./logiqctl ls -h
NAME:
   logiqctl list - List of applications that you can tail

USAGE:
   logiqctl list [command options] [arguments...]

OPTIONS:
   --namespaces  -namespaces (default: false)
   --help, -h    show help (default: false)
```

## List All applications

_**logiqctl**_ cli provides the ls command to show the list of applications from which logs have been ingested. In the below example we see 3 different applications have sent log data from two namespaces

```bash
#./logiqctl ls
Namespace | Application                   | ProcId                                                 | Last Seen      | First Seen
logiq     | kube-state-metrics            | prometheus-kube-state-metrics-78b7d687c5-9k95v         | 6 minutes ago  | 1 day ago
logiq     | beetle-coordinator            | beetle-coordinator-b649dfd86-rs8zt                     | 20 minutes ago | 3 days ago
logiq     | logiq-flash                   | logiq-flash-2                                          | 27 minutes ago | 3 days ago
logiq     | coffee-worker                 | coffee-worker-678d4f4bfd-jpzqg                         | 3 hours ago    | 1 day ago
logiq     | flash-discovery               | flash-discovery-0                                      | 4 hours ago    | 2 days ago
logiq     | prometheus                    | prometheus-prometheus-prometheus-0                     | 5 hours ago    | 1 day ago
logiq     | logiq-flash                   | logiq-flash-0                                          | 7 hours ago    | 3 days ago
logiq     | kubernetes-ingress-controller | logiq-kubernetes-ingress-596dc9bbf6-95m6h              | 1 day ago      | 1 day ago
logiq     | pithos                        | pithos-554d576564-vd7sh                                | 1 day ago      | 3 days ago
logiq     | redis                         | redis-master-0                                         | 1 day ago      | 2 days ago
```

## Filter by Namespace

Using thee --namespaces filter, one can restrict the application listing for one or more namespaces only

```bash
#./logiqctl ls --namespaces kube-logging
Namespace    | Application | ProcId        | Last Seen      | First Seen
kube-logging | fluentd     | fluentd-s9csl | 2 minutes ago  | 2 days ago
kube-logging | fluentd     | fluentd-7dgxz | 3 minutes ago  | 2 days ago
kube-logging | fluentd     | fluentd-p89rq | 14 minutes ago | 2 days ago
```

