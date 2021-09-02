package healthy

//对Pool中的节点进行探测，探测失败则剔除

import (
	"smart_proxy/config"

	"go.uber.org/zap"
)

type HealthCheck struct {
	Logger *zap.Logger
}

func NewHealthCheck(logger *zap.Logger, smpconf *config.SmartProxyConfig) *HealthCheck {
	return &HealthCheck{
		Logger: logger,
	}
}

func (h *HealthCheck) Run() {
	h.Logger.Info("HealthCheck run ...")
}
