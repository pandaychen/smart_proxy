package errcode

import (
	"errors"
)

var (
	ErrNotSupportedMethod = errors.New("Not Supported HTTP method")

	ErrNotSupportedDiscoveryType = errors.New("Not support Discovery Method")

	ErrNoneProxyNodes = errors.New("No Proxy nodes")
)
