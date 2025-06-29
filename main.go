package main

import (
	"flag"
	"net"
	"os"
	"runtime"
	"os/exec"
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

	if nslookup(*host) {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

func exec_ping(host string) bool {
	var cmd *exec.Cmd
	if runtime.GOOS != "windows" {
		cmd = exec.Command("ping", "-c", "1", host)
		} else {
		cmd = exec.Command("ping", "-n", "1", host)
	}

	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
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
