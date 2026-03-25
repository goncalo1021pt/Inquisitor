package main

import (
	"fmt"
)

func main() {
	fmt.Println("Inquisitor - FTP Traffic Interceptor")
	if err := parseArgs(); err != nil {
		fmt.Printf("Error parsing arguments: %v\n", err)
	}
	exec()
}
