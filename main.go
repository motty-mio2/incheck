package main

import (
	"flag"
	"net"
	"os"
	"strings"
)

func main() {
	var (
		host       = flag.String("h", "", "Host to check")
		ipAddrList = []string{}
	)

	var ipaddr string
	flag.Func("i", "IP address to check", func(s string) error {
		ipaddr = s
		ipAddrList = append(ipAddrList, ipaddr)
		return nil
	})

	flag.Parse()

	success := false
	for _, ip := range ipAddrList {
		ip = strings.TrimSpace(ip)
		success = success || check_ipaddr(ip)
	}

	if success && nslookup(*host) {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

// Nslookup
func nslookup(host string) bool {
	_, err := net.LookupHost(host)
	if err != nil {
		return false
	}
	return true
}

// Ipaddress match
func check_ipaddr(ipaddress string) bool {
	ip := net.ParseIP(ipaddress)
	if ip == nil {
		return false
	}

	mask := net.CIDRMask(24, 32)
	networkAddr := ip.Mask(mask)

	ift, err := net.Interfaces()
	if err != nil {
		return false
	}

	for _, ifi := range ift {
		addrs, err := ifi.Addrs()
		if err != nil {
			return false
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			ip := ipNet.IP
			if ip.To4() == nil {
				continue
			}

			localNetworkAddr := ip.Mask(mask)

			if networkAddr.Equal(localNetworkAddr) {
				return true
			}
		}
	}

	return false
}
