package healthy

//对Pool中的节点进行探测，探测失败则剔除

import (
	"context"
	"smart_proxy/backend"
	"smart_proxy/config"
	"smart_proxy/enums"
	sphttp "smart_proxy/pkg/http"
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
	go func() {
		statTicker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-statTicker.C:
				backend_nodes := h.ReverseGroup.GetAllProxys()
				for proxy_name, v := range backend_nodes {
					for _, sip := range v {
						checkret := sphttp.CheckTcpAlive(sip)
						if checkret {
							backendnode := backend.BackendNodeOperator{
								Target: backend.BackendNode{
									ProxyName: proxy_name,
									Addr:      sip},
								Op: enums.BACKEND_UP,
							}
							h.HealthyCheckChan <- backendnode
						}
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
