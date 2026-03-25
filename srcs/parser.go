package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	srcIP  string
	srcMAC string
	dstIP  string
	dstMAC string
)

var parser = &cobra.Command{
	Use:   "inquisitor",
	Short: "FTP Traffic Interceptor",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Source IP: %s\n", srcIP)
		fmt.Printf("Source MAC: %s\n", srcMAC)
		fmt.Printf("Destination IP: %s\n", dstIP)
		fmt.Printf("Destination MAC: %s\n", dstMAC)
	},
}

func init() {
	parser.Flags().StringVarP(&srcIP, "src-ip", "s", "", "Source IP address")
	parser.Flags().StringVarP(&srcMAC, "src-mac", "m", "", "Source MAC address")
	parser.Flags().StringVarP(&dstIP, "dst-ip", "d", "", "Destination IP address")
	parser.Flags().StringVarP(&dstMAC, "dst-mac", "M", "", "Destination MAC address")
}

func parseArgs() error {
	if err := parser.Execute(); err != nil {
		return err
	}
	return nil
}
