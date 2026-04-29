package main

import (
	"fmt"
	"os"
)

func main() {
	if err := parseArgs(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Inquisitor - ARP Poisoning + FTP Interceptor")
	if err := exec(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
