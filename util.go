package main

import (
	"fmt"
	"net"

	"github.com/cmatsuoka/ouidb"
	"github.com/mostlygeek/arp"
)

var db *ouidb.OuiDB

func init() {
	db = ouidb.New("/etc/manuf")
	if db == nil {
		db = ouidb.New("manuf")
	}
}

func getMAC(s string) (string, error) {
	ifaces, err := net.Interfaces()
	checkError(err)
	for _, i := range ifaces {
		if i.Name == s {
			return i.HardwareAddr.String(), nil
		}
	}
	return "", fmt.Errorf("%s: no such interface", s)
}

func getName(addr string) string {
	names, err := net.LookupAddr(addr)
	if err != nil {
		return ""
	}
	return names[0]
}

func getMACFromIP(addr string) string {
	arp.CacheUpdate()
	mac := arp.Search(addr)
	if mac != "" {
		return mac
	}

	ip := net.ParseIP(addr)
	if ip == nil {
		return mac
	}

	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{ip, 0, ""})
	if err != nil {
		return mac
	}
	conn.Write([]byte{0})
	conn.Close()

	arp.CacheUpdate()
	return arp.Search(addr)
}

func getVendor(mac string) string {
	v, _ := db.Lookup(mac)
	return v

}
