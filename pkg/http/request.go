package http

import (
	"net"
	"net/http"
	"strings"
	"time"
)

// GetClientIP acquire the client IP address
func GetClientIP(r *http.Request) (clientIP string) {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ", ")
		clientIP = ips[len(ips)-1]
	} else {
		clientIP, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	return clientIP
}

func CheckTcpAlive(backend_addr string) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", backend_addr, timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
