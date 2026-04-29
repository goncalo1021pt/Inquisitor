package main

import (
	"fmt"
	"net"
	"os"
)

var (
	srcIP   net.IP
	srcMAC  net.HardwareAddr
	dstIP   net.IP
	dstMAC  net.HardwareAddr
	verbose bool
)

func usage() {
	fmt.Fprintf(os.Stderr,
		"Usage: inquisitor [-v] <IP-src> <MAC-src> <IP-target> <MAC-target>\n"+
			"  -v   verbose: print all FTP control traffic (login included)\n")
}

func parseArgs() error {
	var positional []string
	for _, a := range os.Args[1:] {
		switch a {
		case "-v", "--verbose":
			verbose = true
		case "-h", "--help":
			usage()
			os.Exit(0)
		default:
			positional = append(positional, a)
		}
	}
	if len(positional) != 4 {
		usage()
		return fmt.Errorf("expected 4 positional arguments, got %d", len(positional))
	}

	var err error
	if srcIP, err = parseIPv4(positional[0]); err != nil {
		return fmt.Errorf("invalid IP-src: %w", err)
	}
	if srcMAC, err = parseMAC(positional[1]); err != nil {
		return fmt.Errorf("invalid MAC-src: %w", err)
	}
	if dstIP, err = parseIPv4(positional[2]); err != nil {
		return fmt.Errorf("invalid IP-target: %w", err)
	}
	if dstMAC, err = parseMAC(positional[3]); err != nil {
		return fmt.Errorf("invalid MAC-target: %w", err)
	}
	if srcIP.Equal(dstIP) {
		return fmt.Errorf("IP-src and IP-target must differ")
	}
	return nil
}

func parseIPv4(s string) (net.IP, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("%q is not a valid IP", s)
	}
	v4 := ip.To4()
	if v4 == nil {
		return nil, fmt.Errorf("%q is not IPv4", s)
	}
	return v4, nil
}

func parseMAC(s string) (net.HardwareAddr, error) {
	m, err := net.ParseMAC(s)
	if err != nil {
		return nil, err
	}
	if len(m) != 6 {
		return nil, fmt.Errorf("%q is not a 48-bit MAC", s)
	}
	return m, nil
}
