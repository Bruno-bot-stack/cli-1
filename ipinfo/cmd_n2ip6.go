package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

var completionsN2IP6 = &complete.Command{
	Flags: map[string]complete.Predictor{
		"--nocolor": predict.Nothing,
		"-h":        predict.Nothing,
		"--help":    predict.Nothing,
	},
}

func printHelpN2IP6() {
	fmt.Printf(
		`Usage: %s n2ip6 [<opts>] <expr>

Example:
  %[1]s n2ip6 "190.87.89.1"
  %[1]s n2ip6 "2001:0db8:85a3:0000:0000:8a2e:0370:7334
  %[1]s n2ip6 "2001:0db8:85a3::8a2e:0370:7334
  %[1]s n2ip6 "::7334
  %[1]s n2ip6 "7334::""
	

Options:
  General:
    --nocolor
      disable colored output.
    --help, -h
      show help.
`, progBase)
}

func cmdN2IP6() error {
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.BoolVar(&fNoColor, "nocolor", false, "disable colored output.")
	pflag.Parse()

	if fNoColor {
		color.NoColor = true
	}

	if fHelp {
		printHelpDefault()
		return nil
	}

	cmd := ""
	if len(os.Args) > 2 {
		cmd = os.Args[2]
	}

	// If no argument is given, print help.
	if strings.TrimSpace(cmd) == "" {
		printHelpN2IP6()
		return nil
	}

	res, err := lib.CmdN2IP6(cmd)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "err: %v\n", err)
		if err != nil {
			return err
		}
	}

	fmt.Println(res)
	return nil
}
