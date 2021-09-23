package scheduler

import (
	"context"
	"smart_proxy/backend"
	"smart_proxy/enums"
	"smart_proxy/reverseproxy"

	"go.uber.org/zap"
)

var GlobalScheduler *SmartProxyScheduler

// smartproxy controller
type SmartProxyScheduler struct {
	/* Backend operation channel */
	BackendChan  chan backend.BackendNodeOperator
	Logger       *zap.Logger
	ReverseGroup *reverseproxy.SmartReverseProxyGroup
}

func NewSmartProxyScheduler(logger *zap.Logger, reverse_group *reverseproxy.SmartReverseProxyGroup, dis2schChan chan backend.BackendNodeOperator) (*SmartProxyScheduler, error) {
	sch := &SmartProxyScheduler{
		Logger:       logger,
		ReverseGroup: reverse_group,
		BackendChan:  dis2schChan,
	}
	return sch, nil
}

func (s *SmartProxyScheduler) Run(ctx context.Context) {
	s.Logger.Info("Starting SchedulerLoopRun")

	//updates and manages backend nodes
	go func() {
		for {
			select {
			// handle backend operation
			case backend := <-s.BackendChan:
				s.Logger.Info("SchedulerLoopRun handle backendChan", zap.Any("backend", backend))
				s.ProcessBackendNodes(&backend)
			case <-ctx.Done():
				return
			}

		}
	}()
}

func (s *SmartProxyScheduler) ProcessBackendNodes(node *backend.BackendNodeOperator) {
	switch node.Op {
	case enums.BACKEND_DEL:
		proxy := s.ReverseGroup.GetReverseProxy(node.Target.ProxyName)
		if proxy != nil {
			proxy.BackendNodePool.RemoveNode(node.Target.Addr)
		}
	case enums.BACKEND_ADD:
		proxy := s.ReverseGroup.GetReverseProxy(node.Target.ProxyName)
		if proxy != nil {
			proxy.BackendNodePool.AddNode(node.Target.Addr, 1)
		}
	default:
		return
	}
}
