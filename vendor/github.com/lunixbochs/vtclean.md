# README

[![Build Status](https://travis-ci.org/lunixbochs/vtclean.svg?branch=master)](https://travis-ci.org/lunixbochs/vtclean)

## vtclean

Clean up raw terminal output by stripping escape sequences, optionally preserving color.

Get it: `go get github.com/lunixbochs/vtclean/vtclean`

API:

```text
import "github.com/lunixbochs/vtclean"
vtclean.Clean(line string, color bool) string
```

Command line example:

```text
$ echo -e '\x1b[1;32mcolor example
color forced to stop at end of line
backspace is ba\b\bgood
no beeps!\x07\x07' | ./vtclean -color

color example
color forced to stop at end of line
backspace is good
no beeps!
```

Go example:

```text
package main

import (
    "fmt"
    "github.com/lunixbochs/vtclean"
)

func main() {
    line := vtclean.Clean(
        "\033[1;32mcolor, " +
        "curs\033[Aor, " +
        "backspace\b\b\b\b\b\b\b\b\b\b\b\033[K", false)
    fmt.Println(line)
}
```

Output:

```text
color, cursor
```

