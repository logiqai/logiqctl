# Search

## Help

```bash
#logiqctl s -h
NAME:
   Logiqctl search - search sudo

USAGE:
   Logiqctl search command [command options] [search_term, relative time]

COMMANDS:
   next, n  query n
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --st value      Relative start time. (default: "10h")
   --et value      Relative end time. (default: "10h")
   --filter value  --filter 'Hostname=127.0.0.1,10.231.253.255;Message=tito*'
   --help, -h      show help (default: false)
   --version, -v   print the version (default: false)
```

## Providing a search term

```bash
#logiqctl s superuser
2020-03-17T03:52:36Z        |error |docker-desktop|1518|postgres|system daemon|       superuser. For example, "-e POSTGRES_PASSWORD=password" on "docker run".
2020-03-17T03:52:36Z        |error |docker-desktop|1518|postgres|system daemon|Error: Database is uninitialized and superuser password is not specified.
Enter `n` or `next' to continue.
```

