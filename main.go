package main

import (
	"flag"
	"net"
	"os"
	"runtime"

	probing "github.com/prometheus-community/pro-bing"
)

func main() {
	var (
		ipaddr = flag.String("i", "", "IP address to check")
		host   = flag.String("h", "", "Host to check")
	)

	flag.Parse()

	if ping_to_host(*ipaddr) && nslookup(*host) {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

// Ping the host
func ping_to_host(host string) bool {
	pinger, err := probing.NewPinger(host)
	pinger.SetPrivileged(true)
	if err != nil {
		return false
	}
	if runtime.GOOS != "windows" {
		pinger.SetDoNotFragment(true)
	}
	pinger.Count = 1
	pinger.Timeout = 3000000000
	err = pinger.Run()
	if err != nil {
		return false
	}

	stats := pinger.Statistics()

	return stats.PacketLoss == 0
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
