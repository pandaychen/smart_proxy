package reverseproxy

import (
	"fmt"
	"os"
	"time"

	"smart_proxy/enums"
)

// common options setter
type SmartReverseProxyOption func(*SmartReverseProxy) error

func WithProxysvcName(name string) SmartReverseProxyOption {
	return func(s *SmartReverseProxy) error {
		s.ProxyName = name
		return nil
	}
}

func WithProxysvcAddr(addr string) SmartReverseProxyOption {
	return func(s *SmartReverseProxy) error {
		s.ProxyAddress = addr
		return nil
	}
}

//IsSafeHttpSig
func WithHttpSigOn() SmartReverseProxyOption {
	return func(s *SmartReverseProxy) error {
		s.IsSafeHttpSig = true
		return nil
	}
}

func WithGinOn() SmartReverseProxyOption {
	return func(s *SmartReverseProxy) error {
		s.IsGinOn = true
		return nil
	}
}

func WithTlsOn(certfile, keyfile string) SmartReverseProxyOption {
	return func(s *SmartReverseProxy) error {
		if _, err := os.Stat(certfile); err != nil {
			return fmt.Errorf("certfile path %s error", certfile)
		}
		if _, err := os.Stat(keyfile); err != nil {
			return fmt.Errorf("keyfile path %s error", keyfile)
		}

		s.IsTlsOn = true
		s.TlsConfig.CertFile = certfile
		s.TlsConfig.KeyFile = keyfile
		return nil
	}
}

func WithTimeout(timeout time.Duration) SmartReverseProxyOption {
	return func(s *SmartProxyReverse) error {
		s.TimeOut = timeout
		return nil
	}
}

func WithDiscoveryType(dis enums.DISCOVERY_TYPE) SmartReverseProxyOption {
	return func(s *SmartProxyReverse) error {
		s.DiscoveryName = dis
		return nil
	}
}
