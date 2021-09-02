package http

import (
	"net"
	"net/http"
	"strings"
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
