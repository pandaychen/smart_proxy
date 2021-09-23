package reverseproxy

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"smart_proxy/config"
	"smart_proxy/enums"
	"smart_proxy/loadbalancer"
	"smart_proxy/loadbalancer/wrr"
	"smart_proxy/pkg/comms"
	sphttp "smart_proxy/pkg/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// ProxyService for loadbalance to backend pools
type SmartReverseProxy struct {
	sync.RWMutex
	comms.TlsConfig
	IsTlsOn      bool
	IsGinOn      bool // 使用 gin 做反向代理
	ProxyName    string
	ProxyAddress string //ip:port
	TimeOut      time.Duration

	LoadBalancerName enums.LB_TYPE
	DiscoveryName    enums.DISCOVERY_TYPE

	IsSafeHttpSigOn bool // 是否加入 http 转发签名

	// 重要！
	//BackendNodePool 指向后端对应的在线 server 列表
	//BackendNodePool 是一个 interface{}，由具体的 lb 算法来完成实例化
	// 真正的后端连接池以 BackendNodePool 对象实例化
	BackendNodePool loadbalancer.BackendNodePool

	//mapping requests to backends，请求会根据某种策略被代理到不同的后端，key 为后端地址
	ReverseProxyMap     map[string]*httputil.ReverseProxy
	ReverseProxyMapLock *sync.RWMutex

	ProxyServer    *http.Server //self http hanlder
	GinProxyServer *gin.Engine

	Logger *zap.Logger

	//metrcis
	httpRequestCount    *prometheus.CounterVec
	httpRequestDuration *prometheus.SummaryVec
}

// init single reverse proxy
func NewSmartReverseProxy(logger *zap.Logger, conf *config.ReverseProxyConfig) (*SmartReverseProxy, error) {
	var (
		lb_name enums.LB_TYPE
		//discovery_name enums.DISCOVERY_TYPE
		backendNodemap = make(map[string]int)
		//address        string
	)

	switch conf.LbType {
	case "weight-rr":
		lb_name = enums.LB_WEIGHT_RR
	case "consistent-hash":
		lb_name = enums.LB_CONSISTENT_HASH
	case "p2c":
		lb_name = enums.LB_P2C
	default:
		lb_name = enums.LB_WEIGHT_RR
	}

	spr := &SmartReverseProxy{
		ProxyName:        conf.ProxyName,
		ProxyAddress:     conf.BindAddr,
		IsSafeHttpSigOn:  conf.SingnatureOn,
		LoadBalancerName: lb_name,
		TimeOut:          20 * time.Second,
		//	DiscoveryName    :discover_name,
		ReverseProxyMap:     make(map[string]*httputil.ReverseProxy), //key与lb算法有关
		ReverseProxyMapLock: new(sync.RWMutex),
		Logger:              logger,
	}

	spr.httpRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "smartproxy_http_request_count",
			Help: "http request count",
		},
		[]string{"method", "host", "path", "status"},
	)

	spr.httpRequestDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "smartproxy_http_request_duration",
			Help: "http request duration",
		},
		[]string{"method", "host", "path"},
	)

	switch spr.LoadBalancerName {
	case enums.LB_WEIGHT_RR:
		for _, v := range conf.PoolConfList {
			weight := v.Weight
			if v.Address == "" {
				continue
			}
			if v.Weight <= 0 {
				weight = 1
			}
			backendNodemap[v.Address] = weight
		}
		wrrpool, err := wrr.NewWrrBalancerPool(logger, backendNodemap)
		if err != nil {
			logger.Error("NewSmartReverseProxy - NewWrrBalancerPool error", zap.String("errmsg", err.Error()))
			return nil, err
		} else {
			spr.BackendNodePool = wrrpool
		}
	case enums.LB_CONSISTENT_HASH:
		return nil, errors.New("not support")
	case enums.LB_P2C:
		return nil, errors.New("not support")
	}

	if spr.IsGinOn {
		//use gin as a proxy
	} else {
		//use http.Server as a proxy
		//spr 实现了 ServeHTTP 方法，传给 http.Server 的 Handler，作为 ProxyServer
		spr.ProxyServer = &http.Server{
			Addr:    spr.ProxyAddress,
			Handler: spr}
	}

	return spr, nil
}

