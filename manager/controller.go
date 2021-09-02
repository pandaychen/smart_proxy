package manager

import (
	"smart_proxy/config"

	"fmt"
	"go.uber.org/zap"
)

type Controller struct {
	BindAddr string
	Logger   *zap.Logger
}

// NewController 创建一个Controller
func NewController(logger *zap.Logger, smpconf *config.SmartProxyConfig) *Controller {
	return &Controller{
		Logger:   logger,
		BindAddr: fmt.Sprintf("%s:%d", smpconf.ControllerConf.Host, smpconf.ControllerConf.Port),
	}
}

func (c *Controller) Run() {
	c.Logger.Info("Controller run...")
}
