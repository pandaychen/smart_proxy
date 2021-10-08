package healthy

//对Pool中的节点进行探测，探测失败则剔除

import (
	"context"
	"smart_proxy/backend"
	"smart_proxy/config"
	"smart_proxy/enums"
	grpool "smart_proxy/pkg/pool"
	"smart_proxy/reverseproxy"
	"time"

	"go.uber.org/zap"
)

type HealthCheck struct {
	Logger           *zap.Logger
	ReverseGroup     *reverseproxy.SmartReverseProxyGroup
	HealthyCheckChan chan backend.BackendNodeOperator
}

func NewHealthCheck(logger *zap.Logger, smpconf *config.SmartProxyConfig, group *reverseproxy.SmartReverseProxyGroup, healthyCheckChan chan backend.BackendNodeOperator) *HealthCheck {
	return &HealthCheck{
		Logger:           logger,
		ReverseGroup:     group,
		HealthyCheckChan: healthyCheckChan,
	}
}

func (h *HealthCheck) Run(ctx context.Context) {
	h.Logger.Info("HealthCheck run ...")

	//初始化任务池
	healthCheckPool := grpool.NewSPool(10, TcpCheckTask)

	go func() {
		for taskret := range healthCheckPool.GetChanResult() {
			task, ok := taskret.OutputData.(Task)
			if !ok {
				h.Logger.Error("HealthCheck: get task output error", zap.Any("output", taskret))
				continue
			}
			var backendnode *backend.BackendNodeOperator
			if taskret.Err != nil {
				backendnode = &backend.BackendNodeOperator{
					Target: backend.BackendNode{
						ProxyName: task.Name,
						Addr:      task.Addr},
					Op: enums.BACKEND_DOWN,
				}
			} else {
				backendnode = &backend.BackendNodeOperator{
					Target: backend.BackendNode{
						ProxyName: task.Name,
						Addr:      task.Addr},
					Op: enums.BACKEND_UP,
				}
			}

			h.HealthyCheckChan <- *backendnode
		}
	}()

	go func() {
		statTicker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-statTicker.C:
				backend_nodes := h.ReverseGroup.GetAllProxys()
				var tasklist []grpool.TaskInput
				for proxy_name, v := range backend_nodes {
					for _, sip := range v {
						//generate all task
						tasklist = append(tasklist, grpool.TaskInput{
							InputData: Task{
								Name: proxy_name,
								Addr: sip,
							},
						})
					}
				}
				healthCheckPool.PoolWorkers(tasklist)
			case <-ctx.Done():
				return
			}
		}
	}()
}
