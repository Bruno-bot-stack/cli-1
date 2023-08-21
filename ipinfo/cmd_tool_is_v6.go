package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolIs_v6 = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":      predict.Nothing,
		"--help":  predict.Nothing,
		"-q":      predict.Nothing,
		"--quiet": predict.Nothing,
	},
}

func printHelpToolIs_v6() {
	fmt.Printf(
		`Usage: %s tool is_v6 [<opts>] <cidr | ip | ip-range | filepath>

Description:
Checks if the input is an IPv6 address.
Input can be IPs, IP ranges, CIDRs, or filepath to a file

Examples:
# Check two CIDR.
$ %[1]s tool is_v6 1.1.1.0/30

# Check IP range.
$ %[1]s tool is_v6 1.1.1.0-1.1.1.244

# Check for file.
$ %[1]s tool is_v6 /path/to/file.txt 

# Check entries from stdin.
$ cat /path/to/file1.txt | %[1]s tool is_v6

Options:
  --help, -h
    show help.
  --quiet, -q
    quiet mode; suppress additional output.
`, progBase)
}

func cmdToolIs_v6() (err error) {
	f := lib.CmdToolIs_v6Flags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolIs_v6(f, pflag.Args()[2:], printHelpToolIs_v6)
}
