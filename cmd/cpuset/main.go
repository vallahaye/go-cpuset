package main

import (
	"flag"
	"fmt"
	"os"

	"go.vallahaye.net/cpuset"
)

const (
	listFormat    = "list"
	maskFormat    = "mask"
	defaultFormat = listFormat
)

const usageHeader = `Usage: cpuset [flags] command s1 s2

Flags:`

const usageFooter = `
Commands:
  difference
    	compute the difference of the two cpusets
  intersection
    	compute the intersection of the two cpusets
  union
    	compute the union of the two cpusets

Examples:
  cpuset difference 0-32 8-16
  cpuset -format mask difference 00000001,ffffffff 0000ff00

See also:
  man cpuset(7) for more information about cpusets`

func usage() {
	fmt.Fprintln(os.Stderr, usageHeader)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, usageFooter)
}

func fail(text string) {
	fmt.Fprintln(os.Stderr, text)
	usage()
	os.Exit(2)
}

func main() {
	var (
		format       string
		printVersion bool
	)

	flag.Usage = usage
	flag.StringVar(&format, "format", defaultFormat, "use the specified format for parsing the two cpusets and outputing the result")
	flag.BoolVar(&printVersion, "version", false, "print the version and exit")
	flag.Parse()

	if printVersion {
		fmt.Println("cpuset " + cpuset.Version)
		os.Exit(0)
	}

	var (
		parseFn  func(string) (cpuset.CPUSet, error)
		stringFn func(*cpuset.CPUSet) string
	)

	switch format {
	case listFormat:
		parseFn, stringFn = cpuset.ParseList, (*cpuset.CPUSet).ListString
	case maskFormat:
		parseFn, stringFn = cpuset.ParseMask, (*cpuset.CPUSet).MaskString
	default:
		fail("flag provided but invalid: -format")
	}

	args := flag.Args()
	if len(args) != 3 {
		fail("invalid number of arguments")
	}

	var commandFn func(cpuset.CPUSet, cpuset.CPUSet) cpuset.CPUSet

	switch args[0] {
	case "difference":
		commandFn = cpuset.Difference
	case "intersection":
		commandFn = cpuset.Intersection
	case "union":
		commandFn = cpuset.Union
	default:
		fail("command provided but not defined: " + args[0])
	}

	s1, err := parseFn(args[1])
	if err != nil {
		fail("s1 provided but invalid: " + err.Error())
	}

	s2, err := parseFn(args[2])
	if err != nil {
		fail("s2 provided but invalid: " + err.Error())
	}

	s := commandFn(s1, s2)
	fmt.Println(stringFn(&s))
}
