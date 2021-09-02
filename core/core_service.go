package core

//lb 的主入口，一个 lb 包含三个子模块
//1.balancer 实现负载均衡逻辑
//2.controller 提供 cgiserver 及 restfulapi，直接操作 pool
//3.discovery，提供从第三方注册中心，获取在线列表以及实时监控后端服务地址变化

import (
	"os"
	"os/signal"
	"smart_proxy/backend"
	"smart_proxy/config"

	//etcdw "smart_proxy/discovery/etcd"
	"smart_proxy/discovery"
	"smart_proxy/healthy"
	"smart_proxy/manager"
	"smart_proxy/reverseproxy"
	"smart_proxy/scheduler"
	"syscall"

	"github.com/pandaychen/goes-wrapper/zaplog"

	"go.uber.org/zap"
)

// SmartProxyService 定义了反向代理的所有组件（集合）
type SmartProxyService struct {
	Logger            *zap.Logger
	ReverseproxyGroup *reverseproxy.SmartReverseProxyGroup // 提供反向代理 + 连接池（Pool）+ 负载均衡
	Scheduler         *scheduler.SmartProxyScheduler       // 核心调度器
	//Etcder            *etcdw.EtcdDiscoveryClient           // 服务发现模块
	Discoveryer   *discovery.DiscoveryClient
	Controller    *manager.Controller  // 提供后端增删查改 API 管理的 Restful-API 模块
	HealthChecker *healthy.HealthCheck // 健康检查

	//channel
	Discovery2SchedulerChan chan backend.BackendNodeOperator
}

// 创建 SmartProxyService 的所有组件
func NewSmartProxyService(proxy_config *config.SmartProxyConfig) (*SmartProxyService, error) {
	logger, err := zaplog.ZapLoggerInit(proxy_config.ProjectName)
	if err != nil {
		panic(err)
	}

	logger.Info("NewSmartProxyService init...")

	sps := &SmartProxyService{
		Logger:                  logger,
		Discovery2SchedulerChan: make(chan backend.BackendNodeOperator),
	}

	//Init all submodules
	sps.Controller = manager.NewController(logger, proxy_config)

	sps.Discoveryer, _ = discovery.NewDiscoveryClient(logger, proxy_config, sps.Discovery2SchedulerChan)

	sps.HealthChecker = healthy.NewHealthCheck(logger, proxy_config)

	sps.Scheduler, _ = scheduler.NewSmartProxyScheduler(logger, sps.Discovery2SchedulerChan)

	sps.ReverseproxyGroup, _ = reverseproxy.NewSmartReverseProxyGroup(logger, proxy_config)

	return sps, nil
}

// 启动 SmartProxyService 的所有子组件
// 需要注意组件的启动顺序不能调整！
func (s *SmartProxyService) RunLoop() error {
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt, syscall.SIGTERM)

	//start scheduler
	s.Scheduler.SchedulerLoopRun()
	//start discovery
	s.Discoveryer.RealClient.Run()
	//start controller
	s.Controller.Run()

	//start healthychecking
	s.HealthChecker.Run()
	//start reverseproxy

	<-sigC
	return nil
}