// start server with goroutine
func (s *SmartReverseProxy) Run() error {
	var err error
	if !s.IsTlsOn {
		go func() {
			if err = s.ProxyServer.ListenAndServe(); err != nil {
				s.Logger.Error("ListenAndServe error", zap.String("proxyname", s.ProxyName))

			}
		}()
	} else {
		//run https
		go func() {
			if err = s.ProxyServer.ListenAndServeTLS(s.CertFile, s.KeyFile); err != nil {
				s.Logger.Error("ListenAndServeTLS error", zap.String("proxyname", s.ProxyName))
			}
		}()
	}

	return err
}

// stop server
func (s *SmartReverseProxy) Shutdown() error {
	if err := s.ProxyServer.Shutdown(context.Background()); err != nil {
		s.Logger.Error("Shutdown error", zap.String("proxyname", s.ProxyName))
	}
	return nil
}

// 核心方法：分发前置请求到合适的后端节点
// r-- 原始请求
// w-- 回复客户端的响应
func (s *SmartReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//startTime := time.Now()
	origin_host := strings.ToLower(r.Host)
	//pre check

	// 原始请求中的 Host 必须等于代理设置的 ProxyName
	if origin_host != s.ProxyName {
		//sphttp.SmartProxyResponse(w, sphttp.ErrorHostNotMatch)
		//return
	}

	// 根据 lb 算法选择一个合适的 backend（ip+port）
	backend_address, err := s.GetBackendNodeWithLoadbalance(r)
	if err != nil {
		s.Logger.Error("GetBackendNodeWithLoadbalance pick next node error", zap.String("errmgs", err.Error()))
		sphttp.SmartProxyResponse(w, sphttp.ErrorNoneProperlyBackendNode)
		return
	}
	s.Logger.Info("SmartReverseProxy-GetBackendNode info", zap.String("backend_address", backend_address))

	// 选择指定的 httputil.ReverseProxy 处理请求
	rsp, err := s.GetRealReverseProxy(backend_address)
	if err != nil {
		s.Logger.Error("ServeHTTP-GetRealReverseProxy err", zap.String("errmsg", err.Error()))
		sphttp.SmartProxyResponse(w, sphttp.ErrorCreateReverseProxy)
		return
	}

	//forward requests
	rsp.ServeHTTP(w, r)
}

func (s *SmartReverseProxy) GetBackendNodeWithLoadbalance(r *http.Request) (string, error) {
	client_ip := sphttp.GetClientIP(r)
	s.Logger.Info("GetBackendNodeWithLoadbalance", zap.String("client_ip", client_ip))
	next_backend_node, err := s.BackendNodePool.Pick(client_ip)
	if err != nil {
		s.Logger.Error("GetBackendNodeWithLoadbalance pick next node error", zap.String("errmgs", err.Error()))
		return "", err
	}

	return next_backend_node.Addr, nil
}

// 根据后端地址 proxy_addr 选择（新建）reverseproxy
func (s *SmartReverseProxy) GetRealReverseProxy(proxy_addr string) (*httputil.ReverseProxy, error) {
	s.ReverseProxyMapLock.RLock()
	rsproxy, exists := s.ReverseProxyMap[proxy_addr]
	s.ReverseProxyMapLock.RUnlock()
	if !exists {
		// create a new reverse proxy
		if !strings.HasPrefix(proxy_addr, "http://") {
			proxy_addr = fmt.Sprintf("http://%s", proxy_addr)
		}
		target, err := url.Parse(proxy_addr)
		if err != nil {
			s.Logger.Error("GetRealReverseProxy url.Parse error", zap.String("errmsg", err.Error()))
			return nil, err
		}
		rsproxy = httputil.NewSingleHostReverseProxy(target)
		s.ReverseProxyMapLock.Lock()
		s.ReverseProxyMap[proxy_addr] = rsproxy
		s.ReverseProxyMapLock.Unlock()
	}
	return rsproxy, nil
}
