package main

import (
	"fmt"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func sniffFTP() error {
	h, err := pcap.OpenLive(ifaceName, 65536, true, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("pcap open: %w", err)
	}
	defer h.Close()
	if err := h.SetBPFFilter("tcp port 21"); err != nil {
		return fmt.Errorf("bpf: %w", err)
	}
	src := gopacket.NewPacketSource(h, h.LinkType())
	pkts := src.Packets()
	for {
		select {
		case <-shutdownSignal:
			return nil
		case pkt, ok := <-pkts:
			if !ok {
				return nil
			}
			handleFTPPacket(pkt)
		}
	}
}

func handleFTPPacket(pkt gopacket.Packet) {
	tcpL := pkt.Layer(layers.LayerTypeTCP)
	if tcpL == nil {
		return
	}
	tcp := tcpL.(*layers.TCP)
	if len(tcp.Payload) == 0 {
		return
	}
	clientToServer := tcp.DstPort == 21
	for _, line := range strings.Split(strings.TrimRight(string(tcp.Payload), "\r\n"), "\r\n") {
		if line == "" {
			continue
		}
		if verbose {
			tag := "S->C"
			if clientToServer {
				tag = "C->S"
			}
			fmt.Printf("[FTP %s] %s\n", tag, line)
			continue
		}
		if !clientToServer {
			continue
		}
		upper := strings.ToUpper(line)
		switch {
		case strings.HasPrefix(upper, "STOR "),
			strings.HasPrefix(upper, "RETR "),
			strings.HasPrefix(upper, "APPE "),
			strings.HasPrefix(upper, "DELE "):
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				fmt.Printf("[FTP] %s -> %s\n", strings.ToUpper(parts[0]), parts[1])
			}
		}
	}
}
