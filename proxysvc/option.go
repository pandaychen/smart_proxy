package proxysvc

import (
	"fmt"
	"os"
	"time"

	"github.com/pandaychen/smart_proxy/enums"
)

// common options setter
type SmartProxyServiceOption func(*SmartProxyService) error

func WithProxysvcName(name string) SmartProxyServiceOption {
	return func(s *SmartProxyService) error {
		s.ProxyName = name
		return nil
	}
}

func WithProxysvcAddr(addr string) SmartProxyServiceOption {
	return func(s *SmartProxyService) error {
		s.ProxyAddress = addr
		return nil
	}
}

//IsSafeHttpSig
func WithHttpSigOn() SmartProxyServiceOption {
	return func(s *SmartProxyService) error {
		s.IsSafeHttpSig = true
		return nil
	}
}

func WithGinOn() SmartProxyServiceOption {
	return func(s *SmartProxyService) error {
		s.IsGinOn = true
		return nil
	}
}

func WithTlsOn(certfile, keyfile string) SmartProxyServiceOption {
	return func(s *SmartProxyService) error {
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

func WithTimeout(timeout time.Duration) SmartProxyServiceOption {
	return func(s *SmartProxyService) error {
		s.TimeOut = timeout
		return nil
	}
}

func WithDiscoveryType(dis enums.DISCOVERY_TYPE) SmartProxyServiceOption {
	return func(s *SmartProxyService) error {
		s.DiscoveryName = dis
		return nil
	}
}
