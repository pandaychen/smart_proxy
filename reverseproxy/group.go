package reverseproxy

//代理分组（group）定义
import (
	"smart_proxy/backend"
	"smart_proxy/config"
	"smart_proxy/pkg/errcode"
	"sync"

	"go.uber.org/zap"
)

// SmartReverseProxyGroup代表了一组ReverseProxy
// 一个SmartReverseProxy代表一个proxy服务
type SmartReverseProxyGroup struct {
	sync.RWMutex
	Logger *zap.Logger
	//支持不同的代理端口，以配置文件中的reverseproxy_group.name为唯一key
	Proxys       map[string]*SmartReverseProxy
	PeerStatChan chan backend.PeerStateOperator
}

// 初始化反向代理组
func NewSmartReverseProxyGroup(logger *zap.Logger, smpconf *config.SmartProxyConfig, peerStatChan chan backend.PeerStateOperator) (*SmartReverseProxyGroup, error) {
	group := &SmartReverseProxyGroup{
		Logger:       logger,
		Proxys:       make(map[string]*SmartReverseProxy),
		PeerStatChan: peerStatChan,
	}

	group_list := smpconf.ReverseProxyListConf
	//先根据静态配置初始化
	if group_list != nil {
		for _, rproxy_conf := range group_list {
			group.AddReverseProxy(&rproxy_conf)
		}
	}

	return group, nil
}

func (g *SmartReverseProxyGroup) AddReverseProxy(conf *config.ReverseProxyConfig) error {
	rproxy, err := g.NewSmartReverseProxy(g.Logger, conf)
	if err != nil {
		g.Logger.Error("AddReverseProxy -  NewSmartReverseProxy error", zap.String("err", err.Error()))
		return err
	}
	g.Lock()
	defer g.Unlock()
	g.Proxys[rproxy.ProxyName] = rproxy

	return nil
}

func (g *SmartReverseProxyGroup) GetReverseProxy(proxy_name string) *SmartReverseProxy {
	//停止proxy
	g.RLock()
	defer g.RUnlock()
	if _, exists := g.Proxys[proxy_name]; exists {
		return g.Proxys[proxy_name]
	}
	return nil
}

func (g *SmartReverseProxyGroup) DelReverseProxy(proxy_name string) {
	//停止proxy
	g.Lock()
	defer g.Unlock()
	delete(g.Proxys, proxy_name)
}

// Run：启动所有的proxy
func (g *SmartReverseProxyGroup) Run() {
	g.RLock()
	defer g.RUnlock()
	if len(g.Proxys) == 0 {
		panic(errcode.ErrNoneProxyNodes)
	}
	for name, rproxyserver := range g.Proxys {
		g.Logger.Info("SmartReverseProxyGroup Start...", zap.String("Proxy name", name))
		if err := rproxyserver.Run(); err != nil {
			g.Logger.Error("SmartReverseProxyGroup Run Error", zap.String("Proxy name", name), zap.String("errmsg", err.Error()))
		}
	}
}

// Stop: graceful shutdown
func (g *SmartReverseProxyGroup) Stop() {
	g.RLock()
	defer g.RUnlock()
	if len(g.Proxys) == 0 {
		return
	}
	for name, rproxyserver := range g.Proxys {
		g.Logger.Info("SmartReverseProxyGroup Stop...", zap.String("Proxy name", name))
		if err := rproxyserver.Shutdown(); err != nil {
			g.Logger.Error("SmartReverseProxyGroup Shutdown Error", zap.String("Proxy name", name), zap.String("errmsg", err.Error()))
		}
	}
}

func (g *SmartReverseProxyGroup) GetAllProxys() map[string][]string {
	var (
		proxymap = make(map[string][]string)
	)
	g.RLock()
	defer g.RUnlock()
	if len(g.Proxys) == 0 {
		return nil
	}
	for name, rproxyserver := range g.Proxys {
		backendlist := rproxyserver.BackendNodePool.GetAllNodes()
		if backendlist != nil {
			proxymap[name] = backendlist
		}
	}
	return proxymap
}
