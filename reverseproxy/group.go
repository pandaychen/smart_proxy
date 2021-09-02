package reverseproxy

//代理分组（group）定义

// SmartReverseProxyGroup代表了一组ReverseProxy
// 一个SmartReverseProxy代表一个proxy服务
type SmartReverseProxyGroup struct {
	sync.RWMutex
	Logger *zap.Logger
	//支持不同的代理端口，以配置文件中的reverseproxy_group.name为唯一key
	Proxys map[string]*SmartReverseProxy
}
