package proxysvc

import (
	"fmt"
	"os"
	"time"

	"github.com/pandaychen/smart_proxy/enums"
)

// common options setter
type SmartProxyReverseOption func(*SmartProxyReverse) error

func WithProxysvcName(name string) SmartProxyReverseOption {
	return func(s *SmartProxyReverse) error {
		s.ProxyName = name
		return nil
	}
}

func WithProxysvcAddr(addr string) SmartProxyReverseOption {
	return func(s *SmartProxyReverse) error {
		s.ProxyAddress = addr
		return nil
	}
}

//IsSafeHttpSig
func WithHttpSigOn() SmartProxyReverseOption {
	return func(s *SmartProxyReverse) error {
		s.IsSafeHttpSig = true
		return nil
	}
}

func WithGinOn() SmartProxyReverseOption {
	return func(s *SmartProxyReverse) error {
		s.IsGinOn = true
		return nil
	}
}

func WithTlsOn(certfile, keyfile string) SmartProxyReverseOption {
	return func(s *SmartProxyReverse) error {
		if _, err := os.Stat(certfile); err != nil {
			return fmt.Errorf("certfile %s error", certfile)
		}
		if _, err := os.Stat(keyfile); err != nil {
			return fmt.Errorf("keyfile %s error", keyfile)
		}

		s.IsTlsOn = true
		s.TlsConfig.CertFile = certfile
		s.TlsConfig.KeyFile = keyfile
		return nil
	}
}

func WithTimeout(timeout time.Duration) SmartProxyReverseOption {
	return func(s *SmartProxyReverse) error {
		s.TimeOut = timeout
		return nil
	}
}

func WithDiscoveryType(dis enums.DISCOVERY_TYPE) SmartProxyReverseOption {
	return func(s *SmartProxyReverse) error {
		s.DiscoveryName = dis
		return nil
	}
}
