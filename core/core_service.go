package core

//DONE
//lb 的主入口，一个 lb 包含三个子模块
//1.balancer 实现负载均衡逻辑
//2.controller 提供 cgiserver 及 restfulapi，直接操作 pool
//3.discovery，提供从第三方注册中心，获取在线列表以及实时监控后端服务地址变化

import (
	"os"
	"os/signal"
	etcdw "smart_proxy/discovery/etcd"
	"smart_proxy/scheduler"
	"syscall"

	"github.com/pandaychen/goes-wrapper/zaplog"

	"go.uber.org/zap"
)

// SmartProxyService 定义了反向代理的所有组件
type SmartProxyService struct {
	Scheduler *scheduler.SmartProxyScheduler // 核心调度器
	Logger    *zap.Logger
	Etcder    *etcdw.EtcdDiscoveryClient
}

// 创建 SmartProxyService 的所有组件
func NewSmartProxyService() (*SmartProxyService, error) {
	logger, err := zaplog.ZapLoggerInit("smart_proxy")
	if err != nil {
		panic(err)
	}
	scheduler_svc, _ := scheduler.NewSmartProxyScheduler(logger)
	etcder, _ := etcdw.NewEtcdDiscoveryClient(logger, scheduler_svc, "/test/")

	svc := &SmartProxyService{
		Scheduler: scheduler_svc,
		Logger:    logger,
		Etcder:    etcder,
	}

	return svc, nil
}

// 启动 SmartProxyService 的所有子组件
func (s *SmartProxyService) RunLoop() error {
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt, syscall.SIGTERM)

	s.Scheduler.SchedulerLoopRun()
	s.Etcder.Run()

	<-sigC
	return nil
}
