package discovery

import (
	"smart_proxy/backend"
	"smart_proxy/config"
	"smart_proxy/discovery/etcd"
	"smart_proxy/enums"
	"smart_proxy/pkg/errcode"

	"go.uber.org/zap"
)

// DiscoveryClient 的通用结构封装（客户端）
type DiscoveryClient struct {
	Type       enums.DISCOVERY_TYPE
	Cluster    string
	RootPrefix string //Root 前缀
	DnsName    map[string]struct{}

	PasswordAuthOn bool
	Username       string
	Password       string
	Cert           string
	Key            string
	CommonName     string
	TrustedCaCert  string
	ResultChan     chan backend.BackendNodeOperator
	Logger         *zap.Logger
	RealClient     discoveryClient //saving client
}

//discoveryClient 模块是生产方
// 通过 discoveryClient 中解析出的后端 IP 需要通过结果 channel 发送给 scheduler 模块处理，由后者负责操作LB连接池
type discoveryClient interface {
	Run()
	Close()
}

// DiscoveryService 的通用结构封装（服务端）
type DiscoveryService struct {
	Address  string `json:"address"`
	Metadata string `json:"metadata"`
}

//discoveryService 通用的服务注册接口
type discoveryService interface {
	Register()
	Resolver()
	UnRegister()
}

func NewDiscoveryClient(logger *zap.Logger, smpconf *config.SmartProxyConfig, discovery2schedulerChan chan backend.BackendNodeOperator) (*DiscoveryClient, error) {
	var (
		client      *DiscoveryClient
		err         error
		wrapper_cli discoveryClient //uniform saving
	)

	client = &DiscoveryClient{
		//Type:       enums.DISCOVERY_TYPE(smpconf.DiscoveryConf.DiscoveryType),
		Cluster:    smpconf.DiscoveryConf.ClusterAddr,
		RootPrefix: smpconf.DiscoveryConf.RootPrefix,
		//DnsName    map[string]struct{}
		//PasswordAuthOn bool
		//Username       string
		//Password       string
		//Cert           string
		//Key            string
		//CommonName     string
		//TrustedCaCert  string
		ResultChan: discovery2schedulerChan,
		Logger:     logger,
	}

	switch smpconf.DiscoveryConf.DiscoveryType {
	case string(enums.ETCD_DISCOVERY):
		//etcd
		wrapper_cli, _ = etcd.NewEtcdDiscoveryClient(&etcd.EtcdConfig{
			Cluster:    client.Cluster,
			RootPrefix: client.RootPrefix,
			ResultChan: discovery2schedulerChan,
			Logger:     logger,
		})
		client.RealClient = wrapper_cli
	case string(enums.DNS_DISCOVERY):
		//dns
	default:
		panic(errcode.ErrNotSupportedDiscoveryType)
	}
	if err != nil {
		logger.Error("NewDiscoveryClient error", zap.Any("errmsg", err))
		return nil, err
	}

	client.RealClient = wrapper_cli
	return client, nil
}

func (c *DiscoveryClient) Run() {
	c.Logger.Info("DiscoveryClient run...", zap.Any("dtype", c.Type))
	go c.Run()
}
