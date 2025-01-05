# Go cpuset

[![Go Reference](https://pkg.go.dev/badge/go.vallahaye.net/cpuset.svg)](https://pkg.go.dev/go.vallahaye.net/cpuset)
[![Go Report Card](https://goreportcard.com/badge/go.vallahaye.net/cpuset)](https://goreportcard.com/report/go.vallahaye.net/cpuset)

Go cpuset provides functions to parse strings according to the formats
specified in the [Linux cpuset(7) man page](https://man7.org/linux/man-pages/man7/cpuset.7.html#FORMATS).
The parsed representation can then be manipulated using well known set
functions (difference, intersection, ...) and formatted back to string.

### Goals

- Handle all formats specified in the man page ("List" and "Mask")
- Provide an API similar to Golang's upcoming standard set container implementation (golang/go#69230)
- Make it accessible from the command-line
- Well-tested and documented

### Non-goals

- Manage cpusets in the Linux kernel (similar to SUSE's
[cset](https://github.com/SUSE/cpuset) command-line tool)

## Install

Install the `cpuset` command-line tool using the Go toolchain:

```shell
go install go.vallahaye.net/cpuset/cmd/cpuset@latest
```

## Example

Parse two cpusets, compute the intersection and format it back to string.

### Golang

```go
// Decode str1 into a CPUSet.
s1, err := cpuset.ParseList(str1)
if err != nil {
  // Do something with the error.
}

// Decode str2 into a CPUSet.
s2, err := cpuset.ParseList(str2)
if err != nil {
  // Do something with the error.
}

// Compute the intersection of the two cpusets.
s := cpuset.Intersection(s1, s2)

// Encode the intersection and print the result.
fmt.Println(s.String())
```

### Command-line

```shell
cpuset intersection "$STR1" "$STR2"
```
