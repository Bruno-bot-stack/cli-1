package lib

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/pflag"
)

type CmdMatchIPFlags struct {
	Expression []string
	Help       bool
}

func (f *CmdMatchIPFlags) Init() {
	pflag.StringSliceVarP(
		&f.Expression,
		"expression", "e", nil,
		"IPs, subnets to be filtered.",
	)
	pflag.BoolVarP(
		&f.Help,
		"help", "", false,
		"show help.",
	)
}

func CmdMatchIP(
	f CmdMatchIPFlags,
	args []string,
	printHelp func(),
) error {
	if f.Help || len(f.Expression) == 0 || args[0] == "" {
		printHelp()
		return nil
	}

	stat, _ := os.Stdin.Stat()
	isStdin := (stat.Mode() & os.ModeCharDevice) == 0

	var sourceCIDRs, filterCIDRs []SubnetPair
	var sourceIPs, filterIPs []net.IP

	source_op := func(s string, inputType INPUT_TYPE) error {
		switch inputType {
		case INPUT_TYPE_IP:
			sourceIPs = append(sourceIPs, net.ParseIP(s))
		case INPUT_TYPE_CIDR:
			sourceCIDRs = append(sourceCIDRs, s)
		case INPUT_TYPE_IP_RANGE:
			sourceCIDRs = append(sourceCIDRs, s)
		default:
			return ErrInvalidInput
		}
		return nil
	}

	filter_op := func(s string, inputType INPUT_TYPE) error {
		switch inputType {
		case INPUT_TYPE_IP:
			filterIPs = append(filterIPs, net.ParseIP(s))
		case INPUT_TYPE_CIDR:
			filterCIDRs = append(filterCIDRs, s)
		case INPUT_TYPE_IP_RANGE:
			filterCIDRs = append(filterCIDRs, s)
		default:
			return ErrInvalidInput
		}
		return nil
	}

	var err error

	for _, expr := range f.Expression {
		if expr == "-" && isStdin {
			err = ProcessStringsFromStdin(source_op)
			if err != nil {
				return err
			}
		} else {
			err = ProcessStringsFromFile(expr, source_op)
			if err != nil {
				return err
			}
		}
	}

	for _, arg := range args {
		if arg == "-" && isStdin {
			err = ProcessStringsFromStdin(filter_op)
			if err != nil {
				return err
			}
		} else {
			err = ProcessStringsFromFile(arg, filter_op)
			if err != nil {
				return err
			}
		}
	}

	matches := findOverlapping(sourceCIDRs, filterCIDRs, sourceIPs, filterIPs)
	for _, v := range matches {
		fmt.Println(v)
	}

	return nil
}

type SubnetPair struct {
	Raw    string
	Parsed []net.IPNet
}

func findOverlapping(sourceCIDRs, filterCIDRs []SubnetPair, sourceIPs, filterIPs []net.IP) []string {
	var matches []string
	for _, sourceCIDR := range sourceCIDRs {
		foundMatch := false
		for _, v := range sourceCIDR.Parsed {
			for _, filterCIDR := range filterCIDRs {
				for _, fv := range filterCIDR.Parsed {
					if isCIDROverlapping(&v, &fv) {
						if !foundMatch {
							matches = append(matches, sourceCIDR.Raw)
							foundMatch = true
						}
						break
					}
				}
			}
		}
	}

	for _, sourceIP := range sourceIPs {
		foundMatch := false
		for _, filterCIDR := range filterCIDRs {
			for _, fv := range filterCIDR.Parsed {
				if isIPRangeOverlapping(&sourceIP, &fv) {
					if !foundMatch {
						matches = append(matches, sourceIP.String())
						foundMatch = true
					}
					break
				}
			}
		}

		for _, filterIP := range filterIPs {
			if sourceIP.String() == filterIP.String() {
				matches = append(matches, sourceIP.String())
				break
			}
		}
	}

	return matches
}

func isCIDROverlapping(cidr1, cidr2 *net.IPNet) bool {
	return cidr1.Contains(cidr2.IP) || cidr2.Contains(cidr1.IP)
}

func isIPRangeOverlapping(ip *net.IP, cidr *net.IPNet) bool {
	return cidr.Contains(*ip)
}
