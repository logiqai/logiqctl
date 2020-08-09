---
description: This page documents command-line usage for the Query command
---

# Query

## Help

```text
NAME:
   Logiqctl query - query "sudo cron" 2h

USAGE:
   Logiqctl query command [command options] [application names, relative time]

COMMANDS:
   next, n  query n
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --end_time value, --et value    Relative end time. (default: "10h")
   --filter value, -f value        Filter expression e.g. 'Hostname=127.0.0.1,10.231.253.255;Message=tito*'
   --output value, -o value        Set output format to be column|json|raw (default: "column")
   --start_time value, --st value  Relative start time. (default: "10h")
   --tail, -t                      Tail the data without paginating (default: false)
   --help, -h                      show help (default: false)
   --version, -v                   print the version (default: false
```

## Historical queries

Historical log data can be queries using logiqctl. Filters allow the user to specify a time range, filter criteria like application name, matching patterns etc.

For e.g. to query all wordpress logs within the last 2 hr, one would pass the application names and the start time for the query

```text
$logiqctl q --st=2h wordpress 
2020-04-21T13:46:56Z        |debug |fluentd-2k2j8|wordpress-app-5cd484b46c-krkhq|wordpress|user|10.52.1.1 - - [21/Apr/2020:13:46:53 +0000] "GET /wp-login.php HTTP/1.1" 200 1832
2020-04-21T13:46:56Z        |debug |fluentd-2k2j8|wordpress-app-5cd484b46c-krkhq|wordpress|user|10.52.1.1 - - [21/Apr/2020:13:46:51 +0000] "GET /wp-login.php HTTP/1.1" 200 1832
2020-04-21T13:46:47Z        |debug |fluentd-2k2j8|wordpress-app-5cd484b46c-krkhq|wordpress|user|10.52.1.1 - - [21/Apr/2020:13:46:43 +0000] "GET /wp-login.php HTTP/1.1" 200 1832
2020-04-21T13:46:42Z        |debug |fluentd-2k2j8|wordpress-app-5cd484b46c-krkhq|wordpress|user|10.52.1.1 - - [21/Apr/2020:13:46:41 +0000] "GET /wp-login.php HTTP/1.1" 200 1832
2020-04-21T13:46:37Z        |debug |fluentd-2k2j8|wordpress-app-5cd484b46c-krkhq|wordpress|user|10.52.1.1 - - [21/Apr/2020:13:46:33 +0000] "GET /wp-login.php HTTP/1.1" 200 1832
2020-04-21T13:46:16Z        |debug |fluentd-2k2j8|wordpress-app-5cd484b46c-krkhq|wordpress|user|10.52.1.1 - - [21/Apr/2020:13:46:13 +0000] "GET /wp-login.php HTTP/1.1" 200 1832

Atleast 1423 records remaining... Enter `n` or `next' to continue.
```

## Output formatting

The query output formatting can be controlled with the `--output` option. Three values are allowed - `raw`, `column` and `json`

```text
$logiqctl q postgres --st 1h --output json
```

The above command returns the data with each row of data formatted as JSON.

