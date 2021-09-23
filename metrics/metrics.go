package metrics

import (
	"fmt"
	"net/http"
	"smart_proxy/config"
	"smart_proxy/reverseproxy"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Metrics struct {
	Logger       *zap.Logger
	ReverseGroup *reverseproxy.SmartReverseProxyGroup
	BindAddr     string
}

func NewMetrics(logger *zap.Logger, smpconf *config.SmartProxyConfig, group *reverseproxy.SmartReverseProxyGroup) *Metrics {
	return &Metrics{
		Logger:       logger,
		ReverseGroup: group,
		BindAddr:     fmt.Sprintf("%s:%d", smpconf.MetricsConf.Host, smpconf.MetricsConf.Port),
	}
}

func (m *Metrics) Run() {
	m.Logger.Info("Metrics run ...")

	for _, v := range m.ReverseGroup.Proxys {
		//register collector
		prometheus.MustRegister(v.HttpRequestCount, v.HttpRequestDuration)
	}

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(m.BindAddr, nil)
		if err != nil {
			m.Logger.Error("Metrics ListenAndServe error", zap.Any("errmsg", err))
		}
	}()
}
