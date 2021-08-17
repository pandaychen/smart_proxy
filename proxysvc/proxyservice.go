package proxysvc

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	zaplog "github.com/pandaychen/goes-wrapper/zaplog"
	"smart_proxy/backend"
	"smart_proxy/enums"
	"go.uber.org/zap"
)

const (
	DEFAULT_SERVICE_NAME = "smart_proxy"
)

type TlsConfig struct {
	CertFile string
	KeyFile  string
}

// ProxyService for loadbalance to backend pools
type SmartProxyReverse struct {
	sync.RWMutex
	TlsConfig
	IsTlsOn      bool
	IsGinOn      bool //使用gin做反向代理
	ProxyName    string
	ProxyAddress string //ip:port
	TimeOut      time.Duration

	LoadBalancerName enums.LB_TYPE
	DiscoveryName    enums.DISCOVERY_TYPE

	IsSafeHttpSig bool //是否加入http转发签名

	BackendNodePool backend.BackendNodePool //存储后端池

	//mapping requests to backends，请求会根据某种策略被代理到不同的后端，key为后端地址
	ReverseProxyMap     map[string]*httputil.ReverseProxy
	ReverseProxyMapLock sync.RWMutex

	SPSServer    *http.Server //self http hanlder
	SPSGinServer *gin.Engine

	Logger *zap.Logger
}

// 初始化
func NewSmartProxyReverse(options ...SmartProxyReverseOption) (*SmartProxyReverse, error) {
	logger, _ := zaplog.ZapLoggerInit(DEFAULT_SERVICE_NAME)
	sps := &SmartProxyReverse{
		ReverseProxyMap: make(map[string]*httputil.ReverseProxy),
		Logger:          logger,
	}

	for _, opt := range options {
		if err := opt(sps); err != nil {
			sps.Logger.Error("init config error", zap.String("errmsg", err.Error()))
			return nil, err
		}
	}

	if sps.IsGinOn {
		//GIN 代理
		return sps, nil
	}

	sps.SPSServer = &http.Server{
		Addr:    sps.ProxyAddress,
		Handler: sps} //HTTP handler

	return sps, nil
}

// 分发前置请求到合适的后端节点
func (s *SmartProxyReverse) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rsp, err := s.ChooseOneReverseProxy("")
	if err != nil {
		s.Logger.Error("ServeHTTP-ChooseOneReverseProxy err", zap.String("errmsg", err.Error()))
		return
	}

	//forward requests
	rsp.ServeHTTP(w, r)
}

func (s *SmartProxyReverse) ChooseOneReverseProxy(proxy_key string) (*httputil.ReverseProxy, error) {
	defer s.ReverseProxyMapLock.RUnlock()
	s.ReverseProxyMapLock.RLock()
	rsp, exists := s.ReverseProxyMap[proxy_key]
	if !exists {
		// create a new reverse proxy
		return nil, errors.New("error to find reverse proxy")
	}
	return rsp, nil
}
