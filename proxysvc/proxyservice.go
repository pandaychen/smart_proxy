package proxysvc

import (
	"net/http/httputil"
	"sync"
)

// ProxyService for loadbalance to backend pools
type SmartProxyService struct {
	sync.RWMutex
	ProxyName    string
	ProxyAddress string

	//mapping requests to backends，请求会根据某种策略被代理到后端
	ReverseProxyMap map[string]*httputil.ReverseProxy
}
