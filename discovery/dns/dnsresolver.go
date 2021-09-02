package dns

import (
	"fmt"
	"net"
)

func DnsResolver() {
	iprecords, err := net.LookupIP("bf-serv-ca")
	fmt.Println(iprecords, err)
	for _, ip := range iprecords {
		fmt.Println(ip)
	}
}
