# sanitized\_anchor\_name

[![Build Status](https://travis-ci.org/shurcooL/sanitized_anchor_name.svg?branch=master)](https://travis-ci.org/shurcooL/sanitized_anchor_name) [![GoDoc](https://godoc.org/github.com/shurcooL/sanitized_anchor_name?status.svg)](https://godoc.org/github.com/shurcooL/sanitized_anchor_name)

Package sanitized\_anchor\_name provides a func to create sanitized anchor names.

Its logic can be reused by multiple packages to create interoperable anchor names and links to those anchors.

At this time, it does not try to ensure that generated anchor names are unique, that responsibility falls on the caller.

## Installation

```bash
go get -u github.com/shurcooL/sanitized_anchor_name
```

## Example

```go
anchorName := sanitized_anchor_name.Create("This is a header")

fmt.Println(anchorName)

// Output:
// this-is-a-header
```

## License

* [MIT License](https://github.com/logiqai/logiqctl/tree/21d7cfb18f20daa606c2eb00d6236892e77feb04/vendor/github.com/shurcooL/sanitized_anchor_name/LICENSE/README.md)

